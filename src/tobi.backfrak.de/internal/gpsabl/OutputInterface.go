package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"os"
	"strings"
)

// DepthArg -  "Enum" Type that represents the different depth modes
type DepthArg string

// SummaryArg - "Enum" Type that represents the different summary modes
type SummaryArg string

// OutputFormaterType - a string type to implement the enum pattern
type OutputFormaterType string

const (
	// TRACK - analyse into track depth
	TRACK DepthArg = "track"
	// FILE - analyse into file depth
	FILE DepthArg = "file"
	// SEGMENT -  analyse into segment depth
	SEGMENT DepthArg = "segment"
)

const (
	// NONE - add no summary to the output
	NONE SummaryArg = "none"
	// ADDITIONAL - add the summary to the output
	ADDITIONAL SummaryArg = "additional"
	// ONLY - write only the summara as output
	ONLY SummaryArg = "only"
)

// GetValidDepthArgs - The valid args values for the depth parameter
func GetValidDepthArgs() []DepthArg {
	ret := []DepthArg{TRACK, FILE, SEGMENT}
	return ret
}

// GetValidSummaryArgs - The valid args values for the summary parameter
func GetValidSummaryArgs() []SummaryArg {
	ret := []SummaryArg{NONE, ADDITIONAL, ONLY}
	return ret
}

// GetValidDepthArgsString - Get the ValidDepthArgs in one string
func GetValidDepthArgsString() string {
	ret := ""
	for _, arg := range GetValidDepthArgs() {
		ret = fmt.Sprintf("%s %s", arg, ret)
	}
	return ret
}

// CheckValidDepthArg -Check if a string is a valid depth arg
func CheckValidDepthArg(agr string) bool {
	return strings.Contains(GetValidDepthArgsString(), agr)
}

// GetValidSummaryArgsString - Get the ValidSummaryArgs in one string
func GetValidSummaryArgsString() string {
	ret := ""
	for _, arg := range GetValidSummaryArgs() {
		ret = fmt.Sprintf("%s %s", arg, ret)
	}
	return ret
}

// CheckValidSummaryArg -Check if a string is a valid summary arg
func CheckValidSummaryArg(agr string) bool {
	return strings.Contains(GetValidSummaryArgsString(), agr)
}

// OutputFormater - Interface for classes that can format a track output into a file format and write this file
type OutputFormater interface {
	// Get a new OutputFormater of this type
	NewOutputFormater() OutputFormater

	// AddOutPut - Add the output values of a TrackFile to the out file buffer
	AddOutPut(trackFile TrackFile, depth DepthArg, filterDuplicate bool) error

	// WriteOutput - Write the output to the output file
	WriteOutput(outFile *os.File, summary SummaryArg) error

	// Check if this OutputFormater can write the given output file
	CheckFileExtension(filePath string) bool

	// Check if this OutputFormater is responsible for the given OutputFormaterType
	CheckOutputFormaterType(formaterType OutputFormaterType) bool

	// Get the list of file extensions this formater can write
	GetFileExtensions() []string

	// Get the list of OutputFormaterType this formater can write
	GetOutputFormaterTypes() []OutputFormaterType
}
