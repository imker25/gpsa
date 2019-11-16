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
