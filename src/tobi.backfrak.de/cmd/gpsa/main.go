package main

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"tobi.backfrak.de/internal/jsonbl"

	"tobi.backfrak.de/internal/csvbl"
	"tobi.backfrak.de/internal/gpsabl"
	"tobi.backfrak.de/internal/gpxbl"
	"tobi.backfrak.de/internal/tcxbl"
)

// Authors - Information about the authors of the program. You might want to add your name here when contributing to this software
const Authors = "tobi@backfrak.de"

// OutputSeperator - The seperator string for csv output files
const OutputSeperator = "; "

// The version of this program, will be set at compile time by the gradle build script
var version = "undefined"

func main() {

	var fileArgs []string

	// Setup and read-in comandline flags
	handleComandlineOptions()

	// Check flags, that will not process files
	if VerboseFlag {
		args := ""
		for _, arg := range os.Args {
			args = fmt.Sprintf("%s %s", args, arg)
		}
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Call: %s", args))
		if !PrintVersionFlag {
			printVersion()
		}
	}

	if HelpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if PrintVersionFlag {
		printVersion()
		os.Exit(0)
	}

	if PrintLicenseFlag {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("(c) %s - Apache License, Version 2.0( http://www.apache.org/licenses/LICENSE-2.0 )", Authors))
		os.Exit(0)
	}

	// If we don't have input files, we might run with stream input
	if len(flag.Args()) != 0 {
		fileArgs = flag.Args()
	} else {
		// No file input, but might a stream
		fileArgs = processInputStream()
	}

	// There might be files to process
	if len(fileArgs) != 0 {
		// Find out where to write the output. May be a file or STDOUT
		out := getOutPutStream()

		// Get the type that handles the output
		iFormater := getOutPutFormater(*out)
		defer out.Close()

		// Process the files, this will fill the buffer of the output type
		successCount := processFiles(fileArgs, iFormater)

		// Write the output
		errWrite := iFormater.WriteOutput(out, gpsabl.SummaryArg(SummaryParameter))
		if errWrite != nil {
			HandleError(errWrite, OutFileParameter, false, DontPanicFlag)
		}

		if VerboseFlag == true {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%d of %d files processed successfully", successCount, len(fileArgs)))
		}
	} else {
		// No files to process
		if VerboseFlag == true {
			fmt.Fprintln(os.Stdout, "No input files given")
		}
	}

	if ErrorsHandled == false {
		os.Exit(0)
	} else {
		fmt.Fprintln(os.Stderr, "At least one error occurred")
		os.Exit(-1)
	}
}

func processInputStream() []string {

	var fileArgs []string
	// Get stdin stream
	info, errStat := os.Stdin.Stat()
	if errStat != nil {
		HandleError(errStat, "", false, DontPanicFlag)
	}

	// Check if stdin gets data in a pipe
	if info.Mode()&os.ModeNamedPipe != 0 {
		// pipe
		if VerboseFlag == true {
			fmt.Fprintln(os.Stdout, "Input is given as a stream")
		}

		reader := bufio.NewReader(os.Stdin)
		var err error
		fileArgs, err = ReadInputStreamBuffer(reader)
		if err != nil {
			HandleError(err, "os.Stdin", false, DontPanicFlag)
		}
	}

	return fileArgs
}

// processFiles - processes the input files and adds the found content to the output buffer
func processFiles(files []string, iFormater gpsabl.OutputFormater) int {

	if !gpsabl.CheckValidCorrectionParameters(gpsabl.CorrectionParameter(CorrectionParameter)) {
		HandleError(gpsabl.NewCorrectionParameterNotKnownError(gpsabl.CorrectionParameter(CorrectionParameter)), "", false, DontPanicFlag)
	}

	allFiles := len(files)
	successCount := 0
	c := make(chan bool, allFiles)
	countFiles := 0

	// Process the files in a go routine
	for _, filePath := range files {
		go goProcessFile(filePath, iFormater, c)
	}

	// Read back the file processing results
	for ret := range c {
		if ret != false {
			successCount++
		}
		countFiles++
		if countFiles == allFiles {
			close(c)
		}
	}

	// Return how may files were processed fine
	return successCount
}

// goProcessFile - Wraper around, processFile. Use this as go routine
func goProcessFile(filePath string, formater gpsabl.OutputFormater, c chan bool) {
	ret := processFile(filePath, formater)

	c <- ret
}

// processFile - processes one input file and adds the found content to the output buffer
func processFile(filePath string, formater gpsabl.OutputFormater) bool {
	if VerboseFlag == true {
		fmt.Println("Read file: " + filePath)
	}
	// Find out if we can read the file
	reader, readerErr := getReader(filePath)
	if HandleError(readerErr, filePath, SkipErrorExitFlag, DontPanicFlag) == true {
		return false
	}

	// Read the *.gpx into a TrackFile type, using the interface
	file, readErr := reader.ReadTracks(gpsabl.CorrectionParameter(CorrectionParameter), MinimalMovingSpeedParameter, MinimalStepHightParameter)
	if HandleError(readErr, filePath, SkipErrorExitFlag, DontPanicFlag) == true {
		return false
	}

	// Add the file to the out buffer of the formater
	addErr := formater.AddOutPut(file, gpsabl.DepthArg(DepthParameter), SuppressDuplicateOutPutFlag)
	if HandleError(addErr, filePath, SkipErrorExitFlag, DontPanicFlag) == true {
		return false
	}

	if PrintElevationOverDistanceFlag {

		// Get the path of the ElevationOverDistance.csv
		outPath := getElevationOverDistanceFileName(file)
		out, createErr := os.Create(outPath)
		if HandleError(createErr, outPath, SkipErrorExitFlag, DontPanicFlag) == true {
			return false
		}

		// Write the ElevationOverDistance.csv
		fmt.Println(fmt.Sprintf("Create %s", outPath))
		printErr := csvbl.WriteElevationOverDistance(file, out, OutputSeperator)
		if HandleError(printErr, outPath, SkipErrorExitFlag, DontPanicFlag) == true {
			return false
		}
	}

	return true
}

func getElevationOverDistanceFileName(file gpsabl.TrackFile) string {

	dir := os.TempDir()
	_, oldName := filepath.Split(file.FilePath)

	fileName := oldName + ".ElevationOverDistance.csv"

	return path.Join(dir, fileName)
}

// Get the Interface to format the output
func getOutPutFormater(outFile os.File) gpsabl.OutputFormater {
	var iFormater gpsabl.OutputFormater
	if !gpsabl.CheckValidDepthArg(DepthParameter) {
		HandleError(gpsabl.NewDepthParameterNotKnownError(gpsabl.DepthArg(DepthParameter)), "", false, DontPanicFlag)
	}
	if !gpsabl.CheckValidSummaryArg(SummaryParameter) {
		HandleError(gpsabl.NewSummaryParamaterNotKnown(gpsabl.SummaryArg(SummaryParameter)), "", false, DontPanicFlag)
	}

	if outFile == *os.Stdout {
		if checkStdOutFormatParameterValue(StdOutFormatParameter) == false {
			HandleError(newUnKnownFileTypeError(StdOutFormatParameter), "", false, DontPanicFlag)
		}

		if strings.ToLower(StdOutFormatParameter) == strings.ToLower(stdOutFormatParameterValues[0]) {
			iFormater = getCsvOutputFormater()
		}
		if strings.ToLower(StdOutFormatParameter) == strings.ToLower(stdOutFormatParameterValues[1]) {
			jsonFormater := jsonbl.NewJSONOutputFormater()
			iFormater = gpsabl.OutputFormater(jsonFormater)
		}
	}

	if strings.HasSuffix(strings.ToLower(outFile.Name()), ".csv") {
		iFormater = getCsvOutputFormater()
	}

	if strings.HasSuffix(strings.ToLower(outFile.Name()), ".json") {
		jsonFormater := jsonbl.NewJSONOutputFormater()
		iFormater = gpsabl.OutputFormater(jsonFormater)
	}

	if iFormater == nil {
		HandleError(newUnKnownFileTypeError(outFile.Name()), "", false, DontPanicFlag)
	}
	return iFormater
}

func getCsvOutputFormater() *csvbl.CsvOutputFormater {
	csvFormater := csvbl.NewCsvOutputFormater(OutputSeperator, PrintCsvHeaderFlag)
	if !csvbl.CheckTimeFormatIsValid(TimeFormatParameter) {
		HandleError(csvbl.NewTimeFormatNotKnown(csvbl.TimeFormat(TimeFormatParameter)), "", false, DontPanicFlag)
	} else {
		csvFormater.SetTimeFormat(TimeFormatParameter)
	}
	return csvFormater
}

// Get the file interface we are using as output. Maybe a file or STDOUT
func getOutPutStream() *os.File {
	var out *os.File
	var errOpen error
	var errCreate error
	if OutFileParameter == "" {
		out = os.Stdout
	} else {
		if outFileExists(OutFileParameter) {
			errDel := os.Remove(OutFileParameter)
			if errDel != nil {
				HandleError(errDel, OutFileParameter, false, DontPanicFlag)
			}
		}
		out, errCreate = os.Create(OutFileParameter)
		if errCreate != nil {
			HandleError(errCreate, OutFileParameter, false, DontPanicFlag)
		}

		out, errOpen = os.OpenFile(OutFileParameter, os.O_APPEND|os.O_WRONLY, 0600)
		if errOpen != nil {
			HandleError(errOpen, OutFileParameter, false, DontPanicFlag)
		}
	}
	return out
}

// Get the interface that can read a given input file
func getReader(file string) (gpsabl.TrackReader, error) {

	if strings.HasSuffix(file, "gpx") == true { // If the file is a *.gpx, we can read it
		return getGpxReader(file), nil
	}

	if strings.HasSuffix(file, "tcx") == true { // If the file is a *.tcx, we can read it
		return getTcxReader(file), nil
	}

	// We dont know the file type
	return nil, newUnKnownFileTypeError(file)
}

// Get the interface that can read a *.gpx file
func getGpxReader(file string) gpsabl.TrackReader {
	// Get the GpxFile type
	gpx := gpxbl.NewGpxFile(file)

	// Convert the GpxFile to the TrackReader interface
	return gpsabl.TrackReader(&gpx)
}

// Get the interface that can read a *.tcx file
func getTcxReader(file string) gpsabl.TrackReader {
	// Get the TcxFile type
	tcx := tcxbl.NewTcxFile(file)

	// Convert the TcxFile to the TrackReader interface
	return gpsabl.TrackReader(&tcx)
}

// outFileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func outFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	if info.IsDir() {
		err := newOutFileIsDirError(filename)
		HandleError(err, filename, false, DontPanicFlag)
	}
	return true
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// Prints the version string
func printVersion() {
	fmt.Fprintln(os.Stdout, getVersion())
}

// Get the version string
func getVersion() string {
	return fmt.Sprintf("Version: %s", version)
}
