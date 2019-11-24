package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

// DepthParametrNotKnownError - Error when the given depth paramter is not known
type DepthParametrNotKnownError struct {
	err string
	// File - The path to the dir that caused this error
	GivenValue string
}

func (e *DepthParametrNotKnownError) Error() string { // Implement the Error Interface for the DepthParametrNotKnownError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// NewDepthParametrNotKnownError - Get a new DepthParametrNotKnownError struct
func NewDepthParametrNotKnownError(givenValue string) *DepthParametrNotKnownError {
	return &DepthParametrNotKnownError{fmt.Sprintf("The given -depth \"%s\" is not known.", givenValue), givenValue}
}

// CsvOutputFormater - type that formats TrackSummary into csv style
type CsvOutputFormater struct {
	// Seperator - The seperator used to seperate values in csv
	Seperator string

	// ValideDepthArgs - The valide args values for the -depth paramter
	ValideDepthArgs []string

	lineBuffer []string
	mux        sync.Mutex
}

// NewCsvOutputFormater - Get a new CsvOutputFormater
func NewCsvOutputFormater(seperator string) *CsvOutputFormater {
	ret := CsvOutputFormater{}
	ret.Seperator = seperator
	ret.ValideDepthArgs = []string{"track", "file", "segment"}
	ret.lineBuffer = []string{}

	return &ret
}

// AddOutPut - Add the formated output of a TrackFile to the internal buffer, so it can be written out later
func (formater *CsvOutputFormater) AddOutPut(trackFile TrackFile, depth string) error {
	formater.mux.Lock()
	defer formater.mux.Unlock()
	lines, err := formater.FormatOutPut(trackFile, false, depth)
	if err != nil {
		return err
	}
	formater.lineBuffer = append(formater.lineBuffer, lines...)

	return nil
}

// AddHeader - Add the formated header line to the internal buffer, so it can be written out later
func (formater *CsvOutputFormater) AddHeader() {
	formater.mux.Lock()
	defer formater.mux.Unlock()
	formater.lineBuffer = append(formater.lineBuffer, formater.GetHeader())
}

// GetLines - Get the lines stored in the internal buffer
func (formater *CsvOutputFormater) GetLines() []string {
	return formater.lineBuffer
}

// WriteOutput - Write the output to a given file handle object. Make sure the file exists before you call this method!
func (formater *CsvOutputFormater) WriteOutput(outFile *os.File) error {
	formater.mux.Lock()
	defer formater.mux.Unlock()
	for _, line := range formater.lineBuffer {
		_, errWrite := outFile.WriteString(line)
		if errWrite != nil {
			return errWrite
		}
	}

	return nil
}

// FormatOutPut - Create the output for a TrackFile
func (formater *CsvOutputFormater) FormatOutPut(trackFile TrackFile, printHeader bool, depth string) ([]string, error) {
	ret := []string{}
	if printHeader {
		header := formater.GetHeader()
		ret = append(ret, header)
	}

	switch depth {
	case formater.ValideDepthArgs[1]:
		ret = append(ret, formater.FormatTrackSummary(TrackSummaryProvider(trackFile), getLineNameFromTrackFile(trackFile)))
	case formater.ValideDepthArgs[0]:
		addLinesFromTracks(formater, trackFile, &ret)
	case formater.ValideDepthArgs[2]:
		addLinesFromTrackSegments(formater, trackFile, &ret)
	default:
		return ret, NewDepthParametrNotKnownError(depth)
	}

	return ret, nil
}

// GetHeader - Get the header line of a csv output
func (formater *CsvOutputFormater) GetHeader() string {
	ret := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s",
		"Name", formater.Seperator,
		"Distance (km)", formater.Seperator,
		"AtituteRange (m)", formater.Seperator,
		"MinimumAtitute (m)", formater.Seperator,
		"MaximumAtitut (m)", formater.Seperator, GetNewLine())

	return ret
}

// FormatTrackSummary - Create the outputline for a  TrackSummaryProvider
func (formater *CsvOutputFormater) FormatTrackSummary(info TrackSummaryProvider, name string) string {
	ret := fmt.Sprintf("%s%s%f%s%f%s%f%s%f%s%s",
		name, formater.Seperator,
		RoundFloat64To2Digits(info.GetDistance()/1000), formater.Seperator,
		RoundFloat64To2Digits(float64(info.GetAtituteRange())), formater.Seperator,
		RoundFloat64To2Digits(float64(info.GetMinimumAtitute())), formater.Seperator,
		RoundFloat64To2Digits(float64(info.GetMaximumAtitute())), formater.Seperator, GetNewLine())

	return ret
}

// GetVlaideDepthArgsString - Get the VlaideDepthArgs in one string
func (formater *CsvOutputFormater) GetVlaideDepthArgsString() string {
	ret := ""
	for _, arg := range formater.ValideDepthArgs {
		ret = fmt.Sprintf("%s %s", arg, ret)
	}
	return ret
}

// CheckVlaideDepthArg -Check if a string is a valide depth arg
func (formater *CsvOutputFormater) CheckVlaideDepthArg(agr string) bool {
	return strings.Contains(formater.GetVlaideDepthArgsString(), agr)
}

// GetNewLine - Get the new line string depending on the OS
func GetNewLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"

}

func addLinesFromTrackSegments(formater *CsvOutputFormater, trackFile TrackFile, lines *[]string) {
	for iTrack, track := range trackFile.Tracks {
		for iSeg, seg := range track.TrackSegments {
			info := TrackSummaryProvider(seg)
			name := fmt.Sprintf("%s: Segment #%d", getLineNameFromTrack(track, trackFile, iTrack), iSeg+1)
			ret := formater.FormatTrackSummary(info, name)
			*lines = append(*lines, ret)
		}
	}
}

func addLinesFromTracks(formater *CsvOutputFormater, trackFile TrackFile, lines *[]string) {

	for i, track := range trackFile.Tracks {
		info := TrackSummaryProvider(track)
		name := getLineNameFromTrack(track, trackFile, i)
		ret := formater.FormatTrackSummary(info, name)
		*lines = append(*lines, ret)
	}

}

func getLineNameFromTrack(track Track, parent TrackFile, index int) string {
	if track.Name != "" {
		return fmt.Sprintf("%s: %s", getLineNameFromTrackFile(parent), track.Name)
	}

	return fmt.Sprintf("%s: Track #%d", getLineNameFromTrackFile(parent), index+1)

}

func getLineNameFromTrackFile(trackFile TrackFile) string {
	if trackFile.Name != "" {
		return trackFile.Name
	}

	return trackFile.FilePath
}
