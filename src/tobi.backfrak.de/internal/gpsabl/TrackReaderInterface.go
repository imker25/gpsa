package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// TrackReader - The interface for all functions that can read gps data files like *.gpx
type TrackReader interface {
	ReadTracks() (TrackFile, error)
}

// TrackInfoProvider - Interface for classes that provide track info data
type TrackInfoProvider interface {
	GetDistance() float64
	GetAtituteRange() float32
	GetMinimumAtitute() float32
	GetMaximumAtitute() float32
}
