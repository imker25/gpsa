package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import "testing"

func TestNewTrackFile(t *testing.T) {
	file := NewTrackFile("myFile")

	if file.FilePath != "myFile" {
		t.Errorf("The FilePath is %s, but %s was expected", file.FilePath, "myFile")
	}
}

func TestTrackFileIsTrackSummary(t *testing.T) {
	file := TrackFile{}
	sum := TrackSummaryProvider(file)

	if sum == nil {
		t.Errorf("The Track struct does not implement the TrackSummaryProvider interface as expected")
	}

	if sum.GetDistance() != 0.0 {
		t.Errorf("The GetDistance() does return %f, but %f was expected", sum.GetDistance(), 0.0)
	}

	if sum.GetAtituteRange() != 0.0 {
		t.Errorf("The GetAtituteRange() does return %f, but %f was expected", sum.GetAtituteRange(), 0.0)
	}

	if sum.GetMaximumAtitute() != 0.0 {
		t.Errorf("The GetMaximumAtitute() does return %f, but %f was expected", sum.GetMaximumAtitute(), 0.0)
	}

	if sum.GetMinimumAtitute() != 0.0 {
		t.Errorf("The GetMinimumAtitute() does return %f, but %f was expected", sum.GetMinimumAtitute(), 0.0)
	}
}

func TestTrackIsTrackSummary(t *testing.T) {
	trk := Track{}
	sum := TrackSummaryProvider(trk)

	if sum == nil {
		t.Errorf("The Track struct does not implement the TrackSummaryProvider interface as expected")
	}

	if sum.GetDistance() != 0.0 {
		t.Errorf("The GetDistance() does return %f, but %f was expected", sum.GetDistance(), 0.0)
	}

	if sum.GetAtituteRange() != 0.0 {
		t.Errorf("The GetAtituteRange() does return %f, but %f was expected", sum.GetAtituteRange(), 0.0)
	}

	if sum.GetMaximumAtitute() != 0.0 {
		t.Errorf("The GetMaximumAtitute() does return %f, but %f was expected", sum.GetMaximumAtitute(), 0.0)
	}

	if sum.GetMinimumAtitute() != 0.0 {
		t.Errorf("The GetMinimumAtitute() does return %f, but %f was expected", sum.GetMinimumAtitute(), 0.0)
	}
}

func TestTrackSegmentIsTrackSummary(t *testing.T) {
	seg := TrackSegment{}
	sum := TrackSummaryProvider(seg)

	if sum == nil {
		t.Errorf("The Track struct does not implement the TrackSummaryProvider interface as expected")
	}

	if sum.GetDistance() != 0.0 {
		t.Errorf("The GetDistance() does return %f, but %f was expected", sum.GetDistance(), 0.0)
	}

	if sum.GetAtituteRange() != 0.0 {
		t.Errorf("The GetAtituteRange() does return %f, but %f was expected", sum.GetAtituteRange(), 0.0)
	}

	if sum.GetMaximumAtitute() != 0.0 {
		t.Errorf("The GetMaximumAtitute() does return %f, but %f was expected", sum.GetMaximumAtitute(), 0.0)
	}

	if sum.GetMinimumAtitute() != 0.0 {
		t.Errorf("The GetMinimumAtitute() does return %f, but %f was expected", sum.GetMinimumAtitute(), 0.0)
	}
}
