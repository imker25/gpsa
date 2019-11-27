package gpsabl

import "math"

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// TrackSummary - the struct to store track statistic data
type TrackSummary struct {
	Distance          float64
	AtituteRange      float32
	MinimumAtitute    float32
	MaximumAtitute    float32
	ElevationGain     float32
	ElevationLose     float32
	UpwardsDistance   float64
	DownwardsDistance float64
}

// SetValues - Set the Values of a TrackSummary (Implement the TrackSummaryProvider )
func (sum *TrackSummary) SetValues(distance float64, minimumAtitute float32, maximumAtitute float32, elevationGain float32, elevationLose float32, upwardsDistance float64, downwardsDistance float64) {
	sum.MinimumAtitute = minimumAtitute
	sum.MaximumAtitute = maximumAtitute
	sum.AtituteRange = maximumAtitute - minimumAtitute
	sum.Distance = distance
	sum.DownwardsDistance = downwardsDistance
	sum.UpwardsDistance = upwardsDistance
	sum.ElevationGain = elevationGain
	sum.ElevationLose = elevationLose
}

// GetElevationGain - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetElevationGain() float32 {
	return sum.ElevationGain
}

// GetElevationLose - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetElevationLose() float32 {
	return sum.ElevationLose
}

// GetUpwardsDistance - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetUpwardsDistance() float64 {
	return sum.UpwardsDistance
}

// GetDownwardsDistance - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetDownwardsDistance() float64 {
	return sum.DownwardsDistance
}

// GetDistance - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetDistance() float64 {
	return sum.Distance
}

// GetAtituteRange - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetAtituteRange() float32 {
	return sum.AtituteRange
}

// GetMaximumAtitute Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetMaximumAtitute() float32 {
	return sum.MaximumAtitute
}

// GetMinimumAtitute - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetMinimumAtitute() float32 {
	return sum.MinimumAtitute
}

// TrackFile - A struct to handle track files
type TrackFile struct {
	TrackSummary
	FilePath       string
	Name           string
	Description    string
	NumberOfTracks int
	Tracks         []Track
}

// NewTrackFile - Constructor for the TrackFile struct
func NewTrackFile(filePath string) TrackFile {
	ret := TrackFile{}
	ret.FilePath = filePath

	return ret
}

// Track - the struct to handle track info in gpsa
type Track struct {
	TrackSummary
	Name             string
	Description      string
	NumberOfSegments int

	TrackSegments []TrackSegment
}

// TrackSegment - the struct to handle track segment info in gpsa
type TrackSegment struct {
	TrackSummary
	TrackPoints []TrackPoint
}

// TrackPoint - the struct to handle track point info in gpsa
type TrackPoint struct {
	Number                   int
	Elevation                float32
	Latitude                 float32
	Longitude                float32
	HorizontalDistanceBefore float64
	HorizontalDistanceNext   float64
	VerticalDistanceBefore   float32
	VerticalDistanceNext     float32
}

// GetDistance - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetDistance() float64 {
	return pnt.HorizontalDistanceNext
}

// GetAtituteRange - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetAtituteRange() float32 {
	return 0.0
}

// GetMaximumAtitute Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetMaximumAtitute() float32 {
	return pnt.Elevation
}

// GetMinimumAtitute - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetMinimumAtitute() float32 {
	return pnt.Elevation
}

// GetElevationGain - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetElevationGain() float32 {
	if pnt.VerticalDistanceBefore > 0 {
		return pnt.VerticalDistanceBefore
	}
	return 0
}

// GetElevationLose - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetElevationLose() float32 {
	if pnt.VerticalDistanceBefore < 0 {
		return float32(math.Abs(float64(pnt.VerticalDistanceBefore)))
	}
	return 0
}

// GetUpwardsDistance - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetUpwardsDistance() float64 {
	if pnt.VerticalDistanceNext > 0 {
		return pnt.HorizontalDistanceNext
	}
	return 0
}

// GetDownwardsDistance - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetDownwardsDistance() float64 {
	if pnt.VerticalDistanceNext < 0 {
		return pnt.HorizontalDistanceNext
	}
	return 0
}
