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
		t.Errorf("The error messaage of UnKnownFileTypeError does not contain the expected Path")
	}
}

func TestOutFileIsDirErrorStruct(t *testing.T) {
	path := "/some/sample/path"
	err := newOutFileIsDirError(path)

	if err.Dir != path {
		t.Errorf("The Dir was %s, but %s was expected", err.Dir, path)
	}

	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error messaage of OutFileIsDirError does not contain the expected Dir")
	}
}

func TestDepthParametrNotKnownErrorStruct(t *testing.T) {
	val := "asdgfg"
	err := newDepthParametrNotKnownError(val)

	if err.GivenValue != val {
		t.Errorf("The GivenValue was %s, but %s was expected", err.GivenValue, val)
	}

	if strings.Contains(err.Error(), val) == false {
		t.Errorf("The error messaage of DepthParametrNotKnownError does not contain the expected GivenValue")
	}
}
