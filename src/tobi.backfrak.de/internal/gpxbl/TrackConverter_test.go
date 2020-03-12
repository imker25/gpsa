package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"testing"
	"time"
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

func getTrkWithTime() Trk {
	track := getTrk()

	track.TrackSegments[0].TrackPoints[0].Time = "2014-08-22T16:49:07Z"
	track.TrackSegments[0].TrackPoints[1].Time = "2014-08-22T16:49:17Z"
	track.TrackSegments[0].TrackPoints[2].Time = "2014-08-22T16:49:27Z"

	return track
}

func TestConvertTrkTimeInfo(t *testing.T) {
	input := getTrkWithTime()

	track, err := ConvertTrk(input, "none", 0.3, 10.0)
	if err != nil {
		t.Errorf("Got a error, but expected none. The error is: %s", err)
	}

	if track.Distance != 49.32007928467905 {
		t.Errorf("track.Distance  has not the expected value %f but is %f", 49.32007928467905, track.Distance)
	}

	for i := range track.TrackSegments[0].TrackPoints {
		if track.TrackSegments[0].TrackPoints[i].TimeValid == false {
			t.Errorf("track.TrackSegments[0].TrackPoints[%d].TimeValid is false but should not", i)
		}
		if track.TrackSegments[0].TrackPoints[i].Time.Format(time.RFC3339) != input.TrackSegments[0].TrackPoints[i].Time {
			t.Errorf("track.TrackSegments[0].TrackPoints[%d].Time is %s but should be %s", i, track.TrackSegments[0].TrackPoints[i].Time.Format(time.RFC3339), input.TrackSegments[0].TrackPoints[i].Time)
		}
	}
}

func TestConvertSegmentsWithNoPoints(t *testing.T) {
	inPut := getTrk()

	inPut.TrackSegments = append(inPut.TrackSegments, Trkseg{})

	outSegs, err := convertSegments(inPut.TrackSegments, "none", 0.3, 10.0)
	if err != nil {
		t.Errorf("Got a error, but expected none. The error is: %s", err)
	}

	if len(inPut.TrackSegments) != 2 {
		t.Errorf("The input had %d segments, but 2 was expected.", len(inPut.TrackSegments))
	}

	if len(outSegs) != 1 {
		t.Errorf("Got %d segments, but expected 1", len(outSegs))
	}
}

func TestConvertBasicPointValues(t *testing.T) {
	inPnt := getTrk().TrackSegments[0].TrackPoints[0]
	pnt := convertBasicPointValues(inPnt.Longitude, inPnt.Longitude, inPnt.Elevation, inPnt.Time)

	if pnt.TimeValid {
		t.Errorf("The TimeValid is true, but false is expected")
	}

	if inPnt.Elevation != pnt.Elevation {
		t.Errorf("Elevation does not match")
	}

	if inPnt.Longitude != pnt.Longitude {
		t.Errorf("Longitude does not match")
	}

	if inPnt.Latitude != pnt.Latitude {
		t.Errorf("Latitude does not match")
	}
}

func TestConvertBasicPointValuesWithTime(t *testing.T) {
	inPnt := getTrkWithTime().TrackSegments[0].TrackPoints[0]
	pnt := convertBasicPointValues(inPnt.Longitude, inPnt.Longitude, inPnt.Elevation, inPnt.Time)

	if !pnt.TimeValid {
		t.Errorf("The TimeValid is false, but true is expected")
	}

	if inPnt.Time != pnt.Time.Format(time.RFC3339) {
		t.Errorf("Time does not match")
	}

	if inPnt.Elevation != pnt.Elevation {
		t.Errorf("Elevation does not match")
	}

	if inPnt.Longitude != pnt.Longitude {
		t.Errorf("Longitude does not match")
	}

	if inPnt.Latitude != pnt.Latitude {
		t.Errorf("Latitude does not match")
	}
}

func TestConvertBasicPointValuesWithUnFormatedTime(t *testing.T) {
	inPnt := getTrkWithTime().TrackSegments[0].TrackPoints[0]
	inPnt.Time = "lala"
	pnt := convertBasicPointValues(inPnt.Longitude, inPnt.Longitude, inPnt.Elevation, inPnt.Time)

	if pnt.TimeValid {
		t.Errorf("The TimeValid is true, but false is expected")
	}

	if inPnt.Elevation != pnt.Elevation {
		t.Errorf("Elevation does not match")
	}

	if inPnt.Longitude != pnt.Longitude {
		t.Errorf("Longitude does not match")
	}

	if inPnt.Latitude != pnt.Latitude {
		t.Errorf("Latitude does not match")
	}
}

func TestConvertTrkBasicInfo(t *testing.T) {
	input := getTrk()

	track, err := ConvertTrk(input, "none", 0.3, 10.0)
	if err != nil {
		t.Errorf("Got a error, but expected none. The error is: %s", err)
	}

	if track.Name != input.Name {
		t.Errorf("track.Name has not the expected value %s", input.Name)
	}

	if track.Description != input.Description {
		t.Errorf("track.Description has not the expected value %s", input.Description)
	}

	if track.NumberOfSegments != 1 {
		t.Errorf("track.NumberOfSegments has not the expected value %d but is %d", 1, track.NumberOfSegments)
	}

	if track.MinimumAltitude != 99.1 {
		t.Errorf("track.MinimumAltitude has not the expected value %f but is %f", 99.1, track.MinimumAltitude)
	}

	if track.MaximumAltitude != 101.1 {
		t.Errorf("track.MaximumAltitude has not the expected value %f but is %f", 101.1, track.MaximumAltitude)
	}

	if track.Distance != 49.32007928467905 {
		t.Errorf("track.Distance has not the expected value %f but is %f", 49.32007928467905, track.Distance)
	}

	if track.ElevationLose != -2.0 {
		t.Errorf("track.ElevationLose has not the expected value %f but is %f", -2.0, track.ElevationLose)
	}

	if track.ElevationGain != 1.0 {
		t.Errorf("track.ElevationGain has not the expected value %f but is %f", 1.0, track.ElevationGain)
	}

	if track.UpwardsDistance != 27.65582137412336 {
		t.Errorf("track.UpwardsDistance has not the expected value %f but is %f", 27.65582137412336, track.UpwardsDistance)
	}

	if track.DownwardsDistance != 21.664257910555698 {
		t.Errorf("track.DownwardsDistance has not the expected value %f but is %f", 21.664257910555698, track.DownwardsDistance)
	}

	for i := range track.TrackSegments[0].TrackPoints {
		if i > 0 {
			if track.TrackSegments[0].TrackPoints[i].DistanceToThisPoint <= track.TrackSegments[0].TrackPoints[i-1].DistanceToThisPoint {
				t.Errorf("The DistanceToThisPoint for point %d, is %f but the point before had %f", i, track.TrackSegments[0].TrackPoints[i].DistanceToThisPoint, track.TrackSegments[0].TrackPoints[i-1].DistanceToThisPoint)
			}
		}
		if track.TrackSegments[0].TrackPoints[i].TimeValid == true {
			t.Errorf("track.TrackSegments[0].TrackPoints[%d].TimeValid is true but should not", i)
		}
	}
}
