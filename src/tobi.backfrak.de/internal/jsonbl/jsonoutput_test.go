package jsonbl

import (
	"fmt"
	"os"
	"testing"
	"time"

	"tobi.backfrak.de/internal/gpsabl"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file

func TestNewJSONOutputFormater(t *testing.T) {

	sut := NewJSONOutputFormater()

	if len(sut.lineBuffer) != 0 {
		t.Errorf("The new JSONOutputFormater does not have an empty buffer")
	}
}

func TestNewOutputFormater(t *testing.T) {
	orig := NewJSONOutputFormater()
	sut := orig.NewOutputFormater()

	if sut.CheckFileExtension("my/output.json") == false {
		t.Errorf("JSONOutputFormater can not write *.json")
	}

	if sut.CheckFileExtension("my/output.csv") == true {
		t.Errorf("JSONOutputFormater can write *.csv")
	}

	if sut.CheckOutputFormaterType(JSONOutputFormatertype) == false {
		t.Errorf("JSONOutputFormater can not write %s type", JSONOutputFormatertype)
	}

	if sut.CheckOutputFormaterType(gpsabl.OutputFormaterType("abs")) == true {
		t.Errorf("JSONOutputFormater can write %s type", "abs")
	}

	ext := sut.GetFileExtensions()
	if len(ext) != 1 {
		t.Errorf("The number of FileExtensions is not expected")
	}

	if ext[0] != ".json" {
		t.Errorf("The file type \"%s\" is not the expexted \"%s\"", ext[0], ".json")
	}

	if sut.GetNumberOfOutputEntries() != -1 {
		t.Errorf("The initial value of GetNumberOfOutputEntries is %d but should be %d", sut.GetNumberOfOutputEntries(), -1)
	}

	form := sut.GetOutputFormaterTypes()

	if len(form) != 1 {
		t.Errorf("The number of FileExtensionsTypes is not expected")
	}

	if form[0] != gpsabl.OutputFormaterType("JSON") {
		t.Errorf("The file type \"%s\" is not the expexted \"%s\"", form[0], "JSON")
	}

	textOut := sut.GetTextOutputFormater()

	if textOut != nil {
		t.Errorf("JSONOutputFormater should not be a TextOutputFormater")
	}

}

func TestAddOutPutInvalidDepth(t *testing.T) {
	sut := NewJSONOutputFormater()
	trk := getSimpleTrackFileWithTime()
	err := sut.AddOutPut(trk, gpsabl.DepthArg("bla"), false)
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
	switch err.(type) {
	case *gpsabl.DepthParameterNotKnownError:
		fmt.Println("OK")
	default:
		t.Errorf("The error is not from the expected type")
	}
}

func TestGetOutputInvalidSummary(t *testing.T) {
	sut := NewJSONOutputFormater()
	trk1 := getSimpleTrackFileWithTime()
	err1 := sut.AddOutPut(trk1, gpsabl.TRACK, false)
	if err1 != nil {
		t.Errorf("Got an error but expected none")
	}

	_, err3 := sut.GetOutput(gpsabl.SummaryArg("blka"))
	if err3 == nil {
		t.Errorf("Got no error but expected one")
	}
	switch err3.(type) {
	case *gpsabl.SummaryParamaterNotKnown:
		fmt.Println("OK")
	default:
		t.Errorf("The error is not from the expected type")
	}
}

func TestWriteOutputInvalidSummary(t *testing.T) {
	sut := NewJSONOutputFormater()
	trk1 := getSimpleTrackFileWithTime()
	err1 := sut.AddOutPut(trk1, gpsabl.TRACK, false)
	if err1 != nil {
		t.Errorf("Got an error but expected none")
	}

	err3 := sut.WriteOutput(os.Stdout, gpsabl.SummaryArg("blka"))
	if err3 == nil {
		t.Errorf("Got no error but expected one")
	}
	switch err3.(type) {
	case *gpsabl.SummaryParamaterNotKnown:
		fmt.Println("OK")
	default:
		t.Errorf("The error is not from the expected type")
	}
}

func TestAddOutPutSegmentDepth(t *testing.T) {
	sut := NewJSONOutputFormater()
	trk1 := getSimpleTrackFileWithTime()
	err1 := sut.AddOutPut(trk1, gpsabl.SEGMENT, true)
	if err1 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 1 {
		t.Errorf("Expected 1 lines, but got %d", len(sut.lineBuffer))
	}

	trk2 := getTrackFileWithDifferentTime()
	err2 := sut.AddOutPut(trk2, gpsabl.SEGMENT, true)
	if err2 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 2 {
		t.Errorf("Expected 2 lines, but got %d", len(sut.lineBuffer))
	}
}

func TestAddOutPutFileDepth(t *testing.T) {
	sut := NewJSONOutputFormater()
	trk1 := getSimpleTrackFileWithTime()
	err1 := sut.AddOutPut(trk1, gpsabl.FILE, true)
	if err1 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 1 {
		t.Errorf("Expected 1 lines, but got %d", len(sut.lineBuffer))
	}

	trk2 := getTrackFileWithDifferentTime()
	err2 := sut.AddOutPut(trk2, gpsabl.FILE, true)
	if err2 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 2 {
		t.Errorf("Expected 2 lines, but got %d", len(sut.lineBuffer))
	}
}

func TestAddOutPutFilterDuplicate1(t *testing.T) {
	sut := NewJSONOutputFormater()
	trk1 := getSimpleTrackFileWithTime()
	err1 := sut.AddOutPut(trk1, gpsabl.TRACK, true)
	if err1 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 1 {
		t.Errorf("Expected 1 lines, but got %d", len(sut.lineBuffer))
	}

	trk2 := getSimpleTrackFileWithTime()
	err2 := sut.AddOutPut(trk2, gpsabl.TRACK, true)
	if err2 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 1 {
		t.Errorf("Expected 1 lines, but got %d", len(sut.lineBuffer))
	}
}

func TestAddOutPutFilterDuplicate2(t *testing.T) {
	sut := NewJSONOutputFormater()
	trk1 := getSimpleTrackFileWithTime()
	err1 := sut.AddOutPut(trk1, gpsabl.TRACK, true)
	if err1 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 1 {
		t.Errorf("Expected 1 lines, but got %d", len(sut.lineBuffer))
	}

	trk2 := getTrackFileWithDifferentTime()
	err2 := sut.AddOutPut(trk2, gpsabl.TRACK, true)
	if err2 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 2 {
		t.Errorf("Expected 1 lines, but got %d", len(sut.lineBuffer))
	}
}

func TestAddOutPutNotFilterDuplicate(t *testing.T) {
	sut := NewJSONOutputFormater()
	trk1 := getSimpleTrackFileWithTime()
	err1 := sut.AddOutPut(trk1, gpsabl.TRACK, false)
	if err1 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 1 {
		t.Errorf("Expected 1 lines, but got %d", len(sut.lineBuffer))
	}

	trk2 := getSimpleTrackFileWithTime()
	err2 := sut.AddOutPut(trk2, gpsabl.TRACK, false)
	if err2 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 2 {
		t.Errorf("Expected 2 lines, but got %d", len(sut.lineBuffer))
	}
}

func TestGetOutputNoSummary(t *testing.T) {
	sut := NewJSONOutputFormater()
	trk1 := getSimpleTrackFileWithTime()
	err1 := sut.AddOutPut(trk1, gpsabl.TRACK, false)
	if err1 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 1 {
		t.Errorf("Expected 1 lines, but got %d", len(sut.lineBuffer))
	}

	trk2 := getTrackFileWithDifferentTime()
	err2 := sut.AddOutPut(trk2, gpsabl.TRACK, false)
	if err2 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 2 {
		t.Errorf("Expected 2 lines, but got %d", len(sut.lineBuffer))
	}

	res, err3 := sut.GetOutput(gpsabl.NONE)
	if err3 != nil {
		t.Errorf("Got an error but expected none")
	}
	if len(res.Statistics) != len(sut.lineBuffer) {
		t.Errorf("GetOutput did not return the expected result")
	}
	if len(res.Summary) != 0 {
		t.Errorf("GetOutput did not return the expected result")
	}
}

func TestGetOutputOnlySummary(t *testing.T) {
	sut := NewJSONOutputFormater()
	trk1 := getSimpleTrackFileWithTime()
	err1 := sut.AddOutPut(trk1, gpsabl.TRACK, false)
	if err1 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 1 {
		t.Errorf("Expected 1 lines, but got %d", len(sut.lineBuffer))
	}

	trk2 := getTrackFileWithDifferentTime()
	err2 := sut.AddOutPut(trk2, gpsabl.TRACK, false)
	if err2 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 2 {
		t.Errorf("Expected 2 lines, but got %d", len(sut.lineBuffer))
	}

	res, err3 := sut.GetOutput(gpsabl.ONLY)
	if err3 != nil {
		t.Errorf("Got an error but expected none")
	}
	if len(res.Statistics) != 0 {
		t.Errorf("GetOutput did not return the expected result")
	}
	if len(res.Summary) != 4 {
		t.Errorf("GetOutput did not return the expected result")
	}
}

func TestGetOutputAdditionalSummary(t *testing.T) {
	sut := NewJSONOutputFormater()
	trk1 := getSimpleTrackFileWithTime()
	err1 := sut.AddOutPut(trk1, gpsabl.TRACK, false)
	if err1 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 1 {
		t.Errorf("Expected 1 lines, but got %d", len(sut.lineBuffer))
	}

	trk2 := getTrackFileWithDifferentTime()
	err2 := sut.AddOutPut(trk2, gpsabl.TRACK, false)
	if err2 != nil {
		t.Errorf("Got an error but expected none")
	}

	if len(sut.lineBuffer) != 2 {
		t.Errorf("Expected 2 lines, but got %d", len(sut.lineBuffer))
	}

	res, err3 := sut.GetOutput(gpsabl.ADDITIONAL)
	if err3 != nil {
		t.Errorf("Got an error but expected none")
	}
	if len(res.Statistics) != len(sut.lineBuffer) {
		t.Errorf("GetOutput did not return the expected result")
	}
	if len(res.Summary) != 4 {
		t.Errorf("GetOutput did not return the expected result")
	}
}

func TestWriteOutput1(t *testing.T) {
	sut := NewJSONOutputFormater()
	trk1 := getTrackFileTwoTracksWithThreeSegmentsWithTime()
	err1 := sut.AddOutPut(trk1, gpsabl.SEGMENT, false)
	if err1 != nil {
		t.Errorf("Got an error but expected none")
	}

	trk2 := getTrackFileWithDifferentTime()
	err2 := sut.AddOutPut(trk2, gpsabl.SEGMENT, false)
	if err2 != nil {
		t.Errorf("Got an error but expected none")
	}

	err3 := sut.WriteOutput(os.Stdout, gpsabl.ADDITIONAL)
	if err3 != nil {
		t.Errorf("Got an error but expected none")
	}

	if sut.GetNumberOfOutputEntries() != 8 {
		t.Errorf("Error: The number of output entries is %d but should be %d", sut.GetNumberOfOutputEntries(), 8)
	}
}

func TestWriteOutput2(t *testing.T) {
	sut := NewJSONOutputFormater()
	trk1 := getTrackFileTwoTracksWithThreeSegmentsWithTime()
	err1 := sut.AddOutPut(trk1, gpsabl.TRACK, true)
	if err1 != nil {
		t.Errorf("Got an error but expected none")
	}

	trk2 := getTrackFileWithDifferentTime()
	err2 := sut.AddOutPut(trk2, gpsabl.TRACK, true)
	if err2 != nil {
		t.Errorf("Got an error but expected none")
	}

	err3 := sut.WriteOutput(os.Stdout, gpsabl.ADDITIONAL)
	if err3 != nil {
		t.Errorf("Got an error but expected none")
	}
}

func TestGetSummaryEntiresWithEmptyTrackList(t *testing.T) {
	sut := NewJSONOutputFormater()

	lines := sut.getSummaryEntires()

	if len(lines) != 0 {
		t.Errorf("Don't get an empty list when no entries are added")
	}
}

func getTrackFileWithDifferentTime() gpsabl.TrackFile {
	ret := gpsabl.NewTrackFile("/mys/track/file")
	trk := getTrackWithDifferentTime()
	gpsabl.FillTrackValues(&trk)
	ret.Tracks = []gpsabl.Track{trk}
	gpsabl.FillTrackFileValues(&ret)
	ret.NumberOfTracks = 1

	return ret
}

func getTrackPoint(lat, lon, ele float32) gpsabl.TrackPoint {
	pnt := gpsabl.TrackPoint{}
	pnt.Latitude = lat
	pnt.Longitude = lon
	pnt.Elevation = ele
	pnt.TimeValid = false

	return pnt
}

func getTrackPointWithTime(lat, lon, ele float32, time time.Time) gpsabl.TrackPoint {
	pnt := getTrackPoint(lat, lon, ele)
	pnt.TimeValid = true
	pnt.Time = time

	return pnt
}

func getTrackWithDifferentTime() gpsabl.Track {
	t1, _ := time.Parse(time.RFC3339, "2015-08-22T17:19:33Z")
	t2, _ := time.Parse(time.RFC3339, "2015-08-22T17:19:43Z")
	t3, _ := time.Parse(time.RFC3339, "2015-08-22T17:19:53Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 108.0, t2)
	pnt3 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t3)
	points := []gpsabl.TrackPoint{pnt1, pnt2, pnt3}

	gpsabl.FillDistancesTrackPoint(&points[0], gpsabl.TrackPoint{}, points[1])
	gpsabl.FillDistancesTrackPoint(&points[1], points[0], points[2])
	gpsabl.FillDistancesTrackPoint(&points[2], points[1], gpsabl.TrackPoint{})
	gpsabl.FillValuesTrackPointArray(points, "none", 0.3, 10.0)
	seg := gpsabl.TrackSegment{}
	seg.TrackPoints = points
	ret := gpsabl.Track{}
	gpsabl.FillTrackSegmentValues(&seg)
	ret.TrackSegments = []gpsabl.TrackSegment{seg}
	gpsabl.FillTrackValues(&ret)
	ret.NumberOfSegments = 1

	return ret
}

func getTrackFileTwoTracksWithThreeSegments() gpsabl.TrackFile {
	trackFile := getTrackFileTwoTracks()
	trackFile.Tracks[0].TrackSegments = append(trackFile.Tracks[0].TrackSegments, getSimpleTrackFile().Tracks[0].TrackSegments[0])
	gpsabl.FillTrackValues(&trackFile.Tracks[0])

	return trackFile
}

func getTrackFileTwoTracksWithThreeSegmentsWithTime() gpsabl.TrackFile {
	trackFile := getTrackFileTwoTracksWithTime()
	trackFile.Tracks[0].TrackSegments = append(trackFile.Tracks[0].TrackSegments, getSimpleTrackFileWithTime().Tracks[0].TrackSegments[0])
	gpsabl.FillTrackValues(&trackFile.Tracks[0])

	return trackFile
}

func getTrackFileTwoTracks() gpsabl.TrackFile {
	trackFile := getSimpleTrackFile()
	trackFile.Tracks = append(trackFile.Tracks, getSimpleTrackFile().Tracks...)
	gpsabl.FillTrackFileValues(&trackFile)

	return trackFile
}

func getTrackFileTwoTracksWithTime() gpsabl.TrackFile {
	trackFile := getSimpleTrackFileWithTime()
	trackFile.Tracks = append(trackFile.Tracks, getSimpleTrackFileWithTime().Tracks...)
	gpsabl.FillTrackFileValues(&trackFile)

	return trackFile
}

func getTrackFileOneTrackWithTimeOneWithout() gpsabl.TrackFile {
	trackFile := getSimpleTrackFileWithTime()
	trackFile.Tracks = append(trackFile.Tracks, getSimpleTrackFile().Tracks...)
	gpsabl.FillTrackFileValues(&trackFile)

	return trackFile
}

func getSimpleTrackFile() gpsabl.TrackFile {
	ret := gpsabl.NewTrackFile("/mys/track/file")
	trk := getSimpleTrack()
	gpsabl.FillTrackValues(&trk)
	ret.Tracks = []gpsabl.Track{trk}
	gpsabl.FillTrackFileValues(&ret)
	ret.NumberOfTracks = 1

	return ret
}

func getSimpleTrackFileWithTime() gpsabl.TrackFile {
	ret := gpsabl.NewTrackFile("/mys/track/file")
	trk := getSimpleTrackWithTime()
	gpsabl.FillTrackValues(&trk)
	ret.Tracks = []gpsabl.Track{trk}
	gpsabl.FillTrackFileValues(&ret)
	ret.NumberOfTracks = 1

	return ret
}

func getSimpleTrack() gpsabl.Track {
	ret := gpsabl.Track{}
	segs := getSimpleTrackSegment()
	gpsabl.FillTrackSegmentValues(&segs)
	ret.TrackSegments = []gpsabl.TrackSegment{segs}
	gpsabl.FillTrackValues(&ret)
	ret.NumberOfSegments = 1

	return ret
}

func getSimpleTrackWithTime() gpsabl.Track {
	ret := gpsabl.Track{}
	segs := getSimpleTrackSegmentWithTime()
	gpsabl.FillTrackSegmentValues(&segs)
	ret.TrackSegments = []gpsabl.TrackSegment{segs}
	gpsabl.FillTrackValues(&ret)
	ret.NumberOfSegments = 1

	return ret
}

func getSimpleTrackSegment() gpsabl.TrackSegment {
	seg := gpsabl.TrackSegment{}
	points := gerSimpleTrackPointArray()
	seg.TrackPoints = points

	return seg
}

func getSimpleTrackSegmentWithTime() gpsabl.TrackSegment {
	seg := gpsabl.TrackSegment{}
	points := getSimpleTrackPointArrayWithTime()
	seg.TrackPoints = points

	return seg
}

func gerSimpleTrackPointArray() []gpsabl.TrackPoint {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	points := []gpsabl.TrackPoint{pnt1, pnt2, pnt3}

	gpsabl.FillDistancesTrackPoint(&points[0], gpsabl.TrackPoint{}, points[1])
	gpsabl.FillDistancesTrackPoint(&points[1], points[0], points[2])
	gpsabl.FillDistancesTrackPoint(&points[2], points[1], gpsabl.TrackPoint{})
	gpsabl.FillValuesTrackPointArray(points, "none", 0.3, 10.0)

	return points
}

func getSimpleTrackPointArrayWithTime() []gpsabl.TrackPoint {
	t1, _ := time.Parse(time.RFC3339, "2014-08-22T17:19:33Z")
	t2, _ := time.Parse(time.RFC3339, "2014-08-22T17:19:43Z")
	t3, _ := time.Parse(time.RFC3339, "2014-08-22T17:19:53Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 108.0, t2)
	pnt3 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t3)
	points := []gpsabl.TrackPoint{pnt1, pnt2, pnt3}

	gpsabl.FillDistancesTrackPoint(&points[0], gpsabl.TrackPoint{}, points[1])
	gpsabl.FillDistancesTrackPoint(&points[1], points[0], points[2])
	gpsabl.FillDistancesTrackPoint(&points[2], points[1], gpsabl.TrackPoint{})
	gpsabl.FillValuesTrackPointArray(points, "none", 0.3, 10.0)

	return points
}

func getTrackFileWithStandStillPoints(correction string, minimalMovingSpeed float64, minimalStepHight float64) gpsabl.TrackFile {
	var file gpsabl.TrackFile

	t1, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:13Z")
	t2, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:33Z")
	t3, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:53Z")
	t4, _ := time.Parse(time.RFC3339, "2014-08-22T19:20:13Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 108.0, t2)
	pnt3 := getTrackPointWithTime(50.11495751, 8.684874771, 108.0, t3)
	pnt4 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t4)
	points := []gpsabl.TrackPoint{pnt1, pnt2, pnt3, pnt4}

	gpsabl.FillDistancesTrackPoint(&points[0], gpsabl.TrackPoint{}, points[1])
	gpsabl.FillDistancesTrackPoint(&points[1], points[0], points[2])
	gpsabl.FillDistancesTrackPoint(&points[2], points[1], points[3])
	gpsabl.FillDistancesTrackPoint(&points[3], points[2], gpsabl.TrackPoint{})
	gpsabl.FillValuesTrackPointArray(points, gpsabl.CorrectionParameter(correction), minimalMovingSpeed, minimalStepHight)
	laterTrack := gpsabl.Track{}
	seg := gpsabl.TrackSegment{}
	seg.TrackPoints = points
	gpsabl.FillTrackSegmentValues(&seg)
	laterTrack.TrackSegments = append(laterTrack.TrackSegments, seg)
	gpsabl.FillTrackValues(&laterTrack)
	laterTrack.NumberOfSegments = 1

	file.Tracks = append(file.Tracks, laterTrack)

	file.NumberOfTracks = 1
	gpsabl.FillTrackFileValues(&file)

	return file
}

func getTrackFileWithBigVerticalDistance() gpsabl.TrackFile {
	file := getSimpleTrackFileWithTime()

	t1, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:13Z")
	t2, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:33Z")
	t3, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:53Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 142.0, t2)
	pnt3 := getTrackPointWithTime(50.11484790, 8.684885500, 151.0, t3)
	points := []gpsabl.TrackPoint{pnt1, pnt2, pnt3}

	gpsabl.FillDistancesTrackPoint(&points[0], gpsabl.TrackPoint{}, points[1])
	gpsabl.FillDistancesTrackPoint(&points[1], points[0], points[2])
	gpsabl.FillDistancesTrackPoint(&points[2], points[1], gpsabl.TrackPoint{})
	gpsabl.FillValuesTrackPointArray(points, "none", 0.3, 10.0)
	laterTrack := gpsabl.Track{}
	seg := gpsabl.TrackSegment{}
	seg.TrackPoints = points
	gpsabl.FillTrackSegmentValues(&seg)
	laterTrack.TrackSegments = append(laterTrack.TrackSegments, seg)
	gpsabl.FillTrackValues(&laterTrack)
	laterTrack.NumberOfSegments = 1

	file.Tracks = append(file.Tracks, laterTrack)
	file.NumberOfTracks = 2
	gpsabl.FillTrackFileValues(&file)

	return file
}

func getTrackFileWithTimeGaps() gpsabl.TrackFile {
	file := getSimpleTrackFileWithTime()

	t1, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:13Z")
	t2, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:33Z")
	t3, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:53Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 108.0, t2)
	pnt3 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t3)
	points := []gpsabl.TrackPoint{pnt1, pnt2, pnt3}

	gpsabl.FillDistancesTrackPoint(&points[0], gpsabl.TrackPoint{}, points[1])
	gpsabl.FillDistancesTrackPoint(&points[1], points[0], points[2])
	gpsabl.FillDistancesTrackPoint(&points[2], points[1], gpsabl.TrackPoint{})
	gpsabl.FillValuesTrackPointArray(points, "none", 0.3, 10.0)
	laterTrack := gpsabl.Track{}
	seg := gpsabl.TrackSegment{}
	seg.TrackPoints = points
	gpsabl.FillTrackSegmentValues(&seg)
	laterTrack.TrackSegments = append(laterTrack.TrackSegments, seg)
	gpsabl.FillTrackValues(&laterTrack)
	laterTrack.NumberOfSegments = 1

	file.Tracks = append(file.Tracks, laterTrack)
	file.NumberOfTracks = 2
	gpsabl.FillTrackFileValues(&file)

	return file
}
