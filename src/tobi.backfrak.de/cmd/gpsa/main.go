package main

// Copyright 2019 by Tobias Zellner. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"fmt"
	"os"

	"tobi.backfrak.de/internal/gpsabl"
)

func main() {
	gpx, err := gpsabl.ReadGPX("")
	if err != nil {
		gpsabl.HandleError(err)
	}
	fmt.Println(gpx.Name)
	fmt.Println(gpx.Description)
	os.Exit(0)
}
