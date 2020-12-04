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
	"time"
)

// NotValidValue - The value set when values are not valid
const NotValidValue = "not valid"

// OutputLine - Represents one line in the output
type OutputLine struct {
	Name string
	Data TrackSummaryProvider
}

func newOutputLine(name string, data TrackSummaryProvider) *OutputLine {
	ret := OutputLine{}
	ret.Data = data
	ret.Name = name

	return &ret
}

// CsvOutputFormater - type that formats TrackSummary into csv style
type CsvOutputFormater struct {
	// Separator - The separator used to separate values in csv
	Separator string

	// Tell if the CSV header should be added to the output
	AddHeader bool

	// ValidDepthArgs - The valid args values for the -depth parameter
	ValidDepthArgs []string

	lineBuffer []OutputLine
	mux        sync.Mutex
}

// NewCsvOutputFormater - Get a new CsvOutputFormater
func NewCsvOutputFormater(separator string, addHeader bool) *CsvOutputFormater {
	ret := CsvOutputFormater{}
	ret.Separator = separator
	ret.AddHeader = addHeader
	ret.ValidDepthArgs = []string{"track", "file", "segment"}
	ret.lineBuffer = []OutputLine{}

	return &ret
}

// AddOutPut - Add the formated output of a TrackFile to the internal buffer, so it can be written out later
func (formater *CsvOutputFormater) AddOutPut(trackFile TrackFile, depth string, filterDuplicate bool) error {

	var lines []OutputLine
	linesFromFile, err := formater.GetOutPutEntries(trackFile, false, depth)
	if err != nil {
		return err
	}
	if filterDuplicate {
		for _, line := range linesFromFile {
			if outPutContainsLineByTimeStamps(lines, line) == false && outPutContainsLineByTimeStamps(formater.lineBuffer, line) == false {
				lines = append(lines, line)
			}
		}
	} else {
		lines = linesFromFile
	}

	if len(lines) > 0 {
		formater.mux.Lock()
		defer formater.mux.Unlock()
		formater.lineBuffer = append(formater.lineBuffer, lines...)
	}

	return nil
}

// GetLines - Get the lines stored in the internal buffer
func (formater *CsvOutputFormater) GetLines() []string {
	ret := []string{}
	if formater.AddHeader {
		ret = append(ret, formater.GetHeader())
	}
	formater.mux.Lock()
	defer formater.mux.Unlock()
	for _, line := range formater.lineBuffer {
		ret = append(ret, formater.FormatTrackSummary(line.Data, line.Name))
	}
	return ret
}

// WriteOutput - Write the output to a given file handle object. Make sure the file exists before you call this method!
func (formater *CsvOutputFormater) WriteOutput(outFile *os.File) error {
	lines := formater.GetLines()
	for _, line := range lines {
		_, errWrite := outFile.WriteString(line)
		if errWrite != nil {
			return errWrite
		}
	}

	return nil
}

// GetOutPutEntries - Add the output of a TrackFile
func (formater *CsvOutputFormater) GetOutPutEntries(trackFile TrackFile, printHeader bool, depth string) ([]OutputLine, error) {
	ret := []OutputLine{}

	switch depth {
	case formater.ValidDepthArgs[1]:
		ret = append(ret, *newOutputLine(getLineNameFromTrackFile(trackFile), TrackSummaryProvider(trackFile)))
	case formater.ValidDepthArgs[0]:
		ret = append(ret, getLinesFromTracks(formater, trackFile)...)
	case formater.ValidDepthArgs[2]:
		ret = append(ret, getLinesFromTrackSegments(formater, trackFile)...)
	default:
		return nil, NewDepthParameterNotKnownError(depth)
	}

	return ret, nil
}

// GetHeader - Get the header line of a csv output
func (formater *CsvOutputFormater) GetHeader() string {
	ret := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s",
		"Name", formater.Separator,
		"StartTime", formater.Separator,
		"EndTime", formater.Separator,
		"TrackTime (xxhxxmxxs)", formater.Separator,
		"Distance (km)", formater.Separator,
		"HorizontalDistance (km)", formater.Separator,
		"AltitudeRange (m)", formater.Separator,
		"MinimumAltitude (m)", formater.Separator,
		"MaximumAltitude (m)", formater.Separator,
		"ElevationGain (m)", formater.Separator,
		"ElevationLose (m)", formater.Separator,
		"UpwardsDistance (km)", formater.Separator,
		"DownwardsDistance (km)", formater.Separator,
		"MovingTime (xxhxxmxxs)", formater.Separator,
		"UpwardsTime (xxhxxmxxs)", formater.Separator,
		"DownwardsTime (xxhxxmxxs)", formater.Separator,
		"AverageSpeed (km/h)", formater.Separator,
		"UpwardsSpeed (km/h)", formater.Separator,
		"DownwardsSpeed (km/h)", formater.Separator,
		GetNewLine(),
	)

	return ret
}

// FormatTrackSummary - Create the OutputLine for a TrackSummaryProvider
func (formater *CsvOutputFormater) FormatTrackSummary(info TrackSummaryProvider, name string) string {
	var ret string
	if info.GetTimeDataValid() {
		ret = fmt.Sprintf("%s%s%s%s%s%s%s%s%f%s%f%s%f%s%f%s%f%s%f%s%f%s%f%s%f%s%s%s%s%s%s%s%f%s%f%s%f%s%s",
			name, formater.Separator,
			info.GetStartTime().Format(time.RFC3339), formater.Separator,
			info.GetEndTime().Format(time.RFC3339), formater.Separator,
			info.GetEndTime().Sub(info.GetStartTime()).String(), formater.Separator,
			RoundFloat64To2Digits(info.GetDistance()/1000), formater.Separator,
			RoundFloat64To2Digits(info.GetHorizontalDistance()/1000), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetAltitudeRange())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetMinimumAltitude())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetMaximumAltitude())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetElevationGain())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetElevationLose())), formater.Separator,
			RoundFloat64To2Digits(info.GetUpwardsDistance()/1000), formater.Separator,
			RoundFloat64To2Digits(info.GetDownwardsDistance()/1000), formater.Separator,
			info.GetMovingTime().String(), formater.Separator,
			info.GetUpwardsTime().String(), formater.Separator,
			info.GetDownwardsTime().String(), formater.Separator,
			RoundFloat64To2Digits(info.GetAvarageSpeed()*3.6), formater.Separator,
			RoundFloat64To2Digits(info.GetUpwardsSpeed()*3.6), formater.Separator,
			RoundFloat64To2Digits(info.GetDownwardsSpeed()*3.6), formater.Separator,
			GetNewLine(),
		)
	} else {
		ret = fmt.Sprintf("%s%s%s%s%s%s%s%s%f%s%f%s%f%s%f%s%f%s%f%s%f%s%f%s%f%s%s%s%s%s%s%s%s%s%s%s%s%s%s",
			name, formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			RoundFloat64To2Digits(info.GetDistance()/1000), formater.Separator,
			RoundFloat64To2Digits(info.GetHorizontalDistance()/1000), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetAltitudeRange())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetMinimumAltitude())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetMaximumAltitude())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetElevationGain())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetElevationLose())), formater.Separator,
			RoundFloat64To2Digits(info.GetUpwardsDistance()/1000), formater.Separator,
			RoundFloat64To2Digits(info.GetDownwardsDistance()/1000), formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			GetNewLine(),
		)
	}

	return ret
}

// GetValidDepthArgsString - Get the ValidDepthArgs in one string
func (formater *CsvOutputFormater) GetValidDepthArgsString() string {
	ret := ""
	for _, arg := range formater.ValidDepthArgs {
		ret = fmt.Sprintf("%s %s", arg, ret)
	}
	return ret
}

// CheckValidDepthArg -Check if a string is a valid depth arg
func (formater *CsvOutputFormater) CheckValidDepthArg(agr string) bool {
	return strings.Contains(formater.GetValidDepthArgsString(), agr)
}

// GetNewLine - Get the new line string depending on the OS
func GetNewLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"

}

func outPutContainsLineByTimeStamps(output []OutputLine, newLine OutputLine) bool {

	// Don't tread all lines with no valid time values as duplicates
	if newLine.Data.GetTimeDataValid() == false {
		return false
	}
	newLineStartTime := newLine.Data.GetStartTime()
	newLineEndTime := newLine.Data.GetEndTime()

	for _, outLine := range output {
		outLineStartTime := outLine.Data.GetStartTime()
		outLineEndTime := outLine.Data.GetEndTime()

		if outLineStartTime == newLineStartTime && outLineEndTime == newLineEndTime {
			return true
		}
	}

	return false
}

func getLinesFromTrackSegments(formater *CsvOutputFormater, trackFile TrackFile) []OutputLine {
	ret := []OutputLine{}
	for iTrack, track := range trackFile.Tracks {
		for iSeg, seg := range track.TrackSegments {
			info := TrackSummaryProvider(seg)
			name := fmt.Sprintf("%s: Segment #%d", getLineNameFromTrack(track, trackFile, iTrack), iSeg+1)
			entry := newOutputLine(name, info)
			ret = append(ret, *entry)
		}
	}

	return ret
}

func getLinesFromTracks(formater *CsvOutputFormater, trackFile TrackFile) []OutputLine {
	ret := []OutputLine{}
	for i, track := range trackFile.Tracks {
		info := TrackSummaryProvider(track)
		name := getLineNameFromTrack(track, trackFile, i)
		entry := newOutputLine(name, info)
		ret = append(ret, *entry)
	}

	return ret
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
