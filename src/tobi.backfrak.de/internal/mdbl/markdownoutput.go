package mdbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"tobi.backfrak.de/internal/gpsabl"
)

// NotValidValue - The value set when values are not valid
const NotValidValue = "not valid"

// MDOutputFormatertype - The gpsabl.OutputFormaterType this formater is responsible for
const MDOutputFormatertype gpsabl.OutputFormaterType = "MD"
const FileExtension = ".md"

// GetValidTimeFormats -  Get the valid TimeFormat values
func GetValidTimeFormats() []gpsabl.TimeFormat {
	return gpsabl.GetValidTimeFormats()
}

// CheckTimeFormatIsValid - Check if the given format string is a valid TimeFormat
func CheckTimeFormatIsValid(format string) bool {
	return strings.Contains(gpsabl.GetValidTimeFormatsString(), format)
}

// MDOutputFormater - type that formats TrackSummary into a markdown table
type MDOutputFormater struct {
	timeFormater        gpsabl.TimeFormat
	Separator           string
	writtenEntiresCount int
	entriesToWriteCount int
	lineBuffer          []gpsabl.OutputLine
	mux                 sync.Mutex
	TrackListText       string
	SummaryText         string
}

// NewMDOutputFormater - Get a new MDOutputFormater
func NewMDOutputFormater() *MDOutputFormater {
	ret := MDOutputFormater{}
	ret.writtenEntiresCount = -1
	ret.entriesToWriteCount = 0
	ret.timeFormater = gpsabl.RFC3339
	ret.lineBuffer = []gpsabl.OutputLine{}
	ret.Separator = "|"
	ret.TrackListText = "List of Tracks:"
	ret.SummaryText = "Summary table:"

	return &ret
}

// NewOutputFormater -  Get a new gpsabl.OutputFormater of this type
func (formater *MDOutputFormater) NewOutputFormater() gpsabl.OutputFormater {
	ret := NewMDOutputFormater()

	return gpsabl.OutputFormater(ret)
}

// GetTimeFormat - Get the time format string used by this MDOutputFormater
func (formater *MDOutputFormater) GetTimeFormat() string {
	return string(formater.timeFormater)
}

// SetTimeFormat - Set the time format string used by this MDOutputFormater. Will return an error if you want to set an unknown format
func (formater *MDOutputFormater) SetTimeFormat(timeFormat string) error {
	if CheckTimeFormatIsValid(timeFormat) == false {
		return gpsabl.NewTimeFormatNotKnown(gpsabl.TimeFormat(timeFormat))
	}
	formater.timeFormater = gpsabl.TimeFormat(timeFormat)
	return nil
}

// GetTextOutputFormater - Get the gpsabl.TextOutputFormater of ths formater
func (formater *MDOutputFormater) GetTextOutputFormater() gpsabl.TextOutputFormater {
	return nil
}

// CheckTimeFormatIsValid - Check if the given format string is a valid TimeFormat
func (formater *MDOutputFormater) CheckTimeFormatIsValid(format string) bool {
	return strings.Contains(gpsabl.GetValidTimeFormatsString(), format)
}

// AddOutPut - Add the formated output of a TrackFile to the internal buffer, so it can be written out later
func (formater *MDOutputFormater) AddOutPut(trackFile gpsabl.TrackFile, depth gpsabl.DepthArg, filterDuplicate bool) error {

	var lines []gpsabl.OutputLine
	linesFromFile, err := formater.getOutPutEntries(trackFile, depth)
	if err != nil {
		return err
	}
	if filterDuplicate {
		for _, line := range linesFromFile {
			if gpsabl.OutputContainsLineByTimeStamps(lines, line) == false && gpsabl.OutputContainsLineByTimeStamps(formater.lineBuffer, line) == false {
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

// CheckOutputFormaterType - Check if this OutputFormater is responsible for the given gpsabl.OutputFormaterType
func (formater *MDOutputFormater) CheckOutputFormaterType(formaterType gpsabl.OutputFormaterType) bool {
	if formaterType == MDOutputFormatertype {
		return true
	}

	return false
}

// GetFileExtensions - Get the list of file extensions this formater can write
func (formater *MDOutputFormater) GetFileExtensions() []string {
	return []string{FileExtension}
}

// GetOutputFormaterTypes - Get the list of gpsabl.OutputFormaterType this formater can write
func (formater *MDOutputFormater) GetOutputFormaterTypes() []gpsabl.OutputFormaterType {
	return []gpsabl.OutputFormaterType{MDOutputFormatertype}
}

// CheckFileExtension - Check if this OutputFormater can write the given output file
func (formater *MDOutputFormater) CheckFileExtension(filePath string) bool {
	if strings.HasSuffix(strings.ToLower(filePath), FileExtension) {
		return true
	}

	return false
}

// Tells the number if output entries already written to output.
// * -1: When output was not written yet
// * 0: Output was written but contains no entries, may because no entry passes the given filter
// * >0: The number of entries written to the outputs
func (formater *MDOutputFormater) GetNumberOfOutputEntries() int {
	return formater.writtenEntiresCount
}

// GetStatisticSummaryLines - Get a summary of statistic data
func (formater *MDOutputFormater) GetStatisticSummaryLines() []string {
	ret := []string{}

	if len(formater.lineBuffer) > 0 {
		summary := gpsabl.GetStatisticSummaryData(formater.lineBuffer)

		ret = append(ret, formater.formatSumSummary(summary.Sum, summary.AllTimeDataValid))
		ret = append(ret, formater.formatAverageSummary(summary.Average, summary.AllTimeDataValid))
		ret = append(ret, formater.formatMinMaxSummary(summary.Minimum, summary.AllTimeDataValid, "Minimum:"))
		ret = append(ret, formater.formatMinMaxSummary(summary.Maximum, summary.AllTimeDataValid, "Maximum:"))
	}
	return ret
}

// GetLines - Get the lines stored in the internal buffer
func (formater *MDOutputFormater) GetLines() []string {
	ret := []string{}

	formater.mux.Lock()
	defer formater.mux.Unlock()
	sort.Slice(formater.lineBuffer, func(i, j int) bool {
		return formater.lineBuffer[i].Data.GetStartTime().Before(formater.lineBuffer[j].Data.GetStartTime())
	})
	for _, line := range formater.lineBuffer {
		ret = append(ret, formater.FormatTrackSummary(line.Data, line.Name))
	}
	return ret
}

// WriteOutput - Write the output to a given file handle object. Make sure the file exists before you call this method!
func (formater *MDOutputFormater) WriteOutput(outFile *os.File, summary gpsabl.SummaryArg) error {
	lines, getErr := formater.GetOutputLines(summary)
	if getErr != nil {
		return getErr
	}

	if formater.entriesToWriteCount == 0 {
		formater.writtenEntiresCount = formater.entriesToWriteCount
		return nil
	}

	for _, line := range lines {
		_, errWrite := outFile.WriteString(line)
		if errWrite != nil {
			return errWrite
		}
	}

	formater.writtenEntiresCount = formater.entriesToWriteCount
	return nil
}

// GetOutputLines - Get all lines of the output
func (formater *MDOutputFormater) GetOutputLines(summary gpsabl.SummaryArg) ([]string, error) {
	var contentLines []string
	var outputLines []string
	var headerLines []string
	if summary == gpsabl.ADDITIONAL {
		outputLines = append(outputLines, fmt.Sprintf("%s%s", formater.TrackListText, GetNewLine()))
		outputLines = append(outputLines, GetNewLine())
	}
	headerLines = append(headerLines, formater.GetHeader())
	headerLines = append(headerLines, formater.GetHeaderContentSeparator())
	switch summary {
	case gpsabl.NONE:
		contentLines = append(contentLines, formater.GetLines()...)
	case gpsabl.ONLY:
		contentLines = append(contentLines, formater.GetStatisticSummaryLines()...)
	case gpsabl.ADDITIONAL:
		contentLines = append(contentLines, formater.GetLines()...)
	default:
		return nil, gpsabl.NewSummaryParamaterNotKnown(summary)
	}

	formater.entriesToWriteCount = len(contentLines)
	if formater.entriesToWriteCount > 0 {
		outputLines = append(outputLines, headerLines...)
		outputLines = append(outputLines, contentLines...)
		if summary == gpsabl.ADDITIONAL {
			outputLines = append(outputLines, GetNewLine())
			outputLines = append(outputLines, fmt.Sprintf("%s%s", formater.SummaryText, GetNewLine()))
			outputLines = append(outputLines, GetNewLine())
			outputLines = append(outputLines, headerLines...)
			outputLines = append(outputLines, formater.GetStatisticSummaryLines()...)
		}
		return outputLines, nil
	}

	return contentLines, nil
}

// getOutPutEntries - Add the output of a TrackFile
func (formater *MDOutputFormater) getOutPutEntries(trackFile gpsabl.TrackFile, depth gpsabl.DepthArg) ([]gpsabl.OutputLine, error) {

	return gpsabl.GetOutlines(trackFile, depth)
}

// GetHeader - Get the header line of a csv output
func (formater *MDOutputFormater) GetHeader() string {
	trackTimeHeader, _ := formater.getTimeDurationHeader("TrackTime")
	movingTimeHeader, _ := formater.getTimeDurationHeader("MovingTime")
	upwardsTimeHeader, _ := formater.getTimeDurationHeader("UpwardsTime")
	downwardsTimeHeader, _ := formater.getTimeDurationHeader("DownwardsTime")
	ret := fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s%s",
		formater.Separator,
		"Name", formater.Separator,
		"StartTime", formater.Separator,
		"EndTime", formater.Separator,
		trackTimeHeader, formater.Separator,
		"Distance (km)", formater.Separator,
		"HorizontalDistance (km)", formater.Separator,
		"AltitudeRange (m)", formater.Separator,
		"MinimumAltitude (m)", formater.Separator,
		"MaximumAltitude (m)", formater.Separator,
		"ElevationGain (m)", formater.Separator,
		"ElevationLose (m)", formater.Separator,
		"UpwardsDistance (km)", formater.Separator,
		"DownwardsDistance (km)", formater.Separator,
		movingTimeHeader, formater.Separator,
		upwardsTimeHeader, formater.Separator,
		downwardsTimeHeader, formater.Separator,
		"AverageSpeed (km/h)", formater.Separator,
		"UpwardsSpeed (km/h)", formater.Separator,
		"DownwardsSpeed (km/h)", formater.Separator,
		GetNewLine(),
	)

	return ret
}

func (formater *MDOutputFormater) GetHeaderContentSeparator() string {
	ret := fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s%s",
		formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		" :----: ", formater.Separator,
		GetNewLine(),
	)

	return ret
}

// FormatTrackSummary - Create the OutputLine for a TrackSummaryProvider
func (formater *MDOutputFormater) FormatTrackSummary(info gpsabl.TrackSummaryProvider, name string) string {
	var ret string
	if info.GetTimeDataValid() {
		duration, _ := formater.formatTimeDuration(info.GetEndTime().Sub(info.GetStartTime()))
		moveTime, _ := formater.formatTimeDuration(info.GetMovingTime())
		upTime, _ := formater.formatTimeDuration(info.GetUpwardsTime())
		downTime, _ := formater.formatTimeDuration(info.GetDownwardsTime())
		ret = fmt.Sprintf("%s %s %s %s %s %s %s %s %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %s %s %s %s %s %s %.2f %s %.2f %s %.2f %s%s",
			formater.Separator,
			name, formater.Separator,
			info.GetStartTime().Format(string(formater.timeFormater)), formater.Separator,
			info.GetEndTime().Format(string(formater.timeFormater)), formater.Separator,
			duration, formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.GetDistance()/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.GetHorizontalDistance()/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.GetAltitudeRange())), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.GetMinimumAltitude())), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.GetMaximumAltitude())), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.GetElevationGain())), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.GetElevationLose())), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.GetUpwardsDistance()/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.GetDownwardsDistance()/1000), formater.Separator,
			moveTime, formater.Separator,
			upTime, formater.Separator,
			downTime, formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.GetAvarageSpeed()*3.6), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.GetUpwardsSpeed()*3.6), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.GetDownwardsSpeed()*3.6), formater.Separator,
			GetNewLine(),
		)
	} else {
		ret = fmt.Sprintf("%s %s %s %s %s %s %s %s %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %s %s %s %s %s %s %s %s %s %s %s %s%s",
			formater.Separator,
			name, formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.GetDistance()/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.GetHorizontalDistance()/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.GetAltitudeRange())), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.GetMinimumAltitude())), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.GetMaximumAltitude())), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.GetElevationGain())), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.GetElevationLose())), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.GetUpwardsDistance()/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.GetDownwardsDistance()/1000), formater.Separator,
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

func (formater *MDOutputFormater) formatMinMaxSummary(info gpsabl.ExtendedTrackSummary, timeValid bool, name string) string {
	var ret string
	if timeValid {
		duration, _ := formater.formatTimeDuration(info.Duration)
		moveTime, _ := formater.formatTimeDuration(info.MovingTime)
		upTime, _ := formater.formatTimeDuration(info.UpwardsTime)
		downTime, _ := formater.formatTimeDuration(info.DownwardsTime)
		ret = fmt.Sprintf("%s %s %s %s %s %s %s %s %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %s %s %s %s %s %s %.2f %s %.2f %s %.2f %s%s",
			formater.Separator,
			fmt.Sprintf("**%s**", name), formater.Separator,
			info.StartTime.Format(string(formater.timeFormater)), formater.Separator,
			info.EndTime.Format(string(formater.timeFormater)), formater.Separator,
			duration, formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.Distance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.HorizontalDistance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.AltitudeRange)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.MinimumAltitude)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.MaximumAltitude)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.ElevationGain)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.ElevationLose)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.UpwardsDistance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.DownwardsDistance/1000), formater.Separator,
			moveTime, formater.Separator,
			upTime, formater.Separator,
			downTime, formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.AverageSpeed*3.6), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.UpwardsSpeed*3.6), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.DownwardsSpeed*3.6), formater.Separator,
			GetNewLine(),
		)
	} else {
		ret = fmt.Sprintf("%s %s %s %s %s %s %s %s %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %.2f %s %s %s %s %s %s %s %s %s %s %s %s %s%s",
			formater.Separator,
			fmt.Sprintf("**%s**", name), formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.Distance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.HorizontalDistance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.AltitudeRange)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.MinimumAltitude)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.MaximumAltitude)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.ElevationGain)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.ElevationLose)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.UpwardsDistance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.DownwardsDistance/1000), formater.Separator,
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

// FormatTrackSummary - Create the OutputLine for a TrackSummaryProvider
func (formater *MDOutputFormater) formatAverageSummary(info gpsabl.ExtendedTrackSummary, timeValid bool) string {
	var ret string
	if timeValid {
		duration, _ := formater.formatTimeDuration(info.Duration)
		moveTime, _ := formater.formatTimeDuration(info.MovingTime)
		upTime, _ := formater.formatTimeDuration(info.UpwardsTime)
		downTime, _ := formater.formatTimeDuration(info.DownwardsTime)
		ret = fmt.Sprintf("%s %s %s %s %s %s %s %s %s %.2f %s %.2f %s %.2f %s %s %s %s %s %.2f %s %.2f %s %.2f %s %.2f %s %s %s %s %s %s %s %.2f %s %.2f %s %.2f %s%s",
			formater.Separator,
			"**Average**:", formater.Separator,
			"-", formater.Separator,
			"-", formater.Separator,
			duration, formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.Distance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.HorizontalDistance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.AltitudeRange)), formater.Separator,
			"-", formater.Separator,
			"-", formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.ElevationGain)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.ElevationLose)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.UpwardsDistance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.DownwardsDistance/1000), formater.Separator,
			moveTime, formater.Separator,
			upTime, formater.Separator,
			downTime, formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.AverageSpeed*3.6), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.UpwardsSpeed*3.6), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.DownwardsSpeed*3.6), formater.Separator,
			GetNewLine(),
		)
	} else {
		ret = fmt.Sprintf("%s %s %s %s %s %s %s %s %s %.2f %s %.2f %s %.2f %s %s %s %s %s %.2f %s %.2f %s %.2f %s %.2f %s %s %s %s %s %s %s %s %s %s %s %s %s%s",
			formater.Separator,
			"**Average:**", formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.Distance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.HorizontalDistance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.AltitudeRange)), formater.Separator,
			"-", formater.Separator,
			"-", formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.ElevationGain)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.ElevationLose)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.UpwardsDistance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.DownwardsDistance/1000), formater.Separator,
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

func (formater *MDOutputFormater) formatSumSummary(info gpsabl.ExtendedTrackSummary, timeValid bool) string {
	var ret string
	if timeValid {
		duration, _ := formater.formatTimeDuration(info.Duration)
		moveTime, _ := formater.formatTimeDuration(info.MovingTime)
		upTime, _ := formater.formatTimeDuration(info.UpwardsTime)
		downTime, _ := formater.formatTimeDuration(info.DownwardsTime)
		ret = fmt.Sprintf("%s %s %s %s %s %s %s %s %s %.2f %s %.2f %s %s %s %s %s %s %s %.2f %s %.2f %s %.2f %s %.2f %s %s %s %s %s %s %s %s %s %s %s %s %s%s",
			formater.Separator,
			"**Sum:**", formater.Separator,
			"-", formater.Separator,
			"-", formater.Separator,
			duration, formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.Distance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.HorizontalDistance/1000), formater.Separator,
			"-", formater.Separator,
			"-", formater.Separator,
			"-", formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.ElevationGain)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.ElevationLose)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.UpwardsDistance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.DownwardsDistance/1000), formater.Separator,
			moveTime, formater.Separator,
			upTime, formater.Separator,
			downTime, formater.Separator,
			"-", formater.Separator,
			"-", formater.Separator,
			"-", formater.Separator,
			GetNewLine(),
		)
	} else {
		ret = fmt.Sprintf("%s %s %s %s %s %s %s %s %s %.2f %s %.2f %s %s %s %s %s %s %s %.2f %s %.2f %s %.2f %s %.2f %s %s %s %s %s %s %s %s %s %s %s %s %s%s",
			formater.Separator,
			"**Sum:**", formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			NotValidValue, formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.Distance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.HorizontalDistance/1000), formater.Separator,
			"-", formater.Separator,
			"-", formater.Separator,
			"-", formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.ElevationGain)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(float64(info.ElevationLose)), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.UpwardsDistance/1000), formater.Separator,
			gpsabl.RoundFloat64To2Digits(info.DownwardsDistance/1000), formater.Separator,
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

// GetNewLine - Get the new line string depending on the OS
func GetNewLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"

}

func (formater *MDOutputFormater) formatTimeDuration(duration time.Duration) (string, error) {
	switch formater.timeFormater {
	case gpsabl.RFC850:
		str := strings.ReplaceAll(duration.String(), "s", "")
		str = strings.ReplaceAll(str, "m", ":")
		str = strings.ReplaceAll(str, "h", ":")
		return str, nil
	case gpsabl.RFC3339:
		return duration.String(), nil
	case gpsabl.UnixDate:
		return fmt.Sprintf("%.2f", duration.Seconds()), nil
	default:
		return "", gpsabl.NewTimeFormatNotKnown(formater.timeFormater)
	}
}

func (formater *MDOutputFormater) getTimeDurationHeader(prefix string) (string, error) {
	switch formater.timeFormater {
	case gpsabl.RFC850:
		return fmt.Sprintf("%s (%s)", prefix, "hh:mm:ss"), nil
	case gpsabl.RFC3339:
		return fmt.Sprintf("%s (%s)", prefix, "xxhxxmxxs"), nil
	case gpsabl.UnixDate:
		return fmt.Sprintf("%s (%s)", prefix, "s"), nil
	default:
		return "", gpsabl.NewTimeFormatNotKnown(formater.timeFormater)
	}
}
