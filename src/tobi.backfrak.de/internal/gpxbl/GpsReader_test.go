package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"tobi.backfrak.de/internal/testhelper"
)

func TestReadNotExistingGPX(t *testing.T) {
	_, err := ReadGPX(testhelper.GetValidGPX("NotExisting.gpx"))
	switch v := err.(type) {
	case nil:
		t.Errorf("No error when reading a not existing gpx file")
	case *os.PathError:
		fmt.Println("OK")
	default:
		t.Errorf("Expected a *os.PathError, got a %s", reflect.TypeOf(v))
	}
}

func TestReadInValidGPX(t *testing.T) {
	_, err := ReadGPX(testhelper.GetInvalidGPX("01.gpx"))
	switch v := err.(type) {
	case nil:
		t.Errorf("No error, when reading an invalid gpx file")
	case *xml.SyntaxError:
		fmt.Println("OK")
	default:
		t.Errorf("Expected a *xml.SyntaxError, got a %s", reflect.TypeOf(v))
	}

}

func TestReadNotGPX(t *testing.T) {

	gpx, err := ReadGPX(testhelper.GetInvalidGPX("02.gpx"))
	switch v := err.(type) {
	case nil:
		t.Errorf("No error, when reading an invalid gpx file")
	case *GpxFileError:
		checkGpxFileError(v, testhelper.GetInvalidGPX("02.gpx"), t)
	default:
		t.Errorf("Expected a *gpsabl.GpxFileError, got a %s", reflect.TypeOf(v))
	}

	fmt.Println(gpx.Name)
}

func TestGpxFileErrorStruct(t *testing.T) {

	path := "/some/sample/path"
	err := newGpxFileError(path)
	checkGpxFileError(err, path, t)
}

func checkGpxFileError(err *GpxFileError, path string, t *testing.T) {
	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error messaage of GpxFileError does not contain the expected Path")
	}

	if err.File != path {
		t.Errorf("The GpxFileError.File does not match the expected value")
	}
}

func TestReadValidMultiSegmentGPX(t *testing.T) {
	gpx, err := ReadGPX(testhelper.GetValidGPX("02.gpx"))
	if err != nil {
		t.Errorf("Something wrong when reading a valid gpx file: %s", err.Error())
	}

	if len(gpx.Tracks[0].TrackSegments) != 2 {
		t.Errorf("Expected 2 TrackSegments, got %d", len(gpx.Tracks[0].TrackSegments))
	}
}

func TestReadValidMultiTrackGPX(t *testing.T) {
	gpx, err := ReadGPX(testhelper.GetValidGPX("03.gpx"))
	if err != nil {
		t.Errorf("Something wrong when reading a valid gpx file: %s", err.Error())
	}

	if len(gpx.Tracks) != 5 {
		t.Errorf("Expected 5 Tracks, got %d", len(gpx.Tracks))
	}
}

func TestReadAllValidGPX(t *testing.T) {
	files, _ := ioutil.ReadDir(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-gpx"))

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".gpx") {
			if file.IsDir() == false {
				gpx, err := ReadGPX(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-gpx", file.Name()))
				if err != nil {
					t.Errorf("Got the following error while reading file %s: %s", filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-gpx", file.Name()), err.Error())
					return
				}
				if len(gpx.Tracks) < 1 {
					t.Errorf("The can not find tracks in %s.", filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-gpx", file.Name()))
				}
			}
		}
	}
}

func TestReadValidSimpleGPX(t *testing.T) {

	gpx, err := ReadGPX(testhelper.GetValidGPX("01.gpx"))
	if err != nil {
		t.Errorf("Something wrong when reading a valid gpx file: %s", err.Error())
	}

	if gpx.Name != "GPX name" {
		t.Errorf("The GPX name was not expected. Got: %s", gpx.Name)
	}

	if gpx.Description != "A valid GPX Track" {
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
