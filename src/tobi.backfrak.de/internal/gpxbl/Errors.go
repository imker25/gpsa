package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import "fmt"

// GpxFileError - Error when trying to load something that is no gpx file
type GpxFileError struct {
	err string
	// File - The path to the file that caused this error
	File string
}

func (e *GpxFileError) Error() string { // Implement the Error Interface for the GpxFileError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// newGpxFileError - Get a new GpxFileError struct
func newGpxFileError(fileName string) *GpxFileError {
	return &GpxFileError{fmt.Sprintf("The file \"%s\" is not a gpx file", fileName), fileName}
}

// EmptyGpxFileError - Error when trying to load a gpx file that does not contain any valid track
type EmptyGpxFileError struct {
	err string
	// File - The path to the file that caused this error
	File string
}

func (e *EmptyGpxFileError) Error() string { // Implement the Error Interface for the EmptyGpxFileError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// newEmptyGpxFileError - Get a new GpxFileError struct
func newEmptyGpxFileError(fileName string) *EmptyGpxFileError {
	return &EmptyGpxFileError{fmt.Sprintf("The file \"%s\" does not contain any valid tracks.", fileName), fileName}
}
