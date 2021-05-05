package gpsabl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// GetInputFileFromBuffer - Get a new InputFile if there is a reader that supports the data in the buffer
// - validReaders    List of valid TrackReader
// - buffer          The buffer that contains data
// - bufferName      Name of the buffer
// return true and the coresponding InputFile if a valid reader was found
// return false and an empty InputFile if no valid reader was found
func GetInputFileFromBuffer(validReaders []TrackReader, buffer []byte, bufferName string) (bool, InputFile) {
	var retVal InputFile

	for _, reader := range validReaders {
		if reader.CheckBuffer(buffer) == true {
			retVal = *reader.NewInputFileForBuffer(buffer, bufferName)
			return true, retVal
		}
	}

	return false, retVal
}

// GetInputFileFromPath - Get a new InputFile if there is a reader that supports the file extencion of the given file
// - validReaders    List of valid TrackReader
// - path            The path to the data file
// return true and the coresponding InputFile if a valid reader was found
// return false and an empty InputFile if no valid reader was found
func GetInputFileFromPath(validReaders []TrackReader, path string) (bool, InputFile) {
	var retVal InputFile

	for _, reader := range validReaders {
		if reader.CheckFile(path) == true {
			retVal = *NewInputFileWithPath(path)
			return true, retVal
		}
	}

	return false, retVal
}

// GetNewReader - Get a new reader for a given input
// return true and the coresponding TrackReader if a valid reader was found
// return false and nil if no valid reader was found
func GetNewReader(validReaders []TrackReader, input InputFile) (bool, TrackReader) {
	for _, reader := range validReaders {
		if reader.CheckInputFile(input) == true {
			return true, reader.NewReader(input)
		}
	}

	return false, nil
}
