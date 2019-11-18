package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// TrackSummary - the struct to store track statistic data
type TrackSummary struct {
	Distance       float64
	AtituteRange   float32
	MinimumAtitute float32
	MaximumAtitute float32
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
	Elevation                float32
	Latitude                 float32
	Longitude                float32
	HorizontalDistanceBefore float64
	HorizontalDistanceNext   float64
	VerticalDistanceBefore   float32
	VerticalDistanceNext     float32
}
