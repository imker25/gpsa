package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// FillDistancesTrackPoint - Adds the distance values to the basePoint
func FillDistancesTrackPoint(basePoint *TrackPoint, beforePoint TrackPoint, nextPoint TrackPoint) {

	if (beforePoint != TrackPoint{}) {
		basePoint.HorizontalDistanceBefore = Distance(*basePoint, beforePoint)
		basePoint.VerticalDistanceBefore = basePoint.Elevation - beforePoint.Elevation
	}

	if (nextPoint != TrackPoint{}) {
		basePoint.HorizontalDistanceNext = Distance(*basePoint, nextPoint)
		basePoint.VerticalDistanceNext = nextPoint.Elevation - basePoint.Elevation
	}
}

// FillTrackSegmentValues - Fills the distance and atitute fields of a tack segment by adding up all TrackPoint distances
func FillTrackSegmentValues(segment *TrackSegment) {
	iPnts := []TrackSummaryProvider{}
	for i := range segment.TrackPoints {
		iPnt := TrackSummaryProvider(&segment.TrackPoints[i])
		iPnts = append(iPnts, iPnt)
	}

	fillTrackSummaryValues(segment, iPnts)
}

// FillTrackValues - Fills the distance and atitute fields of a tack  by adding up all TrackSegments distances
func FillTrackValues(track *Track) {
	iSegs := []TrackSummaryProvider{}
	for i := range track.TrackSegments {
		iSeg := TrackSummaryProvider(&track.TrackSegments[i])
		iSegs = append(iSegs, iSeg)
	}

	fillTrackSummaryValues(track, iSegs)
}

// FillTrackFileValues - Fills the distance and atitute fields of a tack  by adding up all TrackSegments distances
func FillTrackFileValues(file *TrackFile) {
	iSegs := []TrackSummaryProvider{}
	for i := range file.Tracks {
		iSeg := TrackSummaryProvider(&file.Tracks[i])
		iSegs = append(iSegs, iSeg)
	}

	fillTrackSummaryValues(file, iSegs)

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
