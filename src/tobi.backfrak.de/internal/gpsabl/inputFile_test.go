package gpsabl

import "testing"

func TestNewInputFileWithPath(t *testing.T) {
	path := "my/test/path"
	sut := NewInputFileWithPath(path)

	if sut.Name != path {
		t.Errorf("The inputFile.Name is %s, but %s is expected", sut.Name, path)
	}

	if sut.Type != FilePath {
		t.Errorf("The inputFile.Type is %s, but %s is expected", sut.Type, FilePath)
	}

	if sut.Buffer != nil {
		t.Errorf("The inputFile..Buffer is not nil")
	}

	if sut.InputFileTypeValid() == false {
		t.Errorf("The inputFileTypeValid tells that inputFile.Type %s is not valide", sut.Type)
	}

}

func TestInputFileTypeValid(t *testing.T) {
	ft := InputFileType("myString")

	if inputFileTypeValid(ft) == true {
		t.Errorf("The inputFileTypeValid tells that  %s is a valid type", ft)
	}
}
