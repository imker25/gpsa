package tcxbl

import (
	"sort"
	"time"

	"tobi.backfrak.de/internal/gpsabl"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// ConvertTcx - Convert a tcxbl.Tcx to a gpsabl.TrackFile
func ConvertTcx(tcx Tcx, filePath string, correction string, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {

	res := gpsabl.NewTrackFile(filePath)

	for _, activityArray := range tcx.ActivityArray {
		for _, activity := range activityArray.Activities {
			track, err := convertActivity(activity, correction, minimalMovingSpeed, minimalStepHight)
			if err != nil {
				return gpsabl.TrackFile{}, err
			}
			res.Tracks = append(res.Tracks, track)
		}
	}
	res.NumberOfTracks = len(res.Tracks)
	gpsabl.FillTrackFileValues(&res)

	return res, nil
}

func convertActivity(activity Activity, correction string, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.Track, error) {

	res := gpsabl.Track{}
	for _, lap := range activity.Laps {
		seg, err := convertLap(lap, correction, minimalMovingSpeed, minimalStepHight)
		if err != nil {
			return gpsabl.Track{}, err
		}
		res.TrackSegments = append(res.TrackSegments, seg)
	}
	res.Name = activity.ID
	res.NumberOfSegments = len(res.TrackSegments)
	gpsabl.FillTrackValues(&res)

	return res, nil
}

func convertLap(lap Lap, correction string, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackSegment, error) {
	res := gpsabl.TrackSegment{}
	trackPoints := []Trackpoint{}
	for _, track := range lap.Tracks {
		for _, trackPoint := range track.Trackpoints {
			trackPoints = append(trackPoints, trackPoint)
		}
	}
	retArr, err := convertTrackpoints(trackPoints, correction, minimalMovingSpeed, minimalStepHight)
	if err != nil {
		return gpsabl.TrackSegment{}, err
	}
	res.TrackPoints = retArr
	gpsabl.FillTrackSegmentValues(&res)
	return res, nil
}

func convertTrackpoints(points []Trackpoint, correction string, minimalMovingSpeed float64, minimalStepHight float64) ([]gpsabl.TrackPoint, error) {
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

	err := gpsabl.FillValuesTrackPointArray(ret, correction, minimalMovingSpeed, minimalStepHight)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func goConvertPointDistance(point Trackpoint, i int, pnts *[]Trackpoint, pointCount int, c chan gpsabl.TrackPoint) {

	c <- convertPointDistance(point, i, pnts, pointCount)
}

func convertPointDistance(point Trackpoint, i int, pnts *[]Trackpoint, pointCount int) gpsabl.TrackPoint {
	pnt := convertBasicPointValues(point)
	pnt.Number = i
	points := *pnts

	if i == 0 && pointCount > 1 {
		pntNext := convertBasicPointValues(points[i+1])
		gpsabl.FillDistancesTrackPoint(&pnt, gpsabl.TrackPoint{}, pntNext)
	}

	if i > 0 && i < pointCount-1 {
		pntNext := convertBasicPointValues(points[i+1])
		pntBefore := convertBasicPointValues(points[i-1])
		gpsabl.FillDistancesTrackPoint(&pnt, pntBefore, pntNext)
	}

	if i == pointCount-1 && pointCount > 1 {
		pntBefore := convertBasicPointValues(points[i-1])
		gpsabl.FillDistancesTrackPoint(&pnt, pntBefore, gpsabl.TrackPoint{})
	}

	return pnt
}

func convertBasicPointValues(point Trackpoint) gpsabl.TrackPoint {
	pnt := gpsabl.TrackPoint{}
	pnt.Latitude = point.Position.LatitudeDegrees
	pnt.Longitude = point.Position.LongitudeDegrees
	pnt.Elevation = point.AltitudeMeters

	if point.Time == "" {
		pnt.TimeValid = false
	} else {

		t, err := time.Parse(time.RFC3339, point.Time)

		// In case the time stamp of the track point is not in the specified format, it is not valid
		if err != nil {
			pnt.TimeValid = false
			return pnt
		}

		pnt.TimeValid = true
		pnt.Time = t
	}

	return pnt
}
