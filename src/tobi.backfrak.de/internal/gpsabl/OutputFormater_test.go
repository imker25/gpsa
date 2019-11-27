package gpsabl

import (
	"fmt"
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

	if len(sut.lineBuffer) != 0 {
		t.Errorf("The line buffer is not empty on a new CsvOutputFormater")
	}

	if len(sut.GetLines()) != 0 {
		t.Errorf("The line buffer is not empty on a new CsvOutputFormater")
	}
}

func TestFormatOutPutWithHeader(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	ret, err := formater.FormatOutPut(trackFile, true, "file")
	if err != nil {
		t.Errorf("Got a error, but did not expect one. The error is: %s", err.Error())
	}

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

	ret, err := formater.FormatOutPut(trackFile, true, "file")
	if err != nil {
		t.Errorf("Got a error, but did not expect one. The error is: %s", err.Error())
	}

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

	ret, err := formater.FormatOutPut(trackFile, false, "file")
	if err != nil {
		t.Errorf("Got a error, but did not expect one. The error is: %s", err.Error())
	}

	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Count(ret[0], ";") != 9 {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(ret[0], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderTrackDepth(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	ret, err := formater.FormatOutPut(trackFile, false, "track")
	if err != nil {
		t.Errorf("Got a error, but did not expect one. The error is: %s", err.Error())
	}

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

	ret, err := formater.FormatOutPut(trackFile, false, "track")
	if err != nil {
		t.Errorf("Got a error, but did not expect one. The error is: %s", err.Error())
	}

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

	ret, err := formater.FormatOutPut(trackFile, false, "track")
	if err != nil {
		t.Errorf("Got a error, but did not expect one. The error is: %s", err.Error())
	}

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

	ret, err := formater.FormatOutPut(trackFile, false, "segment")
	if err != nil {
		t.Errorf("Got a error, but did not expect one. The error is: %s", err.Error())
	}

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
	_, err := formater.FormatOutPut(trackFile, false, "abc")

	if err == nil {
		t.Errorf("Did not get an error as expected")
	}

	switch err.(type) {
	case *DepthParametrNotKnownError:
		fmt.Println("OK")
	default:
		t.Errorf("The error is not from the expected type")
	}
}

func TestAddHeader(t *testing.T) {
	frt := NewCsvOutputFormater(";")

	frt.AddHeader()

	lines := frt.GetLines()

	if len(lines) != 1 {
		t.Errorf("The number of lines was not expected. Got %d, expected %d", len(lines), 1)
	}

	if strings.Count(lines[0], ";") != 9 {
		t.Errorf("The Number of semicolons is not the same in each line")
	}
}

func TestAddOutPut(t *testing.T) {
	frt := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	err := frt.AddOutPut(trackFile, "file")
	if err != nil {
		t.Errorf("Got a error, but did not expect one. The error is: %s", err.Error())
	}

	lines := frt.GetLines()

	if len(lines) != 1 {
		t.Errorf("The number of lines was not expected. Got %d, expected %d", len(lines), 1)
	}

	if strings.Count(lines[0], ";") != 9 {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Count(lines[0], "0.020000;") != 1 {
		t.Errorf("The output does not contian the distance as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "1.000000;") != 3 {
		t.Errorf("The output does not contian the ElevationGain as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "0.010000;") != 2 {
		t.Errorf("The output does not contian the UpwardsDistance as expected. It is: %s", lines[0])
	}
}

func TestAddHeaderAndOutPut(t *testing.T) {
	frt := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	frt.AddHeader()
	err := frt.AddOutPut(trackFile, "file")
	if err != nil {
		t.Errorf("Got a error, but did not expect one. The error is: %s", err.Error())
	}

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
	err := frt.AddOutPut(trackFile, "track")
	if err != nil {
		t.Errorf("Got a error, but did not expect one. The error is: %s", err.Error())
	}

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
	err := frt.AddOutPut(trackFile, "segment")
	if err != nil {
		t.Errorf("Got a error, but did not expect one. The error is: %s", err.Error())
	}

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
	errAdd := frt.AddOutPut(trackFile, "segment")
	if errAdd != nil {
		t.Errorf("Got a error, but did not expect one. The error is: %s", errAdd.Error())
	}

	errWrite := frt.WriteOutput(os.Stdout)

	if errWrite != nil {
		t.Errorf("Error while writing the output: %s", errWrite.Error())
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
	errAdd := iFrt.AddOutPut(trackFile, "track")
	if errAdd != nil {
		t.Errorf("Got a error, but did not expect one. The error is: %s", errAdd.Error())
	}

	errWrite := iFrt.WriteOutput(os.Stdout)

	if errWrite != nil {
		t.Errorf("Error while writing the output: %s", errWrite.Error())
	}
}

func getTrackFileTwoTracksWithThreeSegments() TrackFile {
	trackFile := getTrackFileTwoTracks()
	trackFile.Tracks[0].TrackSegments = append(trackFile.Tracks[0].TrackSegments, getSimpleTrackFile().Tracks[0].TrackSegments[0])
	FillTrackValues(&trackFile.Tracks[0])

	return trackFile
}

func getTrackFileTwoTracks() TrackFile {
	trackFile := getSimpleTrackFile()
	trackFile.Tracks = append(trackFile.Tracks, getSimpleTrackFile().Tracks...)
	FillTrackFileValues(&trackFile)

	return trackFile
}
