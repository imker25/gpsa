package main

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"encoding/xml"
	"fmt"
	"os"

	"tobi.backfrak.de/internal/gpxbl"
)

// HandleError - Handles an error
func HandleError(err error) {
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			fmt.Println("Error: The given track file was not found: ", err.Error())
		case *xml.SyntaxError:
			fmt.Println("Error: The given track file is not well formated: ", err.Error())
		case *gpxbl.GpxFileError:
			fmt.Println("Error: The given track file is not a GPX file: ", err.Error())
		default:
			fmt.Println("Error: ", err.Error())
		}
	}

	os.Exit(-1)
}
