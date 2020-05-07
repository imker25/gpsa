package tcxbl

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"tobi.backfrak.de/internal/gpsabl"
	"tobi.backfrak.de/internal/testhelper"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestTrackReader02(t *testing.T) {
	tcx := NewTcxFile(testhelper.GetValidTcx("01.tcx"))

	file, _ := tcx.ReadTracks("none", 0.3, 10.0)

	if file.GetDistance() != 6216.201383825188 {
		t.Errorf("The Distance is %f, but should be %f", file.GetDistance(), 6216.201383825188)
	}

	if file.GetMovingTime() != 1103000000000 {
		t.Errorf("The MovingTime is %d, but should be %d", file.GetMovingTime(), 1103000000000)
	}

	if file.GetStartTime().Format(time.RFC3339) != "2016-06-05T10:45:59Z" {
		t.Errorf("The StartTime is %s, but should be %s", file.GetStartTime().Format(time.RFC3339), "2016-06-05T10:45:59Z")
	}

	if file.GetUpwardsSpeed() != 4.797203586378565 {
		t.Errorf("The UpwardsSpeed is %f, but should be %f", file.GetUpwardsSpeed(), 4.797204)
	}

	if file.GetDownwardsSpeed() != 7.683785245071919 {
		t.Errorf("The UpwardsSpeed is %f, but should be %f", file.GetDownwardsSpeed(), 7.683785)
	}
}

func TestTrackReaderEmptyTrack(t *testing.T) {
	tcx := NewTcxFile(testhelper.GetInvalidTcx("03.tcx"))

	_, err := tcx.ReadTracks("none", 0.3, 10.0)
	if err != nil {
		switch ty := err.(type) {
		case *EmptyTcxFileError:
			fmt.Println("OK")
		default:
			t.Errorf("The Error ReadTracks gave is of the wrong type. The type is %v", ty)
		}
	} else {
		t.Errorf("ReadTracks did not return a error, but was expected")
	}
}

func TestTrackReaderInValidCorrectionParameter(t *testing.T) {
	tcx := NewTcxFile(testhelper.GetValidTcx("02.tcx"))

	_, err := tcx.ReadTracks("asdfg", 0.3, 10.0)
	if err != nil {
		switch ty := err.(type) {
		case *gpsabl.CorrectionParameterNotKnownError:
			fmt.Println("OK")
		default:
			t.Errorf("The Error ReadTracks gave is of the wrong type. The type is %v", ty)
		}
	} else {
		t.Errorf("ReadTracks did not return a error, but was expected")
	}
}

func TestReadAllValidTcxDueInterface(t *testing.T) {
	files, _ := ioutil.ReadDir(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-tcx"))

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tcx") {
			if file.IsDir() == false {
				tcxFile := NewTcxFile(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-tcx", file.Name()))
				iTcx := gpsabl.TrackReader(&tcxFile)

				track, err := iTcx.ReadTracks("none", 0.3, 10.0)
				if err != nil {
					t.Errorf("Got the error \"%s\" while reading file %s.", err.Error(), filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-tcx", file.Name()))

				}
				if track.Distance <= 0.0 {
					t.Errorf("The track.Distance is %f but should not be less then %f", track.Distance, 0.0)
				}
			}
		}
	}
}
