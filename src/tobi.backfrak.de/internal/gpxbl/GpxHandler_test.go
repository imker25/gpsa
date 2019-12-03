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
				trackFile, err := gpxFile.ReadTracks("linear")
				if err != nil {
					t.Errorf("Got the following error while reading file %s: %s", filepath.Join(testhelper.GetProjectRoot(), "testdata", "valide-gpx", file.Name()), err.Error())
					return
				}
				if len(trackFile.Tracks) < 1 {
					t.Errorf("The can not find tracks in %s.", filepath.Join(testhelper.GetProjectRoot(), "testdata", "valide-gpx", file.Name()))
				}

				for _, track := range gpxFile.Tracks {
					for _, seg := range track.TrackSegments {
						for i := range seg.TrackPoints {
							if i > 0 {
								if seg.TrackPoints[i].DistanceToThisPoint <= seg.TrackPoints[i-1].DistanceToThisPoint {
									t.Errorf("File %s: The DistanceToThisPoint for point %d, is %f but the point before had %f", gpxFile.FilePath, i, seg.TrackPoints[i].DistanceToThisPoint, seg.TrackPoints[i-1].DistanceToThisPoint)
								}
							}
						}
					}
				}

				for _, track := range trackFile.Tracks {
					for _, seg := range track.TrackSegments {
						for i := range seg.TrackPoints {
							if i > 0 {
								if seg.TrackPoints[i].DistanceToThisPoint <= seg.TrackPoints[i-1].DistanceToThisPoint {
									t.Errorf("File %s: The DistanceToThisPoint for point %d, is %f but the point before had %f", trackFile.FilePath, i, seg.TrackPoints[i].DistanceToThisPoint, seg.TrackPoints[i-1].DistanceToThisPoint)
								}
							}
						}
					}
				}
			}
		}
	}
}

func TestTrackReaderOnePointTrack(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("06.gpx"))

	file, _ := gpx.ReadTracks("none")

	if file.Tracks[0].Distance != 0.0 {
		t.Errorf("The Distance is %f, but %f was expected", file.Tracks[0].Distance, 0.0)
	}

	if file.Tracks[0].UpwardsDistance != 0.0 {
		t.Errorf("The AtituteRange is %f, but %f was expected", file.Tracks[0].UpwardsDistance, 0.0)
	}
}

func TestTrackReader02(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("02.gpx"))

	file, _ := gpx.ReadTracks("none")

	if file.Tracks[0].Distance != 37823.344979382266 {
		t.Errorf("The Distance is %f, but %f was expected", file.Tracks[0].Distance, 37823.344979382266)
	}
}

func TestTrackReaderUnValideCorrectionParameter(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("02.gpx"))

	_, err := gpx.ReadTracks("asdfg")
	if err != nil {
		switch ty := err.(type) {
		case *gpsabl.CorectionParamterNotKnownError:
			fmt.Println("OK")
		default:
			t.Errorf("The Error ReadTracks gave is of the wrong type. The type is %v", ty)
		}
	} else {
		t.Errorf("ReadTracks did not return a error, but was expected")
	}
}

func TestTrackReaderImpl(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("01.gpx"))

	if gpx.FilePath != testhelper.GetValideGPX("01.gpx") {
		t.Errorf("GpxFile.FilePath was not %s but %s", testhelper.GetValideGPX("01.gpx"), gpx.FilePath)
	}

	file, err := gpx.ReadTracks("linear")

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

	if file.Tracks[0].MinimumAtitute != 298.0 {
		t.Errorf("The MinimumAtitute is %f, but %f was expected", file.Tracks[0].MinimumAtitute, 298.00)
	}

	if file.Tracks[0].MaximumAtitute != 402.0 {
		t.Errorf("The MaximumAtitute is %f, but %f was expected", file.Tracks[0].MaximumAtitute, 402.00)
	}

	if file.Tracks[0].ElevationLose != -257.20975 {
		t.Errorf("The ElevationLose is %f, but %f was expected", file.Tracks[0].ElevationLose, 306.00)
	}

	if file.Tracks[0].ElevationGain != 278.20874 {
		t.Errorf("The ElevationGain is %f, but %f was expected", file.Tracks[0].ElevationGain, 326.999)
	}

	if file.Tracks[0].DownwardsDistance != 9152.075973681809 {
		t.Errorf("The DownwardsDistance is %f, but %f was expected", file.Tracks[0].DownwardsDistance, 9152.075973681809)
	}

	if file.Tracks[0].UpwardsDistance != 8038.332888190817 {
		t.Errorf("The UpwardsDistance is %f, but %f was expected", file.Tracks[0].UpwardsDistance, 8038.332888190817)
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

	if file.MinimumAtitute != file.Tracks[0].MinimumAtitute {
		t.Errorf("The MinimumAtitute is %f, but %f was expected", file.MinimumAtitute, file.Tracks[0].MinimumAtitute)
	}

	if file.MaximumAtitute != file.Tracks[0].MaximumAtitute {
		t.Errorf("The MaximumAtitute is %f, but %f was expected", file.MaximumAtitute, file.Tracks[0].MaximumAtitute)
	}

	if file.ElevationGain != file.Tracks[0].ElevationGain {
		t.Errorf("The ElevationGain is %f, but %f was expected", file.ElevationGain, file.Tracks[0].ElevationGain)
	}

	if file.ElevationLose != file.Tracks[0].ElevationLose {
		t.Errorf("The ElevationLose is %f, but %f was expected", file.ElevationLose, file.Tracks[0].ElevationLose)
	}

	if file.DownwardsDistance != file.Tracks[0].DownwardsDistance {
		t.Errorf("The DownwardsDistance is %f, but %f was expected", file.DownwardsDistance, file.Tracks[0].DownwardsDistance)
	}

	if file.UpwardsDistance != file.Tracks[0].UpwardsDistance {
		t.Errorf("The UpwardsDistance is %f, but %f was expected", file.UpwardsDistance, file.Tracks[0].UpwardsDistance)
	}
}

func TestGpxFileInterfaceImplentaion1(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("01.gpx"))

	reader := gpsabl.TrackReader(&gpx)

	file, err := reader.ReadTracks("none")

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

	if info.GetAtituteRange() != file.Tracks[0].GetAtituteRange() {
		t.Errorf("The AtituteRange is %f, but %f was expected", info.GetAtituteRange(), file.Tracks[0].GetAtituteRange())
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

	file, err := reader.ReadTracks("none")

	if err != nil {
		t.Errorf("Got not expected error:  %s", err.Error())
	}

	info := gpsabl.TrackSummaryProvider(&gpx)

	if info.GetDistance() != file.Tracks[0].Distance {
		t.Errorf("The Distance is %f, but %f was expected", info.GetDistance(), file.Tracks[0].Distance)
	}

	if info.GetAtituteRange() != file.Tracks[0].GetAtituteRange() {
		t.Errorf("The AtituteRange is %f, but %f was expected", info.GetAtituteRange(), file.Tracks[0].GetAtituteRange())
	}

	if info.GetMinimumAtitute() != file.Tracks[0].MinimumAtitute {
		t.Errorf("The MinimumAtitute is %f, but %f was expected", info.GetMinimumAtitute(), file.Tracks[0].MinimumAtitute)
	}

	if info.GetMaximumAtitute() != file.Tracks[0].MaximumAtitute {
		t.Errorf("The MaximumAtitute is %f, but %f was expected", info.GetMaximumAtitute(), file.Tracks[0].MaximumAtitute)
	}

	if info.GetDownwardsDistance() != file.Tracks[0].DownwardsDistance {
		t.Errorf("The GetDownwardsDistance is %f, but %f was expected", info.GetDownwardsDistance(), file.Tracks[0].DownwardsDistance)
	}

	if info.GetElevationGain() != file.Tracks[0].ElevationGain {
		t.Errorf("The GetElevationGain is %f, but %f was expected", info.GetElevationGain(), file.Tracks[0].ElevationGain)
	}

	if info.GetElevationLose() != file.Tracks[0].ElevationLose {
		t.Errorf("The GetElevationLose is %f, but %f was expected", info.GetElevationLose(), file.Tracks[0].ElevationLose)
	}

	if info.GetUpwardsDistance() != file.Tracks[0].UpwardsDistance {
		t.Errorf("The GetUpwardsDistance is %f, but %f was expected", info.GetUpwardsDistance(), file.Tracks[0].UpwardsDistance)
	}

	if gpx.Distance != file.Tracks[0].Distance {
		t.Errorf("The Distance is %f, but %f was expected", gpx.Distance, file.Tracks[0].Distance)
	}

	if gpx.MinimumAtitute != file.Tracks[0].MinimumAtitute {
		t.Errorf("The MinimumAtitute is %f, but %f was expected", gpx.MinimumAtitute, file.Tracks[0].MinimumAtitute)
	}

	if gpx.MaximumAtitute != file.Tracks[0].MaximumAtitute {
		t.Errorf("The MaximumAtitute is %f, but %f was expected", gpx.MaximumAtitute, file.Tracks[0].MaximumAtitute)
	}
}

func TestReadGpxFile(t *testing.T) {
	file, err := ReadGpxFile(testhelper.GetValideGPX("01.gpx"), "none")
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

	if file.MinimumAtitute != file.Tracks[0].MinimumAtitute {
		t.Errorf("The MinimumAtitute is %f, but %f was expected", file.MinimumAtitute, file.Tracks[0].MinimumAtitute)
	}

	if file.MaximumAtitute != file.Tracks[0].MaximumAtitute {
		t.Errorf("The MaximumAtitute is %f, but %f was expected", file.MaximumAtitute, file.Tracks[0].MaximumAtitute)
	}

	if file.ElevationGain != file.Tracks[0].ElevationGain {
		t.Errorf("The ElevationGain is %f, but %f was expected", file.ElevationGain, file.Tracks[0].ElevationGain)
	}

	if file.ElevationLose != file.Tracks[0].ElevationLose {
		t.Errorf("The ElevationLose is %f, but %f was expected", file.ElevationLose, file.Tracks[0].ElevationLose)
	}

	if file.UpwardsDistance != file.Tracks[0].UpwardsDistance {
		t.Errorf("The UpwardsDistance is %f, but %f was expected", file.UpwardsDistance, file.Tracks[0].UpwardsDistance)
	}

	if file.DownwardsDistance != file.Tracks[0].DownwardsDistance {
		t.Errorf("The DownwardsDistance is %f, but %f was expected", file.DownwardsDistance, file.Tracks[0].DownwardsDistance)
	}
}

func TestReadTracksNotExistingGPX(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValideGPX("NotExisting.gpx"))
	_, err := gpx.ReadTracks("none")
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
	_, err := gpx.ReadTracks("none")
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
	file, err := gpx.ReadTracks("none")
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
