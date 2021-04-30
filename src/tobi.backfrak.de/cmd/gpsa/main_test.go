package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"tobi.backfrak.de/internal/jsonbl"

	"tobi.backfrak.de/internal/csvbl"
	"tobi.backfrak.de/internal/gpsabl"
	"tobi.backfrak.de/internal/testhelper"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// Use this Mutex to sync write access to the testdata/test-out.csv file
var outFileMux sync.Mutex

var ComandlineOptionsHandled bool

func TestHandleComandlineOptions(t *testing.T) {
	if !ComandlineOptionsHandled {
		handleComandlineOptions()
		ComandlineOptionsHandled = true
	}

	if HelpFlag == true {
		t.Errorf("The HelpFlag is set to true but false was expected")
	}

	if VerboseFlag == true {
		t.Errorf("The VerboseFlag is set to true but false was expected")
	}

	if SkipErrorExitFlag == true {
		t.Errorf("The SkipErrorExitFlag is set to true but false was expected")
	}

	if PrintVersionFlag == true {
		t.Errorf("The PrintVersionFlag is set to true but false was expected")
	}

	if PrintLicenseFlag == true {
		t.Errorf("The PrintLicenseFlag is set to true but false was expected")
	}

	if DepthParameter != "track" {
		t.Errorf("The DepthParameter is \"%s\" but \"track\" was expected", DepthParameter)
	}

	if OutFileParameter != "" {
		t.Errorf("The DepthParameter is \"%s\" but \"\" was expected", OutFileParameter)
	}

	if DontPanicFlag == false {
		t.Errorf("The DontPanicFlag is set to false but true was expected")
	}

	if PrintCsvHeaderFlag == false {
		t.Errorf("The PrintCsvHeaderFlag is set to false but true was expected")
	}

	if CorrectionParameter != "steps" {
		t.Errorf("The CorrectionParameter is \"%s\" but \"steps\" was expected", DepthParameter)
	}

	if MinimalStepHightParameter != 10.0 {
		t.Errorf("The MinimalStepHightParameter is \"%f\" but \"10.0\" was expected", MinimalStepHightParameter)
	}

	if MinimalMovingSpeedParameter != 0.3 {
		t.Errorf("The MinimalMovingSpeedParameter is \"%f\" but \"10.0\" was expected", MinimalMovingSpeedParameter)
	}
}

func TestCostumHelpMessage(t *testing.T) {
	if !ComandlineOptionsHandled {
		handleComandlineOptions()
		ComandlineOptionsHandled = true
	}
	customHelpMessage()
}

func TestHandleErrorNil(t *testing.T) {
	if HandleError(nil, "my/path", true, true) == true {
		t.Errorf("HandleError returns true, when nil error was given")
	}

}

func TestHandleErrorNotNil(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	if HandleError(newUnKnownFileTypeError("my/path"), "my/path", true, true) == false {
		t.Errorf("HandleError retutns false, when error was given")
	}

	if ErrorsHandled == false {
		t.Errorf("ErrorsHandled should be true after a error was handled")
	}
	SkipErrorExitFlag = oldFlagValue
	ErrorsHandled = false
}

func TestGetReaderGpxFile(t *testing.T) {
	reader, err := getReader("/some/track.gpx")

	if err != nil {
		t.Errorf("Got an error when try to get a reader for a gpx file: %s", err.Error())
	}

	if reader == nil {
		t.Errorf("The reader we got was nil")
	}
}

func TestGetReaderUnknownFile(t *testing.T) {
	reader, err := getReader("/some/track.txt")

	if err == nil {
		t.Errorf("Got no error when trying to get a reader for a txt file.")
	}

	if reader != nil {
		t.Errorf("The reader we got was not nil")
	}
}

func TestProcessValidFiles(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParameter
	DepthParameter = "file"
	oldCorrectionPAr := CorrectionParameter
	CorrectionParameter = "linear"

	formater := csvbl.NewCsvOutputFormater(";", false)
	iFormater := gpsabl.OutputFormater(formater)

	fileStrs := []string{testhelper.GetValidGPX("01.gpx"), testhelper.GetValidGPX("02.gpx"), testhelper.GetValidTcx("02.tcx")}
	var files []gpsabl.InputFile
	for _, file := range fileStrs {
		files = append(files, *gpsabl.NewInputFileWithPath(file))
	}
	successCount := processFiles(files, iFormater)
	if successCount != 3 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	if ErrorsHandled == true {
		t.Errorf("Errors occurred that were not expected")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParameter = oldDepthValue
	CorrectionParameter = oldCorrectionPAr
}

func TestProcessValidFilesWithEmpyElements(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParameter
	DepthParameter = "track"
	oldCorrectionPar := CorrectionParameter
	CorrectionParameter = "linear"
	oldPrintCsvHeaderFlag := PrintCsvHeaderFlag
	PrintCsvHeaderFlag = false

	formater := csvbl.NewCsvOutputFormater(";", false)
	iFormater := gpsabl.OutputFormater(formater)

	fileStrs := []string{testhelper.GetValidGPX("13.gpx"), testhelper.GetValidGPX("14.gpx")}
	var files []gpsabl.InputFile
	for _, file := range fileStrs {
		files = append(files, *gpsabl.NewInputFileWithPath(file))
	}
	successCount := processFiles(files, iFormater)
	if successCount != 2 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	if ErrorsHandled == true {
		t.Errorf("Errors occurred that were not expected")
	}

	if len(formater.GetLines()) != 2 {
		t.Errorf("Got %d lines, but expected 2", len(formater.GetLines()))
	}

	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParameter = oldDepthValue
	CorrectionParameter = oldCorrectionPar
	PrintCsvHeaderFlag = oldPrintCsvHeaderFlag
}

func TestProcessValidFilesWithDuplicateElementsDetectionenabeled(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParameter
	DepthParameter = "track"
	oldCorrectionPar := CorrectionParameter
	CorrectionParameter = "linear"
	oldPrintCsvHeaderFlag := PrintCsvHeaderFlag
	PrintCsvHeaderFlag = false
	oldSuppressDuplicateOutPutFlag := SuppressDuplicateOutPutFlag
	SuppressDuplicateOutPutFlag = true

	formater := csvbl.NewCsvOutputFormater(";", false)
	iFormater := gpsabl.OutputFormater(formater)

	fileStrs := []string{testhelper.GetValidGPX("15.gpx")}
	var files []gpsabl.InputFile
	for _, file := range fileStrs {
		files = append(files, *gpsabl.NewInputFileWithPath(file))
	}
	successCount := processFiles(files, iFormater)
	if successCount != 1 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	if ErrorsHandled == true {
		t.Errorf("Errors occurred that were not expected")
	}

	if len(formater.GetLines()) != 1 {
		t.Errorf("Got %d lines, but expected 1", len(formater.GetLines()))
	}

	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParameter = oldDepthValue
	CorrectionParameter = oldCorrectionPar
	PrintCsvHeaderFlag = oldPrintCsvHeaderFlag
	SuppressDuplicateOutPutFlag = oldSuppressDuplicateOutPutFlag
}

func TestProcessValidFilesWithDuplicateElementsDetectionDisabeled(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParameter
	DepthParameter = "track"
	oldCorrectionPar := CorrectionParameter
	CorrectionParameter = "linear"
	oldPrintCsvHeaderFlag := PrintCsvHeaderFlag
	PrintCsvHeaderFlag = false
	oldSuppressDuplicateOutPutFlag := SuppressDuplicateOutPutFlag
	SuppressDuplicateOutPutFlag = false

	formater := csvbl.NewCsvOutputFormater(";", false)
	iFormater := gpsabl.OutputFormater(formater)

	fileStrs := []string{testhelper.GetValidGPX("15.gpx")}
	var files []gpsabl.InputFile
	for _, file := range fileStrs {
		files = append(files, *gpsabl.NewInputFileWithPath(file))
	}
	successCount := processFiles(files, iFormater)
	if successCount != 1 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	if ErrorsHandled == true {
		t.Errorf("Errors occurred that were not expected")
	}

	if len(formater.GetLines()) != 2 {
		t.Errorf("Got %d lines, but expected 2", len(formater.GetLines()))
	}

	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParameter = oldDepthValue
	CorrectionParameter = oldCorrectionPar
	PrintCsvHeaderFlag = oldPrintCsvHeaderFlag
	SuppressDuplicateOutPutFlag = oldSuppressDuplicateOutPutFlag
}

func TestProcessFilesDifferentCorrection(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParameter
	DepthParameter = "file"
	oldCorrectionPAr := CorrectionParameter
	oldPrintCsvHeaderFlag := PrintCsvHeaderFlag
	PrintCsvHeaderFlag = false
	fileStrs := []string{testhelper.GetValidGPX("01.gpx"), testhelper.GetValidGPX("12.gpx")}
	var files []gpsabl.InputFile
	for _, file := range fileStrs {
		files = append(files, *gpsabl.NewInputFileWithPath(file))
	}
	CorrectionParameter = "none"
	formater1 := csvbl.NewCsvOutputFormater(";", false)
	iFormater1 := gpsabl.OutputFormater(formater1)
	successCount1 := processFiles(files, iFormater1)
	if successCount1 != 2 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	CorrectionParameter = "linear"
	formater2 := csvbl.NewCsvOutputFormater(";", false)
	iFormater2 := gpsabl.OutputFormater(formater2)
	successCount2 := processFiles(files, iFormater2)
	if successCount2 != 2 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	if len(formater2.GetLines()) != len(formater1.GetLines()) {
		t.Errorf("The formater have a different amount of lines")
	}

	if formater2.GetLines()[1] == formater1.GetLines()[1] {
		t.Errorf("Both formaters return the same values")
	}

	if ErrorsHandled == true {
		t.Errorf("Errors occurred that were not expected")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParameter = oldDepthValue
	CorrectionParameter = oldCorrectionPAr
	PrintCsvHeaderFlag = oldPrintCsvHeaderFlag
}

func TestProcessFilesDifferentMovingSpeed(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParameter
	DepthParameter = "file"
	oldCorrectionPAr := CorrectionParameter
	CorrectionParameter = "none"
	oldPrintCsvHeaderFlag := PrintCsvHeaderFlag
	PrintCsvHeaderFlag = false
	oldMinMovingSpeed := MinimalMovingSpeedParameter
	fileStrs := []string{testhelper.GetValidGPX("02.gpx"), testhelper.GetValidGPX("12.gpx")}
	var files []gpsabl.InputFile
	for _, file := range fileStrs {
		files = append(files, *gpsabl.NewInputFileWithPath(file))
	}
	MinimalMovingSpeedParameter = 0.1
	formater1 := csvbl.NewCsvOutputFormater(";", false)
	iFormater1 := gpsabl.OutputFormater(formater1)
	successCount1 := processFiles(files, iFormater1)
	if successCount1 != 2 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	MinimalMovingSpeedParameter = 0.9
	formater2 := csvbl.NewCsvOutputFormater(";", false)
	iFormater2 := gpsabl.OutputFormater(formater2)
	successCount2 := processFiles(files, iFormater2)
	if successCount2 != 2 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	if len(formater2.GetLines()) != len(formater1.GetLines()) {
		t.Errorf("The formater have a different amount of lines")
	}

	if formater2.GetLines()[0] == formater1.GetLines()[0] {
		t.Errorf("Both formaters return the same values for %s", formater1.GetLines()[0])
	}

	if formater2.GetLines()[1] == formater1.GetLines()[1] {
		t.Errorf("Both formaters return the same values for %s", formater1.GetLines()[1])
	}

	if ErrorsHandled == true {
		t.Errorf("Errors occurred that were not expected")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParameter = oldDepthValue
	CorrectionParameter = oldCorrectionPAr
	MinimalMovingSpeedParameter = oldMinMovingSpeed
	PrintCsvHeaderFlag = oldPrintCsvHeaderFlag
}

func TestProcessFilesDifferentStepHight(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParameter
	DepthParameter = "file"
	oldCorrectionPAr := CorrectionParameter
	CorrectionParameter = "steps"
	oldPrintCsvHeaderFlag := PrintCsvHeaderFlag
	PrintCsvHeaderFlag = false
	oldStepHight := MinimalStepHightParameter
	fileStrs := []string{testhelper.GetValidGPX("02.gpx"), testhelper.GetValidGPX("12.gpx")}
	var files []gpsabl.InputFile
	for _, file := range fileStrs {
		files = append(files, *gpsabl.NewInputFileWithPath(file))
	}
	MinimalStepHightParameter = 20.0
	formater1 := csvbl.NewCsvOutputFormater(";", false)
	iFormater1 := gpsabl.OutputFormater(formater1)
	successCount1 := processFiles(files, iFormater1)
	if successCount1 != 2 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	MinimalStepHightParameter = 0.5
	formater2 := csvbl.NewCsvOutputFormater(";", false)
	iFormater2 := gpsabl.OutputFormater(formater2)
	successCount2 := processFiles(files, iFormater2)
	if successCount2 != 2 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	if len(formater2.GetLines()) != len(formater1.GetLines()) {
		t.Errorf("The formater have a different amount of lines")
	}

	if formater2.GetLines()[0] == formater1.GetLines()[0] {
		t.Errorf("Both formaters return the same values for %s", formater1.GetLines()[0])
	}

	if formater2.GetLines()[1] == formater1.GetLines()[1] {
		t.Errorf("Both formaters return the same values for %s", formater1.GetLines()[1])
	}

	if ErrorsHandled == true {
		t.Errorf("Errors occurred that were not expected")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParameter = oldDepthValue
	CorrectionParameter = oldCorrectionPAr
	MinimalStepHightParameter = oldStepHight
	PrintCsvHeaderFlag = oldPrintCsvHeaderFlag
}

func TestProcessFilesStepHightEffectsOther(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParameter
	DepthParameter = "file"
	oldCorrectionPAr := CorrectionParameter
	CorrectionParameter = "none"
	oldPrintCsvHeaderFlag := PrintCsvHeaderFlag
	PrintCsvHeaderFlag = false
	oldStepHight := MinimalStepHightParameter
	fileStrs := []string{testhelper.GetValidGPX("02.gpx"), testhelper.GetValidGPX("12.gpx")}
	var files []gpsabl.InputFile
	for _, file := range fileStrs {
		files = append(files, *gpsabl.NewInputFileWithPath(file))
	}

	MinimalStepHightParameter = 20.0
	formater1 := csvbl.NewCsvOutputFormater(";", false)
	iFormater1 := gpsabl.OutputFormater(formater1)
	successCount1 := processFiles(files, iFormater1)
	if successCount1 != 2 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	MinimalStepHightParameter = 0.5
	formater2 := csvbl.NewCsvOutputFormater(";", false)
	iFormater2 := gpsabl.OutputFormater(formater2)
	successCount2 := processFiles(files, iFormater2)
	if successCount2 != 2 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	if len(formater2.GetLines()) != len(formater1.GetLines()) {
		t.Errorf("The formater have a different amount of lines")
	}

	if formater2.GetLines()[0] != formater1.GetLines()[0] {
		t.Errorf("Both formaters return not the same values: %s", formater1.GetLines()[0])
	}

	if formater2.GetLines()[1] != formater1.GetLines()[1] {
		t.Errorf("Both formaters return the not same values: %s", formater1.GetLines()[1])
	}

	if ErrorsHandled == true {
		t.Errorf("Errors occurred that were not expected")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParameter = oldDepthValue
	CorrectionParameter = oldCorrectionPAr
	MinimalStepHightParameter = oldStepHight
	PrintCsvHeaderFlag = oldPrintCsvHeaderFlag
}

func TestProcessMixedFiles(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParameter
	DepthParameter = "file"
	oldCorrectionPAr := CorrectionParameter
	CorrectionParameter = "linear"

	formater := csvbl.NewCsvOutputFormater(";", false)
	iFormater := gpsabl.OutputFormater(formater)

	fileStrs := []string{testhelper.GetInvalidGPX("01.gpx"), testhelper.GetValidGPX("01.gpx"), testhelper.GetInvalidGPX("02.gpx"), testhelper.GetInvalidGPX("03.gpx")}
	var files []gpsabl.InputFile
	for _, file := range fileStrs {
		files = append(files, *gpsabl.NewInputFileWithPath(file))
	}
	successCount := processFiles(files, iFormater)
	if successCount != 1 {
		t.Errorf("More or less than two files were processed with error - expected exactly two of them")
	}

	if ErrorsHandled == false {
		t.Errorf("No errors occured where errors were expected")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue

	SkipErrorExitFlag = oldFlagValue
	DepthParameter = oldDepthValue
	CorrectionParameter = oldCorrectionPAr
}

func TestProcessInValidFiles(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParameter
	DepthParameter = "file"
	oldCorrectionPAr := CorrectionParameter
	CorrectionParameter = "linear"

	formater := csvbl.NewCsvOutputFormater(";", false)
	iFormater := gpsabl.OutputFormater(formater)

	fileStrs := []string{testhelper.GetInvalidGPX("01.gpx"), testhelper.GetInvalidGPX("02.gpx")}
	var files []gpsabl.InputFile
	for _, file := range fileStrs {
		files = append(files, *gpsabl.NewInputFileWithPath(file))
	}
	successCount := processFiles(files, iFormater)
	if successCount != 0 {
		t.Errorf("Not all files were processed with error as expected")
	}

	if ErrorsHandled == false {
		t.Errorf("No errors occured were errors were expected")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParameter = oldDepthValue
	CorrectionParameter = oldCorrectionPAr
}

func TestGetOutPutStream_StdOut(t *testing.T) {
	ErrorsHandled = false
	oldOutFileParameter := OutFileParameter
	OutFileParameter = ""
	str := getOutPutStream()

	switch os.File(*str) {
	case *os.Stdout:
		fmt.Println("OK")
	default:
		t.Errorf("Did not receive the expected stream")
	}
	ErrorsHandled = false
	OutFileParameter = oldOutFileParameter

}

func TestGetOutPutStream_AFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skip this test on windows")
	} else {
		ErrorsHandled = false
		oldOutFileParameter := OutFileParameter
		var big big.Int
		big.SetInt64(10000)
		max, _ := rand.Int(rand.Reader, &big)
		filePath := filepath.Join(testhelper.GetProjectRoot(), "testdata", fmt.Sprintf("test-out-%d.csv", max))
		OutFileParameter = filePath

		// Make sure the tests can be executed in parallel
		outFileMux.Lock()
		defer outFileMux.Unlock()

		if outFileExists(filePath) {
			err := os.Remove(filePath)
			if err != nil {
				t.Errorf("Test setup was not able to delete %s. Error was: %s", filePath, err.Error())
			}
		}
		str := getOutPutStream()
		str.Close()

		if outFileExists(filePath) {
			err := os.Remove(filePath)
			if err != nil {
				t.Errorf("Test cleanup was not able to delete %s. Error was: %s", filePath, err.Error())
			}
		} else {
			t.Errorf("The outfile \"%s\" was not created as expected", filePath)
		}

		if str.Name() != filePath {
			t.Errorf("The Outstream is not for %s, expected %s ", str.Name(), filePath)
		}

		if ErrorsHandled == true {
			t.Errorf("Got an error where no error was expected")
		}

		ErrorsHandled = false
		OutFileParameter = oldOutFileParameter
	}

}

func TestGetOutPutStream_AExistingFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skip this test on windows")
	} else {
		ErrorsHandled = false
		oldOutFileParameter := OutFileParameter
		var big big.Int
		big.SetInt64(10000)
		max, _ := rand.Int(rand.Reader, &big)
		filePath := filepath.Join(testhelper.GetProjectRoot(), "testdata", fmt.Sprintf("test-out-%d.csv", max))
		OutFileParameter = filePath

		// Make sure the tests can be executed in parallel
		outFileMux.Lock()
		defer outFileMux.Unlock()
		if !outFileExists(filePath) {
			file, _ := os.Create(filePath)
			file.Close()
		}

		if !outFileExists(filePath) {
			t.Errorf("Error while creating out file at test setup")
		}
		str := getOutPutStream()

		str.Sync()
		closeErr := str.Close()
		if closeErr != nil {
			t.Errorf("Test cleanup was not able to close %s. Error was: %s", filePath, closeErr.Error())
		}

		if outFileExists(filePath) {
			err := os.Remove(filePath)
			if err != nil {
				t.Errorf("Test cleanup was not able to delete %s. Error was: %s", filePath, err.Error())
			}
		} else {
			t.Errorf("The outfile \"%s\" was not created, as expected", filePath)
		}

		if str.Name() != filePath {
			t.Errorf("The Outstream is not for %s, expected %s ", str.Name(), filePath)
		}

		if ErrorsHandled == true {
			t.Errorf("Got an error where no error was expected")
		}

		ErrorsHandled = false
		OutFileParameter = oldOutFileParameter
	}
}

func TestGetOutPutFormaterStdOut(t *testing.T) {
	frt := getOutPutFormater(*os.Stdout)

	switch frt.(type) {
	case *csvbl.CsvOutputFormater:
		fmt.Println("OK")
	default:
		t.Errorf("Did not receive the expected formater")
	}
}

func TestGetOutPutFormaterCSV(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skip this test on windows")
	}
	var out *os.File
	filePath := filepath.Join(testhelper.GetProjectRoot(), "testdata", fmt.Sprintf("test-out.csv"))
	out, errCreate := os.Create(filePath)
	if errCreate != nil {
		t.Errorf("%s", errCreate)
	}
	out, errOpen := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
	if errOpen != nil {
		t.Errorf("%s", errOpen)
	}

	frt := getOutPutFormater(*out)
	out.Sync()
	out.Close()
	if outFileExists(filePath) {
		err := os.Remove(filePath)
		if err != nil {
			t.Errorf("Test cleanup was not able to delete %s. Error was: %s", filePath, err.Error())
		}
	} else {
		t.Errorf("The outfile \"%s\" was not created, as expected", filePath)
		return
	}
	switch frt.(type) {
	case *csvbl.CsvOutputFormater:
		fmt.Println("OK")
	default:
		t.Errorf("Did not receive the expected formater")
	}
}

func TestGetOutPutFormaterCSVStdOut(t *testing.T) {
	oldStdOutFormatParameter := StdOutFormatParameter
	StdOutFormatParameter = stdOutFormatParameterValues[0]
	frt := getOutPutFormater(*os.Stdout)

	switch frt.(type) {
	case *csvbl.CsvOutputFormater:
		fmt.Println("OK")
	default:
		t.Errorf("Did not receive the expected formater")
	}
	StdOutFormatParameter = oldStdOutFormatParameter
}

func TestGetOutPutFormaterJSONStdOut(t *testing.T) {
	oldStdOutFormatParameter := StdOutFormatParameter
	StdOutFormatParameter = stdOutFormatParameterValues[1]
	frt := getOutPutFormater(*os.Stdout)

	switch frt.(type) {
	case *jsonbl.JSONOutputFormater:
		fmt.Println("OK")
	default:
		t.Errorf("Did not receive the expected formater")
	}
	StdOutFormatParameter = oldStdOutFormatParameter
}

func TestGetOutPutFormaterJSON(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skip this test on windows")
	}
	var out *os.File
	filePath := filepath.Join(testhelper.GetProjectRoot(), "testdata", fmt.Sprintf("test-out.json"))
	out, errCreate := os.Create(filePath)
	if errCreate != nil {
		t.Errorf("%s", errCreate)
	}
	out, errOpen := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
	if errOpen != nil {
		t.Errorf("%s", errOpen)
	}

	frt := getOutPutFormater(*out)
	out.Sync()
	out.Close()
	if outFileExists(filePath) {
		err := os.Remove(filePath)
		if err != nil {
			t.Errorf("Test cleanup was not able to delete %s. Error was: %s", filePath, err.Error())
		}
	} else {
		t.Errorf("The outfile \"%s\" was not created, as expected", filePath)
		return
	}
	switch frt.(type) {
	case *jsonbl.JSONOutputFormater:
		fmt.Println("OK")
	default:
		t.Errorf("Did not receive the expected formater")
	}
}

func TestProcessInputStreamWithFileNames(t *testing.T) {
	file1, file2, read, err := getTwoValidInputFilePathStream()
	if err != nil {
		t.Fatal(err)
	}

	stdin := os.Stdin
	// Restore stdin right after the test.
	defer func() { os.Stdin = stdin }()
	os.Stdin = read

	inFiles := processInputStream()

	if len(inFiles) != 2 {
		t.Errorf("The number %d of input files is not expected", len(inFiles))
	}
	if inFiles[0].Name != file1 {
		t.Errorf("processInputStream does not return the expected string")
	}
	if inFiles[1].Name != file2 {
		t.Errorf("processInputStream does not return the expected string")
	}

	if inFiles[1].Type != gpsabl.FilePath {
		t.Errorf("The type is %s, but %s is expected", inFiles[1].Type, gpsabl.FilePath)
	}
}

func TestInputStreamBufferWithTwoGPXFileContent(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skip this test on windows")
	}
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParameter
	DepthParameter = "file"
	oldCorrectionPAr := CorrectionParameter
	CorrectionParameter = "linear"
	read, errGet := getValidInputGPXContentStream()
	if errGet != nil {
		t.Fatal(errGet)
	}
	formater := csvbl.NewCsvOutputFormater(";", false)
	stdin := os.Stdin
	// Restore stdin right after the test.
	defer func() { os.Stdin = stdin }()
	os.Stdin = read

	inFiles := processInputStream()
	if inFiles[1].Type != GpxBuffer {
		t.Errorf("The type is %s, but %s is expected", inFiles[1].Type, GpxBuffer)
	}
	if len(inFiles) != 2 {
		t.Errorf("%d files expected, but got %d", 2, len(inFiles))
	}

	successCount := processFiles(inFiles, gpsabl.OutputFormater(formater))

	if successCount != len(inFiles) {
		t.Errorf("only %d files processed successfull, but %d should", successCount, len(inFiles))
	}

	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParameter = oldDepthValue
	CorrectionParameter = oldCorrectionPAr
}
