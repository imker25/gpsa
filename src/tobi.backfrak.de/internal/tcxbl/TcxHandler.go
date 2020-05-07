package tcxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import "tobi.backfrak.de/internal/gpsabl"

// TcxFile - The struct to handle *.gpx data files
type TcxFile struct {
	gpsabl.TrackFile
}

// NewTcxFile - Constructor for the GpxFile struct
func NewTcxFile(filePath string) TcxFile {
	gpx := TcxFile{}
	gpx.FilePath = filePath

	return gpx
}

// ReadTracks - Read the *.tcx from the inputs GpxFile.FilePath, and return a GpxFile struct that contains all information
// Implement the gpsabl.TrackReader interface for *.tcx files
func (gpx *TcxFile) ReadTracks(correction string, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {
	ret, err := ReadTcxFile(gpx.FilePath, correction, minimalMovingSpeed, minimalStepHight)

	if err == nil {
		gpx.TrackFile = ret
	}

	return ret, err
}

// ReadTcxFile - Reads a *.gpx file
func ReadTcxFile(filePath string, correction string, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {
	ret := gpsabl.TrackFile{}
	tcx, fileError := ReadTcx(filePath)

	if fileError != nil {
		return gpsabl.TrackFile{}, fileError
	}

	ret, convertError := ConvertTcx(tcx, filePath, correction, minimalMovingSpeed, minimalStepHight)
	if convertError != nil {
		return gpsabl.TrackFile{}, convertError
	}

	return ret, nil
}
