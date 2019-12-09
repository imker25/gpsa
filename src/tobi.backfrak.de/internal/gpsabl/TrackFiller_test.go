package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"strings"
	"testing"
)

func TestCheckValideCorectionParamters(t *testing.T) {

	if CheckValideCorectionParamters("asd") {
		t.Errorf("The CheckValideCorectionParamters return true for \"asd\"")
	}

	if !CheckValideCorectionParamters("none") {
		t.Errorf("The CheckValideCorectionParamters return false for \"none\"")
	}

	if !CheckValideCorectionParamters("linear") {
		t.Errorf("The CheckValideCorectionParamters return false for \"linear\"")
	}

	if !CheckValideCorectionParamters("steps") {
		t.Errorf("The CheckValideCorectionParamters return false for \"steps\"")
	}

	if len(GetValideCorectionParamters()) != 3 {
		t.Errorf("The number of ValideCorectionParamters is %d, but %d was expected", len(GetValideCorectionParamters()), 2)
	}

}

func TestCheckValideCorectionParamtersString(t *testing.T) {

	valideParms := GetValideCorectionParamtersString()

	if strings.Contains(valideParms, "asd") {
		t.Errorf("The ValideCorectionParamtersString contains \"asd\"")
	}

	if !strings.Contains(valideParms, "none") {
		t.Errorf("The ValideCorectionParamtersString not contains \"none\"")
	}

	if !strings.Contains(valideParms, "steps") {
		t.Errorf("The ValideCorectionParamtersString not contains \"steps\"")
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
		case *CorectionParamterNotKnownError:
			fmt.Println("OK")
		default:
			t.Errorf("The Error FillDistancesTrackPoint gave is of the wrong type. The type is %v", ty)
		}
	} else {
		t.Errorf("FillDistancesTrackPoint did not return a error, but was expected")
	}
}

func TestFillValuesTrackPointArrayValideCorrection(t *testing.T) {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	FillDistancesTrackPoint(&pnt2, pnt1, pnt3)
	err := FillValuesTrackPointArray([]TrackPoint{pnt1, pnt2, pnt3}, GetValideCorectionParamters()[0])
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

	fillCorectedElevationTrackPoint(pnts, "none")
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

	fillCorectedElevationTrackPoint(pnts, "linear")
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

	fillCorectedElevationTrackPoint(pnts, "steps")
	fillElevationGainLoseTrackPoint(pnts)
	fillCountUpDownWards(pnts, "steps")

	if pnts[1].VerticalDistanceBefore != 0.0 {
		t.Errorf("The VerticalDistanceBefore is %f, but %f was expected", pnts[1].VerticalDistanceBefore, 0.0)
	}

	if pnts[1].VerticalDistanceNext != 0.0 {
		t.Errorf("The VerticalDistanceNext is %f, but %f was expected", pnts[1].VerticalDistanceNext, 0.0)
	}
}
func TestFillDistancesThreePointWithUnkonwCorection(t *testing.T) {
	pnts := gerSimpleTrackPointArray()
	err := fillCorectedElevationTrackPoint(pnts, "asd")
	if err != nil {
		switch ty := err.(type) {
		case *CorectionParamterNotKnownError:
			fmt.Println("OK")
		default:
			t.Errorf("The Error FillCorectedElevationTrackPoint gave is of the wrong type. The type is %v", ty)
		}
	} else {
		t.Errorf("FillCorectedElevationTrackPoint did not return a error, but was expected")
	}
}

func TestFillDistancesTwoPointNext(t *testing.T) {
	pnt1 := TrackPoint{}
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	pnts := []TrackPoint{pnt1, pnt2, pnt3}

	FillDistancesTrackPoint(&pnts[1], pnts[0], pnts[2])

	fillCorectedElevationTrackPoint(pnts, "none")
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

	if seg.MaximumAtitute != 109.0 {
		t.Errorf("The MaximumAtitute is %f, but %f expected.", seg.MaximumAtitute, 109.0)
	}

	if seg.MinimumAtitute != 108.0 {
		t.Errorf("The MinimumAtitute is %f, but %f expected.", seg.MinimumAtitute, 108.0)
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

	if track.MaximumAtitute != 109.0 {
		t.Errorf("The MaximumAtitute is %f, but %f expected.", track.MaximumAtitute, 109.0)
	}

	if track.MinimumAtitute != 108.0 {
		t.Errorf("The MinimumAtitute is %f, but %f expected.", track.MinimumAtitute, 108.0)
	}
}

func TestFillTrackFileValuesSimple(t *testing.T) {
	file := TrackFile{}
	file.Tracks = []Track{getSimpleTrack()}
	FillTrackFileValues(&file)

	if file.Distance != 23.885148437468256 {
		t.Errorf("The Distance is %f, but %f expected.", file.Distance, 23.885148437468256)
	}

	if file.MaximumAtitute != 109.0 {
		t.Errorf("The MaximumAtitute is %f, but %f expected.", file.MaximumAtitute, 109.0)
	}

	if file.MinimumAtitute != 108.0 {
		t.Errorf("The MinimumAtitute is %f, but %f expected.", file.MinimumAtitute, 108.0)
	}
}

func getSimpleTrackFile() TrackFile {
	ret := NewTrackFile("/mys/track/file")
	trk := getSimpleTrack()
	FillTrackValues(&trk)
	ret.Tracks = []Track{trk}
	FillTrackFileValues(&ret)

	return ret
}

func getSimpleTrack() Track {
	ret := Track{}
	segs := getSimpleTrackSegment()
	FillTrackSegmentValues(&segs)
	ret.TrackSegments = []TrackSegment{segs}
	FillTrackValues(&ret)

	return ret
}

func getSimpleTrackSegment() TrackSegment {
	seg := TrackSegment{}
	points := gerSimpleTrackPointArray()
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
