package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"sort"

	"tobi.backfrak.de/internal/gpsabl"
)

// ConvertTrk - Convert a gpxbl.Track to a gpsabl.Track
func ConvertTrk(track Trk) gpsabl.Track {

	res := gpsabl.Track{}

	res.Name = track.Name
	res.NumberOfSegments = len(track.TrackSegments)
	res.Description = track.Description

	res.TrackSegments = convertSegments(track.TrackSegments)

	gpsabl.FillTrackValues(&res)

	return res
}

func convertSegments(segments []Trkseg) []gpsabl.TrackSegment {
	var ret []gpsabl.TrackSegment

	for _, seg := range segments {

		segment := gpsabl.TrackSegment{}
		segment.TrackPoints = convertPoints(seg.TrackPoints)

		gpsabl.FillTrackSegmentValues(&segment)
		ret = append(ret, segment)
	}

	return ret
}

func convertPoints(points []Trkpt) []gpsabl.TrackPoint {
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

	gpsabl.FillDistanceToThisPoint(ret)
	gpsabl.FillCorectedElevationTrackPoint(ret)
	gpsabl.FillElevationGainLoseTrackPoint(ret)

	return ret
}

func convertPointDistance(point Trkpt, i int, pnts *[]Trkpt, pointCount int) gpsabl.TrackPoint {
	pnt := convertBasicPointValues(point.Latitude, point.Longitude, point.Elevation)
	pnt.Number = i
	points := *pnts

	if i == 0 && pointCount > 1 {
		pntNext := convertBasicPointValues(points[i+1].Latitude, points[i+1].Longitude, points[i+1].Elevation)
		gpsabl.FillDistancesTrackPoint(&pnt, gpsabl.TrackPoint{}, pntNext)
	}

	if i > 0 && i < pointCount-1 {
		pntNext := convertBasicPointValues(points[i+1].Latitude, points[i+1].Longitude, points[i+1].Elevation)
		pntBefore := convertBasicPointValues(points[i-1].Latitude, points[i-1].Longitude, points[i-1].Elevation)
		gpsabl.FillDistancesTrackPoint(&pnt, pntBefore, pntNext)
	}

	if i == pointCount-1 && pointCount > 1 {
		pntBefore := convertBasicPointValues(points[i-1].Latitude, points[i-1].Longitude, points[i-1].Elevation)
		gpsabl.FillDistancesTrackPoint(&pnt, pntBefore, gpsabl.TrackPoint{})
	}

	return pnt
}

func goConvertPointDistance(point Trkpt, i int, pnts *[]Trkpt, pointCount int, c chan gpsabl.TrackPoint) {

	c <- convertPointDistance(point, i, pnts, pointCount)
}

func convertBasicPointValues(latitude, longitude, elevation float32) gpsabl.TrackPoint {
	pnt := gpsabl.TrackPoint{}
	pnt.Latitude = latitude
	pnt.Longitude = longitude
	pnt.Elevation = elevation

	return pnt
}
