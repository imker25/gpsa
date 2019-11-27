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

func TestTrackSummary(t *testing.T) {
	sum := TrackSummary{}
	iSumSet := TrackSummarySetter(&sum)
	iSumSet.SetValues(100.1, 10.4, 40.6, 40.2, 10.0, 70.1, 30.0)
	iSum := TrackSummaryProvider(sum)

	if iSum.GetDistance() != 100.1 {
		t.Errorf("The GetDistance() rutrns %f, but %f was expected", iSum.GetDistance(), 100.1)
	}

	if !CompareFloat64With4Digits(float64(iSum.GetAtituteRange()), 30.2) {
		t.Errorf("The GetAtituteRange() rutrns %f, but %f was expected", iSum.GetAtituteRange(), 30.2)
	}

	if iSum.GetMaximumAtitute() != 40.6 {
		t.Errorf("The GetMaximumAtitute() rutrns %f, but %f was expected", iSum.GetMaximumAtitute(), 40.6)
	}

	if iSum.GetMinimumAtitute() != 10.4 {
		t.Errorf("The GetMinimumAtitute() rutrns %f, but %f was expected", iSum.GetMinimumAtitute(), 10.4)
	}

	if iSum.GetElevationGain() != 40.2 {
		t.Errorf("The GetElevationGain() rutrns %f, but %f was expected", iSum.GetElevationGain(), 40.2)
	}

	if iSum.GetElevationLose() != 10.0 {
		t.Errorf("The GetElevationLose() rutrns %f, but %f was expected", iSum.GetElevationLose(), 10.0)
	}

	if iSum.GetUpwardsDistance() != 70.1 {
		t.Errorf("The GetUpwardsDistance() rutrns %f, but %f was expected", iSum.GetUpwardsDistance(), 70.1)
	}

	if iSum.GetDownwardsDistance() != 30.0 {
		t.Errorf("The GetDownwardsDistance() rutrns %f, but %f was expected", iSum.GetDownwardsDistance(), 30.0)
	}
}

func TestTrackFileIsTrackSummary(t *testing.T) {
	file := TrackFile{}
	sum := TrackSummaryProvider(&file)

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
	sum := TrackSummaryProvider(&trk)

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
	sum := TrackSummaryProvider(&seg)

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
