package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"testing"
	"time"
)

func TestNewTrackFile(t *testing.T) {
	file := NewTrackFile("myFile")

	if file.FilePath != "myFile" {
		t.Errorf("The FilePath is %s, but %s was expected", file.FilePath, "myFile")
	}
}

func TestTrackSummary(t *testing.T) {
	sum := TrackSummary{}
	iSumSet := TrackSummarySetter(&sum)
	// now := time.Now()
	iSumSet.SetValues(100.1, 10.4, 40.6, 40.2, 10.0, 70.1, 30.0, false, time.Now(), time.Now(), time.Now().Sub(time.Now()), time.Now().Sub(time.Now()), time.Now().Sub(time.Now()))
	iSum := TrackSummaryProvider(sum)

	if iSum.GetDistance() != 100.1 {
		t.Errorf("The GetDistance() returns %f, but %f was expected", iSum.GetDistance(), 100.1)
	}

	if !CompareFloat64With4Digits(float64(iSum.GetAltitudeRange()), 30.2) {
		t.Errorf("The GetAltitudeRange() returns %f, but %f was expected", iSum.GetAltitudeRange(), 30.2)
	}

	if iSum.GetMaximumAltitude() != 40.6 {
		t.Errorf("The GetMaximumAltitude() returns %f, but %f was expected", iSum.GetMaximumAltitude(), 40.6)
	}

	if iSum.GetMinimumAltitude() != 10.4 {
		t.Errorf("The GetMinimumAltitude() returns %f, but %f was expected", iSum.GetMinimumAltitude(), 10.4)
	}

	if iSum.GetElevationGain() != 40.2 {
		t.Errorf("The GetElevationGain() returns %f, but %f was expected", iSum.GetElevationGain(), 40.2)
	}

	if iSum.GetElevationLose() != 10.0 {
		t.Errorf("The GetElevationLose() returns %f, but %f was expected", iSum.GetElevationLose(), 10.0)
	}

	if iSum.GetUpwardsDistance() != 70.1 {
		t.Errorf("The GetUpwardsDistance() returns %f, but %f was expected", iSum.GetUpwardsDistance(), 70.1)
	}

	if iSum.GetDownwardsDistance() != 30.0 {
		t.Errorf("The GetDownwardsDistance() returns %f, but %f was expected", iSum.GetDownwardsDistance(), 30.0)
	}

	if iSum.GetTimeDataValid() == true {
		t.Errorf("The GetTimeDataValid()  is true, but false is expected")
	}
}

func TestTrackSummaryWithTime(t *testing.T) {
	sum := TrackSummary{}
	iSumSet := TrackSummarySetter(&sum)
	endTime, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:33Z")
	startTime, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:13Z")
	movingTime := endTime.Sub(startTime)
	distance := 100.0
	speed := distance / float64(movingTime/1000000000)
	endTimeUp, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:17Z")
	upwardsTime := startTime.Sub(endTimeUp)
	endTimeDown, _ := time.Parse(time.RFC3339, "2014-08-22T19:19:23Z")
	downwarsTime := startTime.Sub(endTimeDown)
	iSumSet.SetValues(distance, 10.4, 40.6, 40.2, 10.0, 70.1, 30.0, true, startTime, endTime, movingTime, upwardsTime, downwarsTime)
	iSum := TrackSummaryProvider(sum)

	if iSum.GetTimeDataValid() == false {
		t.Errorf("The GetTimeDataValid()  is false, but true is expected")
	}

	if iSum.GetStartTime() != startTime {
		t.Errorf("GetStartTime() is not the expected value")
	}

	if iSum.GetEndTime() != endTime {
		t.Errorf("GetEndTime() is not the expected value")
	}

	if iSum.GetMovingTime() != movingTime {
		t.Errorf("GetMovingTime() is not the expected value")
	}

	if iSum.GetAvarageSpeed() != speed {
		t.Errorf("GetAvarageSpeed() is not the expected value")
	}

	if iSum.GetDownwardsTime() != downwarsTime {
		t.Errorf("GetDownwardsTime() is not the expected value")
	}

	if iSum.GetUpwardsTime() != upwardsTime {
		t.Errorf("GetUpwardsTime() is not the expected value")
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

	if sum.GetAltitudeRange() != 0.0 {
		t.Errorf("The GetAltitudeRange() does return %f, but %f was expected", sum.GetAltitudeRange(), 0.0)
	}

	if sum.GetMaximumAltitude() != 0.0 {
		t.Errorf("The GetMaximumAltitude() does return %f, but %f was expected", sum.GetMaximumAltitude(), 0.0)
	}

	if sum.GetMinimumAltitude() != 0.0 {
		t.Errorf("The GetMinimumAltitude() does return %f, but %f was expected", sum.GetMinimumAltitude(), 0.0)
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

	if sum.GetAltitudeRange() != 0.0 {
		t.Errorf("The GetAltitudeRange() does return %f, but %f was expected", sum.GetAltitudeRange(), 0.0)
	}

	if sum.GetMaximumAltitude() != 0.0 {
		t.Errorf("The GetMaximumAltitude() does return %f, but %f was expected", sum.GetMaximumAltitude(), 0.0)
	}

	if sum.GetMinimumAltitude() != 0.0 {
		t.Errorf("The GetMinimumAltitude() does return %f, but %f was expected", sum.GetMinimumAltitude(), 0.0)
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

	if sum.GetAltitudeRange() != 0.0 {
		t.Errorf("The GetAltitudeRange() does return %f, but %f was expected", sum.GetAltitudeRange(), 0.0)
	}

	if sum.GetMaximumAltitude() != 0.0 {
		t.Errorf("The GetMaximumAltitude() does return %f, but %f was expected", sum.GetMaximumAltitude(), 0.0)
	}

	if sum.GetMinimumAltitude() != 0.0 {
		t.Errorf("The GetMinimumAltitude() does return %f, but %f was expected", sum.GetMinimumAltitude(), 0.0)
	}
}
