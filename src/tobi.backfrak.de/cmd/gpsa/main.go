package main

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"fmt"

	"tobi.backfrak.de/internal/gpsabl"
	"tobi.backfrak.de/internal/gpxbl"
)

func main() {
	gpx, err := gpxbl.ReadGPX("")
	if err != nil {
		gpsabl.HandleError(err)
	}
	fmt.Println(gpx.Name)
	fmt.Println(gpx.Description)
	// os.Exit(0)
}
