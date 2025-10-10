package gpsabl

import (
	"os"
	"strings"
	"sync"
	"testing"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestGetOutputFormaterEmptyList(t *testing.T) {
	var formaters []OutputFormater

	res, formater := GetOutputFormater(formaters, *os.Stdout, OutputFormaterType("myType"))

	if res == true {
		t.Errorf("GetOutputFormater find formater in an empty list")
	}

	if formater != nil {
		t.Errorf("GetOutputFormater find formater in an empty list")
	}
}

func TestGetOutputFormater(t *testing.T) {
	formaters := []OutputFormater{&formaterMock{}}

	res, formater := GetOutputFormater(formaters, *os.Stdout, OutputFormaterType("myType"))

	if res == true {
		t.Errorf("GetOutputFormater find formater for myType")
	}

	if formater != nil {
		t.Errorf("GetOutputFormater find formater for myType")
	}

	res, formater = GetOutputFormater(formaters, *os.Stdout, OutputFormaterType("ABS"))

	if res == false {
		t.Errorf("GetOutputFormater find no formater for ABS")
	}

	if formater == nil {
		t.Errorf("GetOutputFormater find no formater for ABS")
	}

	out, errCreate := os.Create("out.abs")
	if errCreate != nil {
		t.Fatalf("Can not create the outfile")
	}
	res, formater = GetOutputFormater(formaters, *out, OutputFormaterType("ABS"))

	if res == false {
		t.Errorf("GetOutputFormater find no formater for ABS")
	}

	if formater == nil {
		t.Errorf("GetOutputFormater find no formater for ABS")
	}

	out, errCreate = os.Create("out.abs")
	if errCreate != nil {
		t.Fatalf("Can not create the outfile")
	}
	res, formater = GetOutputFormater(formaters, *out, OutputFormaterType("CSV"))

	if res == false {
		t.Errorf("GetOutputFormater find no formater for ABS")
	}

	if formater == nil {
		t.Errorf("GetOutputFormater find no formater for ABS")
	}

	out, errCreate = os.Create("out.abc")
	if errCreate != nil {
		t.Fatalf("Can not create the outfile")
	}
	res, formater = GetOutputFormater(formaters, *out, OutputFormaterType("CSV"))

	if res == true {
		t.Errorf("GetOutputFormater find no formater for ABS")
	}

	if formater != nil {
		t.Errorf("GetOutputFormater find no formater for ABS")
	}
}

type formaterMock struct {
	WrittenEntriesCount int
	lineBuffer          []OutputLine
	mux                 sync.Mutex
}

func (formater *formaterMock) NewOutputFormater() OutputFormater {
	ret := formaterMock{}
	ret.WrittenEntriesCount = -1
	ret.lineBuffer = []OutputLine{}

	return OutputFormater(&ret)
}

func (formater *formaterMock) AddOutPut(trackFile TrackFile, depth DepthArg, filterDuplicate bool) error {
	return nil
}

func (formater *formaterMock) WriteOutput(outFile *os.File, summary SummaryArg) error {
	return nil
}

func (formater *formaterMock) CheckFileExtension(filePath string) bool {
	if strings.HasSuffix(strings.ToLower(filePath), ".abs") {
		return true
	}

	return false
}

func (formater *formaterMock) CheckOutputFormaterType(formaterType OutputFormaterType) bool {
	if formaterType == OutputFormaterType("ABS") {
		return true
	}

	return false
}

func (formater *formaterMock) GetFileExtensions() []string {
	return []string{".abs"}
}

func (formater *formaterMock) GetOutputFormaterTypes() []OutputFormaterType {
	return []OutputFormaterType{"ABS"}
}

func (formater *formaterMock) GetTextOutputFormater() TextOutputFormater {
	return nil
}

func (formater *formaterMock) GetNumberOfOutputEntries() int {
	return formater.WrittenEntriesCount
}
