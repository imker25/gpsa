package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"tobi.backfrak.de/internal/csvbl"
	"tobi.backfrak.de/internal/testhelper"
)

func TestReadInputStreamBufferWithFileList(t *testing.T) {

	file1, file2, read, err := getTwoValidInputFilePathStream()
	if err != nil {
		t.Fatal(err)
	}
	result, err := ReadInputStreamBuffer(bufio.NewReader(read))
	if err != nil {
		t.Errorf("Got error \"%s\" but expected none", err)
	}

	if len(result) != 2 {
		t.Errorf("The result list does contain %d files, but %d expected", len(result), 2)
	}

	if result[0] != file1 {
		t.Errorf("The path is %s, but %s is expected", result[0], file1)
	}

	if result[1] != file2 {
		t.Errorf("The path is %s, but %s is expected", result[1], file1)
	}

}

func TestReadInputStreamBufferWithNotExistingFileList(t *testing.T) {

	read, errStream := getInValidInputFilePathStream()
	if errStream != nil {
		t.Fatal(errStream)
	}

	result, err := ReadInputStreamBuffer(bufio.NewReader(read))
	if err == nil {
		t.Errorf("No error, but one expected")
	}

	if result != nil {
		t.Errorf("The file list should be empty")
	}
	switch err.(type) {
	case *UnKnownInputStreamError:
		fmt.Println("OK")
	default:
		t.Errorf("The error is not of the expected type.")
	}
}

func getInValidInputFilePathStream() (*os.File, error) {
	file1 := testhelper.GetValidGPX("12.gpx")
	filenotExist := "myNotExisting.gpx"
	file2 := testhelper.GetValidGPX("10.gpx")
	input := []byte(fmt.Sprintf("%s%s%s%s%s", file1, csvbl.GetNewLine(), filenotExist, csvbl.GetNewLine(), file2))
	read, write, errCreate := os.Pipe()
	if errCreate != nil {
		return nil, errCreate
	}

	_, errWrite := write.Write(input)
	if errWrite != nil {
		return nil, errWrite
	}
	write.Close()

	return read, nil
}

func getTwoValidInputFilePathStream() (string, string, *os.File, error) {
	file1 := testhelper.GetValidGPX("12.gpx")
	file2 := testhelper.GetValidGPX("10.gpx")
	input := []byte(fmt.Sprintf("%s%s%s", file1, csvbl.GetNewLine(), file2))
	read, write, errCreate := os.Pipe()
	if errCreate != nil {
		return "", "", nil, errCreate
	}

	_, errWrite := write.Write(input)
	if errWrite != nil {
		return "", "", nil, errWrite
	}
	write.Close()

	return file1, file2, read, nil
}
