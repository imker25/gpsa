package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"tobi.backfrak.de/internal/gpsabl"
)

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
	for i, point := range points {
		pnt := convertPoint(point, i, points, pointCount)

		ret = append(ret, pnt)
	}

	return ret
}

func convertPoint(point Trkpt, i int, points []Trkpt, pointCount int) gpsabl.TrackPoint {
	pnt := gpsabl.TrackPoint{}
	pnt.Latitude = point.Latitude
	pnt.Longitude = point.Longitude
	pnt.Elevation = point.Elevation

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

	return pnt
}
