package gpsabl

// Copyright 2025 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import "time"

const DEFAULT_START_TIME = "2014-08-22T17:19:33Z"

func getTrackPoint(lat, lon, ele float32) TrackPoint {
	pnt := TrackPoint{}
	pnt.Latitude = lat
	pnt.Longitude = lon
	pnt.Elevation = ele
	pnt.TimeValid = false

	return pnt
}

func getTrackPointWithTime(lat, lon, ele float32, time time.Time) TrackPoint {
	pnt := getTrackPoint(lat, lon, ele)
	pnt.TimeValid = true
	pnt.Time = time

	return pnt
}

func getTrackFileWithStandStillPoints(correction string, minimalMovingSpeed float64, minimalStepHight float64) TrackFile {
	var file TrackFile

	t1, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:13Z")
	t2, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:33Z")
	t3, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:53Z")
	t4, _ := time.Parse(time.RFC3339, "2014-08-22T19:20:13Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 108.0, t2)
	pnt3 := getTrackPointWithTime(50.11495751, 8.684874771, 108.0, t3)
	pnt4 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t4)
	points := []TrackPoint{pnt1, pnt2, pnt3, pnt4}

	FillDistancesTrackPoint(&points[0], TrackPoint{}, points[1])
	FillDistancesTrackPoint(&points[1], points[0], points[2])
	FillDistancesTrackPoint(&points[2], points[1], points[3])
	FillDistancesTrackPoint(&points[3], points[2], TrackPoint{})
	FillValuesTrackPointArray(points, CorrectionParameter(correction), minimalMovingSpeed, minimalStepHight)
	laterTrack := Track{}
	seg := TrackSegment{}
	seg.TrackPoints = points
	FillTrackSegmentValues(&seg)
	laterTrack.TrackSegments = append(laterTrack.TrackSegments, seg)
	FillTrackValues(&laterTrack)
	laterTrack.NumberOfSegments = 1

	file.Tracks = append(file.Tracks, laterTrack)

	file.NumberOfTracks = 1
	FillTrackFileValues(&file)

	return file
}

func getTrackFileWithTimeGaps() TrackFile {
	file := getSimpleTrackFileWithTime()

	t1, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:13Z")
	t2, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:33Z")
	t3, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:53Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 108.0, t2)
	pnt3 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t3)
	points := []TrackPoint{pnt1, pnt2, pnt3}

	FillDistancesTrackPoint(&points[0], TrackPoint{}, points[1])
	FillDistancesTrackPoint(&points[1], points[0], points[2])
	FillDistancesTrackPoint(&points[2], points[1], TrackPoint{})
	FillValuesTrackPointArray(points, "none", 0.3, 10.0)
	laterTrack := Track{}
	seg := TrackSegment{}
	seg.TrackPoints = points
	FillTrackSegmentValues(&seg)
	laterTrack.TrackSegments = append(laterTrack.TrackSegments, seg)
	FillTrackValues(&laterTrack)
	laterTrack.NumberOfSegments = 1

	file.Tracks = append(file.Tracks, laterTrack)
	file.NumberOfTracks = 2
	FillTrackFileValues(&file)

	return file
}

func getSimpleTrackFile() TrackFile {
	ret := NewTrackFile("/mys/track/file")
	trk := getSimpleTrack()
	FillTrackValues(&trk)
	ret.Tracks = []Track{trk}
	FillTrackFileValues(&ret)
	ret.NumberOfTracks = 1

	return ret
}

func getSimpleTrackFileWithTime() TrackFile {
	return getSimpleTrackFileWithStartTime(DEFAULT_START_TIME)
}

func getSimpleTrackFileWithStartTime(startTime string) TrackFile {
	ret := NewTrackFile("/mys/track/file")
	trk := getSimpleTrackWithStartTime(startTime)
	FillTrackValues(&trk)
	ret.Tracks = []Track{trk}
	FillTrackFileValues(&ret)
	ret.NumberOfTracks = 1

	return ret
}

func getSimpleTrack() Track {
	ret := Track{}
	segs := getSimpleTrackSegment()
	FillTrackSegmentValues(&segs)
	ret.TrackSegments = []TrackSegment{segs}
	FillTrackValues(&ret)
	ret.NumberOfSegments = 1

	return ret
}

func getSimpleTrackWithTime() Track {
	return getSimpleTrackWithStartTime(DEFAULT_START_TIME)
}

func getSimpleTrackWithStartTime(startTime string) Track {
	ret := Track{}
	segs := getSimpleTrackSegmentWithStartTime(startTime)
	FillTrackSegmentValues(&segs)
	ret.TrackSegments = []TrackSegment{segs}
	FillTrackValues(&ret)
	ret.NumberOfSegments = 1

	return ret
}

func getSimpleTrackSegment() TrackSegment {
	seg := TrackSegment{}
	points := gerSimpleTrackPointArray()
	seg.TrackPoints = points

	return seg
}

func getSimpleTrackSegmentWithStartTime(startTime string) TrackSegment {
	seg := TrackSegment{}
	points := getSimpleTrackPointArrayWithStartTime(startTime)
	seg.TrackPoints = points

	return seg
}

func gerSimpleTrackPointArray() []TrackPoint {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	points := []TrackPoint{pnt1, pnt2, pnt3}

	FillDistancesTrackPoint(&points[0], TrackPoint{}, points[1])
	FillDistancesTrackPoint(&points[1], points[0], points[2])
	FillDistancesTrackPoint(&points[2], points[1], TrackPoint{})
	FillValuesTrackPointArray(points, "none", 0.3, 10.0)

	return points
}

func getSimpleTrackPointArrayWithTime() []TrackPoint {
	return getSimpleTrackPointArrayWithStartTime(DEFAULT_START_TIME)
}

func getSimpleTrackPointArrayWithStartTime(startTime string) []TrackPoint {
	t1, _ := time.Parse(time.RFC3339, startTime)
	t2 := t1.Add(time.Second * 10)
	t3 := t2.Add(time.Second * 10)
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 108.0, t2)
	pnt3 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t3)
	points := []TrackPoint{pnt1, pnt2, pnt3}

	FillDistancesTrackPoint(&points[0], TrackPoint{}, points[1])
	FillDistancesTrackPoint(&points[1], points[0], points[2])
	FillDistancesTrackPoint(&points[2], points[1], TrackPoint{})
	FillValuesTrackPointArray(points, "none", 0.3, 10.0)

	return points
}

func getSimpleTrackList() []Track {
	tracks := []Track{}
	tracks = append(tracks, getSimpleTrackWithStartTime("2014-08-21T17:19:33Z"))
	tracks = append(tracks, getSimpleTrackWithStartTime("2014-08-22T17:19:33Z"))
	tracks = append(tracks, getSimpleTrackWithStartTime("2014-08-23T17:19:33Z"))
	tracks = append(tracks, getSimpleTrackWithStartTime("2014-08-24T17:19:33Z"))

	return tracks
}

func getTrackFileWithMultipleTracks() TrackFile {
	ret := getSimpleTrackFileWithStartTime("2014-08-21T17:19:33Z")
	ret.Tracks = append(ret.Tracks, getSimpleTrackWithStartTime("2014-08-22T17:19:33Z"))
	ret.Tracks = append(ret.Tracks, getSimpleTrackWithStartTime("2014-08-23T17:19:33Z"))
	ret.Tracks = append(ret.Tracks, getSimpleTrackWithStartTime("2014-08-24T17:19:33Z"))
	FillTrackFileValues(&ret)

	return ret
}
