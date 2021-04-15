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

// NewTcxFile - Constructor for the TcxFile struct
func NewTcxFile(filePath string) TcxFile {
	tcx := TcxFile{}
	tcx.FilePath = filePath

	return tcx
}

// ReadTracks - Read the *.tcx from the inputs TcxFile.FilePath, and return a GpxFile struct that contains all information
// When using this method, the FilePath property may contain any string
// Implement the gpsabl.TrackReader interface for *.tcx files
func (tcx *TcxFile) ReadTracks(correction gpsabl.CorrectionParameter, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {
	ret, err := ReadTcxFile(tcx.FilePath, correction, minimalMovingSpeed, minimalStepHight)

	if err == nil {
		tcx.TrackFile = ret
	}

	return ret, err
}

// ReadBuffer - Read the tcx data from a buffer, and return a gpsabl.TrackFile struct that contains all information
// Implement the gpsabl.TrackReader interface for *.tcx files
func (tcx *TcxFile) ReadBuffer(buffer []byte, correction gpsabl.CorrectionParameter, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {
	content, readErr := readTCXBuffer(buffer, tcx.FilePath)
	if readErr != nil {
		return gpsabl.TrackFile{}, readErr
	}
	ret, convertError := ConvertTcx(content, tcx.FilePath, correction, minimalMovingSpeed, minimalStepHight)
	if convertError != nil {
		return gpsabl.TrackFile{}, convertError
	}

	return ret, nil
}

// ReadTcxFile - Reads a *.gpx file
func ReadTcxFile(filePath string, correction gpsabl.CorrectionParameter, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {
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
