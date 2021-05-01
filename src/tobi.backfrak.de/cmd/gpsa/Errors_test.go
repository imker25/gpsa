package main

import (
	"strings"
	"testing"
)

func TestUnKnownFileTypeErrorStruct(t *testing.T) {
	path := "/some/sample/path"
	err := newUnKnownFileTypeError(path)

	if err.File != path {
		t.Errorf("The File was %s, but %s was expected", err.File, path)
	}

	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error message of UnKnownFileTypeError does not contain the expected path")
	}
}

func TestOutFileIsDirErrorStruct(t *testing.T) {
	path := "/some/sample/path"
	err := newOutFileIsDirError(path)

	if err.Dir != path {
		t.Errorf("The directory was %s, but %s was expected", err.Dir, path)
	}

	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error message of OutFileIsDirError does not contain the expected directory")
	}
}

func TestUnKnownInputStreamError(t *testing.T) {
	path := "<{}>"
	err := newUnKnownInputStreamError(path)

	if err.Line != path {
		t.Errorf("The directory was %s, but %s was expected", err.Line, path)
	}

	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error message of UnKnownInputStreamError does not contain the expected line")
	}
}
