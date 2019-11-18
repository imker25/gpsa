package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"sort"
	"sync"

	"tobi.backfrak.de/internal/gpsabl"
)

var wg sync.WaitGroup

// ConvertTrk - Convert a gpxbl.Trk to a gpsabl.Track
func ConvertTrk(track Trk) gpsabl.Track {

	info := GetTrackInfo(track)

	return ConvertTrackInfo(info)
}

// ConvertTrackInfo - Convert a gpxbl.TrackInfo to a gpsabl.Track
func ConvertTrackInfo(track TrackInfo) gpsabl.Track {

	res := gpsabl.Track{}

	res.Name = track.Name
	res.NumberOfSegments = track.NumberOfSegments
	res.Description = track.Description

	res.TrackSegments = convertSegments(track.Track.TrackSegments)

	res = gpsabl.FillTrackValues(res)

	return res
}

func convertSegments(segments []Trkseg) []gpsabl.TrackSegment {
	var ret []gpsabl.TrackSegment

	for _, seg := range segments {

		segment := gpsabl.TrackSegment{}
		segment.TrackPoints = convertPoints(seg.TrackPoints)

		segment = gpsabl.FillTrackSegmentValues(segment)
		ret = append(ret, segment)
	}

	return ret
}

func convertPoints(points []Trkpt) []gpsabl.TrackPoint {
	var ret []gpsabl.TrackPoint

	pointCount := len(points)
	c := make(chan gpsabl.TrackPoint)
	pointCounter := 0
	for i, point := range points {
		wg.Add(1)
		go goConvertPoint(point, i, &points, pointCount, c)

	}
	for pnt := range c {
		ret = append(ret, pnt)
		pointCounter++
		if pointCounter == pointCount {
			close(c)
		}
	}
	wg.Wait()
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Number < ret[j].Number
	})
	return ret
}

func convertPoint(point Trkpt, i int, pnts *[]Trkpt, pointCount int) gpsabl.TrackPoint {
	pnt := gpsabl.TrackPoint{}
	pnt.Latitude = point.Latitude
	pnt.Longitude = point.Longitude
	pnt.Elevation = point.Elevation
	pnt.Number = i
	points := *pnts

	if i == 0 && pointCount > 1 {
		pntNext := gpsabl.TrackPoint{}
		pntNext.Latitude = points[i+1].Latitude
		pntNext.Longitude = points[i+1].Longitude
		pntNext.Elevation = points[i+1].Elevation
		pnt = gpsabl.FillDistancesTrackPoint(pnt, gpsabl.TrackPoint{}, pntNext)
	}

	if i > 0 && i < pointCount-1 {
		pntNext := gpsabl.TrackPoint{}
		pntNext.Latitude = points[i+1].Latitude
		pntNext.Longitude = points[i+1].Longitude
		pntNext.Elevation = points[i+1].Elevation

		pntBefore := gpsabl.TrackPoint{}
		pntBefore.Latitude = points[i-1].Latitude
		pntBefore.Longitude = points[i-1].Longitude
		pntBefore.Elevation = points[i-1].Elevation

		pnt = gpsabl.FillDistancesTrackPoint(pnt, pntBefore, pntNext)
	}

	if i == pointCount-1 && pointCount > 1 {
		pntBefore := gpsabl.TrackPoint{}
		pntBefore.Latitude = points[i-1].Latitude
		pntBefore.Longitude = points[i-1].Longitude
		pntBefore.Elevation = points[i-1].Elevation
		pnt = gpsabl.FillDistancesTrackPoint(pnt, pntBefore, gpsabl.TrackPoint{})
	}
	defer wg.Done()
	return pnt
}

func goConvertPoint(point Trkpt, i int, pnts *[]Trkpt, pointCount int, c chan gpsabl.TrackPoint) {

	c <- convertPoint(point, i, pnts, pointCount)
}
