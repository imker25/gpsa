package gpsabl

import (
	"strings"
	"testing"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestNewCsvOutputFormater(t *testing.T) {
	sut := NewCsvOutputFormater(";")

	if sut.Seperator != ";" {
		t.Errorf("The Seperator was \"%s\", but \";\" was expected", sut.Seperator)
	}

	if len(sut.ValideDepthArgs) != 2 {
		t.Errorf("The ValideDepthArgs array does not contain the expeced number of values")
	}
}

func TestFormatOutPutWithHeader(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	ret := formater.FormatOutPut(trackFile, true, "file")

	if len(ret) != 2 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Count(ret[0], ";") != strings.Count(ret[1], ";") {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(ret[1], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", ret[1])
	}

	if strings.Contains(ret[1], trackFile.FilePath) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[1])
	}
}

func TestFormatOutPutWithHeaderAndSetName(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()
	trackFile.Name = "My Track File"

	ret := formater.FormatOutPut(trackFile, true, "file")

	if len(ret) != 2 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Count(ret[0], ";") != strings.Count(ret[1], ";") {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(ret[1], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", ret[1])
	}

	if strings.Contains(ret[1], trackFile.Name) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[1])
	}

	if strings.Contains(ret[1], trackFile.FilePath) == true {
		t.Errorf("The output does contian the FilePath but should not. It is: %s", ret[1])
	}
}

func TestFormatOutPutWithOutHeader(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	ret := formater.FormatOutPut(trackFile, false, "file")

	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Count(ret[0], ";") != 5 {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(ret[0], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", ret[1])
	}
}

func TestFormatOutPutWithOutHeaderTrackDepth(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	ret := formater.FormatOutPut(trackFile, false, "track")

	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Contains(ret[0], trackFile.FilePath) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], "#1;") == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderTrackDepthSetTrackName(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()
	trackFile.Tracks[0].Name = "My Track"

	ret := formater.FormatOutPut(trackFile, false, "track")

	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Contains(ret[0], trackFile.FilePath) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], trackFile.Tracks[0].Name) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], "#1;") == true {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderTrackDepthSetTrackFileNameSetTrackName(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()
	trackFile.Tracks[0].Name = "My Track"
	trackFile.Name = "My track file"

	ret := formater.FormatOutPut(trackFile, false, "track")

	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Contains(ret[0], trackFile.FilePath) == true {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], trackFile.Name) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], trackFile.Tracks[0].Name) == false {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}

	if strings.Contains(ret[0], "#1;") == true {
		t.Errorf("The output does not contian the name as expected. It is: %s", ret[0])
	}
}

func TestFormatOutPutWithOutHeaderUnValideDepth(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()
	ret := formater.FormatOutPut(trackFile, false, "abc")

	if len(ret) != 1 {
		t.Errorf("The output has not the expected number of lines")
	}

	if strings.HasPrefix(ret[0], "Error:") == false {
		t.Errorf("The line does not start with \"Error\" as expected")
	}
}
