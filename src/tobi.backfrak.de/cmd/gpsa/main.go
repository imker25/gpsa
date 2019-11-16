package main

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"fmt"
	"os"

	"tobi.backfrak.de/internal/gpsabl"
	"tobi.backfrak.de/internal/gpxbl"
)

func main() {

	if cap(os.Args) > 1 {

		fmt.Println("Argument given:")
		for _, arg := range os.Args[1:] { // Skip the 0s argument, becaus this will always be the program itselfe
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
			fmt.Println("Distance: ", info.GetDistance())
			fmt.Println("AtituteRange: ", info.GetAtituteRange())
			fmt.Println("MinimumAtitute: ", info.GetMinimumAtitute())
			fmt.Println("MaximumAtitute: ", info.GetMaximumAtitute())
		}
	}

}
