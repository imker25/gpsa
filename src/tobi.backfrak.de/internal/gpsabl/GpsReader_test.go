package gpsabl

// Copyright 2019 by Tobias Zellner. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"testing"

	"tobi.backfrak.de/internal/testhelper"
)

func TestReadValideGPX(t *testing.T) {
	file := testhelper.GetValideGPX("01.gpx")

	if file == "" {
		t.Errorf("Test failed, expected not to get an empty string")
	}

	gpx, err := ReadGPX(file)
	if err != nil {
		t.Errorf("Something wrong when reading a valide gpx file: %s", err.Error())
	}

	if gpx.Name != "GPX name" {
		t.Errorf("The GPX name was not expected. Got: %s", gpx.Name)
	}

	if gpx.Description != "A valide GPX Track" {
		t.Errorf("The GPX Description was not expected. Got: %s", gpx.Description)
	}

	if gpx.Track.Name != "Track name" {
		t.Errorf("The Track Name was not expected. Got: %s", gpx.Track.Name)
	}

	if gpx.Track.Number != 1 {
		t.Errorf("The Track Number was not expected. Got: %d", gpx.Track.Number)
	}

	if len(gpx.Track.TrackSegment.TrackPoints) != 637 {
		t.Errorf("The Number of track points was not expected. Got: %d", len(gpx.Track.TrackSegment.TrackPoints))
	}

	if gpx.Track.TrackSegment.TrackPoints[0].Elevation != 308.00100 {
		t.Errorf("The track point 0 Elevation was not expected. Got: %f", gpx.Track.TrackSegment.TrackPoints[0].Elevation)
	}

	if gpx.Track.TrackSegment.TrackPoints[0].Latitude != 49.41594200 {
		t.Errorf("The track point 0 Latitude was not expected. Got: %f", gpx.Track.TrackSegment.TrackPoints[0].Latitude)
	}

	if gpx.Track.TrackSegment.TrackPoints[0].Longitude != 11.01744700 {
		t.Errorf("The track point 0 Longitude was not expected. Got: %f", gpx.Track.TrackSegment.TrackPoints[0].Longitude)
	}
}
