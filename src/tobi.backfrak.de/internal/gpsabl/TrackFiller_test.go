package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestCheckValidCorrectionParameters(t *testing.T) {

	if CheckValidCorrectionParameters("asd") {
		t.Errorf("The CheckValidCorrectionParameters return true for \"asd\"")
	}

	if !CheckValidCorrectionParameters("none") {
		t.Errorf("The CheckValidCorrectionParameters return false for \"none\"")
	}

	if !CheckValidCorrectionParameters("linear") {
		t.Errorf("The CheckValidCorrectionParameters return false for \"linear\"")
	}

	if !CheckValidCorrectionParameters("steps") {
		t.Errorf("The CheckValidCorrectionParameters return false for \"steps\"")
	}

	if len(GetValidCorrectionParameters()) != 3 {
		t.Errorf("The number of ValidCorrectionParameters is %d, but %d was expected", len(GetValidCorrectionParameters()), 2)
	}

}

func TestCheckValidCorrectionParametersString(t *testing.T) {

	validParms := GetValidCorrectionParametersString()

	if strings.Contains(validParms, "asd") {
		t.Errorf("The ValidCorrectionParametersString contains \"asd\"")
	}

	if !strings.Contains(validParms, "none") {
		t.Errorf("The ValidCorrectionParametersString not contains \"none\"")
	}

	if !strings.Contains(validParms, "steps") {
		t.Errorf("The ValidCorrectionParametersString not contains \"steps\"")
	}
}

func TestFillValuesTrackPointArrayWrongCorrection(t *testing.T) {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	FillDistancesTrackPoint(&pnt2, pnt1, pnt3)
	err := FillValuesTrackPointArray([]TrackPoint{pnt1, pnt2, pnt3}, "asd")
	if err != nil {
		switch ty := err.(type) {
		case *CorrectionParameterNotKnownError:
			fmt.Println("OK")
		default:
			t.Errorf("The Error FillDistancesTrackPoint gave is of the wrong type. The type is %v", ty)
		}
	} else {
		t.Errorf("FillDistancesTrackPoint did not return a error, but was expected")
	}
}

func TestFillValuesTrackPointArrayValidCorrection(t *testing.T) {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	FillDistancesTrackPoint(&pnt2, pnt1, pnt3)
	err := FillValuesTrackPointArray([]TrackPoint{pnt1, pnt2, pnt3}, GetValidCorrectionParameters()[0])
	if err != nil {
		t.Errorf("FillDistancesTrackPoint did return a error, but was expected to. The error is %s", err.Error())
	}
}

func TestFillDistancesThreePoints(t *testing.T) {
	pnts := gerSimpleTrackPointArray()

	if pnts[1].VerticalDistanceBefore != -1.0 {
		t.Errorf("The VerticalDistanceBefore is %f, but %f was expected", pnts[1].VerticalDistanceBefore, -1.0)
	}

	if pnts[1].HorizontalDistanceBefore != pnts[1].HorizontalDistanceNext {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was expected", pnts[1].HorizontalDistanceBefore, pnts[1].HorizontalDistanceNext)
	}

	if pnts[1].VerticalDistanceNext != 1.0 {
		t.Errorf("The VerticalDistanceNext is %f, but %f was expected", pnts[1].VerticalDistanceNext, 1.0)
	}

	for i := range pnts {
		if i > 0 {
			if pnts[i].DistanceToThisPoint <= pnts[i-1].DistanceToThisPoint {
				t.Errorf("The DistanceToThisPoint for point %d, is %f but the point before had %f", i, pnts[i].DistanceToThisPoint, pnts[i-1].DistanceToThisPoint)
			}
		}
	}
}

func TestFillDistancesThreePointsBeforeAfter(t *testing.T) {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	lon := pnt2.Longitude
	lat := pnt2.Latitude
	eve := pnt2.Elevation

	FillDistancesTrackPoint(&pnt2, pnt1, pnt3)

	if pnt2.Elevation != eve {
		t.Errorf("The Elevation changed during FillDistancesTrackPoint")
	}

	if pnt2.Longitude != lon {
		t.Errorf("The Longitude changed during FillDistancesTrackPoint")
	}

	if pnt2.Latitude != lat {
		t.Errorf("The Latitude changed during FillDistancesTrackPoint")
	}
}

func TestFillDistancesTwoPointBefore(t *testing.T) {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := TrackPoint{}

	pnts := []TrackPoint{pnt1, pnt2, pnt3}

	FillDistancesTrackPoint(&pnts[1], pnts[0], pnts[2])

	fillCorrectedElevationTrackPoint(pnts, "none")
	fillElevationGainLoseTrackPoint(pnts)

	if pnts[1].VerticalDistanceBefore != -1.0 {
		t.Errorf("The VerticalDistanceBefore is %f, but %f was expected", pnts[1].VerticalDistanceBefore, -1.0)
	}

	if pnts[1].HorizontalDistanceBefore == 0.0 {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was not expected", pnts[1].HorizontalDistanceBefore, 0.0)
	}

	if pnts[1].HorizontalDistanceNext != 0.0 {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was expected", pnts[1].HorizontalDistanceNext, 0.0)
	}

	if pnts[1].VerticalDistanceNext != -108.0 {
		t.Errorf("The VerticalDistanceNext is %f, but %f was expected", pnts[1].VerticalDistanceNext, 0.0)
	}
}

func TestFillDistancesThreePointWithLinearCorection(t *testing.T) {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	pnts := []TrackPoint{pnt1, pnt2, pnt3}

	FillDistancesTrackPoint(&pnts[1], pnts[0], pnts[2])

	fillCorrectedElevationTrackPoint(pnts, "linear")
	fillElevationGainLoseTrackPoint(pnts)

	if pnts[1].VerticalDistanceBefore != 0.0 {
		t.Errorf("The VerticalDistanceBefore is %f, but %f was expected", pnts[1].VerticalDistanceBefore, 0.0)
	}

	if pnts[1].VerticalDistanceNext != 0.0 {
		t.Errorf("The VerticalDistanceNext is %f, but %f was expected", pnts[1].VerticalDistanceNext, 0.0)
	}
}

func TestFillDistancesThreePointWithStepsCorection(t *testing.T) {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	pnts := []TrackPoint{pnt1, pnt2, pnt3}

	FillDistancesTrackPoint(&pnts[1], pnts[0], pnts[2])

	fillCorrectedElevationTrackPoint(pnts, "steps")
	fillElevationGainLoseTrackPoint(pnts)
	fillCountUpDownWards(pnts, "steps")

	if pnts[1].VerticalDistanceBefore != 0.0 {
		t.Errorf("The VerticalDistanceBefore is %f, but %f was expected", pnts[1].VerticalDistanceBefore, 0.0)
	}

	if pnts[1].VerticalDistanceNext != 0.0 {
		t.Errorf("The VerticalDistanceNext is %f, but %f was expected", pnts[1].VerticalDistanceNext, 0.0)
	}
}
func TestFillDistancesThreePointWithUnknownCorrection(t *testing.T) {
	pnts := gerSimpleTrackPointArray()
	err := fillCorrectedElevationTrackPoint(pnts, "asd")
	if err != nil {
		switch ty := err.(type) {
		case *CorrectionParameterNotKnownError:
			fmt.Println("OK")
		default:
			t.Errorf("The Error fillCorrectedElevationTrackPoint gave is of the wrong type. The type is %v", ty)
		}
	} else {
		t.Errorf("fillCorrectedElevationTrackPoint did not return a error, but was expected")
	}
}

func TestFillDistancesTwoPointNext(t *testing.T) {
	pnt1 := TrackPoint{}
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	pnts := []TrackPoint{pnt1, pnt2, pnt3}

	FillDistancesTrackPoint(&pnts[1], pnts[0], pnts[2])

	fillCorrectedElevationTrackPoint(pnts, "none")
	fillElevationGainLoseTrackPoint(pnts)

	if pnts[1].VerticalDistanceBefore != 108.0 {
		t.Errorf("The VerticalDistanceBefore is %f, but %f was expected", pnts[1].VerticalDistanceBefore, 0.0)
	}

	if pnts[1].HorizontalDistanceBefore != 0.0 {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was expected", pnts[1].HorizontalDistanceBefore, 0.0)
	}

	if pnts[1].HorizontalDistanceNext == 0.0 {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was not expected", pnts[1].HorizontalDistanceNext, 0.0)
	}

	if pnts[1].VerticalDistanceNext != 1.0 {
		t.Errorf("The VerticalDistanceNext is %f, but %f was expected", pnts[1].VerticalDistanceNext, 0.0)
	}
}

func TestFillTrackSegmentValuesSimple(t *testing.T) {
	seg := getSimpleTrackSegment()

	oldPointNumber := len(seg.TrackPoints)

	FillTrackSegmentValues(&seg)

	if len(seg.TrackPoints) != oldPointNumber {
		t.Errorf("The number of track points changed during FillTrackSegmentValues() call")
	}

	if seg.Distance != 23.885148437468256 {
		t.Errorf("The Distance is %f, but %f expected.", seg.Distance, 23.885148437468256)
	}

	if seg.MaximumAltitude != 109.0 {
		t.Errorf("The MaximumAltitude is %f, but %f expected.", seg.MaximumAltitude, 109.0)
	}

	if seg.MinimumAltitude != 108.0 {
		t.Errorf("The MinimumAltitude is %f, but %f expected.", seg.MinimumAltitude, 108.0)
	}
}

func TestFillTrackValuesBeforeAfter(t *testing.T) {
	name := "Track"
	description := "My test track"
	track := Track{}

	track.Name = name
	track.Description = description
	track.TrackSegments = []TrackSegment{getSimpleTrackSegment()}
	track.NumberOfSegments = 1

	FillTrackValues(&track)

	if track.Name != name {
		t.Errorf("The Name changed during FillTrackValues")
	}

	if track.Description != description {
		t.Errorf("The Description changed during FillTrackValues")
	}

	if track.NumberOfSegments != 1 {
		t.Errorf("The NumberOfSegments changed during FillTrackValues")
	}

	if len(track.TrackSegments) != 1 {
		t.Errorf("The TrackSegments changed during FillTrackValues")
	}
}

func TestFillTrackFileValuesBeforeAfter(t *testing.T) {
	name := "Track"
	description := "My test track"
	file := TrackFile{}

	file.Name = name
	file.Description = description
	file.Tracks = []Track{getSimpleTrack()}
	file.NumberOfTracks = 1

	FillTrackFileValues(&file)

	if file.Name != name {
		t.Errorf("The Name changed during FillTrackFileValues")
	}

	if file.Description != description {
		t.Errorf("The Description changed during FillTrackFileValues")
	}

	if file.NumberOfTracks != 1 {
		t.Errorf("The NumberOfTracks changed during FillTrackFileValues")
	}

	if len(file.Tracks) != 1 {
		t.Errorf("The Tracks changed during FillTrackFileValues")
	}

	if file.GetTimeDataValid() == true {
		t.Errorf("The Time data for TrackFile is valide, but should not")
	}

	if file.Tracks[0].GetTimeDataValid() == true {
		t.Errorf("The Time data for Track is valide, but should not")
	}

	if file.Tracks[0].TrackSegments[0].GetTimeDataValid() == true {
		t.Errorf("The Time data for TrackSegments is valide, but should not")
	}
}

func TestSimpleTrackFileNoTime(t *testing.T) {

	file := getSimpleTrackFile()

	if file.NumberOfTracks != 1 {
		t.Errorf("The NumberOfTracks changed during FillTrackFileValues")
	}

	if len(file.Tracks) != 1 {
		t.Errorf("The Tracks changed during FillTrackFileValues")
	}

	if file.GetTimeDataValid() == true {
		t.Errorf("The Time data for TrackFile is valide, but should not")
	}

	if file.Tracks[0].GetTimeDataValid() == true {
		t.Errorf("The Time data for Track is valide, but should not")
	}

	if file.Tracks[0].TrackSegments[0].GetTimeDataValid() == true {
		t.Errorf("The Time data for TrackSegments is valide, but should not")
	}
}

func TestFillTrackValuesSimple(t *testing.T) {
	track := Track{}
	segs := getSimpleTrackSegment()
	FillTrackSegmentValues(&segs)
	track.TrackSegments = []TrackSegment{segs}
	FillTrackValues(&track)

	if track.Distance != 23.885148437468256 {
		t.Errorf("The Distance is %f, but %f expected.", track.Distance, 23.885148437468256)
	}

	if track.MaximumAltitude != 109.0 {
		t.Errorf("The MaximumAltitude is %f, but %f expected.", track.MaximumAltitude, 109.0)
	}

	if track.MinimumAltitude != 108.0 {
		t.Errorf("The MinimumAltitude is %f, but %f expected.", track.MinimumAltitude, 108.0)
	}
}

func TestFillTrackFileValuesSimple(t *testing.T) {
	file := TrackFile{}
	file.Tracks = []Track{getSimpleTrack()}
	FillTrackFileValues(&file)

	if file.Distance != 23.885148437468256 {
		t.Errorf("The Distance is %f, but %f expected.", file.Distance, 23.885148437468256)
	}

	if file.MaximumAltitude != 109.0 {
		t.Errorf("The MaximumAltitude is %f, but %f expected.", file.MaximumAltitude, 109.0)
	}

	if file.MinimumAltitude != 108.0 {
		t.Errorf("The MinimumAltitude is %f, but %f expected.", file.MinimumAltitude, 108.0)
	}
}

func TestFillSimpleTrackFileWithTime(t *testing.T) {
	file := getSimpleTrackFileWithTime()

	if len(file.Tracks) != 1 {
		t.Errorf("Expected 1 Tracks, got %d", len(file.Tracks))
	}

	if file.GetStartTime() != file.Tracks[0].GetStartTime() {
		t.Errorf("The StartTime does not match for Track")
	}

	if file.GetEndTime() != file.Tracks[0].GetEndTime() {
		t.Errorf("The EndTime does not match for Track")
	}

	if file.GetStartTime() != file.Tracks[0].TrackSegments[0].GetStartTime() {
		t.Errorf("The StartTime does not match for TrackSegments")
	}

	if file.GetEndTime() != file.Tracks[0].TrackSegments[0].GetEndTime() {
		t.Errorf("The EndTime does not match for TrackSegments")
	}

	if file.GetStartTime() != file.Tracks[0].TrackSegments[0].TrackPoints[0].GetStartTime() {
		t.Errorf("The StartTime does not match for TrackPoints")
	}

	lastPoint := len(file.Tracks[0].TrackSegments[0].TrackPoints) - 1
	if file.GetEndTime() != file.Tracks[0].TrackSegments[0].TrackPoints[lastPoint].GetEndTime() {
		t.Errorf("The EndTime does not match for TrackPoints")
	}

	if file.GetStartTime().Format(time.RFC3339) != "2014-08-22T17:19:33Z" {
		t.Errorf("The start time is %s, but %s was expected", file.GetStartTime().Format(time.RFC3339), "2014-08-22T17:19:33Z")
	}

	if file.GetEndTime().Format(time.RFC3339) != "2014-08-22T17:19:53Z" {
		t.Errorf("The end time is %s, but %s was expected", file.GetEndTime().Format(time.RFC3339), "2014-08-22T17:19:53Z")
	}

	if file.GetTimeDataValid() == false {
		t.Errorf("The Time data for TrackFile is not valide, but should")
	}

	if file.Tracks[0].GetTimeDataValid() == false {
		t.Errorf("The Time data for Track is not valide, but should")
	}

	if file.Tracks[0].TrackSegments[0].GetTimeDataValid() == false {
		t.Errorf("The Time data for TrackSegments is not valide, but should")
	}
}

func getSimpleTrackFile() TrackFile {
	ret := NewTrackFile("/mys/track/file")
	trk := getSimpleTrack()
	FillTrackValues(&trk)
	ret.Tracks = []Track{trk}
	FillTrackFileValues(&ret)
	ret.NumberOfTracks = 1

	return ret
}

func getSimpleTrackFileWithTime() TrackFile {
	ret := NewTrackFile("/mys/track/file")
	trk := getSimpleTrackWithTime()
	FillTrackValues(&trk)
	ret.Tracks = []Track{trk}
	FillTrackFileValues(&ret)
	ret.NumberOfTracks = 1

	return ret
}

func getSimpleTrack() Track {
	ret := Track{}
	segs := getSimpleTrackSegment()
	FillTrackSegmentValues(&segs)
	ret.TrackSegments = []TrackSegment{segs}
	FillTrackValues(&ret)
	ret.NumberOfSegments = 1

	return ret
}

func getSimpleTrackWithTime() Track {
	ret := Track{}
	segs := getSimpleTrackSegmentWithTime()
	FillTrackSegmentValues(&segs)
	ret.TrackSegments = []TrackSegment{segs}
	FillTrackValues(&ret)
	ret.NumberOfSegments = 1

	return ret
}

func getSimpleTrackSegment() TrackSegment {
	seg := TrackSegment{}
	points := gerSimpleTrackPointArray()
	seg.TrackPoints = points

	return seg
}

func getSimpleTrackSegmentWithTime() TrackSegment {
	seg := TrackSegment{}
	points := getSimpleTrackPointArrayWithTime()
	seg.TrackPoints = points

	return seg
}

func gerSimpleTrackPointArray() []TrackPoint {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	points := []TrackPoint{pnt1, pnt2, pnt3}

	FillDistancesTrackPoint(&points[0], TrackPoint{}, points[1])
	FillDistancesTrackPoint(&points[1], points[0], points[2])
	FillDistancesTrackPoint(&points[2], points[1], TrackPoint{})
	FillValuesTrackPointArray(points, "none")

	return points
}

func getSimpleTrackPointArrayWithTime() []TrackPoint {
	t1, _ := time.Parse(time.RFC3339, "2014-08-22T17:19:33Z")
	t2, _ := time.Parse(time.RFC3339, "2014-08-22T17:19:43Z")
	t3, _ := time.Parse(time.RFC3339, "2014-08-22T17:19:53Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 108.0, t2)
	pnt3 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t3)
	points := []TrackPoint{pnt1, pnt2, pnt3}

	FillDistancesTrackPoint(&points[0], TrackPoint{}, points[1])
	FillDistancesTrackPoint(&points[1], points[0], points[2])
	FillDistancesTrackPoint(&points[2], points[1], TrackPoint{})
	FillValuesTrackPointArray(points, "none")

	return points
}
