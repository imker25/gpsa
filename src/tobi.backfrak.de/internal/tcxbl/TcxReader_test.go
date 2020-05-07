package tcxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"tobi.backfrak.de/internal/testhelper"
)

func TestReadInValidTcx(t *testing.T) {
	_, err := ReadTcx(testhelper.GetInvalidTcx("01.tcx"))
	switch v := err.(type) {
	case nil:
		t.Errorf("No error, when reading an invalid tcx file")
	case *xml.SyntaxError:
		fmt.Println("OK")
	default:
		t.Errorf("Expected a *xml.SyntaxError, got a %s", reflect.TypeOf(v))
	}

}

func TestReadAllInValidTcx(t *testing.T) {
	files, _ := ioutil.ReadDir(filepath.Join(testhelper.GetProjectRoot(), "testdata", "invalid-tcx"))

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tcx") {
			if file.IsDir() == false {
				_, err := ReadTcx(filepath.Join(testhelper.GetProjectRoot(), "testdata", "invalid-tcx", file.Name()))
				if err == nil {
					t.Errorf("Got no error while reading file %s.", filepath.Join(testhelper.GetProjectRoot(), "testdata", "invalid-tcx", file.Name()))

				}
			}
		}
	}
}

func TestReadAllValidTcx(t *testing.T) {
	files, _ := ioutil.ReadDir(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-tcx"))

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tcx") {
			if file.IsDir() == false {
				tcx, err := ReadTcx(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-tcx", file.Name()))
				if err != nil {
					t.Errorf("Got the error \"%s\" while reading file %s.", err.Error(), filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-tcx", file.Name()))

				}
				if len(tcx.ActivityArray) <= 0 {
					t.Errorf("ActivityArray array is empty")
				}
				if len(tcx.ActivityArray[0].Activities) <= 0 {
					t.Errorf("Activities array is empty")
				}
				if len(tcx.ActivityArray[0].Activities[0].Laps) <= 0 {
					t.Errorf("Laps array is empty")
				}

				if len(tcx.ActivityArray[0].Activities[0].Laps[0].Tracks) <= 0 {
					t.Errorf("Tracks array is empty")
				}

				if tcx.ActivityArray[0].Activities[0].ID == "" {
					t.Errorf("A Activity ID should never be \"\"")
				}
			}
		}
	}
}

func TestReadValidTcx01(t *testing.T) {
	tcx, err := ReadTcx(testhelper.GetValidTcx("01.tcx"))
	if err != nil {
		t.Errorf("Error, but none expected")
	}
	if len(tcx.ActivityArray) <= 0 {
		t.Errorf("The ActivityArray array is empty, but should not")
	}

	if len(tcx.ActivityArray[0].Activities) != 1 {
		t.Errorf("The Activity array has %d entries but should have %d", len(tcx.ActivityArray[0].Activities), 1)
	}

	if len(tcx.ActivityArray[0].Activities[0].Laps) != 7 {
		t.Errorf("The Laps array has %d entries, but should have %d", len(tcx.ActivityArray[0].Activities[0].Laps), 7)
	}

	if len(tcx.ActivityArray[0].Activities[0].Laps[0].Tracks) != 1 {
		t.Errorf("Track array has %d entries, but should have %d", len(tcx.ActivityArray[0].Activities[0].Laps[0].Tracks), 1)
	}

	if len(tcx.ActivityArray[0].Activities[0].Laps[0].Tracks[0].Trackpoints) != 48 {
		t.Errorf("Trackpoint array has %d entries, but should have %d", len(tcx.ActivityArray[0].Activities[0].Laps[0].Tracks[0].Trackpoints), 48)
	}

	if tcx.ActivityArray[0].Activities[0].ID != "2016-06-23T20:48:36Z" {
		t.Errorf("The Activity.ID \"is\" %s but should be \"2016-06-23T20:48:36Z\"", tcx.ActivityArray[0].Activities[0].ID)
	}

	if tcx.ActivityArray[0].Activities[0].Laps[0].DistanceMeters != "1000.0" {
		t.Errorf("The Laps[0].DistanceMeters \"is\" %s but should be \"1000.0\"", tcx.ActivityArray[0].Activities[0].Laps[0].DistanceMeters)
	}

	if tcx.ActivityArray[0].Activities[0].Laps[0].TotalTimeSeconds != "188.9" {
		t.Errorf("The Laps[0].TotalTimeSeconds is \"%s\" but should be \"188.9\"", tcx.ActivityArray[0].Activities[0].Laps[0].TotalTimeSeconds)
	}

	if tcx.ActivityArray[0].Activities[0].Laps[0].StartTime != "2016-06-05T10:45:18Z" {
		t.Errorf("The Laps[0].StartTime is \"%s\" but should be \"2016-06-05T10:45:18Z", tcx.ActivityArray[0].Activities[0].Laps[0].StartTime)
	}

	if tcx.ActivityArray[0].Activities[0].Laps[0].Tracks[0].Trackpoints[0].Time != "2016-06-05T10:45:59Z" {
		t.Errorf("The firsts Points time is \"%s\" but should be \"2016-06-05T10:45:59Z\"", tcx.ActivityArray[0].Activities[0].Laps[0].Tracks[0].Trackpoints[0].Time)
	}

	if tcx.ActivityArray[0].Activities[0].Laps[0].Tracks[0].Trackpoints[0].Position.LatitudeDegrees != 49.516803938895464 {
		t.Errorf("The firsts Points LatitudeDegrees is \"%f\" but should be \"49.516803938895464\"", tcx.ActivityArray[0].Activities[0].Laps[0].Tracks[0].Trackpoints[0].Position.LatitudeDegrees)
	}

	if tcx.ActivityArray[0].Activities[0].Laps[0].Tracks[0].Trackpoints[0].Position.LongitudeDegrees != 11.374258445575833 {
		t.Errorf("The firsts Points LongitudeDegrees is \"%f\" but should be \"11.374258445575833\"", tcx.ActivityArray[0].Activities[0].Laps[0].Tracks[0].Trackpoints[0].Position.LongitudeDegrees)
	}

	if tcx.ActivityArray[0].Activities[0].Laps[0].Tracks[0].Trackpoints[0].AltitudeMeters != 349.000000 {
		t.Errorf("The firsts Points AltitudeMeters is \"%f\" but should be \"349.000000\"", tcx.ActivityArray[0].Activities[0].Laps[0].Tracks[0].Trackpoints[0].AltitudeMeters)
	}

}
