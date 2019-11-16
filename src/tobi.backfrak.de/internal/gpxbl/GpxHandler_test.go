package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"tobi.backfrak.de/internal/testhelper"
)

func TestTrackReaderAllValideGPX(t *testing.T) {
	files, _ := ioutil.ReadDir(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valide-gpx"))

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".gpx") {
			if file.IsDir() == false {
				gpxFile := NewGpxFile(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valide-gpx", file.Name()))
				tracks, err := gpxFile.ReadTracks()
				if err != nil {
					t.Errorf("Got the following error while reading file %s: %s", filepath.Join(testhelper.GetProjectRoot(), "testdata", "valide-gpx", file.Name()), err.Error())
					return
				}
				if len(tracks) < 1 {
					t.Errorf("The can not find tracks in %s.", filepath.Join(testhelper.GetProjectRoot(), "testdata", "valide-gpx", file.Name()))
				}
			}
		}
	}
}

func TestTrackReaderOnePointTrack(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("06.gpx"))

	tracks, _ := gpx.ReadTracks()

	if tracks[0].Distance != 0.0 {
		t.Errorf("The Distance is %f, but %f was expected", tracks[0].Distance, 0.0)
	}

	if tracks[0].AtituteRange != 0.0 {
		t.Errorf("The AtituteRange is %f, but %f was expected", tracks[0].AtituteRange, 0.0)
	}
}

func TestTrackReader02(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("02.gpx"))

	tracks, _ := gpx.ReadTracks()

	if tracks[0].Distance != 37823.344979382266 {
		t.Errorf("The Distance is %f, but %f was expected", tracks[0].Distance, 18478.293509238614)
	}
}

func TestTrackReaderImpl(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("01.gpx"))

	if gpx.FilePath != testhelper.GetValideGPX("01.gpx") {
		t.Errorf("GpxFile.FilePath was not %s but %s", testhelper.GetValideGPX("01.gpx"), gpx.FilePath)
	}

	tracks, err := gpx.ReadTracks()

	if err != nil {
		t.Errorf("Got not expected error:  %s", err.Error())
	}

	if tracks == nil {
		t.Errorf("Got nil tracks when reading a valide file")
	}

	if len(tracks) != 1 {
		t.Errorf("The number of tracks was not %d, but was %d", 1, len(tracks))
	}

	if tracks[0].Distance != 18478.293509238614 {
		t.Errorf("The Distance is %f, but %f was expected", tracks[0].Distance, 18478.293509238614)
	}

	if tracks[0].AtituteRange != 104.0 {
		t.Errorf("The AtituteRange is %f, but %f was expected", tracks[0].AtituteRange, 104.00)
	}

	if tracks[0].MinimumAtitute != 298.0 {
		t.Errorf("The MinimumAtitute is %f, but %f was expected", tracks[0].MinimumAtitute, 298.00)
	}

	if tracks[0].MaximumAtitute != 402.0 {
		t.Errorf("The MaximumAtitute is %f, but %f was expected", tracks[0].MaximumAtitute, 402.00)
	}
}

func TestReadGpxFile(t *testing.T) {
	tracks, err := ReadGpxFile(testhelper.GetValideGPX("01.gpx"))
	if err != nil {
		t.Errorf("Something wrong when reading a valide gpx file: %s", err.Error())
	}

	if len(tracks) != 1 {
		t.Errorf("The number of tracks was not %d, but was %d", 1, len(tracks))
	}
}
