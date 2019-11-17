package main

import (
	"strings"
	"testing"

	"tobi.backfrak.de/internal/testhelper"
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

func TestHandleErrorNil(t *testing.T) {
	if HandleError(nil, "my/path") == true {
		t.Errorf("HandleError reutrns true, when nil error was given")
	}

}

func TestHandleErrorNotNil(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	if HandleError(newUnKnownFileTypeError("my/path"), "my/path") == false {
		t.Errorf("HandleError reutrns false, when error was given")
	}

	if ErrorsHandled == false {
		t.Errorf("ErrorsHandled should be true, after a error was handeled")
	}
	SkipErrorExitFlag = oldFlagValue
	ErrorsHandled = false
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

func TestProcessValideFiles(t *testing.T) {
	ErrorsHandled = false
	files := []string{testhelper.GetValideGPX("01.gpx"), testhelper.GetValideGPX("02.gpx")}
	if processFiles(files) != 2 {
		t.Errorf("Not all files was proccess successfull as expected")
	}

	if ErrorsHandled == true {
		ErrorsHandled = false
		t.Errorf("Errors occured, but should not")
	}
}

func TestProcessMixedFiles(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	files := []string{testhelper.GetUnValideGPX("01.gpx"), testhelper.GetValideGPX("01.gpx"), testhelper.GetUnValideGPX("02.gpx")}
	if processFiles(files) != 1 {
		t.Errorf("Not two files was proccess with error as expected")
	}

	if ErrorsHandled == false {
		t.Errorf("No errors occured, but should")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
}

func TestProcessUnValideFiles(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	files := []string{testhelper.GetUnValideGPX("01.gpx"), testhelper.GetUnValideGPX("02.gpx")}
	if processFiles(files) != 0 {
		t.Errorf("Not all files was proccess with error as expected")
	}

	if ErrorsHandled == false {
		t.Errorf("No errors occured, but should")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
}
