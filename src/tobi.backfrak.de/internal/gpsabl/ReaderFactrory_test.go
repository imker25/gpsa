package gpsabl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"strings"
	"testing"
)

func TestGetInputFileFromPathWithNoReaders(t *testing.T) {
	var readers []TrackReader
	file := "my/flie.gpx"
	res, _ := GetInputFileFromPath(readers, file)

	if res == true {
		t.Errorf("Can find inputFile without reader")
	}

}

func TestGetInputFileFromBufferWithNoReaders(t *testing.T) {
	var readers []TrackReader
	buffer := []byte{1, 2, 3, 4, 5, 6}
	res, _ := GetInputFileFromBuffer(readers, buffer, "name")

	if res == true {
		t.Errorf("Can find inputFile without reader")
	}

}

func TestGetNewReaderWithNoReaders(t *testing.T) {
	var readers []TrackReader
	input := *NewInputFileWithPath("my/flie.gpx")
	res, _ := GetNewReader(readers, input)

	if res == true {
		t.Errorf("Can find inputFile without reader")
	}

}

func TestGetInputFileFromPath(t *testing.T) {
	var readers []TrackReader
	var mok readerMock
	readers = append(readers, &mok)
	file := "my/flie.abc"
	res, _ := GetInputFileFromPath(readers, file)

	if res == false {
		t.Errorf("Can not find the input file for the mok reader")
	}

}

func TestGetNewReadert(t *testing.T) {
	var readers []TrackReader
	var mok readerMock
	readers = append(readers, &mok)
	input := *NewInputFileWithPath("my/flie.abc")
	res, _ := GetNewReader(readers, input)

	if res == false {
		t.Errorf("Can not get a reader for mok")
	}

}

func TestGetInputFileFromBuffer(t *testing.T) {
	var readers []TrackReader
	var mok readerMock
	readers = append(readers, &mok)
	buffer := []byte{97, 98, 99, 110, 111, 110}
	res, _ := GetInputFileFromBuffer(readers, buffer, "name")

	if res == false {
		t.Errorf("Can find reader for mok buffer")
	}

}

type readerMock struct {
	TrackFile
	input InputFile
}

func (mok *readerMock) NewReader(data InputFile) TrackReader {
	newMok := readerMock{}
	newMok.input = data
	if data.Type == FilePath {
		newMok.FilePath = data.Name
	}

	return &newMok
}

func (mok *readerMock) CheckBuffer(buffer []byte) bool {
	for i, _ := range buffer {
		section := buffer[i : i+3]
		if string(section) == "abc" {
			return true
		}
	}
	return false
}

func (mok *readerMock) NewInputFileForBuffer(buffer []byte, name string) *InputFile {
	file := InputFile{}
	file.Name = name
	file.Type = InputFileType("abcType")
	file.Buffer = buffer

	return &file
}

func (mok *readerMock) CheckInputFile(input InputFile) bool {
	if input.Type == InputFileType("abcType") {
		return true
	} else if input.Type == FilePath && mok.CheckFile(input.Name) {
		return true
	}

	return false
}

func (mok *readerMock) CheckFile(path string) bool {
	if strings.HasSuffix(path, "abc") == true {
		return true
	}

	return false
}

func (mok *readerMock) GetValidFileExtensions() []string {

	extensions := []string{".abc"}

	return extensions
}

func (mok *readerMock) ReadTracks(correction CorrectionParameter, minimalMovingSpeed float64, minimalStepHight float64) (TrackFile, error) {
	return mok.TrackFile, nil
}
