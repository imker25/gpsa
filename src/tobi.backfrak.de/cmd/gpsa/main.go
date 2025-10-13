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

	"tobi.backfrak.de/internal/gpsabl"

	"tobi.backfrak.de/internal/csvbl"
	"tobi.backfrak.de/internal/jsonbl"
	"tobi.backfrak.de/internal/mdbl"

	"tobi.backfrak.de/internal/gpxbl"
	"tobi.backfrak.de/internal/tcxbl"
)

// Authors - Information about the authors of the program. You might want to add your name here when contributing to this software
const Authors = "tobi@backfrak.de"

// OutputSeperator - The seperator string for csv output files
const OutputSeperator = "; "

// The version of this program, will be set at compile time by the gradle build script
var version = "undefined"

var ValidReaders = []gpsabl.TrackReader{&gpxbl.GpxFile{}, &tcxbl.TcxFile{}}
var ValidFormaters = []gpsabl.OutputFormater{&csvbl.CsvOutputFormater{}, &jsonbl.JSONOutputFormater{}, &mdbl.MDOutputFormater{}}
var DefinedFilters = []gpsabl.TrackFilter{}

func main() {

	var fileArgs []gpsabl.InputFile

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
		fileArgs = proccessFileArgs(flag.Args())
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
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%d of %d files processed successfully.", successCount, len(fileArgs)))
		}

		// In case there are no entries in the outfile we remove it from disk
		if iFormater.GetNumberOfOutputEntries() <= 0 {
			deleteOutFile(out)
			if VerboseFlag == true {
				fmt.Fprintln(os.Stdout, fmt.Sprintf("Output is empty. No output file created."))
			}
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

func createFilters() bool {
	if MinStartTime != "" {
		minFilter, minFilterErr := gpsabl.NewMinStartTimeFilter(MinStartTime)
		if minFilterErr != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Can not parse the string \"%s\" as date or date time", MinStartTime))
			return false
		}

		DefinedFilters = append(DefinedFilters, &minFilter)
	}

	if MaxStartTime != "" {
		maxFilter, maxFilterErr := gpsabl.NewMaxStartTimeFilter(MaxStartTime)
		if maxFilterErr != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Can not parse the string \"%s\" as date or date time", MaxStartTime))
			return false
		}

		DefinedFilters = append(DefinedFilters, &maxFilter)
	}

	return true
}

func proccessFileArgs(args []string) []gpsabl.InputFile {
	var fileArgs []gpsabl.InputFile
	for _, file := range args {
		res, input := gpsabl.GetInputFileFromPath(ValidReaders, file)
		if res == true {
			fileArgs = append(fileArgs, input)
		} else {
			HandleError(newUnKnownFileTypeError(file), file, SkipErrorExitFlag, DontPanicFlag)
		}
	}

	return fileArgs
}

// Get the input files from a stream buffer
func processInputStream() []gpsabl.InputFile {

	var fileArgs []gpsabl.InputFile
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
func processFiles(files []gpsabl.InputFile, iFormater gpsabl.OutputFormater) int {

	if !gpsabl.CheckValidCorrectionParameters(gpsabl.CorrectionParameter(CorrectionParameter)) {
		HandleError(gpsabl.NewCorrectionParameterNotKnownError(gpsabl.CorrectionParameter(CorrectionParameter)), "", false, DontPanicFlag)
	}

	if !createFilters() {
		os.Exit(-10)
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
func goProcessFile(file gpsabl.InputFile, formater gpsabl.OutputFormater, c chan bool) {
	ret := processFile(file, formater)

	c <- ret
}

// processFile - processes one input file and adds the found content to the output buffer
func processFile(inFile gpsabl.InputFile, formater gpsabl.OutputFormater) bool {
	if VerboseFlag == true {
		fmt.Println("Read file: " + inFile.Name)
	}
	var file gpsabl.TrackFile
	var readErr error

	// Find out if we can read the file
	reader, readerErr := getReader(inFile)
	if HandleError(readerErr, inFile.Name, SkipErrorExitFlag, DontPanicFlag) == true {
		return false
	}

	// Read the *.gpx into a TrackFile type, using the interface
	file, readErr = reader.ReadTracks(gpsabl.CorrectionParameter(CorrectionParameter), MinimalMovingSpeedParameter, MinimalStepHightParameter)

	if HandleError(readErr, inFile.Name, SkipErrorExitFlag, DontPanicFlag) == true {
		return false
	}

	// Filter the track
	if len(DefinedFilters) > 0 {
		// TODO Write Tests for this functionality
		file = gpsabl.FilterTrackFile(file, DefinedFilters)
	}

	// Add the file to the out buffer of the formater, if it contains tracks
	if len(file.Tracks) < 1 {
		if VerboseFlag {
			fmt.Println(fmt.Sprintf("File \"%s\" does not contain any tracks after applying the given filters", inFile.Name))
		}
		return true
	}

	addErr := formater.AddOutPut(file, gpsabl.DepthArg(DepthParameter), SuppressDuplicateOutPutFlag)
	if HandleError(addErr, inFile.Name, SkipErrorExitFlag, DontPanicFlag) == true {
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
	var res bool
	if !gpsabl.CheckValidDepthArg(DepthParameter) {
		HandleError(gpsabl.NewDepthParameterNotKnownError(gpsabl.DepthArg(DepthParameter)), "", false, DontPanicFlag)
	}
	if !gpsabl.CheckValidSummaryArg(SummaryParameter) {
		HandleError(gpsabl.NewSummaryParamaterNotKnown(gpsabl.SummaryArg(SummaryParameter)), "", false, DontPanicFlag)
	}
	if outFile != *os.Stdout {
		if !checkOutFileExtension(outFile.Name()) {
			HandleError(newUnKnownFileTypeError(outFile.Name()), "", false, DontPanicFlag)
		}
	} else {
		if !checkOutFileType(StdOutFormatParameter) {
			HandleError(newUnKnownFileTypeError(StdOutFormatParameter), "", false, DontPanicFlag)
		}
	}

	res, iFormater = gpsabl.GetOutputFormater(ValidFormaters, outFile, gpsabl.OutputFormaterType(strings.ToUpper(StdOutFormatParameter)))
	if res == false {
		HandleError(newUnKnownFileTypeError(outFile.Name()), "", false, DontPanicFlag)
	}
	switch iFormater.(type) {
	case gpsabl.TextOutputFormater:
		iFormater = setTextutputFormater(iFormater.GetTextOutputFormater())
	}
	if iFormater == nil {
		HandleError(newUnKnownFileTypeError(outFile.Name()), "", false, DontPanicFlag)
	}
	return iFormater
}

func checkOutFileExtension(filePath string) bool {
	for _, formater := range ValidFormaters {
		if formater.CheckFileExtension(filePath) {
			return true
		}
	}

	return false
}

func checkOutFileType(fileType string) bool {
	for _, formater := range ValidFormaters {
		if formater.CheckOutputFormaterType(gpsabl.OutputFormaterType(strings.ToUpper(fileType))) {
			return true
		}
	}

	return false
}

func setTextutputFormater(formater gpsabl.TextOutputFormater) gpsabl.OutputFormater {
	if !formater.CheckTimeFormatIsValid(TimeFormatParameter) {
		HandleError(gpsabl.NewTimeFormatNotKnown(gpsabl.TimeFormat(TimeFormatParameter)), "", false, DontPanicFlag)
	} else {
		formater.SetTimeFormat(TimeFormatParameter)
	}
	formater.SetAddHeader(PrintCsvHeaderFlag)
	formater.SetSeperator(OutputSeperator)

	return formater
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

func deleteOutFile(outFile *os.File) error {

	// In case the output is Stdout there is no temp file that needs deletion
	if outFile != os.Stdout {
		errClose := outFile.Close()
		if errClose != nil {
			return errClose
		}
		errDelete := os.Remove(outFile.Name())
		if errDelete != nil {
			return errDelete
		}
	}

	return nil
}

// Get the interface that can read a given input file
func getReader(file gpsabl.InputFile) (gpsabl.TrackReader, error) {

	res, reader := gpsabl.GetNewReader(ValidReaders, file)
	if res == true {
		return reader, nil
	}

	// We dont know the file type
	return nil, newUnKnownFileTypeError(file.Name)
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

// Prints the version string
func printVersion() {
	fmt.Fprintln(os.Stdout, getVersion())
}

// Get the version string
func getVersion() string {
	return fmt.Sprintf("Version: %s", version)
}
