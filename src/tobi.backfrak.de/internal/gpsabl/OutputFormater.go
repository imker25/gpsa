package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import "fmt"

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
		header := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s",
			"Name", formater.Seperator, "Distance", formater.Seperator, "AtituteRange", formater.Seperator,
			"MinimumAtitute", formater.Seperator, "MaximumAtitut", formater.Seperator)
		ret = append(ret, header)
	}

	ret = append(ret, formater.formatTrackSummary(TrackSummaryProvider(trackFile), trackFile.FilePath))

	return ret
}

func (formater CsvOutputFormater) formatTrackSummary(info TrackSummaryProvider, name string) string {
	ret := fmt.Sprintf("%s%s%f%s%f%s%f%s%f%s",
		name, formater.Seperator, info.GetDistance(), formater.Seperator, info.GetAtituteRange(), formater.Seperator,
		info.GetMinimumAtitute(), formater.Seperator, info.GetMaximumAtitute(), formater.Seperator)

	return ret
}
