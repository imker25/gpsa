package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// FillDistancesTrackPoint - Adds the distance values to the basePoint
func FillDistancesTrackPoint(basePoint, beforePoint, nextPoint TrackPoint) TrackPoint {
	retPoint := TrackPoint{}
	retPoint.Elevation = basePoint.Elevation
	retPoint.Latitude = basePoint.Latitude
	retPoint.Longitude = basePoint.Longitude

	if (beforePoint != TrackPoint{}) {
		retPoint.HorizontalDistanceBefore = Distance(basePoint, beforePoint)
		retPoint.VerticalDistanceBefore = basePoint.Elevation - beforePoint.Elevation
	}

	if (nextPoint != TrackPoint{}) {
		retPoint.HorizontalDistanceNext = Distance(basePoint, nextPoint)
		retPoint.VerticalDistanceNext = nextPoint.Elevation - basePoint.Elevation
	}

	return retPoint
}

// FillTrackSegmentValues - Fills the distance and atitute fields of a tack segment by adding up all TrackPoint distances
func FillTrackSegmentValues(segment TrackSegment) TrackSegment {
	iPnts := []TrackSummaryProvider{}
	for i := range segment.TrackPoints {
		iPnt := TrackSummaryProvider(&segment.TrackPoints[i])
		iPnts = append(iPnts, iPnt)
	}

	ret := TrackSegment{}
	fillTrackSummaryValues(&ret, iPnts)
	ret.TrackPoints = segment.TrackPoints

	return ret
}

func fillTrackSummaryValues(target TrackSummarySetter, input []TrackSummaryProvider) {
	var dist float64
	var minimumAtitute float32
	var maximumAtitute float32

	for i, sum := range input {
		dist = dist + sum.GetDistance()

		if i == 0 || sum.GetMaximumAtitute() > maximumAtitute {
			maximumAtitute = sum.GetMaximumAtitute()
		}

		if i == 0 || sum.GetMinimumAtitute() < minimumAtitute {
			minimumAtitute = sum.GetMinimumAtitute()
		}
	}

	target.SetValues(dist, minimumAtitute, maximumAtitute)
}

// FillTrackValues - Fills the distance and atitute fields of a tack  by adding up all TrackSegments distances
func FillTrackValues(track Track) Track {
	iSegs := []TrackSummaryProvider{}
	for i := range track.TrackSegments {
		iSeg := TrackSummaryProvider(&track.TrackSegments[i])
		iSegs = append(iSegs, iSeg)
	}

	ret := Track{}
	fillTrackSummaryValues(&ret, iSegs)

	ret.Name = track.Name
	ret.NumberOfSegments = track.NumberOfSegments
	ret.Description = track.Description
	ret.TrackSegments = track.TrackSegments

	return ret
}

// FillTrackFileValues - Fills the distance and atitute fields of a tack  by adding up all TrackSegments distances
func FillTrackFileValues(file TrackFile) TrackFile {
	iSegs := []TrackSummaryProvider{}
	for i := range file.Tracks {
		iSeg := TrackSummaryProvider(&file.Tracks[i])
		iSegs = append(iSegs, iSeg)
	}

	ret := TrackFile{}
	fillTrackSummaryValues(&ret, iSegs)

	ret.Name = file.Name
	ret.NumberOfTracks = len(file.Tracks)
	ret.Description = file.Description
	ret.Tracks = file.Tracks
	ret.FilePath = file.FilePath

	return ret
}
