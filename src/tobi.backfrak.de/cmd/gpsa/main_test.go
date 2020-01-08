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

	if DepthParametr != "track" {
		t.Errorf("The DepthParametr is \"%s\" but \"track\" was expected", DepthParametr)
	}

	if OutFileParameter != "" {
		t.Errorf("The DepthParametr is \"%s\" but \"\" was expected", OutFileParameter)
	}

	if DontPanicFlag == false {
		t.Errorf("The DontPanicFlag is set to false but true was expected")
	}

	if PrintCsvHeaderFlag == false {
		t.Errorf("The PrintCsvHeaderFlag is set to false but true was expected")
	}

	if CorrectionParameter != "steps" {
		t.Errorf("The CorrectionParameter is \"%s\" but \"steps\" was expected", DepthParametr)
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
	oldDepthValue := DepthParametr
	DepthParametr = "file"
	oldCorrectionPAr := CorrectionParameter
	CorrectionParameter = "linear"

	formater := gpsabl.NewCsvOutputFormater(";")
	iFormater := gpsabl.OutputFormater(formater)

	files := []string{testhelper.GetValidGPX("01.gpx"), testhelper.GetValidGPX("02.gpx")}
	successCount := processFiles(files, iFormater)
	if successCount != 2 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	if ErrorsHandled == true {
		t.Errorf("Errors occured that were not expected")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParametr = oldDepthValue
	CorrectionParameter = oldCorrectionPAr
}

func TestProcessFilesDifferenCorrection(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParametr
	DepthParametr = "file"
	oldCorrectionPAr := CorrectionParameter
	files := []string{testhelper.GetValidGPX("01.gpx"), testhelper.GetValidGPX("02.gpx")}

	CorrectionParameter = "none"
	formater1 := gpsabl.NewCsvOutputFormater(";")
	iFormater1 := gpsabl.OutputFormater(formater1)
	successCount1 := processFiles(files, iFormater1)
	if successCount1 != 2 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	CorrectionParameter = "linear"
	formater2 := gpsabl.NewCsvOutputFormater(";")
	iFormater2 := gpsabl.OutputFormater(formater2)
	successCount2 := processFiles(files, iFormater2)
	if successCount2 != 2 {
		t.Errorf("Not all files were processed successfully as expected")
	}

	if len(formater2.GetLines()) != len(formater1.GetLines()) {
		t.Errorf("The formater have a different amout of lines")
	}

	if formater2.GetLines()[1] == formater1.GetLines()[1] {
		t.Errorf("Both formaters return the same values")
	}

	if ErrorsHandled == true {
		t.Errorf("Errors occured that were not expected")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParametr = oldDepthValue
	CorrectionParameter = oldCorrectionPAr
}

func TestProcessMixedFiles(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParametr
	DepthParametr = "file"
	oldCorrectionPAr := CorrectionParameter
	CorrectionParameter = "linear"

	formater := gpsabl.NewCsvOutputFormater(";")
	iFormater := gpsabl.OutputFormater(formater)

	files := []string{testhelper.GetInvalidGPX("01.gpx"), testhelper.GetValidGPX("01.gpx"), testhelper.GetInvalidGPX("02.gpx")}
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
	DepthParametr = oldDepthValue
	CorrectionParameter = oldCorrectionPAr
}

func TestProcessInValidFiles(t *testing.T) {
	ErrorsHandled = false
	oldFlagValue := SkipErrorExitFlag
	SkipErrorExitFlag = true
	oldDepthValue := DepthParametr
	DepthParametr = "file"
	oldCorrectionPAr := CorrectionParameter
	CorrectionParameter = "linear"

	formater := gpsabl.NewCsvOutputFormater(";")
	iFormater := gpsabl.OutputFormater(formater)

	files := []string{testhelper.GetInvalidGPX("01.gpx"), testhelper.GetInvalidGPX("02.gpx")}
	successCount := processFiles(files, iFormater)
	if successCount != 0 {
		t.Errorf("Not all files were processed with error as expected")
	}

	if ErrorsHandled == false {
		t.Errorf("No errors occured were errors were expected")
	}
	ErrorsHandled = false
	SkipErrorExitFlag = oldFlagValue
	DepthParametr = oldDepthValue
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

func TestGetOutPutFormater(t *testing.T) {
	frt := getOutPutFormater()

	switch frt.(type) {
	case *gpsabl.CsvOutputFormater:
		fmt.Println("OK")
	default:
		t.Errorf("Did not receive the expected formater")
	}
}
