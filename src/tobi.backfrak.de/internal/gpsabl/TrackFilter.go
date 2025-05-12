package gpsabl

// Copyright 2025 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import "time"

type TrackFilter interface {
	Filter(track TrackSummary) bool
	GetFilterText() string
}

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

func NewMinStartTimeFilter(filterText string) (MinStartTimeFilter, error) {
	ret := MinStartTimeFilter{}
	ret.myFilterText = filterText
	var err error

	ret.MinStartTime, err = time.Parse(time.DateTime, filterText)
	if err == nil {
		return ret, nil
	}

	ret.MinStartTime, err = time.Parse("2006-01-02 15:04", filterText)
	if err == nil {
		return ret, nil
	}

	ret.MinStartTime, err = time.Parse(time.DateOnly, filterText)
	if err == nil {
		return ret, nil
	}

	return ret, err
}

func (filter *MinStartTimeFilter) GetFilterText() string {
	return filter.myFilterText
}
