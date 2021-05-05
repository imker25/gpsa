package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestGetValidDepthArgsString(t *testing.T) {
	str := GetValidDepthArgsString()

	if strings.Contains(str, "blabla") {
		t.Errorf("The GetValidDepthArgsString contains \"blabla\"")
	}

	if !strings.Contains(str, "file") {
		t.Errorf("The GetValidDepthArgsString not contains \"file\"")
	}

	if !strings.Contains(str, "track") {
		t.Errorf("The GetValidDepthArgsString not contains \"track\"")
	}

	if !strings.Contains(str, "segment") {
		t.Errorf("The GetValidDepthArgsString not contains \"segment\"")
	}

	if len(GetValidDepthArgs()) != 3 {
		t.Errorf("The ValidDepthArgs array does not contain the expected number of values")
	}
}

func TestGetValidTimeFormatsString(t *testing.T) {
	str := GetValidTimeFormatsString()
	res := fmt.Sprintf("\"%s\" \"%s\" \"%s\" ", time.UnixDate, time.RFC850, time.RFC3339)
	if str != res {
		t.Errorf("Got \"%s\", but expected \"%s \"", str, res)
	}
}

func TestCheckValidDepthArg(t *testing.T) {

	if CheckValidDepthArg("blabla") {
		t.Errorf("The CheckValidDepthArg contains \"blabla\"")
	}

	if !CheckValidDepthArg("file") {
		t.Errorf("The CheckValidDepthArg not contains \"file\"")
	}

	if !CheckValidDepthArg("track") {
		t.Errorf("The CheckValidDepthArg not contains \"track\"")
	}

	if !CheckValidDepthArg("segment") {
		t.Errorf("The CheckValidDepthArg not contains \"segment\"")
	}
}

func TestGetValidSummaryArgsString(t *testing.T) {
	str := GetValidSummaryArgsString()

	if strings.Contains(str, "blabla") {
		t.Errorf("The GetValidDepthArgsString contains \"blabla\"")
	}

	if !strings.Contains(str, "none") {
		t.Errorf("The GetValidDepthArgsString not contains \"none\"")
	}

	if !strings.Contains(str, "additional") {
		t.Errorf("The GetValidDepthArgsString not contains \"additional\"")
	}

	if !strings.Contains(str, "only") {
		t.Errorf("The GetValidDepthArgsString not contains \"only\"")
	}
}

func TestCheckValidSummaryArg(t *testing.T) {

	if CheckValidSummaryArg("blabla") {
		t.Errorf("The CheckValidDepthArg contains \"blabla\"")
	}

	if !CheckValidSummaryArg("none") {
		t.Errorf("The CheckValidDepthArg not contains \"none\"")
	}

	if !CheckValidSummaryArg("only") {
		t.Errorf("The CheckValidDepthArg not contains \"only\"")
	}

	if !CheckValidSummaryArg("additional") {
		t.Errorf("The CheckValidDepthArg not contains \"additional\"")
	}

	if len(GetValidSummaryArgs()) != 3 {
		t.Errorf("The ValidSummaryArgs array does not contain the expected number of values")
	}
}
