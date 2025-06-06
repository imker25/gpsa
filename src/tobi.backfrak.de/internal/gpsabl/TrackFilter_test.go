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

	if sut.Filter(track.TrackSummary) {
		t.Errorf("The tracks start date \"%s\" is newer than the min date \"%s\"", track.TrackSummary.StartTime, sut.MinStartTime)
	}

}

func TestNewMaxStartTimeFilterNoFilterString(t *testing.T) {

	_, err := NewMaxStartTimeFilter("bla")

	if err == nil {
		t.Errorf("The MinStartTimeFilter can parse the string \"bla\"")
	}

}

func TestNewMaxStartTimeFilterValidDateTimeFilterString(t *testing.T) {

	filterString := "2024-12-13 20:51:42"
	sut, err := NewMaxStartTimeFilter(filterString)

	if err != nil {
		t.Errorf("The MinStartTimeFilter can not parse the string \"%s\"", filterString)
	}

	if sut.GetFilterText() != filterString {
		t.Errorf("The MinStartTimeFilter.GetFilterText() returns \"%s\" but should return \"%s\"", sut.GetFilterText(), filterString)
	}

	if sut.MaxStartTime.Unix() != 1734123102 {
		t.Errorf("The MinStartTimeFilter.MinStartTime is \"%d\" but not the expected \"1734123102\"", sut.MaxStartTime.Unix())
	}

}

func TestNewMaxStartTimeFilterValidDateTimeNoSecondsFilterString(t *testing.T) {

	filterString := "2024-12-13 20:51"
	sut, err := NewMaxStartTimeFilter(filterString)

	if err != nil {
		t.Errorf("The MinStartTimeFilter can not parse the string \"%s\"", filterString)
	}

	if sut.GetFilterText() != filterString {
		t.Errorf("The MinStartTimeFilter.GetFilterText() returns \"%s\" but should return \"%s\"", sut.GetFilterText(), filterString)
	}

	if sut.MaxStartTime.Unix() != 1734123060 {
		t.Errorf("The MinStartTimeFilter.MinStartTime is \"%d\" but not the expected \"1734123060\"", sut.MaxStartTime.Unix())
	}

}

func TestNewMaxStartTimeFilterValidDateFilterString(t *testing.T) {

	filterString := "2024-12-13"
	sut, err := NewMaxStartTimeFilter(filterString)

	if err != nil {
		t.Errorf("The MinStartTimeFilter can not parse the string \"%s\"", filterString)
	}

	if sut.GetFilterText() != filterString {
		t.Errorf("The MinStartTimeFilter.GetFilterText() returns \"%s\" but should return \"%s\"", sut.GetFilterText(), filterString)
	}

	if sut.MaxStartTime.Unix() != 1734048000 {
		t.Errorf("The MinStartTimeFilter.MinStartTime is \"%d\" but not the expected \"1734048000\"", sut.MaxStartTime.Unix())
	}

}

func TestMaxStartTimeFilterFilterMatchingDate1(t *testing.T) {
	filterString := "2014-08-24"
	track := getSimpleTrackWithTime()
	sut, _ := NewMaxStartTimeFilter(filterString)

	if !sut.Filter(track.TrackSummary) {
		t.Errorf("The tracks start date \"%s\" is older than the min date \"%s\"", track.TrackSummary.StartTime, sut.MaxStartTime)
	}

}

func TestMaxStartTimeFilterFilterMatchingDate2(t *testing.T) {
	filterString := "2014-08-23"
	track := getSimpleTrackWithTime()
	sut, _ := NewMaxStartTimeFilter(filterString)

	if !sut.Filter(track.TrackSummary) {
		t.Errorf("The tracks start date \"%s\" is older than the min date \"%s\"", track.TrackSummary.StartTime, sut.MaxStartTime)
	}

}

func TestMaxStartTimeFilterFilterNotMatchingDate1(t *testing.T) {
	filterString := "2014-08-21"
	track := getSimpleTrackWithTime()
	sut, _ := NewMaxStartTimeFilter(filterString)

	if sut.Filter(track.TrackSummary) {
		t.Errorf("The tracks start date \"%s\" is older than the min date \"%s\"", track.TrackSummary.StartTime, sut.MaxStartTime)
	}

}

func TestMaxStartTimeFilterFilterNotValidDateTime(t *testing.T) {
	filterString := "2014-08-23"
	track := getSimpleTrack()
	sut, _ := NewMaxStartTimeFilter(filterString)

	if sut.Filter(track.TrackSummary) {
		t.Errorf("The tracks start date \"%s\" is newer than the min date \"%s\"", track.TrackSummary.StartTime, sut.MaxStartTime)
	}

}

func TestFilterTrackWithMinAndMaxFilterPass(t *testing.T) {
	track := getSimpleTrackWithStartTime("2014-08-23T17:19:33Z")

	minFilterString := "2014-08-22"
	maxFilterString := "2014-08-24"
	filters := []TrackFilter{}
	minFilter, _ := NewMinStartTimeFilter(minFilterString)
	filters = append(filters, &minFilter)
	maxFilter, _ := NewMaxStartTimeFilter(maxFilterString)
	filters = append(filters, &maxFilter)

	if !FilterTrack(track, filters) {
		t.Errorf("The track with StartTime \"%s\" does not pass the FilterTrack with filters %s and %s", track.StartTime.String(), minFilterString, maxFilterString)
	}
}

func TestFilterTrackWithMinAndMaxFilterFailToOld(t *testing.T) {
	track := getSimpleTrackWithStartTime("2014-08-21T17:19:33Z")

	minFilterString := "2014-08-22"
	maxFilterString := "2014-08-24"
	filters := []TrackFilter{}
	minFilter, _ := NewMinStartTimeFilter(minFilterString)
	filters = append(filters, &minFilter)
	maxFilter, _ := NewMaxStartTimeFilter(maxFilterString)
	filters = append(filters, &maxFilter)

	if FilterTrack(track, filters) {
		t.Errorf("The track with StartTime \"%s\" does pass the FilterTrack with filters %s and %s", track.StartTime.String(), minFilterString, maxFilterString)
	}
}

func TestFilterTrackWithMinAndMaxFilterFailToNew(t *testing.T) {
	track := getSimpleTrackWithStartTime("2014-08-25T17:19:33Z")

	minFilterString := "2014-08-22"
	maxFilterString := "2014-08-24"
	filters := []TrackFilter{}
	minFilter, _ := NewMinStartTimeFilter(minFilterString)
	filters = append(filters, &minFilter)
	maxFilter, _ := NewMaxStartTimeFilter(maxFilterString)
	filters = append(filters, &maxFilter)

	if FilterTrack(track, filters) {
		t.Errorf("The track with StartTime \"%s\" does pass the FilterTrack with filters %s and %s", track.StartTime.String(), minFilterString, maxFilterString)
	}
}

func TestFilterTracksAppliesMaxStartTimeFilterCorrect(t *testing.T) {
	tracks := getSimpleTrackList()

	filterString := "2014-08-23"
	filters := []TrackFilter{}
	maxFilter, _ := NewMaxStartTimeFilter(filterString)
	filters = append(filters, &maxFilter)

	filteredTracks := FilterTracks(tracks, filters)

	if len(filteredTracks) != 2 {
		t.Errorf("FilterTracks with maxFilter \"%s\" does return the following Tracks %s", filterString, printTracks(filteredTracks))
	}
}

func TestFilterTracksAppliesMinStartTimeFilterCorrect(t *testing.T) {
	tracks := getSimpleTrackList()

	filterString := "2014-08-23"
	filters := []TrackFilter{}
	minFilter, _ := NewMinStartTimeFilter(filterString)
	filters = append(filters, &minFilter)

	filteredTracks := FilterTracks(tracks, filters)

	if len(filteredTracks) != 2 {
		t.Errorf("FilterTracks with minFilter \"%s\" does return the following Tracks %s", filterString, printTracks(filteredTracks))
	}
}

func TestFilterTracksAppliesMinAndMaxStartTimeFilterCorrect(t *testing.T) {
	tracks := getSimpleTrackList()

	minFilterString := "2014-08-22"
	maxFilterString := "2014-08-24"
	filters := []TrackFilter{}
	minFilter, _ := NewMinStartTimeFilter(minFilterString)
	filters = append(filters, &minFilter)
	maxFilter, _ := NewMaxStartTimeFilter(maxFilterString)
	filters = append(filters, &maxFilter)

	filteredTracks := FilterTracks(tracks, filters)

	if len(filteredTracks) != 2 {
		t.Errorf("FilterTracks with minFilter \"%s\" and maxFilter \"%s\" does return the following Tracks %s", minFilterString, maxFilterString, printTracks(filteredTracks))
	}
}

func TestFilterTrackFilesWithOneTrackPassAll(t *testing.T) {
	file := getSimpleTrackFileWithTime()
	filters := []TrackFilter{}
	minFilterString := "2014-08-22"
	maxFilterString := "2014-08-24"
	minFilter, _ := NewMinStartTimeFilter(minFilterString)
	filters = append(filters, &minFilter)
	maxFilter, _ := NewMaxStartTimeFilter(maxFilterString)
	filters = append(filters, &maxFilter)

	filteredFile := FilterTrackFile(file, filters)
	if len(filteredFile.Tracks) != 1 {
		t.Errorf("The filteredFile contains %d tracks, but %d are expected", len(filteredFile.Tracks), 1)
	}

	if filteredFile.FilePath != file.FilePath {
		t.Errorf("The filteredFile.FilePath is \"%s\", but \"%s\" is expected", filteredFile.FilePath, file.FilePath)
	}

	if filteredFile.Distance != file.Distance {
		t.Errorf("The filteredFile.Distance is \"%f\", but \"%f\" is expected", filteredFile.Distance, file.Distance)
	}

	if filteredFile.ElevationGain != file.ElevationGain {
		t.Errorf("The filteredFile.ElevationGain is \"%f\", but \"%f\" is expected", filteredFile.ElevationGain, file.ElevationGain)
	}

	if filteredFile.NumberOfTracks != 1 {
		t.Errorf("The filteredFile.NumberOfTracks is \"%d\", but \"%d\" is expected", filteredFile.NumberOfTracks, 1)
	}
}

func TestFilterTrackFilesWithFourTracksTwoPassAll(t *testing.T) {
	file := getTrackFileWithMultipleTracks()
	filters := []TrackFilter{}
	minFilterString := "2014-08-22"
	maxFilterString := "2014-08-24"
	minFilter, _ := NewMinStartTimeFilter(minFilterString)
	filters = append(filters, &minFilter)
	maxFilter, _ := NewMaxStartTimeFilter(maxFilterString)
	filters = append(filters, &maxFilter)

	filteredFile := FilterTrackFile(file, filters)
	if len(filteredFile.Tracks) != 2 {
		t.Errorf("The filteredFile contains %d tracks, but %d are expected", len(filteredFile.Tracks), 1)
	}

	if filteredFile.FilePath != file.FilePath {
		t.Errorf("The filteredFile.FilePath is \"%s\", but \"%s\" is expected", filteredFile.FilePath, file.FilePath)
	}

	if filteredFile.Distance != file.Distance/2 {
		t.Errorf("The filteredFile.Distance is \"%f\", but \"%f\" is expected", filteredFile.Distance, file.Distance/2)
	}

	if filteredFile.ElevationGain != file.ElevationGain/2 {
		t.Errorf("The filteredFile.ElevationGain is \"%f\", but \"%f\" is expected", filteredFile.ElevationGain, file.ElevationGain/2)
	}

	if filteredFile.NumberOfTracks != 2 {
		t.Errorf("The filteredFile.NumberOfTracks is \"%d\", but \"%d\" is expected", filteredFile.NumberOfTracks, 2)
	}
}

func TestFilterTrackFilesWithFourTracksZeroPassAll(t *testing.T) {
	file := getTrackFileWithMultipleTracks()
	filters := []TrackFilter{}
	minFilterString := "2014-08-24"
	maxFilterString := "2014-08-22"
	minFilter, _ := NewMinStartTimeFilter(minFilterString)
	filters = append(filters, &minFilter)
	maxFilter, _ := NewMaxStartTimeFilter(maxFilterString)
	filters = append(filters, &maxFilter)

	filteredFile := FilterTrackFile(file, filters)
	if len(filteredFile.Tracks) != 0 {
		t.Errorf("The filteredFile contains %d tracks, but %d are expected", len(filteredFile.Tracks), 1)
	}

	if filteredFile.FilePath != file.FilePath {
		t.Errorf("The filteredFile.FilePath is \"%s\", but \"%s\" is expected", filteredFile.FilePath, file.FilePath)
	}

	if filteredFile.Distance != 0 {
		t.Errorf("The filteredFile.Distance is \"%f\", but \"%f\" is expected", filteredFile.Distance, 0.0)
	}

	if filteredFile.ElevationGain != 0 {
		t.Errorf("The filteredFile.ElevationGain is \"%f\", but \"%f\" is expected", filteredFile.ElevationGain, 0.0)
	}

	if filteredFile.NumberOfTracks != 0 {
		t.Errorf("The filteredFile.NumberOfTracks is \"%d\", but \"%d\" is expected", filteredFile.NumberOfTracks, 0)
	}
}

func printTracks(tracks []Track) string {
	ret := ""

	for _, track := range tracks {
		ret += track.TrackSummary.StartTime.String() + " "
	}

	return ret
}
