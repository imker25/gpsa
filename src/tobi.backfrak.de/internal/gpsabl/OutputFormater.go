package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"fmt"
	"runtime"
)

// CsvOutputFormater - type that formats TrackSummary into csv style
type CsvOutputFormater struct {
	// Seperator - The seperator used to seperate values in csv
	Seperator string

	// ValideDepthArgs - The valide args values for the -depth paramter
	ValideDepthArgs []string
}

// NewCsvOutputFormater - Get a new CsvOutputFormater
func NewCsvOutputFormater(seperator string) CsvOutputFormater {
	ret := CsvOutputFormater{}
	ret.Seperator = seperator
	ret.ValideDepthArgs = []string{"track", "file", "segment"}

	return ret
}

// FormatOutPut - Create the output for a TrackFile
func (formater CsvOutputFormater) FormatOutPut(trackFile TrackFile, printHeader bool, depth string) []string {
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
		ret = append(ret, fmt.Sprintf("Error: Can not handle the given depth value \"%s\"%s", depth, GetNewLine()))
	}

	return ret
}

// GetHeader - Get the header line of a csv output
func (formater CsvOutputFormater) GetHeader() string {
	ret := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s",
		"Name", formater.Seperator,
		"Distance (km)", formater.Seperator,
		"AtituteRange (m)", formater.Seperator,
		"MinimumAtitute (m)", formater.Seperator,
		"MaximumAtitut (m)", formater.Seperator, GetNewLine())

	return ret
}

// GetVlaideDepthArgsString - Get the VlaideDepthArgs in one string
func (formater CsvOutputFormater) GetVlaideDepthArgsString() string {
	ret := ""
	for _, arg := range formater.ValideDepthArgs {
		ret = fmt.Sprintf("%s %s", arg, ret)
	}
	return ret
}

// FormatTrackSummary - Create the outputline for a  TrackSummaryProvider
func (formater CsvOutputFormater) FormatTrackSummary(info TrackSummaryProvider, name string) string {
	ret := fmt.Sprintf("%s%s%f%s%f%s%f%s%f%s%s",
		name, formater.Seperator,
		RoundFloat64To2Digits(info.GetDistance()/1000), formater.Seperator,
		RoundFloat64To2Digits(float64(info.GetAtituteRange())), formater.Seperator,
		RoundFloat64To2Digits(float64(info.GetMinimumAtitute())), formater.Seperator,
		RoundFloat64To2Digits(float64(info.GetMaximumAtitute())), formater.Seperator, GetNewLine())

	return ret
}

// GetNewLine - Get the new line string depending on the OS
func GetNewLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"

}

func addLinesFromTrackSegments(formater CsvOutputFormater, trackFile TrackFile, lines *[]string) {
	for iTrack, track := range trackFile.Tracks {
		for iSeg, seg := range track.TrackSegments {
			info := TrackSummaryProvider(seg)
			name := fmt.Sprintf("%s: Segment #%d", getLineNameFromTrack(track, trackFile, iTrack), iSeg+1)
			ret := formater.FormatTrackSummary(info, name)
			*lines = append(*lines, ret)
		}
	}
}

func addLinesFromTracks(formater CsvOutputFormater, trackFile TrackFile, lines *[]string) {

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
