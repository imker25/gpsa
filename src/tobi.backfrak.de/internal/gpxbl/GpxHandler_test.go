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

	"tobi.backfrak.de/internal/gpsabl"

	"tobi.backfrak.de/internal/testhelper"
)

func TestTrackReaderAllValideGPX(t *testing.T) {
	files, _ := ioutil.ReadDir(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valide-gpx"))

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".gpx") {
			if file.IsDir() == false {
				gpxFile := NewGpxFile(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valide-gpx", file.Name()))
				trackFile, err := gpxFile.ReadTracks()
				if err != nil {
					t.Errorf("Got the following error while reading file %s: %s", filepath.Join(testhelper.GetProjectRoot(), "testdata", "valide-gpx", file.Name()), err.Error())
					return
				}
				if len(trackFile.Tracks) < 1 {
					t.Errorf("The can not find tracks in %s.", filepath.Join(testhelper.GetProjectRoot(), "testdata", "valide-gpx", file.Name()))
				}
			}
		}
	}
}

func TestTrackReaderOnePointTrack(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("06.gpx"))

	file, _ := gpx.ReadTracks()

	if file.Tracks[0].Distance != 0.0 {
		t.Errorf("The Distance is %f, but %f was expected", file.Tracks[0].Distance, 0.0)
	}

	if file.Tracks[0].AtituteRange != 0.0 {
		t.Errorf("The AtituteRange is %f, but %f was expected", file.Tracks[0].AtituteRange, 0.0)
	}
}

func TestTrackReader02(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("02.gpx"))

	file, _ := gpx.ReadTracks()

	if file.Tracks[0].Distance != 37823.344979382266 {
		t.Errorf("The Distance is %f, but %f was expected", file.Tracks[0].Distance, 18478.293509238614)
	}
}

func TestTrackReaderImpl(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("01.gpx"))

	if gpx.FilePath != testhelper.GetValideGPX("01.gpx") {
		t.Errorf("GpxFile.FilePath was not %s but %s", testhelper.GetValideGPX("01.gpx"), gpx.FilePath)
	}

	file, err := gpx.ReadTracks()

	if err != nil {
		t.Errorf("Got not expected error:  %s", err.Error())
	}

	if file.Tracks == nil {
		t.Errorf("Got nil tracks when reading a valide file")
	}

	if len(file.Tracks) != 1 {
		t.Errorf("The number of tracks was not %d, but was %d", 1, len(file.Tracks))
	}

	if file.Tracks[0].Distance != 18478.293509238614 {
		t.Errorf("The Distance is %f, but %f was expected", file.Tracks[0].Distance, 18478.293509238614)
	}

	if file.Tracks[0].AtituteRange != 104.0 {
		t.Errorf("The AtituteRange is %f, but %f was expected", file.Tracks[0].AtituteRange, 104.00)
	}

	if file.Tracks[0].MinimumAtitute != 298.0 {
		t.Errorf("The MinimumAtitute is %f, but %f was expected", file.Tracks[0].MinimumAtitute, 298.00)
	}

	if file.Tracks[0].MaximumAtitute != 402.0 {
		t.Errorf("The MaximumAtitute is %f, but %f was expected", file.Tracks[0].MaximumAtitute, 402.00)
	}

	if file.FilePath != testhelper.GetValideGPX("01.gpx") {
		t.Errorf("The FilePath is %s, but %s was expected", file.FilePath, testhelper.GetValideGPX("01.gpx"))
	}

	if file.NumberOfTracks != 1 {
		t.Errorf("The NumberOfTracks is %d, but %d was expected", file.NumberOfTracks, 1)
	}

	if file.Distance != file.Tracks[0].Distance {
		t.Errorf("The Distance is %f, but %f was expected", file.Distance, file.Tracks[0].Distance)
	}

	if file.AtituteRange != file.Tracks[0].AtituteRange {
		t.Errorf("The AtituteRange is %f, but %f was expected", file.AtituteRange, file.Tracks[0].AtituteRange)
	}

	if file.MinimumAtitute != file.Tracks[0].MinimumAtitute {
		t.Errorf("The MinimumAtitute is %f, but %f was expected", file.MinimumAtitute, file.Tracks[0].MinimumAtitute)
	}

	if file.MaximumAtitute != file.Tracks[0].MaximumAtitute {
		t.Errorf("The MaximumAtitute is %f, but %f was expected", file.MaximumAtitute, file.Tracks[0].MaximumAtitute)
	}
}

func TestGpxFileInterfaceImplentaion1(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("01.gpx"))

	reader := gpsabl.TrackReader(&gpx)

	file, err := reader.ReadTracks()

	if err != nil {
		t.Errorf("Got not expected error:  %s", err.Error())
	}

	if file.Name != "GPX name" {
		t.Errorf("The GPX name was not expected. Got: %s", file.Name)
	}

	info := gpsabl.TrackSummaryProvider(&file)

	if info.GetDistance() != file.Tracks[0].Distance {
		t.Errorf("The Distance is %f, but %f was expected", info.GetDistance(), file.Tracks[0].Distance)
	}

	if info.GetAtituteRange() != file.Tracks[0].AtituteRange {
		t.Errorf("The AtituteRange is %f, but %f was expected", info.GetAtituteRange(), file.Tracks[0].AtituteRange)
	}

	if info.GetMinimumAtitute() != file.Tracks[0].MinimumAtitute {
		t.Errorf("The MinimumAtitute is %f, but %f was expected", info.GetMinimumAtitute(), file.Tracks[0].MinimumAtitute)
	}

	if info.GetMaximumAtitute() != file.Tracks[0].MaximumAtitute {
		t.Errorf("The MaximumAtitute is %f, but %f was expected", info.GetMaximumAtitute(), file.Tracks[0].MaximumAtitute)
	}
}

func TestGpxFileInterfaceImplentaion2(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("01.gpx"))

	reader := gpsabl.TrackReader(&gpx)

	file, err := reader.ReadTracks()

	if err != nil {
		t.Errorf("Got not expected error:  %s", err.Error())
	}

	info := gpsabl.TrackSummaryProvider(&gpx)

	if info.GetDistance() != file.Tracks[0].Distance {
		t.Errorf("The Distance is %f, but %f was expected", info.GetDistance(), file.Tracks[0].Distance)
	}

	if info.GetAtituteRange() != file.Tracks[0].AtituteRange {
		t.Errorf("The AtituteRange is %f, but %f was expected", info.GetAtituteRange(), file.Tracks[0].AtituteRange)
	}

	if info.GetMinimumAtitute() != file.Tracks[0].MinimumAtitute {
		t.Errorf("The MinimumAtitute is %f, but %f was expected", info.GetMinimumAtitute(), file.Tracks[0].MinimumAtitute)
	}

	if info.GetMaximumAtitute() != file.Tracks[0].MaximumAtitute {
		t.Errorf("The MaximumAtitute is %f, but %f was expected", info.GetMaximumAtitute(), file.Tracks[0].MaximumAtitute)
	}

	if gpx.Distance != file.Tracks[0].Distance {
		t.Errorf("The Distance is %f, but %f was expected", gpx.Distance, file.Tracks[0].Distance)
	}

	if gpx.AtituteRange != file.Tracks[0].AtituteRange {
		t.Errorf("The AtituteRange is %f, but %f was expected", gpx.AtituteRange, file.Tracks[0].AtituteRange)
	}

	if gpx.MinimumAtitute != file.Tracks[0].MinimumAtitute {
		t.Errorf("The MinimumAtitute is %f, but %f was expected", gpx.MinimumAtitute, file.Tracks[0].MinimumAtitute)
	}

	if gpx.MaximumAtitute != file.Tracks[0].MaximumAtitute {
		t.Errorf("The MaximumAtitute is %f, but %f was expected", gpx.MaximumAtitute, file.Tracks[0].MaximumAtitute)
	}
}

func TestReadGpxFile(t *testing.T) {
	file, err := ReadGpxFile(testhelper.GetValideGPX("01.gpx"))
	if err != nil {
		t.Errorf("Something wrong when reading a valide gpx file: %s", err.Error())
	}

	if len(file.Tracks) != 1 {
		t.Errorf("The number of tracks was not %d, but was %d", 1, len(file.Tracks))
	}

	if file.Name != "GPX name" {
		t.Errorf("The GPX name was not expected. Got: %s", file.Name)
	}

	if file.Description != "A valide GPX Track" {
		t.Errorf("The GPX Description was not expected. Got: %s", file.Description)
	}

	if file.Tracks[0].Name != "Track name" {
		t.Errorf("The Track Name was not expected. Got: %s", file.Tracks[0].Name)
	}

	if len(file.Tracks[0].TrackSegments[0].TrackPoints) != 637 {
		t.Errorf("The Number of track points was not expected. Got: %d", len(file.Tracks[0].TrackSegments[0].TrackPoints))
	}

	if file.Tracks[0].TrackSegments[0].TrackPoints[0].Elevation != 308.00100 {
		t.Errorf("The track point 0 Elevation was not expected. Got: %f", file.Tracks[0].TrackSegments[0].TrackPoints[0].Elevation)
	}

	if file.Tracks[0].TrackSegments[0].TrackPoints[0].Latitude != 49.41594200 {
		t.Errorf("The track point 0 Latitude was not expected. Got: %f", file.Tracks[0].TrackSegments[0].TrackPoints[0].Latitude)
	}

	if file.Tracks[0].TrackSegments[0].TrackPoints[0].Longitude != 11.01744700 {
		t.Errorf("The track point 0 Longitude was not expected. Got: %f", file.Tracks[0].TrackSegments[0].TrackPoints[0].Longitude)
	}

	if file.FilePath != testhelper.GetValideGPX("01.gpx") {
		t.Errorf("The FilePath is %s, but %s was expected", file.FilePath, testhelper.GetValideGPX("01.gpx"))
	}

	if file.NumberOfTracks != 1 {
		t.Errorf("The NumberOfTracks is %d, but %d was expected", file.NumberOfTracks, 1)
	}

	if file.Distance != file.Tracks[0].Distance {
		t.Errorf("The Distance is %f, but %f was expected", file.Distance, file.Tracks[0].Distance)
	}

	if file.AtituteRange != file.Tracks[0].AtituteRange {
		t.Errorf("The AtituteRange is %f, but %f was expected", file.AtituteRange, file.Tracks[0].AtituteRange)
	}

	if file.MinimumAtitute != file.Tracks[0].MinimumAtitute {
		t.Errorf("The MinimumAtitute is %f, but %f was expected", file.MinimumAtitute, file.Tracks[0].MinimumAtitute)
	}

	if file.MaximumAtitute != file.Tracks[0].MaximumAtitute {
		t.Errorf("The MaximumAtitute is %f, but %f was expected", file.MaximumAtitute, file.Tracks[0].MaximumAtitute)
	}
}

func TestReadTracksNotExistingGPX(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("NotExisting.gpx"))
	_, err := gpx.ReadTracks()
	switch v := err.(type) {
	case nil:
		t.Errorf("No error, when reading a not existing gpx file")
	case *os.PathError:
		fmt.Println("OK")
	default:
		t.Errorf("Expected a *os.PathError, got a %s", reflect.TypeOf(v))
	}
}

func TestReadTracksUnValideGPX(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetUnValideGPX("01.gpx"))
	_, err := gpx.ReadTracks()
	switch v := err.(type) {
	case nil:
		t.Errorf("No error, when reading a unvalide gpx file")
	case *xml.SyntaxError:
		fmt.Println("OK")
	default:
		t.Errorf("Expected a *xml.SyntaxError, got a %s", reflect.TypeOf(v))
	}

}

func TestReadTracksNotGPX(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetUnValideGPX("02.gpx"))
	file, err := gpx.ReadTracks()
	switch v := err.(type) {
	case nil:
		t.Errorf("No error, when reading a unvalide gpx file")
	case *GpxFileError:
		checkGpxFileError(v, testhelper.GetUnValideGPX("02.gpx"), t)
	default:
		t.Errorf("Expected a *gpsabl.GpxFileError, got a %s", reflect.TypeOf(v))
	}

	fmt.Println(file.Name)
}
