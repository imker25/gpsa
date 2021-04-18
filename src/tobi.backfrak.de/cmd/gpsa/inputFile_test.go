package main

import "testing"

func TestNewInputFileWithPath(t *testing.T) {
	path := "my/test/path"
	sut := newInputFileWithPath(path)

	if sut.Name != path {
		t.Errorf("The inputFile.Name is %s, but %s is expected", sut.Name, path)
	}

	if sut.Type != FilePath {
		t.Errorf("The inputFile.Type is %s, but %s is expected", sut.Type, FilePath)
	}

	if sut.Buffer != nil {
		t.Errorf("The inputFile..Buffer is not nil")
	}

	if sut.inputFileTypeValid() == false {
		t.Errorf("The inputFileTypeValid tells that inputFile.Type %s is not valide", sut.Type)
	}

}

func TestNewInputFileGpxBuffer(t *testing.T) {
	name := "stream buffer 1"
	buffer := []byte{1, 2, 3, 4, 5, 6}
	sut := newInputFileGpxBuffer(buffer, name)

	if sut.Name != name {
		t.Errorf("The inputFile.Name is %s, but %s is expected", sut.Name, name)
	}

	if sut.Type != GpxBuffer {
		t.Errorf("The inputFile.Type is %s, but %s is expected", sut.Type, GpxBuffer)
	}

	if sut.Buffer[3] != buffer[3] {
		t.Errorf("The inputFile.Buffer is has not the expected value")
	}

	if sut.inputFileTypeValid() == false {
		t.Errorf("The inputFileTypeValid tells that inputFile.Type %s is not valide", sut.Type)
	}

}

func TestNewInputFileTcxBuffer(t *testing.T) {
	name := "stream buffer 1"
	buffer := []byte{1, 2, 3, 4, 5, 6}
	sut := newInputFileTcxBuffer(buffer, name)

	if sut.Name != name {
		t.Errorf("The inputFile.Name is %s, but %s is expected", sut.Name, name)
	}

	if sut.Type != TcxBuffer {
		t.Errorf("The inputFile.Type is %s, but %s is expected", sut.Type, TcxBuffer)
	}

	if sut.Buffer[3] != buffer[3] {
		t.Errorf("The inputFile.Buffer is has not the expected value")
	}

	if sut.inputFileTypeValid() == false {
		t.Errorf("The inputFileTypeValid tells that inputFile.Type %s is not valide", sut.Type)
	}

}

func TestInputFileTypeValid(t *testing.T) {
	ft := inputFileType("myString")

	if inputFileTypeValid(ft) == true {
		t.Errorf("The inputFileTypeValid tells that  %s is a valid type", ft)
	}
}
