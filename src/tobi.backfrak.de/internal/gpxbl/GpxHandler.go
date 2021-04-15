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
func (gpx *GpxFile) ReadTracks(correction gpsabl.CorrectionParameter, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {
	ret, err := ReadGpxFile(gpx.FilePath, correction, minimalMovingSpeed, minimalStepHight)

	if err == nil {
		gpx.TrackFile = ret
	}

	return ret, err
}

// ReadGpxFile - Reads a *.gpx file
func ReadGpxFile(filePath string, correction gpsabl.CorrectionParameter, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {
	gpx, fileError := ReadGPX(filePath)

	if fileError != nil {
		return gpsabl.TrackFile{}, fileError
	}
	ret, convertError := ConvertGPXFile(gpx, filePath, correction, minimalMovingSpeed, minimalStepHight)
	if convertError != nil {
		return gpsabl.TrackFile{}, convertError
	}

	return ret, nil
}
