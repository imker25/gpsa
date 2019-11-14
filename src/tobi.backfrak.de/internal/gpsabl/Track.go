package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// Track - the struct to handle track info in gpsa
type Track struct {
	Name             string
	Description      string
	NumberOfSegments int
	Distance         float64
	AtituteRange     float32
	MinimumAtitute   float32
	MaximumAtitute   float32

	TrackSegments []TrackSegment
}

// GetDinstance - Get the Distance. TrackInfo interface implementaion
func (trk Track) GetDinstance() float64 {
	return trk.Distance
}

// GetAtituteRange - Get the AtituteRange. TrackInfo interface implementaion
func (trk Track) GetAtituteRange() float32 {
	return trk.AtituteRange
}

// GetMinimumAtitute - Get the MinimumAtitute. TrackInfo interface implementaion
func (trk Track) GetMinimumAtitute() float32 {
	return trk.MinimumAtitute
}

// GetMaximumAtitute - Get the MaximumAtitute. TrackInfo interface implementaion
func (trk Track) GetMaximumAtitute() float32 {
	return trk.MaximumAtitute
}

// TrackSegment - the struct to handle track segment info in gpsa
type TrackSegment struct {
	TrackPoints    []TrackPoint
	Distance       float64
	AtituteRange   float32
	MinimumAtitute float32
	MaximumAtitute float32
}

// GetDinstance - Get the Distance. TrackInfo interface implementaion
func (trk TrackSegment) GetDinstance() float64 {
	return trk.Distance
}

// GetAtituteRange - Get the AtituteRange. TrackInfo interface implementaion
func (trk TrackSegment) GetAtituteRange() float32 {
	return trk.AtituteRange
}

// GetMinimumAtitute - Get the MinimumAtitute. TrackInfo interface implementaion
func (trk TrackSegment) GetMinimumAtitute() float32 {
	return trk.MinimumAtitute
}

// GetMaximumAtitute - Get the MaximumAtitute. TrackInfo interface implementaion
func (trk TrackSegment) GetMaximumAtitute() float32 {
	return trk.MaximumAtitute
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
