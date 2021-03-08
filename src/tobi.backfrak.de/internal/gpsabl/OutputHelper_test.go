package gpsabl

import (
	"fmt"
	"testing"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestGetLineNameFromTrackFile(t *testing.T) {
	file := getTrackFileTwoTracksWithThreeSegmentsWithTime()
	name := "My test track"
	if getLineNameFromTrackFile(file) != file.FilePath {
		t.Errorf("getLineNameFromTrackFile does not return the  expected value")
	}
	file.Name = name
	if getLineNameFromTrackFile(file) != name {
		t.Errorf("getLineNameFromTrackFile does not return the  expected value")
	}
}

func TestGetLineNameFromTrack(t *testing.T) {
	file := getSimpleTrackFileWithTime()
	track := file.Tracks[0]

	trackName := "My test track"
	fileName := "My test file"
	file.Name = fileName
	index := 1
	out1 := getLineNameFromTrack(track, file, index)
	res1 := fmt.Sprintf("%s: Track #%d", fileName, index+1)
	if out1 != res1 {
		t.Errorf("The name is \"%s\", but expected \"%s\"", out1, res1)
	}
	file.Tracks[0].Name = trackName
	track.Name = trackName
	out2 := getLineNameFromTrack(track, file, index)
	res2 := fmt.Sprintf("%s: %s", fileName, trackName)
	if out2 != res2 {
		t.Errorf("The name is \"%s\", but expected \"%s\"", out2, res2)
	}

}

func TestGetOutlinesFileDepth(t *testing.T) {
	file := getTrackFileTwoTracksWithThreeSegmentsWithTime()
	outlines, err := GetOutlines(file, FILE)
	if err != nil {
		t.Errorf("Got error  \"%s\", but expected none", err.Error())
	}
	if len(outlines) != 1 {
		t.Errorf("Got not the expected number of outlines")
	}

	if outlines[0].Name != getLineNameFromTrackFile(file) {
		t.Errorf("The outline.Name is \"%s\", but expected \"%s\"", outlines[0].Name, getLineNameFromTrackFile(file))
	}

	if outlines[0].Data.GetStartTime() != file.GetStartTime() {
		t.Errorf("The outline.Data is not expected")
	}
}

func TestGetOutlinesTrackDepth(t *testing.T) {
	file := getTrackFileTwoTracksWithThreeSegmentsWithTime()
	outlines, err := GetOutlines(file, TRACK)
	if err != nil {
		t.Errorf("Got error  \"%s\", but expected none", err.Error())
	}
	if len(outlines) != 2 {
		t.Errorf("Got not the expected number of outlines")
	}

	if outlines[0].Name != getLineNameFromTrack(file.Tracks[0], file, 0) {
		t.Errorf("The outline.Name is \"%s\", but expected \"%s\"", outlines[0].Name, getLineNameFromTrack(file.Tracks[0], file, 0))
	}

	if outlines[0].Data.GetStartTime() != file.GetStartTime() {
		t.Errorf("The outline.Data is not expected")
	}
}

func TestGetOutlinesSegmentDepth(t *testing.T) {
	file := getTrackFileTwoTracksWithThreeSegmentsWithTime()
	outlines, err := GetOutlines(file, SEGMENT)
	if err != nil {
		t.Errorf("Got error  \"%s\", but expected none", err.Error())
	}
	if len(outlines) != 3 {
		t.Errorf("Got not the expected number of outlines")
	}

	if outlines[0].Name != getLineNameFromTrack(file.Tracks[0], file, 0)+": Segment #1" {
		t.Errorf("The outline.Name is \"%s\", but expected \"%s\"", outlines[0].Name, getLineNameFromTrack(file.Tracks[0], file, 0)+": Segment #1")
	}

	if outlines[0].Data.GetStartTime() != file.GetStartTime() {
		t.Errorf("The outline.Data is not expected")
	}

	if outlines[2].Data.GetEndTime() != file.GetEndTime() {
		t.Errorf("The outline.Data is not expected")
	}
}

func TestStripOutlines(t *testing.T) {
	file := getTrackFileTwoTracksWithThreeSegmentsWithTime()
	outlines, err := GetOutlines(file, SEGMENT)
	if err != nil {
		t.Errorf("Got error  \"%s\", but expected none", err.Error())
	}

	stripped := StripOutlines(outlines)

	if len(stripped) != len(outlines) {
		t.Errorf("The stripped outliens are not the right amount")
	}

	for i := range outlines {
		if stripped[i].Name != outlines[i].Name {
			t.Errorf("The stripped line name does not match")
		}

		if stripped[i].Data.GetAltitudeRange() != outlines[i].Data.GetAltitudeRange() {
			t.Errorf("The stripped line GetAltitudeRange does not match")
		}

		if stripped[i].Data.GetDistance() != outlines[i].Data.GetDistance() {
			t.Errorf("The stripped line GetDistance does not match")
		}

		if stripped[i].Data.GetHorizontalDistance() != outlines[i].Data.GetHorizontalDistance() {
			t.Errorf("The stripped line GetHorizontalDistance does not match")
		}

		if stripped[i].Data.GetMaximumAltitude() != outlines[i].Data.GetMaximumAltitude() {
			t.Errorf("The stripped line GetMaximumAltitude does not match")
		}

		if stripped[i].Data.GetMinimumAltitude() != outlines[i].Data.GetMinimumAltitude() {
			t.Errorf("The stripped line GetMinimumAltitude does not match")
		}

		if stripped[i].Data.GetUpwardsDistance() != outlines[i].Data.GetUpwardsDistance() {
			t.Errorf("The stripped line GetUpwardsDistance does not match")
		}

		if stripped[i].Data.GetAvarageSpeed() != outlines[i].Data.GetAvarageSpeed() {
			t.Errorf("The stripped line GetAvarageSpeed does not match")
		}

		if stripped[i].Data.GetMovingTime() != outlines[i].Data.GetMovingTime() {
			t.Errorf("The stripped line GetMovingTime does not match")
		}
	}
}

func TestGetOutlinesUnkownDepth(t *testing.T) {
	file := getTrackFileTwoTracksWithThreeSegmentsWithTime()
	_, err := GetOutlines(file, "blabla")
	if err == nil {
		t.Errorf("Got no errorbut expected one")
	}
	switch err.(type) {
	case *DepthParameterNotKnownError:
		fmt.Println("OK")
	default:
		t.Errorf("The error is not from the expected type")
	}
}
