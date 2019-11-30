package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"testing"
)

func getTrk() Trk {
	track := Trk{}

	track.Name = "Test"
	track.Description = "A sample track"
	track.Number = 1

	segment := Trkseg{}
	var points = []Trkpt{}
	point1 := Trkpt{}
	point1.Elevation = 100.1
	point1.Latitude = 33.33001
	point1.Longitude = 33.33001
	points = append(points, point1)

	point2 := Trkpt{}
	point2.Elevation = 101.1
	point2.Latitude = 33.3302
	point2.Longitude = 33.3302
	points = append(points, point2)

	point3 := Trkpt{}
	point3.Elevation = 99.1
	point3.Latitude = 33.33009
	point3.Longitude = 33.330009
	points = append(points, point3)

	segment.TrackPoints = points
	var segments = []Trkseg{}
	segments = append(segments, segment)

	track.TrackSegments = segments

	return track

}

func TestConvertTrkBasicInfo(t *testing.T) {
	input := getTrk()

	track := ConvertTrk(input)

	if track.Name != input.Name {
		t.Errorf("track.Name  has not the expected value %s", input.Name)
	}

	if track.Description != input.Description {
		t.Errorf("track.Description  has not the expected value %s", input.Description)
	}

	if track.NumberOfSegments != 1 {
		t.Errorf("track.NumberOfSegments  has not the expected value %d but is %d", 1, track.NumberOfSegments)
	}

	if track.MinimumAtitute != 99.1 {
		t.Errorf("track.MinimumAtitute  has not the expected value %f but is %f", 99.1, track.MinimumAtitute)
	}

	if track.MaximumAtitute != 101.1 {
		t.Errorf("track.MaximumAtitute  has not the expected value %f but is %f", 101.1, track.MaximumAtitute)
	}

	if track.Distance != 49.32007928467905 {
		t.Errorf("track.Distance  has not the expected value %f but is %f", 49.32007928467905, track.Distance)
	}

	if track.ElevationLose != -2.0 {
		t.Errorf("track.ElevationLose  has not the expected value %f but is %f", -2.0, track.ElevationLose)
	}

	if track.ElevationGain != 1.0 {
		t.Errorf("track.ElevationGain  has not the expected value %f but is %f", 1.0, track.ElevationGain)
	}

	if track.UpwardsDistance != 27.65582137412336 {
		t.Errorf("track.UpwardsDistance  has not the expected value %f but is %f", 27.65582137412336, track.UpwardsDistance)
	}

	if track.DownwardsDistance != 21.664257910555698 {
		t.Errorf("track.DownwardsDistance  has not the expected value %f but is %f", 21.664257910555698, track.DownwardsDistance)
	}
}
