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

func TestTrackReaderAllValidGPX(t *testing.T) {
	files, _ := ioutil.ReadDir(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-gpx"))

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".gpx") {
			if file.IsDir() == false {
				gpxFile := NewGpxFile(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-gpx", file.Name()))
				trackFile, err := gpxFile.ReadTracks("linear", 0.3, 10.0)
				if err != nil {
					t.Errorf("Got the following error while reading file %s: %s", filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-gpx", file.Name()), err.Error())
					return
				}
				if len(trackFile.Tracks) < 1 {
					t.Errorf("The can not find tracks in %s.", filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-gpx", file.Name()))
				}

				for _, track := range gpxFile.Tracks {
					for _, seg := range track.TrackSegments {
						for i := range seg.TrackPoints {
							if i > 0 {
								if seg.TrackPoints[i].CountMoving && seg.TrackPoints[i].DistanceToThisPoint <= seg.TrackPoints[i-1].DistanceToThisPoint {
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
								if seg.TrackPoints[i].CountMoving && seg.TrackPoints[i].DistanceToThisPoint <= seg.TrackPoints[i-1].DistanceToThisPoint {
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

func TestCompexTrackWithTimeStampInSomeSegments(t *testing.T) {
	gpxFile := NewGpxFile(filepath.Join(testhelper.GetProjectRoot(), "testdata", "valid-gpx", "04.gpx"))
	fmt.Println(gpxFile.FilePath)
	trackFile, err := gpxFile.ReadTracks("linear", 0.3, 10.0)
	if err != nil {
		t.Errorf("Got the following error while reading file %s: %s", gpxFile.FilePath, err.Error())
		return
	}

	if trackFile.GetTimeDataValid() == true {
		t.Errorf("Track file seems to have time data, but should not")
	}
}

func TestTrackReaderWithTimeStamps(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValidGPX("12.gpx"))

	file, _ := gpx.ReadTracks("none", 0.3, 10.0)

	if len(file.Tracks) != 1 {
		t.Errorf("Expected 1 Tracks, got %d", len(file.Tracks))
	}

	if file.GetStartTime() != file.Tracks[0].GetStartTime() {
		t.Errorf("The StartTime does not match for Track")
	}

	if file.GetEndTime() != file.Tracks[0].GetEndTime() {
		t.Errorf("The EndTime does not match for Track")
	}

	if file.GetStartTime() != file.Tracks[0].TrackSegments[0].GetStartTime() {
		t.Errorf("The StartTime does not match for TrackSegments")
	}

	if file.GetEndTime() != file.Tracks[0].TrackSegments[0].GetEndTime() {
		t.Errorf("The EndTime does not match for TrackSegments")
	}

	if file.GetStartTime() != file.Tracks[0].TrackSegments[0].TrackPoints[0].GetStartTime() {
		t.Errorf("The StartTime does not match for TrackPoints")
	}

	lastPoint := len(file.Tracks[0].TrackSegments[0].TrackPoints) - 1
	if file.GetEndTime() != file.Tracks[0].TrackSegments[0].TrackPoints[lastPoint].GetEndTime() {
		t.Errorf("The EndTime does not match for TrackPoints")
	}
}

func TestTrackReaderOnePointTrack(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValidGPX("06.gpx"))

	file, _ := gpx.ReadTracks("none", 0.3, 10.0)

	if file.Tracks[0].Distance != 0.0 {
		t.Errorf("The Distance is %f, but %f was expected", file.Tracks[0].Distance, 0.0)
	}

	if file.Tracks[0].UpwardsDistance != 0.0 {
		t.Errorf("The AltitudeRange is %f, but %f was expected", file.Tracks[0].UpwardsDistance, 0.0)
	}
}

func TestTrackReader02(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValidGPX("02.gpx"))

	file, _ := gpx.ReadTracks("none", 0.3, 10.0)

	if file.Tracks[0].Distance != 37823.344979382266 {
		t.Errorf("The Distance is %f, but %f was expected", file.Tracks[0].Distance, 37823.344979382266)
	}
}

func TestTrackReaderInValidCorrectionParameter(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValidGPX("02.gpx"))

	_, err := gpx.ReadTracks("asdfg", 0.3, 10.0)
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

func TestTrackReaderEmptyTrack(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetInvalidGPX("03.gpx"))

	_, err := gpx.ReadTracks("none", 0.3, 10.0)
	if err != nil {
		switch ty := err.(type) {
		case *EmptyGpxFileError:
			fmt.Println("OK")
		default:
			t.Errorf("The Error ReadTracks gave is of the wrong type. The type is %v", ty)
		}
	} else {
		t.Errorf("ReadTracks did not return a error, but was expected")
	}
}

func TestTrackReaderTrackWithEmptySegment(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValidGPX("13.gpx"))

	trk, err := gpx.ReadTracks("none", 0.3, 10.0)
	if err != nil {
		t.Errorf("Got a error, but expected none. The error is: %s", err)
	}

	if len(trk.Tracks) != 1 {
		t.Errorf("Got %d tracks, but expected 1", len(trk.Tracks))
	}

	if len(trk.Tracks[0].TrackSegments) != 1 {
		t.Errorf("Got %d track segments, but expected 1", len(trk.Tracks))
	}

}

func TestTrackReaderOneEmptyTrack(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValidGPX("14.gpx"))

	trk, err := gpx.ReadTracks("none", 0.3, 10.0)
	if err != nil {
		t.Errorf("Got a error, but expected none. The error is: %s", err)
	}

	if len(trk.Tracks) != 1 {
		t.Errorf("Got %d tracks, but expected 1", len(trk.Tracks))
	}

	if len(trk.Tracks[0].TrackSegments) != 1 {
		t.Errorf("Got %d track segments, but expected 1", len(trk.Tracks))
	}

}

func TestTrackReaderImpl(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValidGPX("01.gpx"))

	if gpx.FilePath != testhelper.GetValidGPX("01.gpx") {
		t.Errorf("GpxFile.FilePath was not %s but %s", testhelper.GetValidGPX("01.gpx"), gpx.FilePath)
	}

	file, err := gpx.ReadTracks("linear", 0.3, 10.0)

	if err != nil {
		t.Errorf("Got not expected error:  %s", err.Error())
	}

	if file.Tracks == nil {
		t.Errorf("Got nil tracks when reading a valid file")
	}

	if len(file.Tracks) != 1 {
		t.Errorf("The number of tracks was not %d, but was %d", 1, len(file.Tracks))
	}

	if file.Tracks[0].Distance != 18478.293509238614 {
		t.Errorf("The Distance is %f, but %f was expected", file.Tracks[0].Distance, 18478.293509238614)
	}

	if file.Tracks[0].MinimumAltitude != 298.0 {
		t.Errorf("The MinimumAltitude is %f, but %f was expected", file.Tracks[0].MinimumAltitude, 298.00)
	}

	if file.Tracks[0].MaximumAltitude != 402.0 {
		t.Errorf("The MaximumAltitude is %f, but %f was expected", file.Tracks[0].MaximumAltitude, 402.00)
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

	if file.FilePath != testhelper.GetValidGPX("01.gpx") {
		t.Errorf("The FilePath is %s, but %s was expected", file.FilePath, testhelper.GetValidGPX("01.gpx"))
	}

	if file.NumberOfTracks != 1 {
		t.Errorf("The NumberOfTracks is %d, but %d was expected", file.NumberOfTracks, 1)
	}

	if file.Distance != file.Tracks[0].Distance {
		t.Errorf("The Distance is %f, but %f was expected", file.Distance, file.Tracks[0].Distance)
	}

	if file.MinimumAltitude != file.Tracks[0].MinimumAltitude {
		t.Errorf("The MinimumAltitude is %f, but %f was expected", file.MinimumAltitude, file.Tracks[0].MinimumAltitude)
	}

	if file.MaximumAltitude != file.Tracks[0].MaximumAltitude {
		t.Errorf("The MaximumAltitude is %f, but %f was expected", file.MaximumAltitude, file.Tracks[0].MaximumAltitude)
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
	gpx := NewGpxFile(testhelper.GetValidGPX("01.gpx"))

	reader := gpsabl.TrackReader(&gpx)

	file, err := reader.ReadTracks("none", 0.3, 10.0)

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

	if info.GetAltitudeRange() != file.Tracks[0].GetAltitudeRange() {
		t.Errorf("The AltitudeRange is %f, but %f was expected", info.GetAltitudeRange(), file.Tracks[0].GetAltitudeRange())
	}

	if info.GetMinimumAltitude() != file.Tracks[0].MinimumAltitude {
		t.Errorf("The MinimumAltitude is %f, but %f was expected", info.GetMinimumAltitude(), file.Tracks[0].MinimumAltitude)
	}

	if info.GetMaximumAltitude() != file.Tracks[0].MaximumAltitude {
		t.Errorf("The MaximumAltitude is %f, but %f was expected", info.GetMaximumAltitude(), file.Tracks[0].MaximumAltitude)
	}
}

func TestGpxFileInterfaceImplentaion2(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValidGPX("01.gpx"))

	reader := gpsabl.TrackReader(&gpx)

	file, err := reader.ReadTracks("none", 0.3, 10.0)

	if err != nil {
		t.Errorf("Got not expected error:  %s", err.Error())
	}

	info := gpsabl.TrackSummaryProvider(&gpx)

	if info.GetDistance() != file.Tracks[0].Distance {
		t.Errorf("The Distance is %f, but %f was expected", info.GetDistance(), file.Tracks[0].Distance)
	}

	if info.GetAltitudeRange() != file.Tracks[0].GetAltitudeRange() {
		t.Errorf("The AltitudeRange is %f, but %f was expected", info.GetAltitudeRange(), file.Tracks[0].GetAltitudeRange())
	}

	if info.GetMinimumAltitude() != file.Tracks[0].MinimumAltitude {
		t.Errorf("The MinimumAltitude is %f, but %f was expected", info.GetMinimumAltitude(), file.Tracks[0].MinimumAltitude)
	}

	if info.GetMaximumAltitude() != file.Tracks[0].MaximumAltitude {
		t.Errorf("The MaximumAltitude is %f, but %f was expected", info.GetMaximumAltitude(), file.Tracks[0].MaximumAltitude)
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

	if gpx.MinimumAltitude != file.Tracks[0].MinimumAltitude {
		t.Errorf("The MinimumAltitude is %f, but %f was expected", gpx.MinimumAltitude, file.Tracks[0].MinimumAltitude)
	}

	if gpx.MaximumAltitude != file.Tracks[0].MaximumAltitude {
		t.Errorf("The MaximumAltitude is %f, but %f was expected", gpx.MaximumAltitude, file.Tracks[0].MaximumAltitude)
	}
}

func TestReadGpxFile(t *testing.T) {
	file, err := ReadGpxFile(testhelper.GetValidGPX("01.gpx"), "none", 0.3, 10.0)
	if err != nil {
		t.Errorf("Something wrong when reading a valid gpx file: %s", err.Error())
	}

	if len(file.Tracks) != 1 {
		t.Errorf("The number of tracks was not %d, but was %d", 1, len(file.Tracks))
	}

	if file.Name != "GPX name" {
		t.Errorf("The GPX name was not expected. Got: %s", file.Name)
	}

	if file.Description != "A valid GPX Track" {
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

	if file.FilePath != testhelper.GetValidGPX("01.gpx") {
		t.Errorf("The FilePath is %s, but %s was expected", file.FilePath, testhelper.GetValidGPX("01.gpx"))
	}

	if file.NumberOfTracks != 1 {
		t.Errorf("The NumberOfTracks is %d, but %d was expected", file.NumberOfTracks, 1)
	}

	if file.Distance != file.Tracks[0].Distance {
		t.Errorf("The Distance is %f, but %f was expected", file.Distance, file.Tracks[0].Distance)
	}

	if file.MinimumAltitude != file.Tracks[0].MinimumAltitude {
		t.Errorf("The MinimumAltitude is %f, but %f was expected", file.MinimumAltitude, file.Tracks[0].MinimumAltitude)
	}

	if file.MaximumAltitude != file.Tracks[0].MaximumAltitude {
		t.Errorf("The MaximumAltitude is %f, but %f was expected", file.MaximumAltitude, file.Tracks[0].MaximumAltitude)
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
	gpx := NewGpxFile(testhelper.GetValidGPX("NotExisting.gpx"))
	_, err := gpx.ReadTracks("none", 0.3, 10.0)
	switch v := err.(type) {
	case nil:
		t.Errorf("No error, when reading a not existing gpx file")
	case *os.PathError:
		fmt.Println("OK")
	default:
		t.Errorf("Expected a *os.PathError, got a %s", reflect.TypeOf(v))
	}
}

func TestReadTracksInvalidGPX(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetInvalidGPX("01.gpx"))
	_, err := gpx.ReadTracks("none", 0.3, 10.0)
	switch v := err.(type) {
	case nil:
		t.Errorf("No error, when reading an invalid gpx file")
	case *xml.SyntaxError:
		fmt.Println("OK")
	default:
		t.Errorf("Expected a *xml.SyntaxError, got a %s", reflect.TypeOf(v))
	}

}

func TestReadTracksNotGPX(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetInvalidGPX("02.gpx"))
	file, err := gpx.ReadTracks("none", 0.3, 10.0)
	switch v := err.(type) {
	case nil:
		t.Errorf("No error, when reading an invalid gpx file")
	case *GpxFileError:
		checkGpxFileError(v, testhelper.GetInvalidGPX("02.gpx"), t)
	default:
		t.Errorf("Expected a *gpsabl.GpxFileError, got a %s", reflect.TypeOf(v))
	}

	fmt.Println(file.Name)
}

func TestTrackReaderAlpineSkiTrack(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValidGPX("16.gpx"))
	file, err := gpx.ReadTracks("linear", 0.3, 10.0)

	if err != nil {
		t.Errorf("Something wrong when reading a valid gpx file: %s", err.Error())
	}

	if file.GetMovingTime() == file.GetEndTime().Sub(file.GetStartTime()) {
		t.Errorf("The GetMovingTime is the same as the speed calculated from start and end time")
	}

	if file.GetMovingTime() != file.Tracks[0].GetMovingTime() {
		t.Errorf("The tracks moving time is not the same as the files moving time")
	}

	if file.GetMovingTime() != file.Tracks[0].TrackSegments[0].GetMovingTime() {
		t.Errorf("The track segments moving time is not the same as the files moving time")
	}

	if file.GetUpwardsTime() >= file.GetMovingTime() {
		t.Errorf("The UpwardsTime %d is bigger than the moving time %d", file.GetUpwardsTime(), file.GetMovingTime())
	}

	if file.GetDownwardsTime() >= file.GetMovingTime() {
		t.Errorf("The DownwardsTime %d is bigger than the moving time %d", file.GetDownwardsTime(), file.GetMovingTime())
	}

	if file.GetUpwardsSpeed() >= file.GetDownwardsSpeed() {
		t.Errorf("The UpwardsSpeed %f is bigger than the DownwardsSpeed %f", file.GetUpwardsSpeed(), file.GetDownwardsSpeed())
	}

	if file.GetDownwardsTime()+file.GetUpwardsTime() > file.GetMovingTime() {
		t.Errorf("The MovingTime %d is smaller than DownwardsTime + UpwardsTime %d", file.GetMovingTime(), file.GetDownwardsTime()+file.GetUpwardsTime())
	}

	if file.GetAvarageSpeed() >= file.GetDownwardsSpeed() {
		t.Errorf("The GetAvarageSpeed %f is bigger than the DownwardsSpeed %f", file.GetAvarageSpeed(), file.GetDownwardsSpeed())
	}

	if file.GetDownwardsTime() != file.Tracks[0].GetDownwardsTime() {
		t.Errorf("The file DownwardsTime %d is not the same the the tracks DownwardsTime %d", file.GetDownwardsTime(), file.Tracks[0].GetDownwardsTime())
	}

	if file.GetDownwardsTime() != file.Tracks[0].TrackSegments[0].GetDownwardsTime() {
		t.Errorf("The file DownwardsTime %d is not the same the the segments DownwardsTime %d", file.GetDownwardsTime(), file.Tracks[0].GetDownwardsTime())
	}

	if file.GetUpwardsTime() != file.Tracks[0].GetUpwardsTime() {
		t.Errorf("The file GetUpwardsTime %d is not the same the the tracks GetUpwardsTime %d", file.GetUpwardsTime(), file.Tracks[0].GetUpwardsTime())
	}

	if file.GetUpwardsTime() != file.Tracks[0].TrackSegments[0].GetUpwardsTime() {
		t.Errorf("The file GetUpwardsTime %d is not the same the the segments GetUpwardsTime %d", file.GetUpwardsTime(), file.Tracks[0].GetUpwardsTime())
	}

	if file.GetUpwardsSpeed() != file.Tracks[0].GetUpwardsSpeed() {
		t.Errorf("The file GetUpwardsSpeed %f is not the same the the tracks GetUpwardsSpeed %f", file.GetUpwardsSpeed(), file.Tracks[0].GetUpwardsSpeed())
	}

	if file.GetUpwardsSpeed() != file.Tracks[0].TrackSegments[0].GetUpwardsSpeed() {
		t.Errorf("The file GetUpwardsSpeed %f is not the same the the segments GetUpwardsSpeed %f", file.GetUpwardsSpeed(), file.Tracks[0].GetUpwardsSpeed())
	}

	if file.GetDownwardsSpeed() != file.Tracks[0].GetDownwardsSpeed() {
		t.Errorf("The file GetDownwardsSpeed %f is not the same the the tracks GetDownwardsSpeed %f", file.GetDownwardsSpeed(), file.Tracks[0].GetDownwardsSpeed())
	}

	if file.GetDownwardsSpeed() != file.Tracks[0].TrackSegments[0].GetDownwardsSpeed() {
		t.Errorf("The file GetDownwardsSpeed %f is not the same the the segments GetDownwardsSpeed %f", file.GetDownwardsSpeed(), file.Tracks[0].GetDownwardsSpeed())
	}
}

func TestMinimalStepHightValues(t *testing.T) {
	gpx := NewGpxFile(testhelper.GetValidGPX("02.gpx"))
	file1, err1 := gpx.ReadTracks("none", 0.3, 10.0)
	file2, err2 := gpx.ReadTracks("steps", 0.3, 0.0)
	file3, err3 := gpx.ReadTracks("steps", 0.3, 20.0)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("Got a error but none expected")
	}

	if file1.GetElevationGain() != file2.GetElevationGain() {
		t.Errorf("The file1.GetElevationGain() \"%f\" is not the same as file2.GetElevationGain() \"%f\"", file1.GetElevationGain(), file2.GetElevationGain())
	}

	if file1.GetElevationLose() != file2.GetElevationLose() {
		t.Errorf("The file1.GetElevationLose() \"%f\" is not the same as file2.GetElevationLose() \"%f\"", file1.GetElevationLose(), file2.GetElevationLose())
	}

	if file2.GetElevationGain() <= file3.GetElevationGain() {
		t.Errorf("The file2.GetElevationGain() \"%f\" is the same as file3.GetElevationGain() \"%f\"", file2.GetElevationGain(), file3.GetElevationGain())
	}

	if file2.GetElevationLose() >= file3.GetElevationLose() {
		t.Errorf("The file2.GetElevationLose() \"%f\" is the same as file3.GetElevationLose() \"%f\"", file2.GetElevationLose(), file3.GetElevationLose())
	}
}

func TestMinimalMovingSpeedValues(t *testing.T) {

	gpx := NewGpxFile(testhelper.GetValidGPX("16.gpx"))
	file1, err1 := gpx.ReadTracks("none", 0.3, 10.0)
	file2, err2 := gpx.ReadTracks("none", 0.0, 10.0)

	if err1 != nil || err2 != nil {
		t.Errorf("Got a error but none expected")
	}
	if file1.GetMovingTime() == file2.GetMovingTime() {
		t.Errorf("The MovingTime is the same, no matter whats the MinimalMovingSpeed")
	}

	if file2.GetMovingTime() != file2.GetEndTime().Sub(file2.GetStartTime()) {
		t.Errorf("The MovingTime %d is not the same as EndTime - StartTime %d when called with zero MinimalMovingSpeed", file2.GetMovingTime(), file2.GetEndTime().Sub(file2.GetStartTime()))
	}

}

func TestParameterErrors(t *testing.T) {

	gpx := NewGpxFile(testhelper.GetValidGPX("01.gpx"))
	_, err1 := gpx.ReadTracks("none", -0.3, 10.0)

	if err1 != nil {
		switch err1.(type) {
		case *gpsabl.MinimalMovingSpeedLessThenZero:
			fmt.Println("OK")
		default:
			t.Errorf("Expected a MinimalMovingSpeedLessThenZero, got a %s", reflect.TypeOf(err1))
		}

	} else {
		t.Errorf("Got no error when a MinimalMovingSpeedLessThenZero error is expected")
	}

	_, err2 := gpx.ReadTracks("none", 0.3, -10.0)
	if err2 != nil {
		switch err2.(type) {
		case *gpsabl.MinimalStepHightLessThenZero:
			fmt.Println("OK")
		default:
			t.Errorf("Expected a MinimalStepHightLessThenZero, got a %s", reflect.TypeOf(err2))
		}

	} else {
		t.Errorf("Got no error when a MinimalStepHightLessThenZero error is expected")
	}
}
