package main

import (
	"strings"
	"testing"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestHandleComandlineOptions(t *testing.T) {
	handleComandlineOptions()

	if HelpFlag == true {
		t.Errorf("The HelpFlag is true, but should not")
	}
}

func TestHandleError(t *testing.T) {
	HandleError(nil, "my/path")

}

func TestGetReaderGpxFile(t *testing.T) {
	reader, err := getReader("/some/track.gpx")

	if err != nil {
		t.Errorf("Got an error when try to get a reader for a gpx file: %s", err.Error())
	}

	if reader == nil {
		t.Errorf("The reader we got was nil")
	}
}

func TestGetReaderUnkonwnFile(t *testing.T) {
	reader, err := getReader("/some/track.txt")

	if err == nil {
		t.Errorf("Got no error when try to get a reader for a txt file.")
	}

	if reader != nil {
		t.Errorf("The reader we got was not nil")
	}
}

func TestUnKnownFileTypeErrorStruct(t *testing.T) {
	path := "/some/sample/path"
	err := newUnKnownFileTypeError(path)

	if err.File != path {
		t.Errorf("The File was %s, but %s was expected", err.File, path)
	}

	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error messaage of GpxFileError does not contain the expected Path")
	}
}
