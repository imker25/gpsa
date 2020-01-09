package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"sort"
	"time"

	"tobi.backfrak.de/internal/gpsabl"
)

// ConvertTrk - Convert a gpxbl.Track to a gpsabl.Track
func ConvertTrk(track Trk, correction string) (gpsabl.Track, error) {

	res := gpsabl.Track{}
	var err error
	res.Name = track.Name
	res.NumberOfSegments = len(track.TrackSegments)
	res.Description = track.Description

	res.TrackSegments, err = convertSegments(track.TrackSegments, correction)
	if err != nil {
		return res, err
	}

	gpsabl.FillTrackValues(&res)

	return res, err
}

func convertSegments(segments []Trkseg, correction string) ([]gpsabl.TrackSegment, error) {
	var ret []gpsabl.TrackSegment
	var err error
	for _, seg := range segments {

		segment := gpsabl.TrackSegment{}
		segment.TrackPoints, err = convertPoints(seg.TrackPoints, correction)
		if err != nil {
			return nil, err
		}

		// Add only segments, that contain points
		if len(seg.TrackPoints) > 0 {
			gpsabl.FillTrackSegmentValues(&segment)
			ret = append(ret, segment)
		}
	}

	return ret, nil
}

func convertPoints(points []Trkpt, correction string) ([]gpsabl.TrackPoint, error) {
	var ret []gpsabl.TrackPoint

	pointCount := len(points)
	if pointCount > 10 { // I think it makes only sense to use go routines for tracks with more then 10 points
		c := make(chan gpsabl.TrackPoint, pointCount)
		pointCounter := 0
		for i, point := range points {
			go goConvertPointDistance(point, i, &points, pointCount, c)

		}
		for pnt := range c {
			ret = append(ret, pnt)
			pointCounter++
			if pointCounter == pointCount {
				close(c)
			}
		}
	} else {
		for i, point := range points {
			pnt := convertPointDistance(point, i, &points, pointCount)
			ret = append(ret, pnt)
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Number < ret[j].Number
	})

	err := gpsabl.FillValuesTrackPointArray(ret, correction)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func convertPointDistance(point Trkpt, i int, pnts *[]Trkpt, pointCount int) gpsabl.TrackPoint {
	pnt := convertBasicPointValues(point.Latitude, point.Longitude, point.Elevation, point.Time)
	pnt.Number = i
	points := *pnts

	if i == 0 && pointCount > 1 {
		pntNext := convertBasicPointValues(points[i+1].Latitude, points[i+1].Longitude, points[i+1].Elevation, points[i+1].Time)
		gpsabl.FillDistancesTrackPoint(&pnt, gpsabl.TrackPoint{}, pntNext)
	}

	if i > 0 && i < pointCount-1 {
		pntNext := convertBasicPointValues(points[i+1].Latitude, points[i+1].Longitude, points[i+1].Elevation, points[i+1].Time)
		pntBefore := convertBasicPointValues(points[i-1].Latitude, points[i-1].Longitude, points[i-1].Elevation, points[i-1].Time)
		gpsabl.FillDistancesTrackPoint(&pnt, pntBefore, pntNext)
	}

	if i == pointCount-1 && pointCount > 1 {
		pntBefore := convertBasicPointValues(points[i-1].Latitude, points[i-1].Longitude, points[i-1].Elevation, points[i-1].Time)
		gpsabl.FillDistancesTrackPoint(&pnt, pntBefore, gpsabl.TrackPoint{})
	}

	return pnt
}

func goConvertPointDistance(point Trkpt, i int, pnts *[]Trkpt, pointCount int, c chan gpsabl.TrackPoint) {

	c <- convertPointDistance(point, i, pnts, pointCount)
}

func convertBasicPointValues(latitude, longitude, elevation float32, timeStamp string) gpsabl.TrackPoint {
	pnt := gpsabl.TrackPoint{}
	pnt.Latitude = latitude
	pnt.Longitude = longitude
	pnt.Elevation = elevation

	if timeStamp == "" {
		pnt.TimeValid = false
	} else {

		t, err := time.Parse(time.RFC3339, timeStamp)

		// In case the time stamp of the track point is not in the specifyed format, it is not valid
		if err != nil {
			pnt.TimeValid = false
			return pnt
		}

		pnt.TimeValid = true
		pnt.Time = t
	}

	return pnt
}
