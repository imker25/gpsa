package tcxbl

import (
	"strings"
	"testing"
)

func TestGpxFileErrorStruct(t *testing.T) {

	path := "/some/sample/path"
	err := newTcxFileError(path)
	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error message of GpxFileError does not contain the expected Path")
	}

	if err.File != path {
		t.Errorf("The GpxFileError.File does not match the expected value")
	}
}

func TestEmptyTcxFileError(t *testing.T) {
	path := "/some/sample/path"
	err := newEmptyTcxFileError(path)
	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error message of EmptyGpxFileError does not contain the expected Path")
	}

	if err.File != path {
		t.Errorf("The EmptyGpxFileError.File does not match the expected value")
	}
}
