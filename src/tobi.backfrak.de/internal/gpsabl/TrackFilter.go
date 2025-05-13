package gpsabl

// Copyright 2025 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import "time"

// TrackFilter - An Interface for all 'classes' used to filter tracks from the output
type TrackFilter interface {
	// Check the given track with the filter
	Filter(track TrackSummary) bool
	// Get the string this filter was created with
	GetFilterText() string
}

// MinStartTimeFilter - An implementation of TrackFilter to filter tracks by there start time.
// All Tracks started after the given time pass the filter
type MinStartTimeFilter struct {
	MinStartTime time.Time
	myFilterText string
}

func (filter *MinStartTimeFilter) Filter(track TrackSummary) bool {

	if !track.TimeDataValid {
		return true
	}

	if track.StartTime.Unix() >= filter.MinStartTime.Unix() {
		return true
	}

	return false
}

// NewMinStartTimeFilter - Get a new instance of the MinStartTimeFilter or an error if the `filterText` can not be parsed as a timestamp
func NewMinStartTimeFilter(filterText string) (MinStartTimeFilter, error) {
	ret := MinStartTimeFilter{}
	ret.myFilterText = filterText
	var err error
	ret.MinStartTime, err = filterTextToTime(filterText)

	return ret, err
}

func (filter *MinStartTimeFilter) GetFilterText() string {
	return filter.myFilterText
}

// MaxStartTimeFilter - An implementation of TrackFilter to filter tracks by there start time.
// All Tracks started before the given time pass the filter
type MaxStartTimeFilter struct {
	MaxStartTime time.Time
	myFilterText string
}

func (filter *MaxStartTimeFilter) Filter(track TrackSummary) bool {

	if !track.TimeDataValid {
		return true
	}

	if track.StartTime.Unix() <= filter.MaxStartTime.Unix() {
		return true
	}

	return false
}

// NewMaxStartTimeFilter - Get a new instance of the MaxStartTimeFilter or an error if the `filterText` can not be parsed as a timestamp
func NewMaxStartTimeFilter(filterText string) (MaxStartTimeFilter, error) {
	ret := MaxStartTimeFilter{}
	ret.myFilterText = filterText
	var err error
	ret.MaxStartTime, err = filterTextToTime(filterText)

	return ret, err
}

func (filter *MaxStartTimeFilter) GetFilterText() string {
	return filter.myFilterText
}

func FilterTracks(tracks []Track, filters []TrackFilter) []Track {
	ret := []Track{}
	for _, track := range tracks {
		trackPassedAll := true
		for _, filter := range filters {
			if !filter.Filter(track.TrackSummary) {
				trackPassedAll = false
				break
			}
		}
		if trackPassedAll {
			ret = append(ret, track)
		}
	}

	return ret
}

func filterTextToTime(filterText string) (time.Time, error) {
	var ret time.Time
	var err error

	ret, err = time.Parse(time.DateTime, filterText)
	if err == nil {
		return ret, nil
	}

	ret, err = time.Parse("2006-01-02 15:04", filterText)
	if err == nil {
		return ret, nil
	}

	ret, err = time.Parse(time.DateOnly, filterText)
	if err == nil {
		return ret, nil
	}

	return ret, err
}
