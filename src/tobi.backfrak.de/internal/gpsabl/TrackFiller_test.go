package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import "testing"

func TestFillDistancesThreePoints(t *testing.T) {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	pnt2 = FillDistancesTrackPoint(pnt2, pnt1, pnt3)

	if pnt2.VerticalDistanceBefore != -1.0 {
		t.Errorf("The VerticalDistanceBefore is %f, but %f was expected", pnt2.VerticalDistanceBefore, -1.0)
	}

	if pnt2.HorizontalDistanceBefore != pnt2.HorizontalDistanceNext {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was expected", pnt2.HorizontalDistanceBefore, pnt2.HorizontalDistanceNext)
	}

	if pnt2.VerticalDistanceNext != 1.0 {
		t.Errorf("The VerticalDistanceNext is %f, but %f was expected", pnt2.VerticalDistanceNext, 1.0)
	}
}

func TestFillDistancesThreePointsBeforeAfter(t *testing.T) {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	lon := pnt2.Longitude
	lat := pnt2.Latitude
	eve := pnt2.Elevation

	pnt2 = FillDistancesTrackPoint(pnt2, pnt1, pnt3)

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

	pnt2 = FillDistancesTrackPoint(pnt2, pnt1, pnt3)

	if pnt2.VerticalDistanceBefore != -1.0 {
		t.Errorf("The VerticalDistanceBefore is %f, but %f was expected", pnt2.VerticalDistanceBefore, -1.0)
	}

	if pnt2.HorizontalDistanceBefore == 0.0 {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was not expected", pnt2.HorizontalDistanceBefore, 0.0)
	}

	if pnt2.HorizontalDistanceNext != 0.0 {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was expected", pnt2.HorizontalDistanceNext, 0.0)
	}

	if pnt2.VerticalDistanceNext != 0.0 {
		t.Errorf("The VerticalDistanceNext is %f, but %f was expected", pnt2.VerticalDistanceNext, 0.0)
	}
}

func TestFillDistancesTwoPointNext(t *testing.T) {
	pnt1 := TrackPoint{}
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	pnt2 = FillDistancesTrackPoint(pnt2, pnt1, pnt3)

	if pnt2.VerticalDistanceBefore != 0.0 {
		t.Errorf("The VerticalDistanceBefore is %f, but %f was expected", pnt2.VerticalDistanceBefore, -1.0)
	}

	if pnt2.HorizontalDistanceBefore != 0.0 {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was expected", pnt2.HorizontalDistanceBefore, 0.0)
	}

	if pnt2.HorizontalDistanceNext == 0.0 {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was not expected", pnt2.HorizontalDistanceNext, 0.0)
	}

	if pnt2.VerticalDistanceNext != 1.0 {
		t.Errorf("The VerticalDistanceNext is %f, but %f was expected", pnt2.VerticalDistanceNext, 0.0)
	}
}

func TestFillTrackSegmentValuesSimple(t *testing.T) {
	seg := getSimpleTrackSegment()

	oldPointNumber := len(seg.TrackPoints)

	seg = FillTrackSegmentValues(seg)

	if len(seg.TrackPoints) != oldPointNumber {
		seg = FillTrackSegmentValues(seg)
		t.Errorf("The number of track points changed during FillTrackSegmentValues() call")
	}

	if seg.Distance != 23.885148437468256 {
		t.Errorf("The Distance is %f, but %f expected.", seg.Distance, 23.885148437468256)
	}

	if seg.AtituteRange != 1.0 {
		t.Errorf("The AtituteRange is %f, but %f expected.", seg.AtituteRange, 1.0)
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

	track = FillTrackValues(track)

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

func TestFillTrackValuesSimple(t *testing.T) {
	track := Track{}
	track.TrackSegments = []TrackSegment{FillTrackSegmentValues(getSimpleTrackSegment())}
	track = FillTrackValues(track)

	if track.AtituteRange != 1.0 {
		t.Errorf("The AtituteRange is %f, but %f expected.", track.AtituteRange, 1.0)
	}

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

func getSimpleTrackSegment() TrackSegment {
	seg := TrackSegment{}
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	points := []TrackPoint{FillDistancesTrackPoint(pnt1, TrackPoint{}, pnt2), FillDistancesTrackPoint(pnt2, pnt1, pnt3), FillDistancesTrackPoint(pnt3, pnt2, TrackPoint{})}
	seg.TrackPoints = points

	return seg
}