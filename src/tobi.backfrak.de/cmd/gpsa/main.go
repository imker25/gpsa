package main

// Copyright 2019 by Tobias Zellner. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"os"

	"tobi.backfrak.de/internal/gpsabl"
)

func main() {
	gpsabl.ReadGPX("")

	os.Exit(0)
}
