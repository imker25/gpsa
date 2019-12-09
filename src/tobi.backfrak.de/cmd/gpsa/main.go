package main

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"tobi.backfrak.de/internal/gpsabl"
	"tobi.backfrak.de/internal/gpxbl"
)

// Authors - Information about the authors of the program. You might want to add your name here when contribute to this software
const Authors = "tobi@backfrak.de"

// The Version of this program, will be set at compile time by the gradle build script
var version = "undefined"

// HelpFlag - Tell if the program was called with -help
var HelpFlag bool

// VerboseFlag - Tell if the program was called with -verbose
var VerboseFlag bool

// SkipErrorExitFlag - Tell if the program was called with -skip-error-exit
var SkipErrorExitFlag bool

// PrintCsvHeaderFlag - Tell if the program was called with  -print-csv-header
var PrintCsvHeaderFlag bool

// OutFileParameter - Tell if and where we should write the output to ( -out-file )
var OutFileParameter string

// DontPanicFlag - Tell if the prgramm was called with -dont-panic
var DontPanicFlag bool

// DepthParametr - Tell in what depth we should mak the analyses ( -depth )
var DepthParametr string

// CorrectionParameter - Tell how we should correct Elevation data from the track ( -correction )
var CorrectionParameter string

// PrintVersionFlag - tell if the program was called with the -version flag
var PrintVersionFlag bool

// PrintLicenseFlag - tell if the program was called with the -license flag
var PrintLicenseFlag bool

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

		// If don't have intput files, we do nothing
		if len(flag.Args()) != 0 {
			// Find out where to write the output. May a file, may STDOUT
			out := getOutPutStream()

			// Get the type that handles the output
			iFormater := getOutPutFormater()
			defer out.Close()

			// Process the files, this will fill the buffer of the output type
			successCount := processFiles(flag.Args(), iFormater)

			// Write the output
			errWrite := iFormater.WriteOutput(out)
			if errWrite != nil {
				HandleError(errWrite, OutFileParameter, false, DontPanicFlag)
			}

			if VerboseFlag == true {
				fmt.Fprintln(os.Stdout, fmt.Sprintf("%d of %d files process successfull", successCount, len(flag.Args())))
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
		fmt.Fprintln(os.Stderr, "At least one error occured")
		os.Exit(-1)
	}
}

// handleComandlineOptions - Setup and parse the comandline options.
// Defines the usage function as well
func handleComandlineOptions() {

	outFormater := gpsabl.NewCsvOutputFormater(";")

	// Setup the valide comandline flags
	flag.BoolVar(&HelpFlag, "help", false, "Prints this help message and exit")
	flag.BoolVar(&PrintVersionFlag, "version", false, "Print the version of the program and exit")
	flag.BoolVar(&PrintLicenseFlag, "license", false, "Print the license information of the program and exit")
	flag.BoolVar(&VerboseFlag, "verbose", false, "Run the program with verbose output")
	flag.BoolVar(&SkipErrorExitFlag, "skip-error-exit", false, "Don't exit the program on track file processing errors")
	flag.BoolVar(&PrintCsvHeaderFlag, "print-csv-header", true, "Print out a csv header line")
	flag.StringVar(&OutFileParameter, "out-file", "", "Tell where to write the output. StdOut is used when not set")
	flag.BoolVar(&DontPanicFlag, "dont-panic", true, "Tell if the prgramm will exit with panic, or with negiatv exit code in error cases")
	flag.StringVar(&DepthParametr, "depth", outFormater.ValideDepthArgs[0],
		fmt.Sprintf("Tell how depth the program should analyse the files. Possible values are [%s]", outFormater.GetVlaideDepthArgsString()))
	flag.StringVar(&CorrectionParameter, "correction", gpsabl.GetValideCorectionParamters()[2],
		fmt.Sprintf("Tell how the programm should correct the elevation data read from the track. Possible values are [%s]", gpsabl.GetValideCorectionParamtersString()))

	// Overwrite the std Usage function with some costum stuff
	flag.Usage = costumHelpMessage

	// Read the given flags
	flag.Parse()
}

func costumHelpMessage() {
	fmt.Fprintln(os.Stdout, fmt.Sprintf("%s: Reads in GPS track files, and writes out basic statistic data found in the track as a CSV style report", os.Args[0]))
	fmt.Fprintln(os.Stdout, fmt.Sprintf("Program %s", getVersion()))
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, fmt.Sprintf("Usage: %s [options] [files]", os.Args[0]))
	fmt.Fprintln(os.Stdout, "  files")
	fmt.Fprintln(os.Stdout, "        One or more track files (only *.gpx) supported at the moment")
	fmt.Fprintln(os.Stdout, "Options:")
	flag.PrintDefaults()
}

// processFiles - prosses the input files and add the found coneted to the output buffer
func processFiles(files []string, iFormater gpsabl.OutputFormater) int {

	if !gpsabl.CheckValideCorectionParamters(CorrectionParameter) {
		HandleError(gpsabl.NewCorectionParamterNotKnownError(CorrectionParameter), "", false, DontPanicFlag)
	}

	allFiles := len(files)
	successCount := 0
	c := make(chan bool, allFiles)
	countFiles := 0

	// Add the header to the output, when needed
	if PrintCsvHeaderFlag {
		iFormater.AddHeader()
	}

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

	// Return how may files was processed fine
	return successCount
}

// goProcessFile - Wraper around, processFile. Use this as go routine
func goProcessFile(filePath string, formater gpsabl.OutputFormater, c chan bool) {
	ret := processFile(filePath, formater)

	c <- ret
}

// processFile - prosess one input file and add the found coneted to the output buffer
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
	file, readErr := reader.ReadTracks(CorrectionParameter)
	if HandleError(readErr, filePath, SkipErrorExitFlag, DontPanicFlag) == true {
		return false
	}

	// Add the file to the out buffer of the formater
	addErr := formater.AddOutPut(file, DepthParametr)
	if HandleError(addErr, filePath, SkipErrorExitFlag, DontPanicFlag) == true {
		return false
	}

	return true
}

// Get the Interface to format the output
func getOutPutFormater() gpsabl.OutputFormater {
	formater := gpsabl.NewCsvOutputFormater(";")
	if !formater.CheckVlaideDepthArg(DepthParametr) {
		HandleError(gpsabl.NewDepthParametrNotKnownError(DepthParametr), "", false, DontPanicFlag)
	}
	iFormater := gpsabl.OutputFormater(formater)

	return iFormater
}

// Get the file interface, we use as output. A file or STDOUT
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

// Get the interface, that can read a given input file
func getReader(file string) (gpsabl.TrackReader, error) {

	if strings.HasSuffix(file, "gpx") == true { // If the file is a *.gpx, we can read it
		return getGpxReader(file), nil
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
