package gpsabl

import (
	"testing"
)

// Copyright 2025 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestNewMinStartTimeFilterNoFilterString(t *testing.T) {

	_, err := NewMinStartTimeFilter("bla")

	if err == nil {
		t.Errorf("The MinStartTimeFilter can parse the string \"bla\"")
	}

}

func TestNewMinStartTimeFilterValidDateTimeFilterString(t *testing.T) {

	filterString := "2024-12-13 20:51:42"
	sut, err := NewMinStartTimeFilter(filterString)

	if err != nil {
		t.Errorf("The MinStartTimeFilter can not parse the string \"%s\"", filterString)
	}

	if sut.GetFilterText() != filterString {
		t.Errorf("The MinStartTimeFilter.GetFilterText() returns \"%s\" but should return \"%s\"", sut.GetFilterText(), filterString)
	}

	if sut.MinStartTime.Unix() != 1734123102 {
		t.Errorf("The MinStartTimeFilter.MinStartTime is \"%d\" but not the expected \"1734123102\"", sut.MinStartTime.Unix())
	}

}

func TestNewMinStartTimeFilterValidDateTimeNoSecondsFilterString(t *testing.T) {

	filterString := "2024-12-13 20:51"
	sut, err := NewMinStartTimeFilter(filterString)

	if err != nil {
		t.Errorf("The MinStartTimeFilter can not parse the string \"%s\"", filterString)
	}

	if sut.GetFilterText() != filterString {
		t.Errorf("The MinStartTimeFilter.GetFilterText() returns \"%s\" but should return \"%s\"", sut.GetFilterText(), filterString)
	}

	if sut.MinStartTime.Unix() != 1734123060 {
		t.Errorf("The MinStartTimeFilter.MinStartTime is \"%d\" but not the expected \"1734123060\"", sut.MinStartTime.Unix())
	}

}

func TestNewMinStartTimeFilterValidDateFilterString(t *testing.T) {

	filterString := "2024-12-13"
	sut, err := NewMinStartTimeFilter(filterString)

	if err != nil {
		t.Errorf("The MinStartTimeFilter can not parse the string \"%s\"", filterString)
	}

	if sut.GetFilterText() != filterString {
		t.Errorf("The MinStartTimeFilter.GetFilterText() returns \"%s\" but should return \"%s\"", sut.GetFilterText(), filterString)
	}

	if sut.MinStartTime.Unix() != 1734048000 {
		t.Errorf("The MinStartTimeFilter.MinStartTime is \"%d\" but not the expected \"1734048000\"", sut.MinStartTime.Unix())
	}

}
func TestMinStartTimeFilterFilterMatchingDate1(t *testing.T) {
	filterString := "2014-08-21"
	track := getSimpleTrackWithTime()
	sut, _ := NewMinStartTimeFilter(filterString)

	if !sut.Filter(track.TrackSummary) {
		t.Errorf("The tracks start date \"%s\" is newer than the min date \"%s\"", track.TrackSummary.StartTime, sut.MinStartTime)
	}

}

func TestMinStartTimeFilterFilterMatchingDate2(t *testing.T) {
	filterString := "2014-08-22"
	track := getSimpleTrackWithTime()
	sut, _ := NewMinStartTimeFilter(filterString)

	if !sut.Filter(track.TrackSummary) {
		t.Errorf("The tracks start date \"%s\" is newer than the min date \"%s\"", track.TrackSummary.StartTime, sut.MinStartTime)
	}

}

func TestMinStartTimeFilterFilterNotMatchingDate1(t *testing.T) {
	filterString := "2014-08-23"
	track := getSimpleTrackWithTime()
	sut, _ := NewMinStartTimeFilter(filterString)

	if sut.Filter(track.TrackSummary) {
		t.Errorf("The tracks start date \"%s\" is newer than the min date \"%s\"", track.TrackSummary.StartTime, sut.MinStartTime)
	}

}

func TestMinStartTimeFilterFilterNotValidDateTime(t *testing.T) {
	filterString := "2014-08-23"
	track := getSimpleTrack()
	sut, _ := NewMinStartTimeFilter(filterString)

	if !sut.Filter(track.TrackSummary) {
		t.Errorf("The tracks start date \"%s\" is newer than the min date \"%s\"", track.TrackSummary.StartTime, sut.MinStartTime)
	}

}
