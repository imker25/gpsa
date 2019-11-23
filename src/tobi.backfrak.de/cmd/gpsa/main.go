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

// PrintVersionFlag - tell if the program was called with the -version flag
var PrintVersionFlag bool

// PrintLicenseFlag - tell if the program was called with the -license flag
var PrintLicenseFlag bool

func main() {

	if cap(os.Args) > 1 {

		handleComandlineOptions()

		if HelpFlag {
			flag.Usage()
			os.Exit(0)
		}

		if PrintVersionFlag {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("Version: %s", version))
			os.Exit(0)
		}

		if PrintLicenseFlag {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("(c) %s - Apache License, Version 2.0( http://www.apache.org/licenses/LICENSE-2.0 )", Authors))
			os.Exit(0)
		}

		out := getOutPutStream()
		defer out.Close()

		successCount, report := processFiles(flag.Args())
		for _, ret := range report {
			_, errWrite := out.WriteString(ret)
			if errWrite != nil {
				HandleError(errWrite, OutFileParameter, false, DontPanicFlag)
			}
		}
		if VerboseFlag == true {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%d of %d files process successfull", successCount, len(flag.Args())))
		}

	}

	if ErrorsHandled == false {
		os.Exit(0)
	} else {
		fmt.Fprintln(os.Stderr, "At least one error occured")
		os.Exit(-1)
	}
}

func getOutPutStream() *os.File {
	var out *os.File
	var errOpen error
	var errCreate error
	if OutFileParameter == "" {
		out = os.Stdout
	} else {
		if fileExists(OutFileParameter) {
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

// handleComandlineOptions - Setup and parse the comandline options.
// Defines the usage function as well
func handleComandlineOptions() {

	outFormater := gpsabl.NewCsvOutputFormater(";")

	flag.BoolVar(&HelpFlag, "help", false, "Prints this message")
	flag.BoolVar(&PrintVersionFlag, "version", false, "Print the version of the program")
	flag.BoolVar(&PrintLicenseFlag, "license", false, "Print the license information of the program")
	flag.BoolVar(&VerboseFlag, "verbose", false, "Run the program with verbose output")
	flag.BoolVar(&SkipErrorExitFlag, "skip-error-exit", false, "Don't exit the program on track file processing errors")
	flag.BoolVar(&PrintCsvHeaderFlag, "print-csv-header", true, "Print out a csv header line")
	flag.StringVar(&OutFileParameter, "out-file", "", "Tell where to write the output. StdOut is used when not set")
	flag.BoolVar(&DontPanicFlag, "dont-panic", true, "Tell if the prgramm will exit with panic, or with negiatv exit code in error cases")
	flag.StringVar(&DepthParametr, "depth", outFormater.ValideDepthArgs[0],
		fmt.Sprintf("Tell how depth the program should analyse the files. Possible values are [%s]", outFormater.GetVlaideDepthArgsString()))

	// Overwrite the std Usage function with some costum stuff
	flag.Usage = func() {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("%s: Reads in GPS track files, and writes out basic statistic data found in the track", os.Args[0]))
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Program Version: %s", version))
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Usage: %s [options] [files]", os.Args[0]))
		fmt.Fprintln(os.Stdout, "  files")
		fmt.Fprintln(os.Stdout, "        One or more track files (only *.gpx) supported at the moment")
		fmt.Fprintln(os.Stdout, "Options:")
		flag.PrintDefaults()
	}
	// fmt.Println("Call: ", os.Args)
	flag.Parse()

	if !strings.Contains(outFormater.GetVlaideDepthArgsString(), DepthParametr) {
		HandleError(newDepthParametrNotKnownError(DepthParametr), "", false, DontPanicFlag)
	}
}

func processFiles(files []string) (int, []string) {
	allFiles := len(files)
	successCount := 0
	c := make(chan string, allFiles)
	countFiles := 0
	retVals := []string{}
	formater := gpsabl.NewCsvOutputFormater(";")
	if PrintCsvHeaderFlag {
		retVals = append(retVals, formater.GetHeader())
	}

	for _, filePath := range files {
		go goProcessFile(filePath, formater, c)
	}

	for ret := range c {
		if ret != "" {
			successCount++
			retVals = append(retVals, ret)
		}
		countFiles++
		if countFiles == allFiles {
			close(c)
		}
	}
	return successCount, retVals
}

func goProcessFile(filePath string, formater gpsabl.CsvOutputFormater, c chan string) {
	rets := processFile(filePath, formater)

	for _, ret := range rets {
		c <- ret
	}
}

func processFile(filePath string, formater gpsabl.CsvOutputFormater) []string {
	if VerboseFlag == true {
		fmt.Println("Read file: " + filePath)
	}
	// Find out if we can read the file
	reader, readerErr := getReader(filePath)
	if HandleError(readerErr, filePath, SkipErrorExitFlag, DontPanicFlag) == true {
		return []string{""}
	}

	// Read the *.gpx into a TrackFile type, using the interface
	file, readErr := reader.ReadTracks()
	if HandleError(readErr, filePath, SkipErrorExitFlag, DontPanicFlag) == true {
		return []string{""}
	}

	// Convert the TrackFile into the TrackSummaryProvider interface
	// info := gpsabl.TrackSummaryProvider(file)

	return formater.FormatOutPut(file, false, DepthParametr)
}

func getReader(file string) (gpsabl.TrackReader, error) {

	if strings.HasSuffix(file, "gpx") == true { // If the file is a *.gpx, we can read it
		return getGpxReader(file), nil
	}

	// We dont know the file type
	return nil, newUnKnownFileTypeError(file)
}

func getGpxReader(file string) gpsabl.TrackReader {
	// Get the GpxFile type
	gpx := gpxbl.NewGpxFile(file)

	// Convert the GpxFile to the TrackReader interface
	return gpsabl.TrackReader(&gpx)
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
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
