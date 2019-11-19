package gpsabl

import (
	"strings"
	"testing"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestFormatOutPutWithHeader(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	ret := formater.FormatOutPut(trackFile, true)

	if len(ret) != 2 {
		t.Errorf("The output has not the expected number of files")
	}

	if strings.Count(ret[0], ";") != strings.Count(ret[1], ";") {
		t.Errorf("The Number of semicolons is not the same in each line")
	}

	if strings.Contains(ret[1], "0.0200") == false {
		t.Errorf("The output does not contian the distance as expected. It is: %s", ret[1])
	}
}

func TestFormatOutPutWithOutHeader(t *testing.T) {
	formater := NewCsvOutputFormater(";")
	trackFile := getSimpleTrackFile()

	ret := formater.FormatOutPut(trackFile, false)

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
