package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"encoding/xml"
	"fmt"
	"os"
	"reflect"
	"testing"

	"tobi.backfrak.de/internal/testhelper"
)

func TestReadNotExistingGPX(t *testing.T) {
	file := testhelper.GetValideGPX("NotExisting.gpx")

	if file == "" {
		t.Errorf("Test failed, expected not to get an empty string")
	}

	_, err := ReadGPX(file)
	switch v := err.(type) {
	case nil:
		t.Errorf("No error, when reading a unvalide gpx file")
	case *os.PathError:
		fmt.Println("OK")
	default:
		t.Errorf("Expected a *os.PathError, got a %s", reflect.TypeOf(v))
	}
}

func TestReadUnValideGPX(t *testing.T) {
	file := testhelper.GetUnValideGPX("01.gpx")

	if file == "" {
		t.Errorf("Test failed, expected not to get an empty string")
	}

	_, err := ReadGPX(file)
	switch v := err.(type) {
	case nil:
		t.Errorf("No error, when reading a unvalide gpx file")
	case *xml.SyntaxError:
		fmt.Println("OK")
	default:
		t.Errorf("Expected a *xml.SyntaxError, got a %s", reflect.TypeOf(v))
	}

}

func TestReadValideSimpleGPX(t *testing.T) {
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

	if gpx.Tracks[0].Name != "Track name" {
		t.Errorf("The Track Name was not expected. Got: %s", gpx.Tracks[0].Name)
	}

	if gpx.Tracks[0].Number != 1 {
		t.Errorf("The Track Number was not expected. Got: %d", gpx.Tracks[0].Number)
	}

	if len(gpx.Tracks[0].TrackSegments[0].TrackPoints) != 637 {
		t.Errorf("The Number of track points was not expected. Got: %d", len(gpx.Tracks[0].TrackSegments[0].TrackPoints))
	}

	if gpx.Tracks[0].TrackSegments[0].TrackPoints[0].Elevation != 308.00100 {
		t.Errorf("The track point 0 Elevation was not expected. Got: %f", gpx.Tracks[0].TrackSegments[0].TrackPoints[0].Elevation)
	}

	if gpx.Tracks[0].TrackSegments[0].TrackPoints[0].Latitude != 49.41594200 {
		t.Errorf("The track point 0 Latitude was not expected. Got: %f", gpx.Tracks[0].TrackSegments[0].TrackPoints[0].Latitude)
	}

	if gpx.Tracks[0].TrackSegments[0].TrackPoints[0].Longitude != 11.01744700 {
		t.Errorf("The track point 0 Longitude was not expected. Got: %f", gpx.Tracks[0].TrackSegments[0].TrackPoints[0].Longitude)
	}
}
