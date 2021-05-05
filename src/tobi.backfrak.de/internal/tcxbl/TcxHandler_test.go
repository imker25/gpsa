package tcxbl

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"tobi.backfrak.de/internal/gpsabl"
	"tobi.backfrak.de/internal/testhelper"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestTrackReader02(t *testing.T) {
	tcx := NewTcxFile(testhelper.GetValidTcx("01.tcx"))

	file, _ := tcx.ReadTracks("none", 0.3, 10.0)

	if file.GetDistance() != 6216.201383825188 {
		t.Errorf("The Distance is %f, but should be %f", file.GetDistance(), 6216.201383825188)
	}

	if file.GetMovingTime() != 1103000000000 {
		t.Errorf("The MovingTime is %d, but should be %d", file.GetMovingTime(), 1103000000000)
	}

	if file.GetStartTime().Format(time.RFC3339) != "2016-06-05T10:45:59Z" {
		t.Errorf("The StartTime is %s, but should be %s", file.GetStartTime().Format(time.RFC3339), "2016-06-05T10:45:59Z")
	}

	if file.GetUpwardsSpeed() != 4.797203586378565 {
		t.Errorf("The UpwardsSpeed is %f, but should be %f", file.GetUpwardsSpeed(), 4.797204)
	}

	if file.GetDownwardsSpeed() != 7.683785245071919 {
		t.Errorf("The UpwardsSpeed is %f, but should be %f", file.GetDownwardsSpeed(), 7.683785)
	}
}

func TestTrackReaderEmptyTrack(t *testing.T) {
	tcx := NewTcxFile(testhelper.GetInvalidTcx("03.tcx"))

	_, err := tcx.ReadTracks("none", 0.3, 10.0)
	if err != nil {
		switch ty := err.(type) {
		case *EmptyTcxFileError:
			fmt.Println("OK")
		default:
			t.Errorf("The Error ReadTracks gave is of the wrong type. The type is %v", ty)
		}
	} else {
		t.Errorf("ReadTracks did not return a error, but was expected")
	}
}

func TestTrackReaderInValidCorrectionParameter(t *testing.T) {
	tcx := NewTcxFile(testhelper.GetValidTcx("02.tcx"))

	_, err := tcx.ReadTracks("asdfg", 0.3, 10.0)
	if err != nil {
		switch ty := err.(type) {
		case *gpsabl.CorrectionParameterNotKnownError:
			fmt.Println("OK")
		default:
			t.Errorf("The Error ReadTracks gave is of the wrong type. The type is %v", ty)
		}
	} else {
		t.Errorf("ReadTracks did not return a error, but was expected")
	}
}

func TestReadAllValidTcxDueInterface(t *testing.T) {
	files, _ := ioutil.ReadDir(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-tcx"))

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tcx") {
			if file.IsDir() == false {
				tcxFile := NewTcxFile(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-tcx", file.Name()))
				iTcx := gpsabl.TrackReader(&tcxFile)

				track, err := iTcx.ReadTracks("none", 0.3, 10.0)
				if err != nil {
					t.Errorf("Got the error \"%s\" while reading file %s.", err.Error(), filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-tcx", file.Name()))

				}
				if track.Distance <= 0.0 {
					t.Errorf("The track.Distance is %f but should not be less then %f", track.Distance, 0.0)
				}
			}
		}
	}
}

func TestNewReaderWithValidFilePath(t *testing.T) {
	tcx := TcxFile{}
	file := testhelper.GetValidTcx("01.tcx")
	checkRes := tcx.CheckFile(file)
	input := *gpsabl.NewInputFileWithPath(file)

	sut := tcx.NewReader(input)

	if checkRes != true {
		t.Errorf("TcxFile can not read %s", file)
	}

	trk, err := sut.ReadTracks(gpsabl.STEPS, 0.01, 10)
	if err != nil {
		t.Errorf("Got error \"%s\" but expect none", err)
	}

	if trk.FilePath != file {
		t.Errorf("The trk.FilePath %s is not the given path %s", trk.FilePath, file)
	}

	if trk.GetDistance() != 6216.201383825188 {
		t.Errorf("The Distance is %f, but should be %f", trk.GetDistance(), 6216.201383825188)
	}

}

func TestNewReaderWithValidBuffer(t *testing.T) {
	tcx := TcxFile{}
	buffer, createErr := testhelper.GetValidTcxBuffer("01.tcx")
	name := "Buffer 1"
	if createErr != nil {
		t.Fatalf("Got error \"%s\" while creating the input buffer", createErr)
	}
	checkRes := tcx.CheckBuffer(buffer)
	input := *tcx.NewInputFileForBuffer(buffer, name)

	sut := tcx.NewReader(input)

	if checkRes != true {
		t.Errorf("TcxFile can not read from buffer")
	}

	trk, err := sut.ReadTracks(gpsabl.STEPS, 0.01, 10)
	if err != nil {
		t.Errorf("Got error \"%s\" but expect none", err)
	}

	if trk.FilePath != name {
		t.Errorf("The trk.FilePath %s is not the given path %s", trk.FilePath, name)
	}

	if trk.GetDistance() != 6216.201383825188 {
		t.Errorf("The Distance is %f, but should be %f", trk.GetDistance(), 6216.201383825188)
	}

}

func TestCheckInputFile(t *testing.T) {
	tcx := TcxFile{}
	file1 := testhelper.GetValidTcx("01.tcx")
	input1 := *gpsabl.NewInputFileWithPath(file1)
	if tcx.CheckInputFile(input1) == false {
		t.Errorf("TcxFile can not read %s", file1)
	}

	file2 := testhelper.GetValidGPX("01.gpx")
	input2 := *gpsabl.NewInputFileWithPath(file2)
	if tcx.CheckInputFile(input2) == true {
		t.Errorf("TcxFile can read %s", file2)
	}

	buffer1, createErr1 := testhelper.GetValidTcxBuffer("01.tcx")
	name := "Buffer 1"
	if createErr1 != nil {
		t.Fatalf("Got error \"%s\" while creating the input buffer", createErr1)
	}
	input3 := *tcx.NewInputFileForBuffer(buffer1, name)
	if tcx.CheckInputFile(input3) == false {
		t.Errorf("TcxFile can not read %s from buffer", file1)
	}
}

func TestCheckBuffer(t *testing.T) {
	tcx := TcxFile{}
	buffer1, createErr1 := testhelper.GetValidTcxBuffer("01.tcx")
	if createErr1 != nil {
		t.Fatalf("Got error \"%s\" while creating the input buffer", createErr1)
	}

	if tcx.CheckBuffer(buffer1) == false {
		t.Errorf("TcxFile can not read %s from buffer", testhelper.GetValidTcx("01.tcx"))
	}

	buffer2, createErr2 := testhelper.GetValidGpxBuffer("01.gpx")
	if createErr2 != nil {
		t.Fatalf("Got error \"%s\" while creating the input buffer", createErr1)
	}

	if tcx.CheckBuffer(buffer2) == true {
		t.Errorf("TcxFile can  read %s from buffer", testhelper.GetValidGPX("01.gpx"))
	}
}

func TestReadBufferWithInvalidBuffer(t *testing.T) {
	buffer1, createErr1 := testhelper.GetInvalidTcxBuffer("01.tcx")
	if createErr1 != nil {
		t.Fatalf("Got error \"%s\" while creating the input buffer", createErr1)
	}

	_, err := ReadBuffer(buffer1, "bla", gpsabl.STEPS, 0.3, 10.0)
	if err == nil {
		t.Errorf("Got no error, but expected one")
	}

}

func TestGetValidFileExtensions(t *testing.T) {
	tcx := TcxFile{}
	file := testhelper.GetValidTcx("01.tcx")
	input := *gpsabl.NewInputFileWithPath(file)

	sut := tcx.NewReader(input)

	extensions := sut.GetValidFileExtensions()

	if len(extensions) != 1 {
		t.Errorf("The GetValidFileExtensions does not give the expected amount of data")
	}

	if extensions[0] != ".tcx" {
		t.Errorf("The GetValidFileExtensions does not give the expected value")
	}

}
