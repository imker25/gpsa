package main

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"fmt"
	"os"

	"tobi.backfrak.de/internal/gpxbl"
)

func main() {

	if cap(os.Args) > 1 {

		fmt.Println("Argument given:")
		for _, arg := range os.Args[1:] { // Skip the 0s argument, becaus this will always be the program itselfe
			// fmt.Println(arg)
			gpx := gpxbl.NewGpxFile(arg)
			file, err := gpx.ReadTracks()
			if err != nil {
				HandleError(err)
			}
			fmt.Println("Name: ", file.Name)
			fmt.Println("Description: ", file.Description)
			fmt.Println("NumberOfTracks: ", file.NumberOfTracks)
			fmt.Println("Distance: ", file.Distance)
			fmt.Println("AtituteRange: ", file.AtituteRange)
			fmt.Println("MinimumAtitute: ", file.MinimumAtitute)
			fmt.Println("MaximumAtitute: ", file.MaximumAtitute)
		}
	}

}
