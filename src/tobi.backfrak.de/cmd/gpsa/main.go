package main

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"flag"
	"fmt"
	"os"

	"tobi.backfrak.de/internal/gpsabl"
	"tobi.backfrak.de/internal/gpxbl"
)

var HelpFlag bool

func main() {

	if cap(os.Args) > 1 {

		handleComandlineOptions()

		if HelpFlag == true {
			flag.Usage()
			os.Exit(0)
		}

		for _, arg := range flag.Args() {
			fmt.Println("Read file: " + arg)
			// Get the GpxFile type
			gpx := gpxbl.NewGpxFile(arg)

			// Convert the GpxFile to the TrackReader interface
			reader := gpsabl.TrackReader(&gpx)

			// Read the *.gpx into a TrackFile type, using the interface
			file, err := reader.ReadTracks()

			if err != nil {
				HandleError(err, arg)
			}
			// Convert the TrackFile into the TrackInfoProvider interface
			info := gpsabl.TrackInfoProvider(file)

			// Read Properties from the TrackFile
			fmt.Println("Name: ", file.Name)
			fmt.Println("Description: ", file.Description)

			// Read Properties from the GpxFile
			fmt.Println("NumberOfTracks: ", gpx.NumberOfTracks)

			// Read properties troutgh the interface
			fmt.Println("Distance: ", info.GetDistance(), "m")
			fmt.Println("AtituteRange: ", info.GetAtituteRange(), "m")
			fmt.Println("MinimumAtitute: ", info.GetMinimumAtitute(), "m")
			fmt.Println("MaximumAtitute: ", info.GetMaximumAtitute(), "m")
		}

		os.Exit(0)
	}

}

// handleComandlineOptions - Setup and parse the Comandline Options
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
