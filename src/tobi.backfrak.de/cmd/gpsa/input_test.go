package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
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

	if result[0].Name != file1 {
		t.Errorf("The path is %s, but %s is expected", result[0], file1)
	}

	if result[0].Type != FilePath {
		t.Errorf("The type is %s, but %s is expected", result[0].Type, FilePath)
	}

	if result[1].Name != file2 {
		t.Errorf("The path is %s, but %s is expected", result[1], file1)
	}

	if result[1].Type != FilePath {
		t.Errorf("The type is %s, but %s is expected", result[1].Type, FilePath)
	}

	if result[1].Buffer != nil {
		t.Errorf("The buffer is expected to be nil")
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

func TestReadInputStreamBufferWithTwoGPXFileContent(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skip this test on windows")
	}

	read, errGet := getValidInputGPXContentStream()
	if errGet != nil {
		t.Fatal(errGet)
	}

	input, err := ReadInputStreamBuffer(bufio.NewReader(read))
	if err != nil {
		t.Errorf("No error, but one expected")
	}

	if len(input) != 2 {
		t.Errorf("The input has %d files, but %d files are expected", len(input), 2)
	}

	if input[0].Type != GpxBuffer {
		t.Errorf("The type is %s, but %s is expected", input[0].Type, GpxBuffer)
	}

	if input[0].Name == "" {
		t.Errorf(("The name is \"\""))
	}

	if input[0].Buffer == nil {
		t.Errorf("The buffer is nil")
	}

	if input[1].Type != GpxBuffer {
		t.Errorf("The type is %s, but %s is expected", input[1].Type, GpxBuffer)
	}

	if input[1].Buffer == nil {
		t.Errorf("The buffer is nil")
	}

	if input[0].Name == input[1].Name {
		t.Errorf(("The names are the same"))
	}

}

func getValidInputGPXContentStream() (*os.File, error) {
	filePath1 := testhelper.GetValidGPX("05.gpx")
	file1, _ := os.Open(filePath1)
	filePath2 := testhelper.GetValidGPX("04.gpx")
	file2, _ := os.Open(filePath2)

	var inputBytes []byte
	reader1 := bufio.NewReader(file1)
	for {
		input, errRead1 := reader1.ReadByte()
		if errRead1 != nil {
			if errRead1 == io.EOF {
				break
			} else {
				return nil, errRead1
			}
		}

		inputBytes = append(inputBytes, input)
	}

	reader2 := bufio.NewReader(file2)
	for {
		input, errRead2 := reader2.ReadByte()
		if errRead2 != nil {
			if errRead2 == io.EOF {
				break
			} else {
				return nil, errRead2
			}
		}

		inputBytes = append(inputBytes, input)
	}
	read, write, errCreate := os.Pipe()
	if errCreate != nil {
		return nil, errCreate
	}

	_, errWrite := write.Write(inputBytes)
	if errWrite != nil {
		return nil, errWrite
	}
	write.Close()

	return read, nil
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
