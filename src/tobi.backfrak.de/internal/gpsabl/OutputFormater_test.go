package gpsabl

import (
	"os"
	"strings"
	"testing"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestDepthParametrNotKnownErrorStruct(t *testing.T) {
	val := "asdgfg"
	err := NewDepthParametrNotKnownError(val)

	if err.GivenValue != val {
		t.Errorf("The GivenValue was %s, but %s was expected", err.GivenValue, val)
	}

	if strings.Contains(err.Error(), val) == false {
		t.Errorf("The error messaage of DepthParametrNotKnownError does not contain the expected GivenValue")
	}
}

func TestNewCsvOutputFormater(t *testing.T) {
	sut := NewCsvOutputFormater(";")

	if sut.Seperator != ";" {
		t.Errorf("The Seperator was \"%s\", but \";\" was expected", sut.Seperator)
	}

	if len(sut.ValideDepthArgs) != 3 {
		t.Errorf("The ValideDepthArgs array does not contain the expeced number of values")
	}
}

func TestFormatOutPutWithHeader(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	ret := formater.FormatOutPut(trackFile, true, "file")

	if len(ret) != 2 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Count(ret[0], ";") != strings.Count(ret[1], ";") {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(ret[1], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", ret[1])
	}

	if strings.Contains(ret[1], trackFile.FilePath) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[1])
	}
}

func TestFormatOutPutWithHeaderAndSetName(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()
	trackFile.Name = "My Track File"

	ret := formater.FormatOutPut(trackFile, true, "file")

	if len(ret) != 2 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Count(ret[0], ";") != strings.Count(ret[1], ";") {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(ret[1], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", ret[1])
	}

	if strings.Contains(ret[1], trackFile.Name) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[1])
	}

	if strings.Contains(ret[1], trackFile.FilePath) == true {
		t.Errorf("The output does contian the FilePath but should not. It is: %s", ret[1])
	}
}

func TestFormatOutPutWithOutHeader(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	ret := formater.FormatOutPut(trackFile, false, "file")

	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Count(ret[0], ";") != 5 {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(ret[0], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderTrackDepth(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	ret := formater.FormatOutPut(trackFile, false, "track")

	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Contains(ret[0], trackFile.FilePath) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], "#1;") == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderTrackDepthSetTrackName(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()
	trackFile.Tracks[0].Name = "My Track"

	ret := formater.FormatOutPut(trackFile, false, "track")

	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Contains(ret[0], trackFile.FilePath) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], trackFile.Tracks[0].Name) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], "#1;") == true {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderTrackDepthSetTrackFileNameSetTrackName(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()
	trackFile.Tracks[0].Name = "My Track"
	trackFile.Name = "My track file"

	ret := formater.FormatOutPut(trackFile, false, "track")

	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Contains(ret[0], trackFile.FilePath) == true {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], trackFile.Name) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], trackFile.Tracks[0].Name) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], "#1;") == true {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderTrackDepthSegmentSetTrackFileNameSetTrackName(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()
	trackFile.Tracks[0].Name = "My Track"
	trackFile.Name = "My track file"

	ret := formater.FormatOutPut(trackFile, false, "segment")

	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Contains(ret[0], trackFile.FilePath) == true {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], trackFile.Name) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], trackFile.Tracks[0].Name) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], "Segment #1;") == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderUnValideDepth(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()
	ret := formater.FormatOutPut(trackFile, false, "abc")

	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of lines")
	}

	if strings.HasPrefix(ret[0], "Error:") == false {
		t.Errorf("The line does not start with \"Error\" as expected")
	}
}

func TestAddHeader(t *testing.T) {
	frt := NewCsvOutputFormater(";")

	frt.AddHeader()

	lines := frt.GetLines()

	if len(lines) != 1 {
		t.Errorf("The number of lines was not expected. Got %d, expected %d", len(lines), 1)
	}

	if strings.Count(lines[0], ";") != 5 {
		t.Errorf("The Number of semicolons is not the same in each line")
	}
}

func TestAddOutPut(t *testing.T) {
	frt := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	frt.AddOutPut(trackFile, "file")

	lines := frt.GetLines()

	if len(lines) != 1 {
		t.Errorf("The number of lines was not expected. Got %d, expected %d", len(lines), 1)
	}

	if strings.Count(lines[0], ";") != 5 {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(lines[0], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", lines[0])
	}
}

func TestAddHeaderAndOutPut(t *testing.T) {
	frt := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	frt.AddHeader()
	frt.AddOutPut(trackFile, "file")

	lines := frt.GetLines()

	if len(lines) != 2 {
		t.Errorf("The number of lines was not expected. Got %d, expected %d", len(lines), 2)
	}

	if strings.Count(lines[0], ";") != strings.Count(lines[1], ";") {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(lines[1], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", lines[0])
	}
}

func TestAddHeaderAndOutPutFileTwoTracksFileDepth(t *testing.T) {
	frt := NewCsvOutputFormater(";")
	trackFile := getTrackFileTwoTracks()

	frt.AddHeader()
	frt.AddOutPut(trackFile, "file")

	lines := frt.GetLines()

	if len(lines) != 2 {
		t.Errorf("The number of lines was not expected. Got %d, expected %d", len(lines), 2)
	}

	if strings.Count(lines[0], ";") != strings.Count(lines[1], ";") {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(lines[1], "0.0500") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", lines[1])
	}
}

func TestAddHeaderAndOutPutFileTwoTracksTrackDepth(t *testing.T) {
	frt := NewCsvOutputFormater(";")
	trackFile := getTrackFileTwoTracks()

	frt.AddHeader()
	frt.AddOutPut(trackFile, "track")

	lines := frt.GetLines()

	if len(lines) != 3 {
		t.Errorf("The number of lines was not expected. Got %d, expected %d", len(lines), 3)
	}

	if strings.Count(lines[0], ";") != strings.Count(lines[1], ";") {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(lines[1], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", lines[1])
	}

	if strings.Contains(lines[2], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", lines[2])
	}
}

func TestAddHeaderAndOutPutFileTwoTracksSegmentDepth(t *testing.T) {
	frt := NewCsvOutputFormater(";")
	trackFile := getTrackFileTwoTracksWithThreeSegments()

	frt.AddHeader()
	frt.AddOutPut(trackFile, "segment")

	lines := frt.GetLines()

	if len(lines) != 4 {
		t.Errorf("The number of lines was not expected. Got %d, expected %d", len(lines), 4)
	}

	if strings.Count(lines[0], ";") != strings.Count(lines[3], ";") {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(lines[1], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", lines[1])
	}

	if strings.Contains(lines[3], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", lines[2])
	}
}

func TestWriteOutputSegmentDepth(t *testing.T) {
	frt := NewCsvOutputFormater(";")
	trackFile := getTrackFileTwoTracksWithThreeSegments()

	frt.AddHeader()
	frt.AddOutPut(trackFile, "segment")

	err := frt.WriteOutput(os.Stdout)

	if err != nil {
		t.Errorf("Error while writing the output: %s", err.Error())
	}
}

func TestCheckVlaideDepthArg(t *testing.T) {
	frt := NewCsvOutputFormater(";")

	if frt.CheckVlaideDepthArg("asfd") == true {
		t.Errorf("The CheckVlaideDepthArg returns true for \"asfd\"")
	}

	if frt.CheckVlaideDepthArg("file") == false {
		t.Errorf("The CheckVlaideDepthArg returns false for \"file\"")
	}
}

func TestCsvOutputFormaterIsOutputFormater(t *testing.T) {
	frt := NewCsvOutputFormater(";")
	trackFile := getTrackFileTwoTracksWithThreeSegments()

	iFrt := OutputFormater(frt)
	iFrt.AddHeader()
	iFrt.AddOutPut(trackFile, "track")

	err := iFrt.WriteOutput(os.Stdout)

	if err != nil {
		t.Errorf("Error while writing the output: %s", err.Error())
	}
}

func getTrackFileTwoTracksWithThreeSegments() TrackFile {
	trackFile := getTrackFileTwoTracks()
	trackFile.Tracks[0].TrackSegments = append(trackFile.Tracks[0].TrackSegments, getSimpleTrackFile().Tracks[0].TrackSegments[0])
	trackFile.Tracks[0] = FillTrackValues(trackFile.Tracks[0])

	return trackFile
}

func getTrackFileTwoTracks() TrackFile {
	trackFile := getSimpleTrackFile()
	trackFile.Tracks = append(trackFile.Tracks, getSimpleTrackFile().Tracks...)
	trackFile = FillTrackFileValues(trackFile)

	return trackFile
}
