package main

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"tobi.backfrak.de/internal/gpsabl"
	"tobi.backfrak.de/internal/gpxbl"
)

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

func main() {

	if cap(os.Args) > 1 {

		handleComandlineOptions()

		if HelpFlag == true {
			flag.Usage()
			os.Exit(0)
		}

		successCount, report := processFiles(flag.Args())
		for _, ret := range report {
			fmt.Fprintln(os.Stdout, ret)
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

// handleComandlineOptions - Setup and parse the comandline options.
// Defines the usage function as well
func handleComandlineOptions() {
	flag.BoolVar(&HelpFlag, "help", false, "Prints this message")
	flag.BoolVar(&VerboseFlag, "verbose", false, "Run the programm with verbose output")
	flag.BoolVar(&SkipErrorExitFlag, "skip-error-exit", false, "Don't exit the programm with first error")

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
	c := make(chan string)
	countFiles := 0
	retVals := []string{}

	for _, filePath := range files {
		go goProcessFile(filePath, c)
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

func goProcessFile(filePath string, c chan string) {
	ret := processFile(filePath)
	c <- ret
}

func processFile(filePath string) string {
	retVal := ""
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
	info := gpsabl.TrackSummaryProvider(&file)

	// Read Properties from the TrackFile
	retVal = retVal + "File name: " + file.FilePath + getNewLine()
	retVal = retVal + "Name: " + file.Name + getNewLine()
	retVal = retVal + "Description: " + file.Description + getNewLine()
	retVal = retVal + "NumberOfTracks: " + fmt.Sprintf(" %d ", file.NumberOfTracks) + getNewLine()

	// Read properties troutgh the interface
	retVal = retVal + "Distance:" + fmt.Sprintf(" %f ", info.GetDistance()) + "m" + getNewLine()
	retVal = retVal + "AtituteRange:" + fmt.Sprintf(" %f ", info.GetAtituteRange()) + "m" + getNewLine()
	retVal = retVal + "MinimumAtitute:" + fmt.Sprintf(" %f ", info.GetMinimumAtitute()) + "m" + getNewLine()
	retVal = retVal + "MaximumAtitute:" + fmt.Sprintf(" %f ", info.GetMaximumAtitute()) + "m" + getNewLine()

	return retVal
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

func getNewLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	} else {
		return "\n"
	}

}
