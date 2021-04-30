package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"strings"

	"tobi.backfrak.de/internal/gpsabl"
)

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

// ReadTracks - Read the *.gpx from the inputs GpxFile.FilePath, and return a gpsabl.TrackFile struct that contains all information
// Implement the gpsabl.TrackReader interface for *.gpx files
func (gpx *GpxFile) ReadTracks(correction gpsabl.CorrectionParameter, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {
	ret, err := ReadGpxFile(gpx.FilePath, correction, minimalMovingSpeed, minimalStepHight)

	if err == nil {
		gpx.TrackFile = ret
	}

	return ret, err
}

// ReadBuffer - Read the *.gpx data from a buffer, and return a gpsabl.TrackFile struct that contains all information
// When using this method, the FilePath property may contain any string
// Implement the gpsabl.TrackReader interface for *.gpx files
func (gpx *GpxFile) ReadBuffer(buffer []byte, correction gpsabl.CorrectionParameter, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {
	content, readErr := readGPXBuffer(buffer, gpx.FilePath)
	if readErr != nil {
		return gpsabl.TrackFile{}, readErr
	}
	ret, convertError := ConvertGPXFile(content, gpx.FilePath, correction, minimalMovingSpeed, minimalStepHight)
	if convertError != nil {
		return gpsabl.TrackFile{}, convertError
	}

	return ret, nil
}

// CheckFile - Check if a file can be read by the GpxFile "class"
func (gpx *GpxFile) CheckFile(path string) bool {
	if strings.HasSuffix(path, "gpx") == true { // If the file is a *.gpx, we can read it
		return true
	}

	return false
}

// CheckBuffer - Check if a buffer can be read by he GpxFile "class"
func (gpx *GpxFile) CheckBuffer(buffer []byte) bool {
	for i, _ := range buffer {
		section := buffer[i : i+4]
		if string(section) == "<gpx" {
			return true
		}
	}
	return false
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
