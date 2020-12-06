package main

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

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

// HelpFlag - Tells if the program was called with -help
var HelpFlag bool

// VerboseFlag - Tells if the program was called with -verbose
var VerboseFlag bool

// SkipErrorExitFlag - Tells if the program was called with -skip-error-exit
var SkipErrorExitFlag bool

// PrintCsvHeaderFlag - Tells if the program was called with -print-csv-header
var PrintCsvHeaderFlag bool

// OutFileParameter - Tells if and where we should write the output to ( -out-file )
var OutFileParameter string

// DontPanicFlag - Tells if the program was called with -dont-panic
var DontPanicFlag bool

// DepthParameter - Tells for which depth we should perform the analyses ( -depth )
var DepthParameter string

// SummaryParameter - Tells if we should add summary to the output ( -summary )
var SummaryParameter string

// CorrectionParameter - Tells how we should correct Elevation data from the track ( -correction )
var CorrectionParameter string

// PrintVersionFlag - Tells if the program was called with the -version flag
var PrintVersionFlag bool

// PrintLicenseFlag - Tells if the program was called with the -license flag
var PrintLicenseFlag bool

// SuppressDuplicateOutPutFlag - Tells if the program was called with the -suppressDuplicateOutPut flag
var SuppressDuplicateOutPutFlag bool

// MinimalMovingSpeedParameter - Tells the minimal moving speed for moving time and speed calculation
var MinimalMovingSpeedParameter float64

// MinimalStepHightParameter - Tells the minimal step hight, when "steps" correction is used
var MinimalStepHightParameter float64

// PrintElevationOverDistanceFlag - Tell if the program was called with the -print-elevation-over-distance flag
var PrintElevationOverDistanceFlag bool

func main() {

	if cap(os.Args) > 1 {

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

		// If we don't have input files, we do nothing
		if len(flag.Args()) != 0 {
			// Find out where to write the output. May be a file or STDOUT
			out := getOutPutStream()

			// Get the type that handles the output
			iFormater := getOutPutFormater()
			defer out.Close()

			// Process the files, this will fill the buffer of the output type
			successCount := processFiles(flag.Args(), iFormater)

			// Write the output
			errWrite := iFormater.WriteOutput(out, SummaryParameter)
			if errWrite != nil {
				HandleError(errWrite, OutFileParameter, false, DontPanicFlag)
			}

			if VerboseFlag == true {
				fmt.Fprintln(os.Stdout, fmt.Sprintf("%d of %d files processed successfully", successCount, len(flag.Args())))
			}
		} else {
			if VerboseFlag == true {
				fmt.Fprintln(os.Stdout, "No input files given")
			}
		}

	}

	if ErrorsHandled == false {
		os.Exit(0)
	} else {
		fmt.Fprintln(os.Stderr, "At least one error occurred")
		os.Exit(-1)
	}
}

// handleComandlineOptions - Setup and parse the commandline options.
// Defines the usage function as well
func handleComandlineOptions() {

	outFormater := gpsabl.NewCsvOutputFormater(OutputSeperator, false)

	// Setup the valid comandline flags
	flag.Float64Var(&MinimalStepHightParameter, "minimal-step-hight", 10.0, "The minimal step hight. Only in use when \"steps\"  elevation correction is used. In [m]")
	flag.Float64Var(&MinimalMovingSpeedParameter, "minimal-moving-speed", 0.3, "The minimal speed. Distances traveled with less speed are not counted. In [m/s]")
	flag.BoolVar(&SuppressDuplicateOutPutFlag, "suppress-duplicate-out-put", false, "Suppress the output of duplicate lines. Duplicates are detected by timestamps. Output with non valid time data may still contains duplicates.")
	flag.BoolVar(&HelpFlag, "help", false, "Print help message and exit")
	flag.BoolVar(&PrintVersionFlag, "version", false, "Print version of the program and exit")
	flag.BoolVar(&PrintLicenseFlag, "license", false, "Print license information of the program and exit")
	flag.BoolVar(&VerboseFlag, "verbose", false, "Run the program with verbose output")
	flag.BoolVar(&SkipErrorExitFlag, "skip-error-exit", false, "Don't exit the program on track file processing errors")
	flag.BoolVar(&PrintCsvHeaderFlag, "print-csv-header", true, "Print out a csv header line. Possible values are [true false]")
	flag.StringVar(&OutFileParameter, "out-file", "", "Decide where to write the output. StdOut is used when not explicitly set")
	flag.BoolVar(&DontPanicFlag, "dont-panic", true, "Decide if the program will exit with panic or with negative exit code in error cases. Possible values are [true false]")
	flag.StringVar(&DepthParameter, "depth", outFormater.ValidDepthArgs[0],
		fmt.Sprintf("Define the way the program should analyse the files. Possible values are [%s]", outFormater.GetValidDepthArgsString()))
	flag.StringVar(&CorrectionParameter, "correction", gpsabl.GetValidCorrectionParameters()[2],
		fmt.Sprintf("Define how to correct the elevation data read in from the track. Possible values are [%s]", gpsabl.GetValidCorrectionParametersString()))
	flag.BoolVar(&PrintElevationOverDistanceFlag, "print-elevation-over-distance", false, "Tell if \"ElevationOverDistance.csv\" should be created for each track. The files will be locate in tmp dir.")
	flag.StringVar(&SummaryParameter, "summary", outFormater.ValidSummaryArgs[0],
		fmt.Sprintf("Define the way the program should analyse the files. Possible values are [%s]", outFormater.GetValidSummaryArgsString()))
	// Overwrite the std Usage function with some custom stuff
	flag.Usage = customHelpMessage

	// Read the given flags
	flag.Parse()
}

func customHelpMessage() {
	fmt.Fprintln(os.Stdout, fmt.Sprintf("%s: Reads in GPS track files, and writes out basic statistic data found in the track as a CSV style report", os.Args[0]))
	fmt.Fprintln(os.Stdout, fmt.Sprintf("Program %s", getVersion()))
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, fmt.Sprintf("Usage: %s [options] [files]", os.Args[0]))
	fmt.Fprintln(os.Stdout, "  files")
	fmt.Fprintln(os.Stdout, "        One or more track files (only *.gpx and *.tcx) supported at the moment")
	fmt.Fprintln(os.Stdout, "Options:")
	flag.PrintDefaults()
}

// processFiles - processes the input files and adds the found content to the output buffer
func processFiles(files []string, iFormater gpsabl.OutputFormater) int {

	if !gpsabl.CheckValidCorrectionParameters(CorrectionParameter) {
		HandleError(gpsabl.NewCorrectionParameterNotKnownError(CorrectionParameter), "", false, DontPanicFlag)
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
	file, readErr := reader.ReadTracks(CorrectionParameter, MinimalMovingSpeedParameter, MinimalStepHightParameter)
	if HandleError(readErr, filePath, SkipErrorExitFlag, DontPanicFlag) == true {
		return false
	}

	// Add the file to the out buffer of the formater
	addErr := formater.AddOutPut(file, DepthParameter, SuppressDuplicateOutPutFlag)
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
		printErr := gpsabl.WriteElevationOverDistance(file, out, OutputSeperator)
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
func getOutPutFormater() gpsabl.OutputFormater {
	formater := gpsabl.NewCsvOutputFormater(OutputSeperator, PrintCsvHeaderFlag)
	if !formater.CheckValidDepthArg(DepthParameter) {
		HandleError(gpsabl.NewDepthParameterNotKnownError(DepthParameter), "", false, DontPanicFlag)
	}
	if !formater.CheckValidSummaryArg(SummaryParameter) {
		HandleError(gpsabl.NewSummaryParamaterNotKnown(SummaryParameter), "", false, DontPanicFlag)
	}
	iFormater := gpsabl.OutputFormater(formater)

	return iFormater
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

// Prints the version string
func printVersion() {
	fmt.Fprintln(os.Stdout, getVersion())
}

// Get the version string
func getVersion() string {
	return fmt.Sprintf("Version: %s", version)
}
