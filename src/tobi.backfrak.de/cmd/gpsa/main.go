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

// OutFileIsDirError - Error when trying to write the output to a directory and not a file
type OutFileIsDirError struct {
	err string
	// File - The path to the dir that caused this error
	Dir string
}

func (e *OutFileIsDirError) Error() string { // Implement the Error Interface for the OutFileIsDirError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// newOutFileIsDirError- Get a new OutFileIsDirError struct
func newOutFileIsDirError(dirName string) *OutFileIsDirError {
	return &OutFileIsDirError{fmt.Sprintf("The given -out-file \"%s\" is a directory.", dirName), dirName}
}

// UnKnownFileTypeError - Error when trying to load not known file type
type UnKnownFileTypeError struct {
	err string
	// File - The path to the file that caused this error
	File string
}

func (e *UnKnownFileTypeError) Error() string { // Implement the Error Interface for the UnKnownFileTypeError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// newUnKnownFileTypeError - Get a new UnKnownFileTypeError struct
func newUnKnownFileTypeError(fileName string) *UnKnownFileTypeError {
	return &UnKnownFileTypeError{fmt.Sprintf("The type of the file \"%s\" is not known.", fileName), fileName}
}

// HelpFlag - Tell if the programm was called with -help
var HelpFlag bool

// VerboseFlag - Tell if the programm was called with -verbose
var VerboseFlag bool

// SkipErrorExitFlag - Tell if the programm was called with -skip-error-exit
var SkipErrorExitFlag bool

// PrintCsvHeaderFlag - Tell if the programm was called with  -print-csv-header
var PrintCsvHeaderFlag bool

// OutFileParameter - Tell if and where we should write the output to ( -out-file )
var OutFileParameter string

func main() {

	if cap(os.Args) > 1 {

		handleComandlineOptions()

		if HelpFlag == true {
			flag.Usage()
			os.Exit(0)
		}

		out := getOutPutStream()
		defer out.Close()

		successCount, report := processFiles(flag.Args())
		for _, ret := range report {
			_, errWrite := out.WriteString(ret)
			if errWrite != nil {
				HandleError(errWrite, OutFileParameter)
				panic(errWrite)
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
				HandleError(errDel, OutFileParameter)
				panic(errDel)
			}
		}
		out, errCreate = os.Create(OutFileParameter)
		if errCreate != nil {
			HandleError(errCreate, OutFileParameter)
			panic(errCreate)
		}

		out, errOpen = os.OpenFile(OutFileParameter, os.O_APPEND|os.O_WRONLY, 0600)
		if errOpen != nil {
			HandleError(errOpen, OutFileParameter)
			panic(errOpen)
		}
	}
	return out
}

// handleComandlineOptions - Setup and parse the comandline options.
// Defines the usage function as well
func handleComandlineOptions() {
	flag.BoolVar(&HelpFlag, "help", false, "Prints this message")
	flag.BoolVar(&VerboseFlag, "verbose", false, "Run the programm with verbose output")
	flag.BoolVar(&SkipErrorExitFlag, "skip-error-exit", false, "Don't exit the programm with first error")
	flag.BoolVar(&PrintCsvHeaderFlag, "print-csv-header", true, "Print out a csv header line (default is true)")
	flag.StringVar(&OutFileParameter, "out-file", "", "Tell where to write the output. StdOut is used when not set")

	// Overwrite the std Usage function with some costum stuff
	flag.Usage = func() {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("%s: Reads in track files, and writes out basic statistic data found in the track", os.Args[0]))
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Usage: %s [options] [files]", os.Args[0]))
		fmt.Fprintln(os.Stdout, "  files")
		fmt.Fprintln(os.Stdout, "        One or more track files (only *.gpx) supported at the moment")
		fmt.Fprintln(os.Stdout, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

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
	ret := processFile(filePath, formater)
	c <- ret
}

func processFile(filePath string, formater gpsabl.CsvOutputFormater) string {
	if VerboseFlag == true {
		fmt.Println("Read file: " + filePath)
	}
	// Find out if we can read the file
	reader, readerErr := getReader(filePath)
	if HandleError(readerErr, filePath) == true {
		return ""
	}

	// Read the *.gpx into a TrackFile type, using the interface
	file, readErr := reader.ReadTracks()
	if HandleError(readErr, filePath) == true {
		return ""
	}

	// Convert the TrackFile into the TrackSummaryProvider interface
	info := gpsabl.TrackSummaryProvider(file)

	return formater.FormatTrackSummary(info, file.FilePath)
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
		HandleError(err, filename)
		panic(err)
	}
	return true
}
