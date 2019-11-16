package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import "tobi.backfrak.de/internal/gpsabl"

// GpxFile - The struct to handle *.gpx data files
type GpxFile struct {
	gpsabl.TrackFile
}

// NewGpxFile - Constructor for the GpxFile struct
func NewGpxFile(filePath string) GpxFile {
	gpx := GpxFile{}
	gpx.FilePath = filePath

	return gpx
}

// ReadTracks - implement the gpsabl.TrackReader interface for *.gpx files
func (gpx GpxFile) ReadTracks() ([]gpsabl.Track, error) {
	return ReadGpxFile(gpx.FilePath)
}

// ReadGpxFile - Reads a *.gpx file
func ReadGpxFile(filepath string) ([]gpsabl.Track, error) {
	var ret []gpsabl.Track
	gpx, err := ReadGPX(filepath)

	if err != nil {
		return ret, err
	}

	var tracks []gpsabl.Track
	for _, trk := range gpx.Tracks {
		track := ConvertTrk(trk)
		tracks = append(tracks, track)

	}

	return tracks, nil
}
