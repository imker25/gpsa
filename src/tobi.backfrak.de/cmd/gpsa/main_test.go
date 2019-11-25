package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"tobi.backfrak.de/internal/gpsabl"
	"tobi.backfrak.de/internal/testhelper"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// Use this Mux to sync write access to the testdata/test-out.csv file
var outFileMux sync.Mutex

var ComandlineOptionsHandled bool

func TestHandleComandlineOptions(t *testing.T) {
	if !ComandlineOptionsHandled {
		handleComandlineOptions()
		ComandlineOptionsHandled = true
	}

	if HelpFlag == true {
		t.Errorf("The HelpFlag is true, but should not")
	}

	if VerboseFlag == true {
		t.Errorf("The VerboseFlag is true, but should not")
	}

	if SkipErrorExitFlag == true {
		t.Errorf("The SkipErrorExitFlag is true, but should not")
	}

	if PrintVersionFlag == true {
		t.Errorf("The PrintVersionFlag is true, but should not")
	}

	if PrintLicenseFlag == true {
		t.Errorf("The PrintLicenseFlag is true, but should not")
	}

	if DepthParametr != "track" {
		t.Errorf("The DepthParametr is \"%s\" but \"track\" was expected", DepthParametr)
	}

	if OutFileParameter != "" {
		t.Errorf("The DepthParametr is \"%s\" but \"\" was expected", OutFileParameter)
	}

	if DontPanicFlag == false {
		t.Errorf("The DontPanicFlag is false, but should not")
	}

	if PrintCsvHeaderFlag == false {
		t.Errorf("The PrintCsvHeaderFlag is false, but should not")
	}
}

func TestCostumHelpMessage(t *testing.T) {
	if !ComandlineOptionsHandled {
		handleComandlineOptions()
		ComandlineOptionsHandled = true
	}
	costumHelpMessage()
}

func TestHandleErrorNil(t *testing.T) {
	if HandleError(nil, "my/path", true, true) == true {
		t.Errorf("HandleError reutrns true, when nil error was given")
	}

}

func TestHandleErrorNotNil(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	if HandleError(newUnKnownFileTypeError("my/path"), "my/path", true, true) == false {
		t.Errorf("HandleError reutrns false, when error was given")
	}

	if ErrorsHandled == false {
		t.Errorf("ErrorsHandled should be true, after a error was handeled")
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

func TestGetReaderUnkonwnFile(t *testing.T) {
	reader, err := getReader("/some/track.txt")

	if err == nil {
		t.Errorf("Got no error when try to get a reader for a txt file.")
	}

	if reader != nil {
		t.Errorf("The reader we got was not nil")
	}
}

func TestProcessValideFiles(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParametr
	DepthParametr = "file"
	formater := gpsabl.NewCsvOutputFormater(";")
	iFormater := gpsabl.OutputFormater(formater)

	files := []string{testhelper.GetValideGPX("01.gpx"), testhelper.GetValideGPX("02.gpx")}
	successCount := processFiles(files, iFormater)
	if successCount != 2 {
		t.Errorf("Not all files was proccess successfull as expected")
	}

	if ErrorsHandled == true {
		t.Errorf("Errors occured, but should not")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParametr = oldDepthValue
}

func TestProcessMixedFiles(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParametr
	DepthParametr = "file"

	formater := gpsabl.NewCsvOutputFormater(";")
	iFormater := gpsabl.OutputFormater(formater)

	files := []string{testhelper.GetUnValideGPX("01.gpx"), testhelper.GetValideGPX("01.gpx"), testhelper.GetUnValideGPX("02.gpx")}
	successCount := processFiles(files, iFormater)
	if successCount != 1 {
		t.Errorf("Not two files was proccess with error as expected")
	}

	if ErrorsHandled == false {
		t.Errorf("No errors occured, but should")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue

	SkipErrorExitFlag = oldFlagValue
	DepthParametr = oldDepthValue
}

func TestProcessUnValideFiles(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParametr
	DepthParametr = "file"

	formater := gpsabl.NewCsvOutputFormater(";")
	iFormater := gpsabl.OutputFormater(formater)

	files := []string{testhelper.GetUnValideGPX("01.gpx"), testhelper.GetUnValideGPX("02.gpx")}
	successCount := processFiles(files, iFormater)
	if successCount != 0 {
		t.Errorf("Not all files was proccess with error as expected")
	}

	if ErrorsHandled == false {
		t.Errorf("No errors occured, but should")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParametr = oldDepthValue
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
		t.Errorf("Got not the expected stream")
	}
	ErrorsHandled = false
	OutFileParameter = oldOutFileParameter

}

func TestGetOutPutStream_AFile(t *testing.T) {
	ErrorsHandled = false
	oldOutFileParameter := OutFileParameter
	filePath := filepath.Join(testhelper.GetProjectRoot(), "testdata", "test-out.csv")
	OutFileParameter = filePath

	// Make sure the tests can be executed in parallel
	outFileMux.Lock()
	defer outFileMux.Unlock()

	str := getOutPutStream()
	str.Close()

	if outFileExists(filePath) {
		os.Remove(filePath)
	} else {
		t.Errorf("The outfile \"%s\" was not created, as expected", filePath)
	}

	if str.Name() != filePath {
		t.Errorf("The Outstream is not for %s, expected %s ", str.Name(), filePath)
	}

	if ErrorsHandled == true {
		t.Errorf("Got an error, but expected none")
	}

	ErrorsHandled = false
	OutFileParameter = oldOutFileParameter

}

func TestGetOutPutStream_AExistingFile(t *testing.T) {
	ErrorsHandled = false
	oldOutFileParameter := OutFileParameter
	filePath := filepath.Join(testhelper.GetProjectRoot(), "testdata", "test-out.csv")
	OutFileParameter = filePath

	// Make sure the tests can be executed in parallel
	outFileMux.Lock()
	defer outFileMux.Unlock()

	os.Create(filePath)
	if !outFileExists(filePath) {
		t.Errorf("Error while creating out file at test setup")
	}
	str := getOutPutStream()
	str.Close()

	if outFileExists(filePath) {
		os.Remove(filePath)
	} else {
		t.Errorf("The outfile \"%s\" was not created, as expected", filePath)
	}

	if str.Name() != filePath {
		t.Errorf("The Outstream is not for %s, expected %s ", str.Name(), filePath)
	}

	if ErrorsHandled == true {
		t.Errorf("Got an error, but expected none")
	}

	ErrorsHandled = false
	OutFileParameter = oldOutFileParameter

}
func TestGetOutPutFormater(t *testing.T) {
	frt := getOutPutFormater()

	switch frt.(type) {
	case *gpsabl.CsvOutputFormater:
		fmt.Println("OK")
	default:
		t.Errorf("Got not the expected formater")
	}
}
