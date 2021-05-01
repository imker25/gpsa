package gpsabl

import "time"

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// TrackReader - The interface for all functions that can read gps data files like *.gpx
type TrackReader interface {
	// Get a new reader of this type for a given InputFile
	NewReader(data InputFile) TrackReader
	// ReadTracks - Read the tracks that are realted to this instance
	ReadTracks(correction CorrectionParameter, minimalMovingSpeed float64, minimalStepHight float64) (TrackFile, error)
	// Check if a reader of this type can read this buffer
	CheckBuffer(buffer []byte) bool
	// Check if a reader of this type can read this a file with the extencion of the given path
	CheckFile(path string) bool
	// Check if a given InputFile ca be read by a reader of this type
	CheckInputFile(input InputFile) bool
	// Get a new InputFile instance for a buffer that can be read by a reader of this type
	NewInputFileForBuffer(buffer []byte, name string) *InputFile
}

// TrackSummaryProvider - Interface for classes that provide track summary data
type TrackSummaryProvider interface {
	GetDistance() float64
	GetHorizontalDistance() float64
	GetAltitudeRange() float32
	GetMinimumAltitude() float32
	GetMaximumAltitude() float32
	GetElevationGain() float32
	GetElevationLose() float32
	GetUpwardsDistance() float64
	GetDownwardsDistance() float64
	GetTimeDataValid() bool
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetMovingTime() time.Duration
	GetUpwardsTime() time.Duration
	GetDownwardsTime() time.Duration
	GetAvarageSpeed() float64
	GetUpwardsSpeed() float64
	GetDownwardsSpeed() float64
}

// TrackSummarySetter - Interface for classes that can set track summary data
type TrackSummarySetter interface {
	SetValues(distance float64,
		horizontalDistance float64,
		minimumAltitude float32,
		maximumAltitude float32,
		elevationGain float32,
		elevationLose float32,
		upwardsDistance float64,
		downwardsDistance float64,
		timeDataValid bool,
		startTime time.Time,
		endTime time.Time,
		movingTime time.Duration,
		upwardsTime time.Duration,
		downwards time.Duration)
}
