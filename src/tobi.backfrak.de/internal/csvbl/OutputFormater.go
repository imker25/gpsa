package csvbl

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

// CSVOutputFormatertype - The gpsabl.OutputFormaterType this formater is responsible for
const CSVOutputFormatertype gpsabl.OutputFormaterType = "CSV"
const FileExtension = ".csv"

// DefaultOutputSeperator - The seperator string for csv output files
const DefaultOutputSeperator = "; "

// TimeFormat - Represents a go Time format string for the enum pattern
type TimeFormat string

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

// TimeFormatNotKnown - Error when the given -summary is not known
type TimeFormatNotKnown struct {
	err string
	// File - The path to the dir that caused this error
	GivenValue TimeFormat
}

func (e *TimeFormatNotKnown) Error() string { // Implement the Error Interface for the TimeFormatNotKnown struct
	return fmt.Sprintf("%s", e.err)
}

// NewTimeFormatNotKnown - Get a new TimeFormatNotKnown struct
func NewTimeFormatNotKnown(givenValue TimeFormat) *TimeFormatNotKnown {
	return &TimeFormatNotKnown{fmt.Sprintf("The given -summary \"%s\" is not known.", givenValue), givenValue}
}

// CheckTimeFormatIsValid - Check if the given format string is a valid TimeFormat
func CheckTimeFormatIsValid(format string) bool {
	return strings.Contains(GetValidTimeFormatsString(), format)
}

// CsvOutputFormater - type that formats TrackSummary into csv style
type CsvOutputFormater struct {
	// Separator - The separator used to separate values in csv
	Separator string

	// Tell if the CSV header should be added to the output
	AddHeader bool

	timeFormater TimeFormat

	lineBuffer []gpsabl.OutputLine
	mux        sync.Mutex
}

// NewCsvOutputFormater - Get a new CsvOutputFormater
func NewCsvOutputFormater(separator string, addHeader bool) *CsvOutputFormater {
	ret := CsvOutputFormater{}
	ret.Separator = separator
	ret.AddHeader = addHeader
	ret.timeFormater = RFC3339
	ret.lineBuffer = []gpsabl.OutputLine{}

	return &ret
}

// NewOutputFormater -  Get a new gpsabl.OutputFormater of this type
func (formater *CsvOutputFormater) NewOutputFormater() gpsabl.OutputFormater {
	ret := NewCsvOutputFormater(DefaultOutputSeperator, true)

	return gpsabl.OutputFormater(ret)
}

// GetTimeFormat - Get the time format string used by this CsvOutputFormater
func (formater *CsvOutputFormater) GetTimeFormat() string {
	return string(formater.timeFormater)
}

// SetTimeFormat - Set the time format string used by this CsvOutputFormater. Will return an error if you want to set an unknown format
func (formater *CsvOutputFormater) SetTimeFormat(timeFormat string) error {
	if CheckTimeFormatIsValid(timeFormat) == false {
		return NewTimeFormatNotKnown(TimeFormat(timeFormat))
	}
	formater.timeFormater = TimeFormat(timeFormat)
	return nil
}

// AddOutPut - Add the formated output of a TrackFile to the internal buffer, so it can be written out later
func (formater *CsvOutputFormater) AddOutPut(trackFile gpsabl.TrackFile, depth gpsabl.DepthArg, filterDuplicate bool) error {

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
func (formater *CsvOutputFormater) CheckOutputFormaterType(formaterType gpsabl.OutputFormaterType) bool {
	if formaterType == CSVOutputFormatertype {
		return true
	}

	return false
}

// GetFileExtensions - Get the list of file extensions this formater can write
func (formater *CsvOutputFormater) GetFileExtensions() []string {
	return []string{FileExtension}
}

// GetOutputFormaterTypes - Get the list of gpsabl.OutputFormaterType this formater can write
func (formater *CsvOutputFormater) GetOutputFormaterTypes() []gpsabl.OutputFormaterType {
	return []gpsabl.OutputFormaterType{CSVOutputFormatertype}
}

// CheckFileExtension - Check if this OutputFormater can write the given output file
func (formater *CsvOutputFormater) CheckFileExtension(filePath string) bool {
	if strings.HasSuffix(strings.ToLower(filePath), FileExtension) {
		return true
	}

	return false
}

// GetStatisticSummaryLines - Get a summary of statistic data
func (formater *CsvOutputFormater) GetStatisticSummaryLines() []string {
	ret := []string{}
	summary := gpsabl.GetStatisticSummaryData(formater.lineBuffer)

	ret = append(ret, formater.formatSumSummary(summary.Sum, summary.AllTimeDataValid))
	ret = append(ret, formater.formatAverageSummary(summary.Average, summary.AllTimeDataValid))
	ret = append(ret, formater.formatMinMaxSummary(summary.Minimum, summary.AllTimeDataValid, "Minimum:"))
	ret = append(ret, formater.formatMinMaxSummary(summary.Maximum, summary.AllTimeDataValid, "Maximum:"))

	return ret
}

// GetLines - Get the lines stored in the internal buffer
func (formater *CsvOutputFormater) GetLines() []string {
	ret := []string{}
	if formater.AddHeader {
		ret = append(ret, formater.GetHeader())
	}
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
func (formater *CsvOutputFormater) WriteOutput(outFile *os.File, summary gpsabl.SummaryArg) error {
	lines, getErr := formater.GetOutputLines(summary)
	if getErr != nil {
		return getErr
	}

	for _, line := range lines {
		_, errWrite := outFile.WriteString(line)
		if errWrite != nil {
			return errWrite
		}
	}

	return nil
}

// GetOutputLines - Get all lines of the output
func (formater *CsvOutputFormater) GetOutputLines(summary gpsabl.SummaryArg) ([]string, error) {
	var lines []string
	switch summary {
	case gpsabl.NONE:
		lines = formater.GetLines()
	case gpsabl.ONLY:
		if formater.AddHeader {
			lines = append(lines, formater.GetHeader())
		}
		lines = append(lines, formater.GetStatisticSummaryLines()...)
	case gpsabl.ADDITIONAL:
		lines = formater.GetLines()
		sepaeratorLine := fmt.Sprintf("%s%s%s", "Statistics:", formater.Separator, GetNewLine())
		lines = append(lines, sepaeratorLine)
		lines = append(lines, formater.GetStatisticSummaryLines()...)
	default:
		return nil, gpsabl.NewSummaryParamaterNotKnown(summary)
	}

	return lines, nil
}

// getOutPutEntries - Add the output of a TrackFile
func (formater *CsvOutputFormater) getOutPutEntries(trackFile gpsabl.TrackFile, depth gpsabl.DepthArg) ([]gpsabl.OutputLine, error) {

	return gpsabl.GetOutlines(trackFile, depth)
}

// GetHeader - Get the header line of a csv output
func (formater *CsvOutputFormater) GetHeader() string {
	trackTimeHeader, _ := formater.getTimeDurationHeader("TrackTime")
	movingTimeHeader, _ := formater.getTimeDurationHeader("MovingTime")
	upwardsTimeHeader, _ := formater.getTimeDurationHeader("UpwardsTime")
	downwardsTimeHeader, _ := formater.getTimeDurationHeader("DownwardsTime")
	ret := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s",
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

// FormatTrackSummary - Create the OutputLine for a TrackSummaryProvider
func (formater *CsvOutputFormater) FormatTrackSummary(info gpsabl.TrackSummaryProvider, name string) string {
	var ret string
	if info.GetTimeDataValid() {
		duration, _ := formater.formatTimeDuration(info.GetEndTime().Sub(info.GetStartTime()))
		moveTime, _ := formater.formatTimeDuration(info.GetMovingTime())
		upTime, _ := formater.formatTimeDuration(info.GetUpwardsTime())
		downTime, _ := formater.formatTimeDuration(info.GetDownwardsTime())
		ret = fmt.Sprintf("%s%s%s%s%s%s%s%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%s%s%s%s%s%s%.2f%s%.2f%s%.2f%s%s",
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
		ret = fmt.Sprintf("%s%s%s%s%s%s%s%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%s%s%s%s%s%s%s%s%s%s%s%s%s",
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

func (formater *CsvOutputFormater) formatMinMaxSummary(info gpsabl.ExtendedTrackSummary, timeValid bool, name string) string {
	var ret string
	if timeValid {
		duration, _ := formater.formatTimeDuration(info.Duration)
		moveTime, _ := formater.formatTimeDuration(info.MovingTime)
		upTime, _ := formater.formatTimeDuration(info.UpwardsTime)
		downTime, _ := formater.formatTimeDuration(info.DownwardsTime)
		ret = fmt.Sprintf("%s%s%s%s%s%s%s%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%s%s%s%s%s%s%.2f%s%.2f%s%.2f%s%s",
			name, formater.Separator,
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
		ret = fmt.Sprintf("%s%s%s%s%s%s%s%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%.2f%s%s%s%s%s%s%s%s%s%s%s%s%s%s",
			name, formater.Separator,
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
func (formater *CsvOutputFormater) formatAverageSummary(info gpsabl.ExtendedTrackSummary, timeValid bool) string {
	var ret string
	if timeValid {
		duration, _ := formater.formatTimeDuration(info.Duration)
		moveTime, _ := formater.formatTimeDuration(info.MovingTime)
		upTime, _ := formater.formatTimeDuration(info.UpwardsTime)
		downTime, _ := formater.formatTimeDuration(info.DownwardsTime)
		ret = fmt.Sprintf("%s%s%s%s%s%s%s%s%.2f%s%.2f%s%.2f%s%s%s%s%s%.2f%s%.2f%s%.2f%s%.2f%s%s%s%s%s%s%s%.2f%s%.2f%s%.2f%s%s",
			"Average:", formater.Separator,
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
		ret = fmt.Sprintf("%s%s%s%s%s%s%s%s%.2f%s%.2f%s%.2f%s%s%s%s%s%.2f%s%.2f%s%.2f%s%.2f%s%s%s%s%s%s%s%s%s%s%s%s%s%s",
			"Average:", formater.Separator,
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

func (formater *CsvOutputFormater) formatSumSummary(info gpsabl.ExtendedTrackSummary, timeValid bool) string {
	var ret string
	if timeValid {
		duration, _ := formater.formatTimeDuration(info.Duration)
		moveTime, _ := formater.formatTimeDuration(info.MovingTime)
		upTime, _ := formater.formatTimeDuration(info.UpwardsTime)
		downTime, _ := formater.formatTimeDuration(info.DownwardsTime)
		ret = fmt.Sprintf("%s%s%s%s%s%s%s%s%.2f%s%.2f%s%s%s%s%s%s%s%.2f%s%.2f%s%.2f%s%.2f%s%s%s%s%s%s%s%s%s%s%s%s%s%s",
			"Sum:", formater.Separator,
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
		ret = fmt.Sprintf("%s%s%s%s%s%s%s%s%.2f%s%.2f%s%s%s%s%s%s%s%.2f%s%.2f%s%.2f%s%.2f%s%s%s%s%s%s%s%s%s%s%s%s%s%s",
			"Sum:", formater.Separator,
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

func (formater *CsvOutputFormater) formatTimeDuration(duration time.Duration) (string, error) {
	switch formater.timeFormater {
	case RFC850:
		str := strings.ReplaceAll(duration.String(), "s", "")
		str = strings.ReplaceAll(str, "m", ":")
		str = strings.ReplaceAll(str, "h", ":")
		return str, nil
	case RFC3339:
		return duration.String(), nil
	case UnixDate:
		return fmt.Sprintf("%.2f", duration.Seconds()), nil
	default:
		return "", NewTimeFormatNotKnown(formater.timeFormater)
	}
}

func (formater *CsvOutputFormater) getTimeDurationHeader(prefix string) (string, error) {
	switch formater.timeFormater {
	case RFC850:
		return fmt.Sprintf("%s (%s)", prefix, "hh:mm:ss"), nil
	case RFC3339:
		return fmt.Sprintf("%s (%s)", prefix, "xxhxxmxxs"), nil
	case UnixDate:
		return fmt.Sprintf("%s (%s)", prefix, "s"), nil
	default:
		return "", NewTimeFormatNotKnown(formater.timeFormater)
	}
}
