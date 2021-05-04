package tcxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"strings"

	"tobi.backfrak.de/internal/gpsabl"
)

const TcxBuffer gpsabl.InputFileType = "TcxBuffer"

// The file extension this Reader can read
const FileExtension string = ".tcx"

// TcxFile - The struct to handle *.gpx data files
type TcxFile struct {
	gpsabl.TrackFile
	input gpsabl.InputFile
}

// NewTcxFile - Constructor for the TcxFile struct
func NewTcxFile(filePath string) TcxFile {
	tcx := TcxFile{}
	tcx.FilePath = filePath
	tcx.input = *gpsabl.NewInputFileWithPath(filePath)

	return tcx
}

// NewReader - Get a new reader for TCX files that will read the data in the given gpsabl.InputFile
// Implement the gpsabl.TrackReader interface for *.tcx files
func (tcx *TcxFile) NewReader(data gpsabl.InputFile) gpsabl.TrackReader {
	newTcx := TcxFile{}
	newTcx.input = data
	if data.Type == gpsabl.FilePath {
		newTcx.FilePath = data.Name
	}

	return &newTcx
}

// ReadTracks - Read the *.tcx from the input from FilePath or InputFile.Buffer, and return a gpsabl.TrackFile struct that contains all information
// Implement the gpsabl.TrackReader interface for *.tcx files
func (tcx *TcxFile) ReadTracks(correction gpsabl.CorrectionParameter, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {
	var err error
	var ret gpsabl.TrackFile
	if tcx.input.Type == gpsabl.FilePath {
		ret, err = ReadTcxFile(tcx.FilePath, correction, minimalMovingSpeed, minimalStepHight)
	} else if tcx.input.Type == TcxBuffer {
		ret, err = ReadBuffer(tcx.input.Buffer, tcx.input.Name, correction, minimalMovingSpeed, minimalStepHight)
	} else {
		err = gpsabl.NewUnKnownInputFileTypeError(tcx.input.Name)
	}

	if err == nil {
		tcx.TrackFile = ret
	}

	return ret, err
}

// CheckInputFile - Check if a gpsabl.InputFile can be handled by the TcxFile reader
// Implement the gpsabl.TrackReader interface for *.tcx files
func (tcx *TcxFile) CheckInputFile(input gpsabl.InputFile) bool {
	if input.Type == TcxBuffer {
		return true
	} else if input.Type == gpsabl.FilePath && tcx.CheckFile(input.Name) {
		return true
	}

	return false
}

// ReadBuffer - Read the tcx data from a buffer, and return a gpsabl.TrackFile struct that contains all information
func ReadBuffer(buffer []byte, name string, correction gpsabl.CorrectionParameter, minimalMovingSpeed float64, minimalStepHight float64) (gpsabl.TrackFile, error) {
	content, readErr := readTCXBuffer(buffer, name)
	if readErr != nil {
		return gpsabl.TrackFile{}, readErr
	}
	ret, convertError := ConvertTcx(content, name, correction, minimalMovingSpeed, minimalStepHight)
	if convertError != nil {
		return gpsabl.TrackFile{}, convertError
	}

	return ret, nil
}

// CheckFile - Check if a file can be read by the TcxFile "class"
// Implement the gpsabl.TrackReader interface for *.tcx files
func (gpx *TcxFile) CheckFile(path string) bool {
	if strings.HasSuffix(strings.ToLower(path), FileExtension) == true { // If the file is a *.gpx, we can read it
		return true
	}

	return false
}

// CheckBuffer - Check if a buffer can be read by he TcxFile "class"
// Implement the gpsabl.TrackReader interface for *.tcx files
func (gpx *TcxFile) CheckBuffer(buffer []byte) bool {
	for i, _ := range buffer {
		section := buffer[i : i+23]
		if string(section) == "<TrainingCenterDatabase" {
			return true
		}
	}
	return false
}

// NewInputFileForBuffer - Get a new InputFile for a buffer containing a tcx files content
// Implement the gpsabl.TrackReader interface for *.tcx files
func (gpx *TcxFile) NewInputFileForBuffer(buffer []byte, name string) *gpsabl.InputFile {
	file := gpsabl.InputFile{}
	file.Name = name
	file.Type = TcxBuffer
	file.Buffer = buffer

	return &file
}

// GetValidFileExtensions - Get a list of file extensions this reader can read
// Implement the gpsabl.TrackReader interface for *.tcx files
func (gpx *TcxFile) GetValidFileExtensions() []string {

	extensions := []string{FileExtension}

	return extensions
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
