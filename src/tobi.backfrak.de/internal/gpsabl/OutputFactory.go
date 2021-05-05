package gpsabl

import "os"

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// GetOutputFormater - Get the OutoutFormater that matches the given output configuration
func GetOutputFormater(validFormaters []OutputFormater, outFile os.File, formaterType OutputFormaterType) (bool, OutputFormater) {

	for _, validFormater := range validFormaters {
		if outFile == *os.Stdout {
			if validFormater.CheckOutputFormaterType(formaterType) {
				return true, validFormater.NewOutputFormater()
			}
		} else {
			if validFormater.CheckFileExtension(outFile.Name()) {
				return true, validFormater.NewOutputFormater()
			}
		}
	}

	return false, nil
}
