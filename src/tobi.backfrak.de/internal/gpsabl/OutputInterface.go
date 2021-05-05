package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// DepthArg -  "Enum" Type that represents the different depth modes
type DepthArg string

// SummaryArg - "Enum" Type that represents the different summary modes
type SummaryArg string

// OutputFormaterType - a string type to implement the enum pattern
type OutputFormaterType string

// TimeFormat - Represents a go Time format string for the enum pattern
type TimeFormat string

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

const (
	// RFC3339 - Internal representation of gos time.RFC3339
	RFC3339 TimeFormat = time.RFC3339

	// RFC850 -  Internal representation of gos time.RFC850
	RFC850 TimeFormat = time.RFC850

	// UnixDate -  Internal representation of gos time.UnixDate
	UnixDate TimeFormat = time.UnixDate
)

// GetValidTimeFormats -  Get the valid TimeFormat values
func GetValidTimeFormats() []TimeFormat {
	return []TimeFormat{RFC3339, RFC850, UnixDate}
}

// GetValidTimeFormatsString - Get a string that contains all valid TimeFormat values
func GetValidTimeFormatsString() string {
	ret := ""
	for _, arg := range GetValidTimeFormats() {
		ret = fmt.Sprintf("\"%s\" %s", arg, ret)
	}
	return ret
}

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

	// Get the TextOutputFormater version of this formater or nil if the formater is not a TextOutputFormater
	GetTextOutputFormater() TextOutputFormater
}

// TextOutputFormater - Interface for classes that can format a track output into a text style file format like csv
type TextOutputFormater interface {
	OutputFormater

	// Set the time format string used by this CsvOutputFormater. Will return an error if you want to set an unknown format
	SetTimeFormat(timeFormat string) error

	// Set the value of formater.AddHeader
	SetAddHeader(value bool)

	// Set the value of formater.Separator
	SetSeperator(value string)

	// Check if the given format string is a valid TimeFormat
	CheckTimeFormatIsValid(format string) bool
}
