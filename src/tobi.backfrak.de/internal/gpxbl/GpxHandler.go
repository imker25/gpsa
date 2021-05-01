package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"strings"

	"tobi.backfrak.de/internal/gpsabl"
)

const GpxBuffer gpsabl.InputFileType = "GpxBuffer"

// GpxFile - The struct to handle *.gpx data files
type GpxFile struct {
	gpsabl.TrackFile
	input gpsabl.InputFile
}

// NewGpxFile - Constructor for the GpxFile struct
func NewGpxFile(filePath string) GpxFile {
	gpx := GpxFile{}
	gpx.FilePath = filePath
	gpx.input = *gpsabl.NewInputFileWithPath(filePath)

	return gpx
}

// NewReader - Get a new reader for GPX files that will read the data in the given gpsabl.InputFile
// Implement the gpsabl.TrackReader interface for *.gpx files
func (gpx *GpxFile) NewReader(data gpsabl.InputFile) gpsabl.TrackReader {
	newGpx := GpxFile{}
	newGpx.input = data
	if data.Type == gpsabl.FilePath {
		newGpx.FilePath = data.Name
	}

	return &newGpx
}

// ReadTracks - Read the *.gpx from the given FilePath or Buffer, and return a gpsabl.TrackFile struct that contains all information
// Implement the gpsabl.TrackReader interface for *.gpx files
func (gpx *GpxFile) ReadTracks(correction gpsabl.CorrectionParameter, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {
	var err error
	var ret gpsabl.TrackFile
	if gpx.input.Type == gpsabl.FilePath {
		ret, err = ReadGpxFile(gpx.FilePath, correction, minimalMovingSpeed, minimalStepHight)
	} else if gpx.input.Type == GpxBuffer {
		ret, err = ReadBuffer(gpx.input.Buffer, gpx.input.Name, correction, minimalMovingSpeed, minimalStepHight)
	} else {
		err = gpsabl.NewUnKnownInputFileTypeError(gpx.input.Name)
	}
	if err == nil {
		gpx.TrackFile = ret
	}

	return ret, err
}

// ReadBuffer - Read the *.gpx data from a buffer, and return a gpsabl.TrackFile struct that contains all information
// When using this method, the FilePath property may contain any string
func ReadBuffer(buffer []byte, name string, correction gpsabl.CorrectionParameter, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {
	content, readErr := readGPXBuffer(buffer, name)
	if readErr != nil {
		return gpsabl.TrackFile{}, readErr
	}
	ret, convertError := ConvertGPXFile(content, name, correction, minimalMovingSpeed, minimalStepHight)
	if convertError != nil {
		return gpsabl.TrackFile{}, convertError
	}

	return ret, nil
}

// CheckFile - Check if a file can be read by the GpxFile "class"
// Implement the gpsabl.TrackReader interface for *.gpx files
func (gpx *GpxFile) CheckFile(path string) bool {
	if strings.HasSuffix(path, "gpx") == true { // If the file is a *.gpx, we can read it
		return true
	}

	return false
}

// CheckBuffer - Check if a buffer can be read by he GpxFile "class"
// Implement the gpsabl.TrackReader interface for *.gpx files
func (gpx *GpxFile) CheckBuffer(buffer []byte) bool {
	for i, _ := range buffer {
		section := buffer[i : i+4]
		if string(section) == "<gpx" {
			return true
		}
	}
	return false
}

// NewInputFileForBuffer - Get a new gpsabl.InputFile for a buffer containing a gpx files content
// Implement the gpsabl.TrackReader interface for *.gpx files
func (gpx *GpxFile) NewInputFileForBuffer(buffer []byte, name string) *gpsabl.InputFile {
	file := gpsabl.InputFile{}
	file.Name = name
	file.Type = GpxBuffer
	file.Buffer = buffer

	return &file
}

// CheckInputFile - Check if a gpsabl.InputFile can be handled by the GpxFile reader
// Implement the gpsabl.TrackReader interface for *.gpx  files
func (gpx *GpxFile) CheckInputFile(input gpsabl.InputFile) bool {
	if input.Type == GpxBuffer {
		return true
	} else if input.Type == gpsabl.FilePath && gpx.CheckFile(input.Name) {
		return true
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
