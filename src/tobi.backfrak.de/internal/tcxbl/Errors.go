package tcxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import "fmt"

// TcxFileError - Error when trying to load something that is no gpx file
type TcxFileError struct {
	err string
	// File - The path to the file that caused this error
	File string
}

func (e *TcxFileError) Error() string { // Implement the Error Interface for the TcxFileError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// newTcxFileError - Get a new TcxFileError struct
func newTcxFileError(fileName string) *TcxFileError {
	return &TcxFileError{fmt.Sprintf("The file \"%s\" is not a gpx file", fileName), fileName}
}

// EmptyTcxFileError - Error when trying to load a gpx file that does not contain any valid track
type EmptyTcxFileError struct {
	err string
	// File - The path to the file that caused this error
	File string
}

func (e *EmptyTcxFileError) Error() string { // Implement the Error Interface for the EmptyTcxFileError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// newEmptyTcxFileError - Get a new GpxFileError struct
func newEmptyTcxFileError(fileName string) *EmptyTcxFileError {
	return &EmptyTcxFileError{fmt.Sprintf("The file \"%s\" does not contain any valid tracks.", fileName), fileName}
}
