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

func main() {

	if cap(os.Args) > 1 {

		handleComandlineOptions()

		if HelpFlag == true {
			flag.Usage()
			os.Exit(0)
		}

		processFiles(flag.Args())

	}
	os.Exit(0)
}

func processFiles(files []string) {
	for _, filePath := range files {
		fmt.Println("Read file: " + filePath)

		// Find out if we can read the file
		reader, readerErr := getReader(filePath)
		HandleError(readerErr, filePath)

		// Read the *.gpx into a TrackFile type, using the interface
		file, readErr := reader.ReadTracks()
		HandleError(readErr, filePath)

		// Convert the TrackFile into the TrackInfoProvider interface
		info := gpsabl.TrackInfoProvider(file)

		// Read Properties from the TrackFile
		fmt.Println("Name:", file.Name)
		fmt.Println("Description:", file.Description)

		// Read Properties from the GpxFile
		fmt.Println("NumberOfTracks:", file.NumberOfTracks)

		// Read properties troutgh the interface
		fmt.Println("Distance:", info.GetDistance(), "m")
		fmt.Println("AtituteRange:", info.GetAtituteRange(), "m")
		fmt.Println("MinimumAtitute:", info.GetMinimumAtitute(), "m")
		fmt.Println("MaximumAtitute:", info.GetMaximumAtitute(), "m")
	}
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

// handleComandlineOptions - Setup and parse the comandline options.
// Defines the usage function as well
func handleComandlineOptions() {
	flag.BoolVar(&HelpFlag, "help", false, "Prints this message")

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
