package gpsabl

import (
	"time"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// TrackSummary - the struct to store track statistic data
type TrackSummary struct {
	Distance          float64
	MinimumAltitude   float32
	MaximumAltitude   float32
	ElevationGain     float32
	ElevationLose     float32
	UpwardsDistance   float64
	DownwardsDistance float64
	TimeDataValid     bool
	StartTime         time.Time
	EndTime           time.Time
	MovingTime        time.Duration
	AvarageSpeed      float64
}

// SetValues - Set the Values of a TrackSummary (Implement the TrackSummaryProvider )
func (sum *TrackSummary) SetValues(distance float64,
	minimumAltitude float32,
	maximumAltitude float32,
	elevationGain float32,
	elevationLose float32,
	upwardsDistance float64,
	downwardsDistance float64,
	timeDataValid bool,
	startTime time.Time,
	endTime time.Time,
	movigTime time.Duration) {

	sum.MinimumAltitude = minimumAltitude
	sum.MaximumAltitude = maximumAltitude
	sum.Distance = distance
	sum.DownwardsDistance = downwardsDistance
	sum.UpwardsDistance = upwardsDistance
	sum.ElevationGain = elevationGain
	sum.ElevationLose = elevationLose
	sum.StartTime = startTime
	sum.EndTime = endTime
	sum.TimeDataValid = timeDataValid
	sum.MovingTime = movigTime
	if timeDataValid && movigTime > 0 {
		sum.AvarageSpeed = sum.Distance / float64(sum.MovingTime/1000000000)
	}
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

// GetAltitudeRange - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetAltitudeRange() float32 {
	return sum.MaximumAltitude - sum.MinimumAltitude
}

// GetMaximumAltitude Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetMaximumAltitude() float32 {
	return sum.MaximumAltitude
}

// GetMinimumAltitude - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetMinimumAltitude() float32 {
	return sum.MinimumAltitude
}

// GetStartTime - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetStartTime() time.Time {
	return sum.StartTime
}

// GetEndTime - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetEndTime() time.Time {
	return sum.EndTime
}

// GetTimeDataValid - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetTimeDataValid() bool {
	return sum.TimeDataValid
}

// GetMovingTime - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetMovingTime() time.Duration {
	return sum.MovingTime
}

// GetAvarageSpeed - Implement the TrackSummaryProvider interface for TrackSummary
func (sum TrackSummary) GetAvarageSpeed() float64 {
	return sum.AvarageSpeed
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
	Time                     time.Time
	TimeValid                bool
	HorizontalDistanceBefore float64
	HorizontalDistanceNext   float64
	DistanceNext             float64
	DistanceBefore           float64
	DistanceToThisPoint      float64
	CorectedElevation        float32
	VerticalDistanceBefore   float32
	VerticalDistanceNext     float32
	CountUpwards             bool
	CountDownwards           bool
	CountMoving              bool
	MovingTime               time.Duration
	TimeDurationBefore       time.Duration
	TimeDurationNext         time.Duration
	AvarageSpeed             float64
	SpeedBefore              float64
	SpeedNext                float64
}

// GetDistance - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetDistance() float64 {
	return pnt.DistanceBefore
}

// GetAltitudeRange - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetAltitudeRange() float32 {
	return 0.0
}

// GetMaximumAltitude Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetMaximumAltitude() float32 {
	return pnt.Elevation
}

// GetMinimumAltitude - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetMinimumAltitude() float32 {
	return pnt.Elevation
}

// GetElevationGain - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetElevationGain() float32 {
	if pnt.VerticalDistanceNext > 0 {
		return pnt.VerticalDistanceNext
	}
	return 0
}

// GetElevationLose - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetElevationLose() float32 {
	if pnt.VerticalDistanceNext < 0 {
		return pnt.VerticalDistanceNext
	}
	return 0
}

// GetUpwardsDistance - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetUpwardsDistance() float64 {
	if pnt.CountUpwards {
		return pnt.DistanceNext
	}
	return 0
}

// GetDownwardsDistance - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetDownwardsDistance() float64 {
	if pnt.CountDownwards {
		return pnt.DistanceNext
	}
	return 0
}

// GetStartTime - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetStartTime() time.Time {
	return pnt.Time
}

// GetEndTime - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetEndTime() time.Time {
	return pnt.Time
}

// GetTimeDataValid - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetTimeDataValid() bool {
	return pnt.TimeValid
}

// GetMovingTime - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetMovingTime() time.Duration {
	return pnt.MovingTime
}

// GetAvarageSpeed - Implement the TrackSummaryProvider interface for TrackPoint
func (pnt TrackPoint) GetAvarageSpeed() float64 {
	return pnt.AvarageSpeed
}
