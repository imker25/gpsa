package main

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import "fmt"

// OutFileIsDirError - Error when trying to write the output to a directory and not a file
type OutFileIsDirError struct {
	err string
	// File - The path to the dir that caused this error
	Dir string
}

func (e *OutFileIsDirError) Error() string { // Implement the Error Interface for the OutFileIsDirError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// newOutFileIsDirError- Get a new OutFileIsDirError struct
func newOutFileIsDirError(dirName string) *OutFileIsDirError {
	return &OutFileIsDirError{fmt.Sprintf("The given -out-file \"%s\" is a directory.", dirName), dirName}
}

// UnKnownFileTypeError - Error when trying to load not known file type
type UnKnownFileTypeError struct {
	err string
	// File - The path to the file that caused this error
	File string
}

func (e *UnKnownFileTypeError) Error() string { // Implement the Error Interface for the UnKnownFileTypeError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// newUnKnownFileTypeError - Get a new UnKnownFileTypeError struct
func newUnKnownFileTypeError(fileName string) *UnKnownFileTypeError {
	return &UnKnownFileTypeError{fmt.Sprintf("The type of the file \"%s\" is not known.", fileName), fileName}
}
