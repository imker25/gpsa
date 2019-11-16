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

	var dist float64
	var minimumAtitute float32
	var maximumAtitute float32

	for i, track := range segment.TrackPoints {
		dist = dist + track.HorizontalDistanceNext

		if i == 0 || track.Elevation > maximumAtitute {
			maximumAtitute = track.Elevation
		}

		if i == 0 || track.Elevation < minimumAtitute {
			minimumAtitute = track.Elevation
		}
	}

	ret := TrackSegment{}
	ret.AtituteRange = maximumAtitute - minimumAtitute
	ret.MaximumAtitute = maximumAtitute
	ret.MinimumAtitute = minimumAtitute
	ret.Distance = dist
	ret.TrackPoints = segment.TrackPoints

	return ret
}

// FillTrackValues - Fills the distance and atitute fields of a tack  by adding up all TrackSegments distances
func FillTrackValues(track Track) Track {
	var dist float64
	var minimumAtitute float32
	var maximumAtitute float32

	for i, seg := range track.TrackSegments {
		dist = dist + seg.Distance

		if i == 0 || seg.MaximumAtitute > maximumAtitute {
			maximumAtitute = seg.MaximumAtitute
		}

		if i == 0 || seg.MinimumAtitute < minimumAtitute {
			minimumAtitute = seg.MinimumAtitute
		}
	}

	ret := Track{}
	ret.AtituteRange = maximumAtitute - minimumAtitute
	ret.MaximumAtitute = maximumAtitute
	ret.MinimumAtitute = minimumAtitute
	ret.Distance = dist

	ret.Name = track.Name
	ret.NumberOfSegments = track.NumberOfSegments
	ret.Description = track.Description
	ret.TrackSegments = track.TrackSegments

	return ret
}

// FillTrackFileValues - Fills the distance and atitute fields of a tack  by adding up all TrackSegments distances
func FillTrackFileValues(file TrackFile) TrackFile {
	var dist float64
	var minimumAtitute float32
	var maximumAtitute float32

	for i, trk := range file.Tracks {
		dist = dist + trk.Distance

		if i == 0 || trk.MaximumAtitute > maximumAtitute {
			maximumAtitute = trk.MaximumAtitute
		}

		if i == 0 || trk.MinimumAtitute < minimumAtitute {
			minimumAtitute = trk.MinimumAtitute
		}
	}

	ret := TrackFile{}
	ret.AtituteRange = maximumAtitute - minimumAtitute
	ret.MaximumAtitute = maximumAtitute
	ret.MinimumAtitute = minimumAtitute
	ret.Distance = dist

	ret.Name = file.Name
	ret.NumberOfTracks = len(file.Tracks)
	ret.Description = file.Description
	ret.Tracks = file.Tracks
	ret.FilePath = file.FilePath

	return ret
}
