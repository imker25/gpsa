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

// HandleError - Handles an error. Will do nothing if the given error is nil
func HandleError(err error, filePath string) {
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: The given track file \"%s\" was not found.", filePath))
		case *xml.SyntaxError:
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: The given track file \"%s\" is not well formated: %s", filePath, err.Error()))
		case *gpxbl.GpxFileError:
			fmt.Fprintln(os.Stderr, err.Error())
		case *UnKnownFileTypeError:
			fmt.Fprintln(os.Stderr, err.Error())
		default:
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: %s", err.Error()))
		}

		os.Exit(-1)
	}
}
