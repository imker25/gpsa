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

// CsvOutputFormater - type that formats TrackSummary into csv style
type CsvOutputFormater struct {
	// Separator - The separator used to separate values in csv
	Separator string

	// ValidDepthArgs - The valid args values for the -depth parameter
	ValidDepthArgs []string

	lineBuffer []string
	mux        sync.Mutex
}

// NewCsvOutputFormater - Get a new CsvOutputFormater
func NewCsvOutputFormater(separator string) *CsvOutputFormater {
	ret := CsvOutputFormater{}
	ret.Separator = separator
	ret.ValidDepthArgs = []string{"track", "file", "segment"}
	ret.lineBuffer = []string{}

	return &ret
}

// AddOutPut - Add the formated output of a TrackFile to the internal buffer, so it can be written out later
func (formater *CsvOutputFormater) AddOutPut(trackFile TrackFile, depth string, filterDuplicate bool) error {

	var lines []string
	linesFromFile, err := formater.FormatOutPut(trackFile, false, depth)
	if err != nil {
		return err
	}
	if filterDuplicate {
		for _, line := range linesFromFile {
			if outPutContainsLineByTimeStamps(lines, line) == false && outPutContainsLineByTimeStamps(formater.GetLines(), line) == false {
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
	case formater.ValidDepthArgs[1]:
		ret = append(ret, formater.FormatTrackSummary(TrackSummaryProvider(trackFile), getLineNameFromTrackFile(trackFile)))
	case formater.ValidDepthArgs[0]:
		addLinesFromTracks(formater, trackFile, &ret)
	case formater.ValidDepthArgs[2]:
		addLinesFromTrackSegments(formater, trackFile, &ret)
	default:
		return ret, NewDepthParameterNotKnownError(depth)
	}

	return ret, nil
}

// GetHeader - Get the header line of a csv output
func (formater *CsvOutputFormater) GetHeader() string {
	ret := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s",
		"Name", formater.Separator,
		"StartTime", formater.Separator,
		"EndTime", formater.Separator,
		"Distance (km)", formater.Separator,
		"AltitudeRange (m)", formater.Separator,
		"MinimumAltitude (m)", formater.Separator,
		"MaximumAltitude (m)", formater.Separator,
		"ElevationGain (m)", formater.Separator,
		"ElevationLose (m)", formater.Separator,
		"UpwardsDistance (km)", formater.Separator,
		"DownwardsDistance (km)", formater.Separator,
		"MovingTime (xxhxxmxxs)", formater.Separator,
		"AverageSpeed (km/h)", formater.Separator,
		GetNewLine(),
	)

	return ret
}

// FormatTrackSummary - Create the outputline for a TrackSummaryProvider
func (formater *CsvOutputFormater) FormatTrackSummary(info TrackSummaryProvider, name string) string {
	var ret string
	if info.GetTimeDataValid() {
		ret = fmt.Sprintf("%s%s%s%s%s%s%f%s%f%s%f%s%f%s%f%s%f%s%f%s%f%s%s%s%f%s%s",
			name, formater.Separator,
			info.GetStartTime().Format(time.RFC3339), formater.Separator,
			info.GetEndTime().Format(time.RFC3339), formater.Separator,
			RoundFloat64To2Digits(info.GetDistance()/1000), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetAltitudeRange())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetMinimumAltitude())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetMaximumAltitude())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetElevationGain())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetElevationLose())), formater.Separator,
			RoundFloat64To2Digits(info.GetUpwardsDistance()/1000), formater.Separator,
			RoundFloat64To2Digits(info.GetDownwardsDistance()/1000), formater.Separator,
			info.GetMovingTime().String(), formater.Separator,
			RoundFloat64To2Digits(info.GetAvarageSpeed()*3.6), formater.Separator,
			GetNewLine(),
		)
	} else {
		ret = fmt.Sprintf("%s%s%s%s%s%s%f%s%f%s%f%s%f%s%f%s%f%s%f%s%f%s%s%s%s%s%s",
			name, formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			RoundFloat64To2Digits(info.GetDistance()/1000), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetAltitudeRange())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetMinimumAltitude())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetMaximumAltitude())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetElevationGain())), formater.Separator,
			RoundFloat64To2Digits(float64(info.GetElevationLose())), formater.Separator,
			RoundFloat64To2Digits(info.GetUpwardsDistance()/1000), formater.Separator,
			RoundFloat64To2Digits(info.GetDownwardsDistance()/1000), formater.Separator,
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

func outPutContainsLineByTimeStamps(output []string, newLine string) bool {

	newLineStartTime := getStartTimeFormOutPutLine(newLine)
	newLineEndTime := getEndTimeFormOutPutLine(newLine)

	// Don't tread all lines with no valid time values as duplicates
	if newLineStartTime == NotValidValue {
		return false
	}

	for _, outLine := range output {
		outLineStartTime := getStartTimeFormOutPutLine(outLine)
		outLineEndTime := getEndTimeFormOutPutLine(outLine)

		if outLineStartTime == newLineStartTime && outLineEndTime == newLineEndTime {
			return true
		}
	}

	return false
}

func getStartTimeFormOutPutLine(line string) string {
	lineFields := strings.Split(line, ";")

	return lineFields[1]
}

func getEndTimeFormOutPutLine(line string) string {
	lineFields := strings.Split(line, ";")

	return lineFields[2]
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
