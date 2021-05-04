package csvbl

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"tobi.backfrak.de/internal/gpsabl"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

const numberOfSemicolonExpected = 19
const numberOfNotValideExpected = 9

func TestNewTimeFormatNotKnown(t *testing.T) {
	val := "asdgfg"
	err := NewTimeFormatNotKnown(TimeFormat(val))

	if err.GivenValue != TimeFormat(val) {
		t.Errorf("The GivenValue was %s, but %s was expected", err.GivenValue, val)
	}

	if strings.Contains(err.Error(), val) == false {
		t.Errorf("The error message of DepthParameterNotKnownError does not contain the expected GivenValue")
	}
}

func TestNewOutputFormater(t *testing.T) {
	var orig CsvOutputFormater
	sut := orig.NewOutputFormater()

	if sut.CheckFileExtension("my/output.csv") == false {
		t.Errorf("JSONOutputFormater can not write *.csv")
	}

	if sut.CheckFileExtension("my/output.json") == true {
		t.Errorf("JSONOutputFormater can write *.json")
	}

	if sut.CheckOutputFormaterType(CSVOutputFormatertype) == false {
		t.Errorf("JSONOutputFormater can not write %s type", CSVOutputFormatertype)
	}

	if sut.CheckOutputFormaterType(gpsabl.OutputFormaterType("abs")) == true {
		t.Errorf("JSONOutputFormater can write %s type", "abs")
	}

	ext := sut.GetFileExtensions()
	if len(ext) != 1 {
		t.Errorf("The number of FileExtensions is not expected")
	}

	if ext[0] != ".csv" {
		t.Errorf("The file type \"%s\" is not the expexted \"%s\"", ext[0], ".csv")
	}

	form := sut.GetOutputFormaterTypes()

	if len(form) != 1 {
		t.Errorf("The number of FileExtensionsTypes is not expected")
	}

	if form[0] != gpsabl.OutputFormaterType("CSV") {
		t.Errorf("The file type \"%s\" is not the expexted \"%s\"", form[0], "CSV")
	}

}

func TestSetTimeFormat(t *testing.T) {
	sut := NewCsvOutputFormater(";", false)

	if sut.GetTimeFormat() != time.RFC3339 {
		t.Errorf("The TimeFormat does not have the expected default value")
	}

	err1 := sut.SetTimeFormat(time.UnixDate)
	if err1 != nil {
		t.Errorf("Got an error but expected none")
	}

	if sut.GetTimeFormat() != time.UnixDate {
		t.Errorf("The TimeFormat does not have the expected default value")
	}

	str := "blablub"
	err2 := sut.SetTimeFormat(str)
	if err2 == nil {
		t.Errorf("Got no error but expected one")
	}

	if err2 == nil {
		t.Errorf("Got no errorbut expected one")
	}
	switch err2.(type) {
	case *TimeFormatNotKnown:
		fmt.Println("OK")
	default:
		t.Errorf("The error is not from the expected type")
	}
}

func TestGetValidTimeFormatsString(t *testing.T) {
	str := GetValidTimeFormatsString()
	res := fmt.Sprintf("\"%s\" \"%s\" \"%s\" ", time.UnixDate, time.RFC850, time.RFC3339)
	if str != res {
		t.Errorf("Got \"%s\", but expected \"%s \"", str, res)
	}
}

func TestNewCsvOutputFormater(t *testing.T) {
	sut := NewCsvOutputFormater(";", false)

	if sut.Separator != ";" {
		t.Errorf("The Separator was \"%s\", but \";\" was expected", sut.Separator)
	}

	if len(sut.lineBuffer) != 0 {
		t.Errorf("The line buffer is not empty on a new CsvOutputFormater")
	}

	if len(sut.GetLines()) != 0 {
		t.Errorf("The line buffer is not empty on a new CsvOutputFormater")
	}

	if sut.GetTimeFormat() != string(sut.timeFormater) {
		t.Errorf("GetTimeFormat does not have the expected value")
	}

	if sut.timeFormater != time.RFC3339 {
		t.Errorf("The TimeFormat does not have the expected default value")
	}
}

func TestFormatOutPutWithHeader(t *testing.T) {
	formater := NewCsvOutputFormater(";", true)
	trackFile := getSimpleTrackFile()

	err := formater.AddOutPut(trackFile, "track", false)
	if err != nil {
		t.Errorf("Got a error but did not expect one. The error is: %s", err.Error())
	}
	ret := formater.GetLines()
	if len(ret) != 2 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Count(ret[0], ";") != strings.Count(ret[1], ";") {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(ret[1], "0.02") == false {
		t.Errorf("The output does not contain the distance as expected. It is: %s", ret[1])
	}

	if strings.Contains(ret[1], trackFile.FilePath) == false {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[1])
	}
}

func TestFormatOutPutWithHeaderAndSetName(t *testing.T) {
	formater := NewCsvOutputFormater(";", true)
	trackFile := getSimpleTrackFile()
	trackFile.Name = "My Track File"

	err := formater.AddOutPut(trackFile, "track", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	ret := formater.GetLines()
	if len(ret) != 2 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Count(ret[0], ";") != strings.Count(ret[1], ";") {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(ret[1], "0.02") == false {
		t.Errorf("The output does not contain the distance as expected. It is: %s", ret[1])
	}

	if strings.Contains(ret[1], trackFile.Name) == false {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[1])
	}

	if strings.Contains(ret[1], trackFile.FilePath) == true {
		t.Errorf("The output does contain the FilePath but should not. It is: %s", ret[1])
	}
}

func TestFormatOutPutWithOutHeader(t *testing.T) {
	formater := NewCsvOutputFormater(";", false)
	trackFile := getSimpleTrackFile()

	entries, err := formater.getOutPutEntries(trackFile, "file")
	if err != nil {
		t.Errorf("Got a error but did not expect one. The error is: %s", err.Error())
	}
	ret := getLinesFormOutputLines(entries)
	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Count(ret[0], ";") != numberOfSemicolonExpected {
		t.Errorf("The Number of semicolons in the line is %d but %d was expected", strings.Count(ret[0], ";"), numberOfSemicolonExpected)
	}

	if strings.Contains(ret[0], "0.02") == false {
		t.Errorf("The output does not contain the distance as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderTrackDepth(t *testing.T) {
	formater := NewCsvOutputFormater(";", false)
	trackFile := getSimpleTrackFile()

	entries, err := formater.getOutPutEntries(trackFile, "track")
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	ret := getLinesFormOutputLines(entries)
	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Contains(ret[0], trackFile.FilePath) == false {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], "#1;") == false {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderTrackDepthSetTrackName(t *testing.T) {
	formater := NewCsvOutputFormater(";", false)
	trackFile := getSimpleTrackFile()
	trackFile.Tracks[0].Name = "My Track"

	entries, err := formater.getOutPutEntries(trackFile, "track")
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	ret := getLinesFormOutputLines(entries)
	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Contains(ret[0], trackFile.FilePath) == false {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], trackFile.Tracks[0].Name) == false {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], "#1;") == true {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderTrackDepthSetTrackFileNameSetTrackName(t *testing.T) {
	formater := NewCsvOutputFormater(";", false)
	trackFile := getSimpleTrackFile()
	trackFile.Tracks[0].Name = "My Track"
	trackFile.Name = "My track file"

	entries, err := formater.getOutPutEntries(trackFile, "track")
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	ret := getLinesFormOutputLines(entries)
	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Contains(ret[0], trackFile.FilePath) == true {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], trackFile.Name) == false {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], trackFile.Tracks[0].Name) == false {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], "#1;") == true {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderTrackDepthSegmentSetTrackFileNameSetTrackName(t *testing.T) {
	formater := NewCsvOutputFormater(";", false)
	trackFile := getSimpleTrackFile()
	trackFile.Tracks[0].Name = "My Track"
	trackFile.Name = "My track file"

	entries, err := formater.getOutPutEntries(trackFile, "segment")
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	ret := getLinesFormOutputLines(entries)
	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Contains(ret[0], trackFile.FilePath) == true {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], trackFile.Name) == false {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], trackFile.Tracks[0].Name) == false {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], "Segment #1;") == false {
		t.Errorf("The output does not contain the name as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderInvalidDepth(t *testing.T) {
	formater := NewCsvOutputFormater(";", false)
	trackFile := getSimpleTrackFile()
	_, err := formater.getOutPutEntries(trackFile, "abc")

	if err == nil {
		t.Errorf("Did not get an error as expected")
	}

	switch err.(type) {
	case *gpsabl.DepthParameterNotKnownError:
		fmt.Println("OK")
	default:
		t.Errorf("The error is not from the expected type")
	}
}

func TestAddOutPut(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile := getSimpleTrackFile()

	err := frt.AddOutPut(trackFile, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}

	lines := frt.GetLines()

	if len(lines) != 1 {
		t.Errorf("The number of lines was not expected. Got %d, expected %d", len(lines), 1)
	}

	if strings.Count(lines[0], ";") != numberOfSemicolonExpected {
		t.Errorf("The Number of semicolons in the line is %d but %d was expected", strings.Count(lines[0], ";"), numberOfSemicolonExpected)
	}

	if strings.Count(lines[0], "0.02;") != 2 {
		t.Errorf("The output does not contain the distance as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "1.00;") != 3 {
		t.Errorf("The output does not contain the ElevationGain as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "0.01;") != 2 {
		t.Errorf("The output does not contain the UpwardsDistance as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "not valid;") != numberOfNotValideExpected {
		t.Errorf("The output does not contain the Time values as expected. It is: %s", lines[0])
	}
}

func TestAddOutPutWithTimeStamp(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile := getSimpleTrackFileWithTime()

	err := frt.AddOutPut(trackFile, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}

	lines := frt.GetLines()

	if len(lines) != 1 {
		t.Errorf("The number of lines was not expected. Got %d, expected %d", len(lines), 1)
	}

	if strings.Count(lines[0], ";") != numberOfSemicolonExpected {
		t.Errorf("The Number of semicolons in the line is %d but %d was expected", strings.Count(lines[0], ";"), numberOfSemicolonExpected)
	}

	if strings.Count(lines[0], "0.02;") != 2 {
		t.Errorf("The output does not contain the distance as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "1.00;") != 3 {
		t.Errorf("The output does not contain the ElevationGain as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "0.01;") != 2 {
		t.Errorf("The output does not contain the UpwardsDistance as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "2014-08-22T17:19:33Z;") != 1 {
		t.Errorf("The output does not contain the StartTime as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "2014-08-22T17:19:53Z;") != 1 {
		t.Errorf("The output does not contain the EndTime as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "20s;") != 2 {
		t.Errorf("The output does not contain the MovingTime as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "4.30;") != 3 {
		t.Errorf("The output does not contain the AvarageSpeed as expected. It is: %s", lines[0])
	}
}

func TestWriteOutputSegmentDepth(t *testing.T) {
	frt := NewCsvOutputFormater(";", true)
	trackFile := getTrackFileTwoTracksWithThreeSegments()

	errAdd := frt.AddOutPut(trackFile, "segment", false)
	if errAdd != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", errAdd.Error())
	}

	errWrite := frt.WriteOutput(os.Stdout, "none")

	if errWrite != nil {
		t.Errorf("Error while writing the output: %s", errWrite.Error())
	}
}

func TestWriteOutputSummaryUnknown(t *testing.T) {
	frt := NewCsvOutputFormater(";", true)
	trackFile := getTrackFileTwoTracksWithThreeSegments()

	errAdd := frt.AddOutPut(trackFile, "segment", false)
	if errAdd != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", errAdd.Error())
	}

	errOut := frt.WriteOutput(os.Stdout, "bla")
	if errOut == nil {
		t.Errorf("Got no error, but was expected")
	}
	switch errOut.(type) {
	case *gpsabl.SummaryParamaterNotKnown:
		fmt.Println("OK")
	default:
		t.Errorf("The error is not of the expected type.")
	}
}

func TestCsvOutputFormaterIsOutputFormater(t *testing.T) {
	frt := NewCsvOutputFormater(";", true)
	trackFile := getTrackFileTwoTracksWithThreeSegments()

	iFrt := gpsabl.OutputFormater(frt)
	errAdd := iFrt.AddOutPut(trackFile, "track", false)
	if errAdd != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", errAdd.Error())
	}

	errWrite := iFrt.WriteOutput(os.Stdout, "none")

	if errWrite != nil {
		t.Errorf("Error while writing the output: %s", errWrite.Error())
	}
}

func TestCsvOutputFormaterDuplicateFilterWithTimeStamp(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile := getTrackFileTwoTracksWithThreeSegmentsWithTime()

	errAdd := frt.AddOutPut(trackFile, "track", true)
	if errAdd != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", errAdd.Error())
	}

	if len(frt.GetLines()) != 1 {
		t.Errorf("Got %d lines, but expected 1", len(frt.GetLines()))
	}
}

func TestCsvOutputFormaterDuplicateFilterWithOutTime(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile := getTrackFileTwoTracksWithThreeSegments()

	errAdd := frt.AddOutPut(trackFile, "track", true)
	if errAdd != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", errAdd.Error())
	}

	if len(frt.GetLines()) != 2 {
		t.Errorf("Got %d lines, but expected 1", len(frt.GetLines()))
	}
}

func TestAddOutPutWithUnValidFilter(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile := getTrackFileTwoTracksWithThreeSegments()

	errAdd := frt.AddOutPut(trackFile, "ba", false)
	if errAdd == nil {
		t.Errorf("Got no error but did expect one.")
	}

}
func TestAddOutPutWithUnValidFilterAndDuplicateFilter(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile := getTrackFileTwoTracksWithThreeSegments()

	errAdd := frt.AddOutPut(trackFile, "ba", true)
	if errAdd == nil {
		t.Errorf("Got no error but did expect one.")
	}

}

func TestAddOutPutMixedTimeAndNoTime(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile := getTrackFileOneTrackWithTimeOneWithout()

	err := frt.AddOutPut(trackFile, "track", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}

	lines := frt.GetLines()

	if len(lines) != 2 {
		t.Errorf("The number of lines was not expected. Got %d, expected %d", len(lines), 1)
	}

	if strings.Count(lines[1], ";") != numberOfSemicolonExpected {
		t.Errorf("The Number of semicolons in the line is %d but %d was expected", strings.Count(lines[0], ";"), numberOfSemicolonExpected)
	}

	if strings.Count(lines[1], "0.02;") != 2 {
		t.Errorf("The output does not contain the distance as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[1], "1.00;") != 3 {
		t.Errorf("The output does not contain the ElevationGain as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[1], "0.01;") != 2 {
		t.Errorf("The output does not contain the UpwardsDistance as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[1], "2014-08-22T17:19:33Z;") != 1 {
		t.Errorf("The output does not contain the StartTime as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[1], "2014-08-22T17:19:53Z;") != 1 {
		t.Errorf("The output does not contain the EndTime as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[1], "20s;") != 2 {
		t.Errorf("The output does not contain the MovingTime as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[1], "4.30;") != 3 {
		t.Errorf("The output does not contain the AvarageSpeed as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "not valid;") != numberOfNotValideExpected {
		t.Errorf("The output does not contain the Time values as often as expected. Found it %d times in: %s", strings.Count(lines[1], "not valid;"), lines[1])
	}

	if strings.Count(lines[1], "10s") != 2 {
		t.Errorf("The output does not contain the Time values as often as expected. Found it %d times in: %s", strings.Count(lines[0], "10s;"), lines[1])
	}
}

func TestOutPutTrackTimeAndMovingTimeIsDifferent(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile := getTrackFileWithTimeGaps()

	err := frt.AddOutPut(trackFile, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}

	lines := frt.GetLines()

	if strings.Count(lines[0], "2h0m20s;") != 1 {
		t.Errorf("The output does not contain the TrackTime as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "1m0s;") != 1 {
		t.Errorf("The output does not contain the MovingTime as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "30s;") != 2 {
		t.Errorf("The output does not contain the Upwards / Downwards Time as expected. It is: %s", lines[0])
	}
}

func TestOutPutDistanceAndHorizontalDistanceIsDifferent(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile := getTrackFileWithBigVerticalDistance()

	err := frt.AddOutPut(trackFile, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}

	lines := frt.GetLines()

	if strings.Count(lines[0], "0.07;") != 1 {
		t.Errorf("The output does not contain the Distance as expected. It is: %s", lines[0])
	}

	if strings.Count(lines[0], "0.05;") != 1 {
		t.Errorf("The output does not contain the HorizontalDistance as expected. It is: %s", lines[0])
	}
}

func TestOutPutContainsLineByTimeStamps1(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile := getTrackFileTwoTracksWithThreeSegmentsWithTime()
	entries, err := frt.getOutPutEntries(trackFile, "track")
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}

	if gpsabl.OutputContainsLineByTimeStamps(entries, entries[0]) == false {
		t.Errorf("Got false, but expect true")
	}

	if gpsabl.OutputContainsLineByTimeStamps(entries, *gpsabl.NewOutputLine("bla", getTrackWithDifferentTime())) == true {
		t.Errorf("Got true, but expect false")
	}
}

func TestOutPutContainsLineByTimeStamps2(t *testing.T) {

	frt := NewCsvOutputFormater(";", false)
	trackFile := getTrackFileTwoTracksWithThreeSegments()
	entries, err := frt.getOutPutEntries(trackFile, "track")
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}

	if gpsabl.OutputContainsLineByTimeStamps(entries, entries[0]) == true {
		t.Errorf("Got true, but expect false")
	}
}

func TestOutputIsSorted(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile1 := getTrackFileWithDifferentTime()
	err := frt.AddOutPut(trackFile1, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	trackFile2 := getSimpleTrackFileWithTime()
	err = frt.AddOutPut(trackFile2, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	lines := frt.GetLines()

	slpitLineOne := strings.Split(lines[0], ";")
	slpitLineTwo := strings.Split(lines[1], ";")
	if slpitLineOne[1] != "2014-08-22T17:19:33Z" {
		t.Errorf("The lines are not in the right order")
	}
	if slpitLineTwo[1] != "2015-08-22T17:19:33Z" {
		t.Errorf("The lines are not in the right order")
	}

}

func TestGetStatisticSummaryLinesWithTime(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile1 := getTrackFileWithDifferentTime()
	err := frt.AddOutPut(trackFile1, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	trackFile2 := getSimpleTrackFileWithTime()
	err = frt.AddOutPut(trackFile2, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	lines := frt.GetStatisticSummaryLines()

	if strings.Count(lines[0], ";") != numberOfSemicolonExpected {
		t.Errorf("In \"%s\" The number of semicolons is %d, but expected %d", lines[0], strings.Count(lines[0], ";"), numberOfSemicolonExpected)
	}

	if strings.Count(lines[1], ";") != numberOfSemicolonExpected {
		t.Errorf("In \"%s\" The number of semicolons is %d, but expected %d", lines[1], strings.Count(lines[1], ";"), numberOfSemicolonExpected)
	}
	if strings.Count(lines[2], ";") != numberOfSemicolonExpected {
		t.Errorf("In \"%s\" The number of semicolons is %d, but expected %d", lines[2], strings.Count(lines[2], ";"), numberOfSemicolonExpected)
	}
	if strings.Count(lines[3], ";") != numberOfSemicolonExpected {
		t.Errorf("In \"%s\" The number of semicolons is %d, but expected %d", lines[3], strings.Count(lines[3], ";"), numberOfSemicolonExpected)
	}
}

func TestGetOutputLinesSummaryNone(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile1 := getTrackFileWithDifferentTime()
	err := frt.AddOutPut(trackFile1, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	trackFile2 := getSimpleTrackFileWithTime()
	err = frt.AddOutPut(trackFile2, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	lines, errOut := frt.GetOutputLines("none")
	if errOut != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	if len(lines) != 2 {
		t.Errorf("Got an unexpected number of lines")
	}
}

func TestGetOutputLinesSummaryOnly(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile1 := getTrackFileWithDifferentTime()
	err := frt.AddOutPut(trackFile1, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	trackFile2 := getSimpleTrackFileWithTime()
	err = frt.AddOutPut(trackFile2, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	lines, errOut := frt.GetOutputLines("only")
	if errOut != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	if len(lines) != 4 {
		t.Errorf("Got an unexpected number of lines")
	}
}

func TestGetOutputLinesSummaryOnlyWithHeader(t *testing.T) {
	frt := NewCsvOutputFormater(";", true)
	trackFile1 := getTrackFileWithDifferentTime()
	err := frt.AddOutPut(trackFile1, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	trackFile2 := getSimpleTrackFileWithTime()
	err = frt.AddOutPut(trackFile2, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	lines, errOut := frt.GetOutputLines("only")
	if errOut != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	if len(lines) != 5 {
		t.Errorf("Got an unexpected number of lines")
	}
}

func TestGetOutputLinesSummaryAdditional(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile1 := getTrackFileWithDifferentTime()
	err := frt.AddOutPut(trackFile1, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	trackFile2 := getSimpleTrackFileWithTime()
	err = frt.AddOutPut(trackFile2, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	lines, errOut := frt.GetOutputLines("additional")
	if errOut != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	if len(lines) != 7 {
		t.Errorf("Got an unexpected number of lines")
	}
}

func TestGetOutputLinesSummaryUnValid(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile1 := getTrackFileWithDifferentTime()
	err := frt.AddOutPut(trackFile1, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	trackFile2 := getSimpleTrackFileWithTime()
	err = frt.AddOutPut(trackFile2, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	_, errOut := frt.GetOutputLines("bla")
	if errOut == nil {
		t.Errorf("Got no error, but was expected")
	}
	switch errOut.(type) {
	case *gpsabl.SummaryParamaterNotKnown:
		fmt.Println("OK")
	default:
		t.Errorf("The error is not of the expected type.")
	}

}

func TestGetStatisticSummaryLinesWithoutTime(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	trackFile1 := getTrackFileTwoTracks()
	err := frt.AddOutPut(trackFile1, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	trackFile2 := getSimpleTrackFileWithTime()
	err = frt.AddOutPut(trackFile2, "file", false)
	if err != nil {
		t.Errorf("Got an error but did not expect one. The error is: %s", err.Error())
	}
	lines := frt.GetStatisticSummaryLines()

	if strings.Count(lines[0], ";") != numberOfSemicolonExpected {
		t.Errorf("In \"%s\" The number of semicolons is %d, but expected %d", lines[0], strings.Count(lines[0], ";"), numberOfSemicolonExpected)
	}

	if strings.Count(lines[1], ";") != numberOfSemicolonExpected {
		t.Errorf("In \"%s\" The number of semicolons is %d, but expected %d", lines[1], strings.Count(lines[1], ";"), numberOfSemicolonExpected)
	}
	if strings.Count(lines[2], ";") != numberOfSemicolonExpected {
		t.Errorf("In \"%s\" The number of semicolons is %d, but expected %d", lines[2], strings.Count(lines[2], ";"), numberOfSemicolonExpected)
	}
	if strings.Count(lines[3], ";") != numberOfSemicolonExpected {
		t.Errorf("In \"%s\" The number of semicolons is %d, but expected %d", lines[3], strings.Count(lines[3], ";"), numberOfSemicolonExpected)
	}
}

func TestFormatTimeDurationUnknownFormat(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	frt.timeFormater = TimeFormat("blabla")
	now := time.Now()
	dur := now.Sub(now.Add(-2 * time.Hour))
	_, err := frt.formatTimeDuration(dur)

	if err == nil {
		t.Errorf("Got no error, but expected one")
	}
	switch err.(type) {
	case *TimeFormatNotKnown:
		fmt.Println("OK")
	default:
		t.Errorf("The error is not from the expected type")
	}
}

func TestFormatTimeDurationRFC3339Format(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	err := frt.SetTimeFormat(string(RFC3339))
	if err != nil {
		t.Errorf("Got an error, but expected none")
	}

	now := time.Now()
	dur := now.Sub(now.Add(-2 * time.Hour))
	str, _ := frt.formatTimeDuration(dur)
	res := dur.String()
	if res != str {
		t.Errorf("Expected %s, but got %s", res, str)
	}

	now = time.Now()
	dur = now.Sub(now.Add(-48 * time.Hour))
	str, _ = frt.formatTimeDuration(dur)
	res = dur.String()
	if res != str {
		t.Errorf("Expected %s, but got %s", res, str)
	}

	now = time.Now()
	dur = now.Sub(now.Add(-37*time.Hour + -2*time.Minute + -20*time.Second))
	str, _ = frt.formatTimeDuration(dur)
	res = dur.String()
	if res != str {
		t.Errorf("Expected %s, but got %s", res, str)
	}

	now = time.Now()
	dur = now.Sub(now.Add(-3*time.Minute + -21*time.Second))
	str, _ = frt.formatTimeDuration(dur)
	res = dur.String()
	if res != str {
		t.Errorf("Expected %s, but got %s", res, str)
	}

}

func TestFormatTimeDurationRFC850Format(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	err := frt.SetTimeFormat(string(RFC850))
	if err != nil {
		t.Errorf("Got an error, but expected none")
	}

	now := time.Now()
	dur := now.Sub(now.Add(-2 * time.Hour))
	str, _ := frt.formatTimeDuration(dur)
	res := "2:0:0"
	if res != str {
		t.Errorf("Expected %s, but got %s", res, str)
	}

	now = time.Now()
	dur = now.Sub(now.Add(-48 * time.Hour))
	str, _ = frt.formatTimeDuration(dur)
	res = "48:0:0"
	if res != str {
		t.Errorf("Expected %s, but got %s", res, str)
	}

	now = time.Now()
	dur = now.Sub(now.Add(-37*time.Hour + -2*time.Minute + -20*time.Second))
	str, _ = frt.formatTimeDuration(dur)
	res = "37:2:20"
	if res != str {
		t.Errorf("Expected %s, but got %s", res, str)
	}

	now = time.Now()
	dur = now.Sub(now.Add(-3*time.Minute + -21*time.Second))
	str, _ = frt.formatTimeDuration(dur)
	res = "3:21"
	if res != str {
		t.Errorf("Expected %s, but got %s", res, str)
	}

}

func TestFormatTimeDurationUnixFormat(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	err := frt.SetTimeFormat(string(UnixDate))
	if err != nil {
		t.Errorf("Got an error, but expected none")
	}

	now := time.Now()
	dur := now.Sub(now.Add(-2 * time.Hour))
	str, _ := frt.formatTimeDuration(dur)
	res := fmt.Sprintf("%.2f", float64(2*3600))
	if res != str {
		t.Errorf("Expected %s, but got %s", res, str)
	}

	now = time.Now()
	dur = now.Sub(now.Add(-48 * time.Hour))
	str, _ = frt.formatTimeDuration(dur)
	res = fmt.Sprintf("%.2f", float64(48*3600))
	if res != str {
		t.Errorf("Expected %s, but got %s", res, str)
	}

	now = time.Now()
	dur = now.Sub(now.Add(-37*time.Hour + -2*time.Minute + -20*time.Second))
	str, _ = frt.formatTimeDuration(dur)
	res = fmt.Sprintf("%.2f", float64(37*3600+140))
	if res != str {
		t.Errorf("Expected %s, but got %s", res, str)
	}

	now = time.Now()
	dur = now.Sub(now.Add(-3*time.Minute + -21*time.Second))
	str, _ = frt.formatTimeDuration(dur)
	res = fmt.Sprintf("%.2f", float64(201))
	if res != str {
		t.Errorf("Expected %s, but got %s", res, str)
	}

}

func TestGetTimeHeader(t *testing.T) {
	frt := NewCsvOutputFormater(";", false)
	frt.timeFormater = TimeFormat("blabla")
	_, err := frt.getTimeDurationHeader("bla")

	if err == nil {
		t.Errorf("Got no error, but expected one")
	}
	switch err.(type) {
	case *TimeFormatNotKnown:
		fmt.Println("OK")
	default:
		t.Errorf("The error is not from the expected type")
	}

	frt.timeFormater = RFC3339
	str, _ := frt.getTimeDurationHeader("bla")
	res := "bla (xxhxxmxxs)"
	if str != res {
		t.Errorf("Get  \"%s \" but expected \"%s\"", str, res)
	}

	frt.timeFormater = UnixDate
	str, _ = frt.getTimeDurationHeader("bla")
	res = "bla (s)"
	if str != res {
		t.Errorf("Get  \"%s \" but expected \"%s\"", str, res)
	}

	frt.timeFormater = RFC850
	str, _ = frt.getTimeDurationHeader("bla")
	res = "bla (hh:mm:ss)"
	if str != res {
		t.Errorf("Get  \"%s \" but expected \"%s\"", str, res)
	}
}

func getLinesFormOutputLines(lines []gpsabl.OutputLine) []string {
	ret := []string{}
	formater := NewCsvOutputFormater(";", true)
	for _, line := range lines {
		ret = append(ret, formater.FormatTrackSummary(line.Data, line.Name))
	}

	return ret
}

func getTrackFileWithDifferentTime() gpsabl.TrackFile {
	ret := gpsabl.NewTrackFile("/mys/track/file")
	trk := getTrackWithDifferentTime()
	gpsabl.FillTrackValues(&trk)
	ret.Tracks = []gpsabl.Track{trk}
	gpsabl.FillTrackFileValues(&ret)
	ret.NumberOfTracks = 1

	return ret
}

func getTrackPoint(lat, lon, ele float32) gpsabl.TrackPoint {
	pnt := gpsabl.TrackPoint{}
	pnt.Latitude = lat
	pnt.Longitude = lon
	pnt.Elevation = ele
	pnt.TimeValid = false

	return pnt
}

func getTrackPointWithTime(lat, lon, ele float32, time time.Time) gpsabl.TrackPoint {
	pnt := getTrackPoint(lat, lon, ele)
	pnt.TimeValid = true
	pnt.Time = time

	return pnt
}

func getTrackWithDifferentTime() gpsabl.Track {
	t1, _ := time.Parse(time.RFC3339, "2015-08-22T17:19:33Z")
	t2, _ := time.Parse(time.RFC3339, "2015-08-22T17:19:43Z")
	t3, _ := time.Parse(time.RFC3339, "2015-08-22T17:19:53Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 108.0, t2)
	pnt3 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t3)
	points := []gpsabl.TrackPoint{pnt1, pnt2, pnt3}

	gpsabl.FillDistancesTrackPoint(&points[0], gpsabl.TrackPoint{}, points[1])
	gpsabl.FillDistancesTrackPoint(&points[1], points[0], points[2])
	gpsabl.FillDistancesTrackPoint(&points[2], points[1], gpsabl.TrackPoint{})
	gpsabl.FillValuesTrackPointArray(points, "none", 0.3, 10.0)
	seg := gpsabl.TrackSegment{}
	seg.TrackPoints = points
	ret := gpsabl.Track{}
	gpsabl.FillTrackSegmentValues(&seg)
	ret.TrackSegments = []gpsabl.TrackSegment{seg}
	gpsabl.FillTrackValues(&ret)
	ret.NumberOfSegments = 1

	return ret
}

func getTrackFileTwoTracksWithThreeSegments() gpsabl.TrackFile {
	trackFile := getTrackFileTwoTracks()
	trackFile.Tracks[0].TrackSegments = append(trackFile.Tracks[0].TrackSegments, getSimpleTrackFile().Tracks[0].TrackSegments[0])
	gpsabl.FillTrackValues(&trackFile.Tracks[0])

	return trackFile
}

func getTrackFileTwoTracksWithThreeSegmentsWithTime() gpsabl.TrackFile {
	trackFile := getTrackFileTwoTracksWithTime()
	trackFile.Tracks[0].TrackSegments = append(trackFile.Tracks[0].TrackSegments, getSimpleTrackFileWithTime().Tracks[0].TrackSegments[0])
	gpsabl.FillTrackValues(&trackFile.Tracks[0])

	return trackFile
}

func getTrackFileTwoTracks() gpsabl.TrackFile {
	trackFile := getSimpleTrackFile()
	trackFile.Tracks = append(trackFile.Tracks, getSimpleTrackFile().Tracks...)
	gpsabl.FillTrackFileValues(&trackFile)

	return trackFile
}

func getTrackFileTwoTracksWithTime() gpsabl.TrackFile {
	trackFile := getSimpleTrackFileWithTime()
	trackFile.Tracks = append(trackFile.Tracks, getSimpleTrackFileWithTime().Tracks...)
	gpsabl.FillTrackFileValues(&trackFile)

	return trackFile
}

func getTrackFileOneTrackWithTimeOneWithout() gpsabl.TrackFile {
	trackFile := getSimpleTrackFileWithTime()
	trackFile.Tracks = append(trackFile.Tracks, getSimpleTrackFile().Tracks...)
	gpsabl.FillTrackFileValues(&trackFile)

	return trackFile
}

func getSimpleTrackFile() gpsabl.TrackFile {
	ret := gpsabl.NewTrackFile("/mys/track/file")
	trk := getSimpleTrack()
	gpsabl.FillTrackValues(&trk)
	ret.Tracks = []gpsabl.Track{trk}
	gpsabl.FillTrackFileValues(&ret)
	ret.NumberOfTracks = 1

	return ret
}

func getSimpleTrackFileWithTime() gpsabl.TrackFile {
	ret := gpsabl.NewTrackFile("/mys/track/file")
	trk := getSimpleTrackWithTime()
	gpsabl.FillTrackValues(&trk)
	ret.Tracks = []gpsabl.Track{trk}
	gpsabl.FillTrackFileValues(&ret)
	ret.NumberOfTracks = 1

	return ret
}

func getSimpleTrack() gpsabl.Track {
	ret := gpsabl.Track{}
	segs := getSimpleTrackSegment()
	gpsabl.FillTrackSegmentValues(&segs)
	ret.TrackSegments = []gpsabl.TrackSegment{segs}
	gpsabl.FillTrackValues(&ret)
	ret.NumberOfSegments = 1

	return ret
}

func getSimpleTrackWithTime() gpsabl.Track {
	ret := gpsabl.Track{}
	segs := getSimpleTrackSegmentWithTime()
	gpsabl.FillTrackSegmentValues(&segs)
	ret.TrackSegments = []gpsabl.TrackSegment{segs}
	gpsabl.FillTrackValues(&ret)
	ret.NumberOfSegments = 1

	return ret
}

func getSimpleTrackSegment() gpsabl.TrackSegment {
	seg := gpsabl.TrackSegment{}
	points := gerSimpleTrackPointArray()
	seg.TrackPoints = points

	return seg
}

func getSimpleTrackSegmentWithTime() gpsabl.TrackSegment {
	seg := gpsabl.TrackSegment{}
	points := getSimpleTrackPointArrayWithTime()
	seg.TrackPoints = points

	return seg
}

func gerSimpleTrackPointArray() []gpsabl.TrackPoint {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	points := []gpsabl.TrackPoint{pnt1, pnt2, pnt3}

	gpsabl.FillDistancesTrackPoint(&points[0], gpsabl.TrackPoint{}, points[1])
	gpsabl.FillDistancesTrackPoint(&points[1], points[0], points[2])
	gpsabl.FillDistancesTrackPoint(&points[2], points[1], gpsabl.TrackPoint{})
	gpsabl.FillValuesTrackPointArray(points, "none", 0.3, 10.0)

	return points
}

func getSimpleTrackPointArrayWithTime() []gpsabl.TrackPoint {
	t1, _ := time.Parse(time.RFC3339, "2014-08-22T17:19:33Z")
	t2, _ := time.Parse(time.RFC3339, "2014-08-22T17:19:43Z")
	t3, _ := time.Parse(time.RFC3339, "2014-08-22T17:19:53Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 108.0, t2)
	pnt3 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t3)
	points := []gpsabl.TrackPoint{pnt1, pnt2, pnt3}

	gpsabl.FillDistancesTrackPoint(&points[0], gpsabl.TrackPoint{}, points[1])
	gpsabl.FillDistancesTrackPoint(&points[1], points[0], points[2])
	gpsabl.FillDistancesTrackPoint(&points[2], points[1], gpsabl.TrackPoint{})
	gpsabl.FillValuesTrackPointArray(points, "none", 0.3, 10.0)

	return points
}

func getTrackFileWithStandStillPoints(correction string, minimalMovingSpeed float64, minimalStepHight float64) gpsabl.TrackFile {
	var file gpsabl.TrackFile

	t1, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:13Z")
	t2, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:33Z")
	t3, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:53Z")
	t4, _ := time.Parse(time.RFC3339, "2014-08-22T19:20:13Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 108.0, t2)
	pnt3 := getTrackPointWithTime(50.11495751, 8.684874771, 108.0, t3)
	pnt4 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t4)
	points := []gpsabl.TrackPoint{pnt1, pnt2, pnt3, pnt4}

	gpsabl.FillDistancesTrackPoint(&points[0], gpsabl.TrackPoint{}, points[1])
	gpsabl.FillDistancesTrackPoint(&points[1], points[0], points[2])
	gpsabl.FillDistancesTrackPoint(&points[2], points[1], points[3])
	gpsabl.FillDistancesTrackPoint(&points[3], points[2], gpsabl.TrackPoint{})
	gpsabl.FillValuesTrackPointArray(points, gpsabl.CorrectionParameter(correction), minimalMovingSpeed, minimalStepHight)
	laterTrack := gpsabl.Track{}
	seg := gpsabl.TrackSegment{}
	seg.TrackPoints = points
	gpsabl.FillTrackSegmentValues(&seg)
	laterTrack.TrackSegments = append(laterTrack.TrackSegments, seg)
	gpsabl.FillTrackValues(&laterTrack)
	laterTrack.NumberOfSegments = 1

	file.Tracks = append(file.Tracks, laterTrack)

	file.NumberOfTracks = 1
	gpsabl.FillTrackFileValues(&file)

	return file
}

func getTrackFileWithBigVerticalDistance() gpsabl.TrackFile {
	file := getSimpleTrackFileWithTime()

	t1, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:13Z")
	t2, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:33Z")
	t3, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:53Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 142.0, t2)
	pnt3 := getTrackPointWithTime(50.11484790, 8.684885500, 151.0, t3)
	points := []gpsabl.TrackPoint{pnt1, pnt2, pnt3}

	gpsabl.FillDistancesTrackPoint(&points[0], gpsabl.TrackPoint{}, points[1])
	gpsabl.FillDistancesTrackPoint(&points[1], points[0], points[2])
	gpsabl.FillDistancesTrackPoint(&points[2], points[1], gpsabl.TrackPoint{})
	gpsabl.FillValuesTrackPointArray(points, "none", 0.3, 10.0)
	laterTrack := gpsabl.Track{}
	seg := gpsabl.TrackSegment{}
	seg.TrackPoints = points
	gpsabl.FillTrackSegmentValues(&seg)
	laterTrack.TrackSegments = append(laterTrack.TrackSegments, seg)
	gpsabl.FillTrackValues(&laterTrack)
	laterTrack.NumberOfSegments = 1

	file.Tracks = append(file.Tracks, laterTrack)
	file.NumberOfTracks = 2
	gpsabl.FillTrackFileValues(&file)

	return file
}

func getTrackFileWithTimeGaps() gpsabl.TrackFile {
	file := getSimpleTrackFileWithTime()

	t1, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:13Z")
	t2, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:33Z")
	t3, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:53Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 108.0, t2)
	pnt3 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t3)
	points := []gpsabl.TrackPoint{pnt1, pnt2, pnt3}

	gpsabl.FillDistancesTrackPoint(&points[0], gpsabl.TrackPoint{}, points[1])
	gpsabl.FillDistancesTrackPoint(&points[1], points[0], points[2])
	gpsabl.FillDistancesTrackPoint(&points[2], points[1], gpsabl.TrackPoint{})
	gpsabl.FillValuesTrackPointArray(points, "none", 0.3, 10.0)
	laterTrack := gpsabl.Track{}
	seg := gpsabl.TrackSegment{}
	seg.TrackPoints = points
	gpsabl.FillTrackSegmentValues(&seg)
	laterTrack.TrackSegments = append(laterTrack.TrackSegments, seg)
	gpsabl.FillTrackValues(&laterTrack)
	laterTrack.NumberOfSegments = 1

	file.Tracks = append(file.Tracks, laterTrack)
	file.NumberOfTracks = 2
	gpsabl.FillTrackFileValues(&file)

	return file
}
