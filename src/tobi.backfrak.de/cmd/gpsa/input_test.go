package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"tobi.backfrak.de/internal/csvbl"
	"tobi.backfrak.de/internal/gpsabl"
	"tobi.backfrak.de/internal/gpxbl"
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

	if result[0].Type != gpsabl.FilePath {
		t.Errorf("The type is %s, but %s is expected", result[0].Type, gpsabl.FilePath)
	}

	if result[1].Name != file2 {
		t.Errorf("The path is %s, but %s is expected", result[1], file1)
	}

	if result[1].Type != gpsabl.FilePath {
		t.Errorf("The type is %s, but %s is expected", result[1].Type, gpsabl.FilePath)
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

	if input[0].Type != gpxbl.GpxBuffer {
		t.Errorf("The type is %s, but %s is expected", input[0].Type, gpxbl.GpxBuffer)
	}

	if input[0].Name == "" {
		t.Errorf(("The name is \"\""))
	}

	if input[0].Buffer == nil {
		t.Errorf("The buffer is nil")
	}

	if input[1].Type != gpxbl.GpxBuffer {
		t.Errorf("The type is %s, but %s is expected", input[1].Type, gpxbl.GpxBuffer)
	}

	if input[1].Buffer == nil {
		t.Errorf("The buffer is nil")
	}

	if input[0].Name == input[1].Name {
		t.Errorf(("The names are the same"))
	}

}

func TestGetValidTrackExtensions(t *testing.T) {
	sut := getValidTrackExtensions()

	if strings.Contains(sut, ".gpx") == false {
		t.Errorf("\"%s\" does not contain \".gpx\"", sut)
	}

	if strings.Contains(sut, ".tcx") == false {
		t.Errorf("\"%s\" does not contain \".tcx\"", sut)
	}
}

func getValidInputGPXContentStream() (*os.File, error) {
	buffer1, err1 := testhelper.GetValidGpxBuffer("05.gpx")
	if err1 != nil {
		return nil, err1
	}
	buffer2, err2 := testhelper.GetValidGpxBuffer("04.gpx")
	if err2 != nil {
		return nil, err1
	}
	var inputBytes []byte
	inputBytes = append(inputBytes, buffer1...)
	inputBytes = append(inputBytes, buffer2...)
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
