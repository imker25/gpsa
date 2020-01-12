package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"strings"
	"testing"
)

func TestGpxFileErrorStruct(t *testing.T) {

	path := "/some/sample/path"
	err := newGpxFileError(path)
	checkGpxFileError(err, path, t)
}

func TestEmptyGpxFileError(t *testing.T) {
	path := "/some/sample/path"
	err := newEmptyGpxFileError(path)
	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error messaage of EmptyGpxFileError does not contain the expected Path")
	}

	if err.File != path {
		t.Errorf("The EmptyGpxFileError.File does not match the expected value")
	}
}

func checkGpxFileError(err *GpxFileError, path string, t *testing.T) {
	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error messaage of GpxFileError does not contain the expected Path")
	}

	if err.File != path {
		t.Errorf("The GpxFileError.File does not match the expected value")
	}
}
