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

// ReadTracks - Read the *.gpx from the inputs GpxFile.FilePath, and return a GpxFile struct that contains all information
// Implement the gpsabl.TrackReader interface for *.gpx files
func (gpx *GpxFile) ReadTracks(correction string) (gpsabl.TrackFile, error) {
	ret, err := ReadGpxFile(gpx.FilePath, correction)

	if err == nil {
		gpx.TrackFile = ret
	}

	return ret, err
}

// ReadGpxFile - Reads a *.gpx file
func ReadGpxFile(filePath string, correction string) (gpsabl.TrackFile, error) {
	ret := gpsabl.TrackFile{}
	ret.FilePath = filePath

	gpx, fileError := ReadGPX(filePath)

	if fileError != nil {
		return ret, fileError
	}

	var tracks []gpsabl.Track
	for _, trk := range gpx.Tracks {

		// Add only tracks that contain segments
		if len(trk.TrackSegments) > 0 {
			track, convertError := ConvertTrk(trk, correction)
			if convertError != nil {
				return ret, convertError
			}
			tracks = append(tracks, track)
		}

	}

	if len(tracks) > 0 {
		ret.Tracks = tracks
		ret.Name = gpx.Name
		ret.Description = gpx.Description
		ret.NumberOfTracks = len(tracks)

		gpsabl.FillTrackFileValues(&ret)
	} else {
		return ret, newEmptyGpxFileError(filePath)
	}
	return ret, nil
}
