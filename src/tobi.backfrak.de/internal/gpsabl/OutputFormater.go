package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"fmt"
)

// CsvOutputFormater - type that formats TrackSummary into csv style
type CsvOutputFormater struct {
	// Seperator - The seperator used to seperate values in csv
	Seperator string
}

// NewCsvOutputFormater - Get a new CsvOutputFormater
func NewCsvOutputFormater(seperator string) CsvOutputFormater {
	return CsvOutputFormater{seperator}
}

// FormatOutPut - Create the output for a TrackFile
func (formater CsvOutputFormater) FormatOutPut(trackFile TrackFile, printHeader bool) []string {
	ret := []string{}
	if printHeader {
		header := formater.GetHeader()
		ret = append(ret, header)
	}

	ret = append(ret, formater.FormatTrackSummary(TrackSummaryProvider(trackFile), trackFile.FilePath))

	return ret
}

// GetHeader - Get the header line of a csv output
func (formater CsvOutputFormater) GetHeader() string {
	ret := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s",
		"Name", formater.Seperator,
		"Distance (km)", formater.Seperator,
		"AtituteRange (m)", formater.Seperator,
		"MinimumAtitute (m)", formater.Seperator,
		"MaximumAtitut (m)", formater.Seperator)

	return ret
}

// FormatTrackSummary - Create the outputline for a  TrackSummaryProvider
func (formater CsvOutputFormater) FormatTrackSummary(info TrackSummaryProvider, name string) string {
	ret := fmt.Sprintf("%s%s%f%s%f%s%f%s%f%s",
		name, formater.Seperator,
		RoundFloat64To2Digits(info.GetDistance()/1000), formater.Seperator,
		RoundFloat64To2Digits(float64(info.GetAtituteRange())), formater.Seperator,
		RoundFloat64To2Digits(float64(info.GetMinimumAtitute())), formater.Seperator,
		RoundFloat64To2Digits(float64(info.GetMaximumAtitute())), formater.Seperator)

	return ret
}
