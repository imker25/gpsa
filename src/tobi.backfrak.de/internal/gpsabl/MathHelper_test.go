package gpsabl

import (
	"fmt"
	"testing"
	"time"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestRoundFloat64To2Digits(t *testing.T) {
	if RoundFloat64To2Digits(23.123) != 23.12 {
		t.Errorf("Round down does not work")
	}

	if RoundFloat64To2Digits(23.127) != 23.13 {
		t.Errorf("Round up does not work")
	}

	if RoundFloat64To2Digits(23.1200) != 23.12 {
		t.Errorf("Round 0 does not work")
	}
}

func TestRoundFloat64To4Digits(t *testing.T) {
	if RoundFloat64To4Digits(23.12344) != 23.1234 {
		t.Errorf("Round down does not work")
	}

	if RoundFloat64To4Digits(23.12717) != 23.1272 {
		t.Errorf("Round up does not work")
	}

	if RoundFloat64To4Digits(23.123300) != 23.1233 {
		t.Errorf("Round 0 does not work")
	}
}

func TestCompareFloat64With4Digits(t *testing.T) {

	if !CompareFloat64With4Digits(23.12344, 23.1234) {
		t.Errorf("Round down does not work")
	}

	if !CompareFloat64With4Digits(23.12348, 23.1235) {
		t.Errorf("Round up does not work")
	}

	if !CompareFloat64With4Digits(23.123700, 23.1237) {
		t.Errorf("Round 0 does not work")
	}
}

func TestSumFloats(t *testing.T) {
	arr := []float64{1.0, 2.3, 3.8}
	sum := sumFloat64Array(arr)
	res := 7.1
	if sum != res {
		t.Errorf("The sum is %f, but expected %f", sum, res)
	}

	arr = []float64{-1.0, -2.3, -3.8}
	sum = sumFloat64Array(arr)
	res = -7.1
	if sum != res {
		t.Errorf("The sum is %f, but expected %f", sum, res)
	}

	arr = []float64{-1.0, 2.3, 3.8}
	sum = sumFloat64Array(arr)
	res = 5.1
	if sum != res {
		t.Errorf("The sum is %f, but expected %f", sum, res)
	}
}

func TestMinMaxFloats(t *testing.T) {
	expMin := -2.1
	expMax := 3.8
	arr := []float64{1.0, 2.3, expMax, expMin}

	min := minFloat64Array(arr)
	if min != expMin {
		t.Errorf("The minimal value is %f, but %f was expected", min, expMin)
	}

	max := maxFloat64Array(arr)
	if max != expMax {
		t.Errorf("The maximal value is %f, but %f was expected", max, expMax)
	}
}

func TestSumTimeDuration(t *testing.T) {
	now := time.Now()
	val1 := now.Sub(now.Add(-1 * time.Minute))
	val2 := now.Sub(now.Add(-3 * time.Minute))
	val3 := now.Sub(now.Add(-25 * time.Second))
	val4 := now.Sub(now.Add(-36 * time.Second))
	arr := []time.Duration{val1, val2, val3, val4}
	sum := sumTimeDurationArray(arr)
	res := val1 + val2 + val3 + val4
	if sum != res {
		t.Errorf("The time duration sum is %s, but expected %s", sum.String(), res.String())
	}
}

func TestMinMaxDuration(t *testing.T) {
	now := time.Now()
	val1 := now.Sub(now.Add(-1 * time.Minute))
	val2 := now.Sub(now.Add(-3 * time.Minute))
	val3 := now.Sub(now.Add(-25 * time.Second))
	val4 := now.Sub(now.Add(-36 * time.Second))
	arr := []time.Duration{val1, val2, val3, val4}

	min := minTimeDurationArray(arr)
	if min != val3 {
		t.Errorf("The minimal value is %s, but %s was expected", min.String(), val3.String())
	}

	max := maxTimeDurationArray(arr)
	if max != val2 {
		t.Errorf("The maximal value is %s, but %s was expected", max.String(), val2.String())
	}
}

func TestMinMaxTime(t *testing.T) {
	now := time.Now()
	val1 := now.Add(-1 * time.Minute)
	val2 := now.Add(-3 * time.Minute)
	val3 := now.Add(-25 * time.Second)
	val4 := now.Add(-36 * time.Second)
	arr := []time.Time{val1, val2, val3, val4}

	min := minTimeArray(arr)
	if min != val2 {
		t.Errorf("The minimal value is %s, but %s was expected", min.String(), val2.String())
	}

	max := maxTimeArray(arr)
	if max != val3 {
		t.Errorf("The maximal value is %s, but %s was expected", max.String(), val3.String())
	}
}

func TestAverageDuration(t *testing.T) {
	now := time.Now()
	val1 := now.Sub(now.Add(-2 * time.Minute))
	val2 := now.Sub(now.Add(-3 * time.Minute))
	val3 := now.Sub(now.Add(-3 * time.Minute))
	val4 := now.Sub(now.Add(-2 * time.Minute))
	arr := []time.Duration{val1, val2, val3, val4}
	sum := sumTimeDurationArray(arr)
	avr := averageDuration(sum, 4)
	res := "2m30s"
	if avr.String() != res {
		t.Errorf("The time duration sum is %s, but expected %s", avr.String(), res)
	}
}

func TestGetTrackDataArraysWithTime(t *testing.T) {
	file := getTrackFileTwoTracksWithThreeSegmentsWithTime()

	lineBuffer := []OutputLine{}
	for i, track := range file.Tracks {
		for j, seg := range track.TrackSegments {
			lineBuffer = append(lineBuffer, *NewOutputLine(fmt.Sprintf("%d-%d", i, j), seg))
		}
	}

	arrays := GetTrackDataArrays(lineBuffer)

	if arrays.AllTimeDataValid == false {
		t.Errorf("Not all time data is valid, but should")
	}

	if len(arrays.DownwardsSpeeds) != len(arrays.AltitudeRanges) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.AverageSpeeds) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.Distances) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.DownwardsDistances) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.DownwardsTimes) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.Durations) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.ElevationGains) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.ElevationLoses) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.EndTimes) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.HorizontalDistances) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.MaximumAltitudes) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.MinimumAltitudes) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.MovingTimes) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.StartTimes) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.UpwardsDistances) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.UpwardsSpeeds) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.DownwardsSpeeds) != len(arrays.UpwardsTimes) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
}

func TestGetTrackDataArraysWithOutTime(t *testing.T) {
	file := getTrackFileTwoTracksWithThreeSegments()
	lineBuffer := []OutputLine{}
	for i, track := range file.Tracks {
		for j, seg := range track.TrackSegments {
			lineBuffer = append(lineBuffer, *NewOutputLine(fmt.Sprintf("%d-%d", i, j), seg))
		}
	}

	arrays := GetTrackDataArrays(lineBuffer)

	if arrays.AllTimeDataValid == true {
		t.Errorf("All time data is valid, but should not")
	}

	if len(arrays.DownwardsSpeeds) != 0 {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if 0 != len(arrays.AverageSpeeds) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if len(arrays.AltitudeRanges) != len(arrays.Distances) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.AltitudeRanges) != len(arrays.DownwardsDistances) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if 0 != len(arrays.DownwardsTimes) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if 0 != len(arrays.Durations) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if len(arrays.AltitudeRanges) != len(arrays.ElevationGains) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.AltitudeRanges) != len(arrays.ElevationLoses) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if 0 != len(arrays.EndTimes) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if len(arrays.AltitudeRanges) != len(arrays.HorizontalDistances) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.AltitudeRanges) != len(arrays.MaximumAltitudes) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.AltitudeRanges) != len(arrays.MinimumAltitudes) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if 0 != len(arrays.MovingTimes) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if 0 != len(arrays.StartTimes) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if len(arrays.AltitudeRanges) != len(arrays.UpwardsDistances) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if 0 != len(arrays.UpwardsSpeeds) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if 0 != len(arrays.UpwardsTimes) {
		t.Errorf("The number of entries is not 0, nut should")
	}
}

func TestGetTrackDataArraysWithMixedTime(t *testing.T) {
	file1 := getTrackFileTwoTracksWithThreeSegments()
	file2 := getTrackFileTwoTracksWithThreeSegmentsWithTime()
	lineBuffer := []OutputLine{}
	for i, track := range file1.Tracks {
		for j, seg := range track.TrackSegments {
			lineBuffer = append(lineBuffer, *NewOutputLine(fmt.Sprintf("%d-%d", i, j), seg))
		}
	}
	for i, track := range file2.Tracks {
		for j, seg := range track.TrackSegments {
			lineBuffer = append(lineBuffer, *NewOutputLine(fmt.Sprintf("%d-%d", i, j), seg))
		}
	}

	arrays := GetTrackDataArrays(lineBuffer)

	if arrays.AllTimeDataValid == true {
		t.Errorf("All time data is valid, but should not")
	}

	if len(arrays.DownwardsSpeeds) != 0 {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if 0 != len(arrays.AverageSpeeds) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if len(arrays.AltitudeRanges) != len(arrays.Distances) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.AltitudeRanges) != len(arrays.DownwardsDistances) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if 0 != len(arrays.DownwardsTimes) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if 0 != len(arrays.Durations) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if len(arrays.AltitudeRanges) != len(arrays.ElevationGains) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.AltitudeRanges) != len(arrays.ElevationLoses) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if 0 != len(arrays.EndTimes) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if len(arrays.AltitudeRanges) != len(arrays.HorizontalDistances) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.AltitudeRanges) != len(arrays.MaximumAltitudes) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if len(arrays.AltitudeRanges) != len(arrays.MinimumAltitudes) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if 0 != len(arrays.MovingTimes) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if 0 != len(arrays.StartTimes) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if len(arrays.AltitudeRanges) != len(arrays.UpwardsDistances) {
		t.Errorf("The number of entries is not the same for all arrays")
	}
	if 0 != len(arrays.UpwardsSpeeds) {
		t.Errorf("The number of entries is not 0, nut should")
	}
	if 0 != len(arrays.UpwardsTimes) {
		t.Errorf("The number of entries is not 0, nut should")
	}
}

func TestGetStatisticSummaryDataWithoutTime(t *testing.T) {
	file := getTrackFileTwoTracksWithThreeSegments()
	lineBuffer := []OutputLine{}
	for i, track := range file.Tracks {
		for j, seg := range track.TrackSegments {
			lineBuffer = append(lineBuffer, *NewOutputLine(fmt.Sprintf("%d-%d", i, j), seg))
		}
	}

	summaries := GetStatisticSummaryData(lineBuffer)

	if summaries.AllTimeDataValid == true {
		t.Errorf("All time data is valid, but should not")
	}
}

func TestOutPutContainsLineByTimeStamps1(t *testing.T) {
	trackFile := getTrackFileTwoTracksWithThreeSegmentsWithTime()
	entries := []OutputLine{}
	for i, track := range trackFile.Tracks {
		for j, seg := range track.TrackSegments {
			entries = append(entries, *NewOutputLine(fmt.Sprintf("%d-%d", i, j), seg))
		}
	}

	if OutputContainsLineByTimeStamps(entries, entries[0]) == false {
		t.Errorf("Got false, but expect true")
	}

	if OutputContainsLineByTimeStamps(entries, *NewOutputLine("bla", getTrackWithDifferentTime())) == true {
		t.Errorf("Got true, but expect false")
	}
}

func TestOutPutContainsLineByTimeStamps2(t *testing.T) {

	trackFile := getTrackFileTwoTracksWithThreeSegments()
	entries := []OutputLine{}
	for i, track := range trackFile.Tracks {
		for j, seg := range track.TrackSegments {
			entries = append(entries, *NewOutputLine(fmt.Sprintf("%d-%d", i, j), seg))
		}
	}

	if OutputContainsLineByTimeStamps(entries, entries[0]) == true {
		t.Errorf("Got true, but expect false")
	}
}

func TestGetStatisticSummaryDataWithTime(t *testing.T) {
	file := getTrackFileTwoTracksWithThreeSegmentsWithTime()
	lineBuffer := []OutputLine{}
	for i, track := range file.Tracks {
		for j, seg := range track.TrackSegments {
			lineBuffer = append(lineBuffer, *NewOutputLine(fmt.Sprintf("%d-%d", i, j), seg))
		}
	}

	summaries := GetStatisticSummaryData(lineBuffer)

	if summaries.AllTimeDataValid && summaries.Sum.TimeDataValid && summaries.Average.TimeDataValid && summaries.Maximum.TimeDataValid && summaries.Minimum.TimeDataValid == false {
		t.Errorf("Not all time data is valid, but should")
	}

	if summaries.Sum.Duration.String() != "1m0s" {
		t.Errorf("The value is %s, but expected was %s", summaries.Sum.Duration.String(), "1m0s")
	}
	if summaries.Average.Duration.String() != "20s" {
		t.Errorf("The value is %s, but expected was %s", summaries.Average.Duration.String(), "20s")
	}
	if summaries.Maximum.Duration.String() != "20s" {
		t.Errorf("The value is %s, but expected was %s", summaries.Average.Duration.String(), "20s")
	}
	if summaries.Minimum.Duration.String() != "20s" {
		t.Errorf("The value is %s, but expected was %s", summaries.Average.Duration.String(), "20s")
	}
}

func getTrackWithDifferentTime() Track {
	t1, _ := time.Parse(time.RFC3339, "2015-08-22T17:19:33Z")
	t2, _ := time.Parse(time.RFC3339, "2015-08-22T17:19:43Z")
	t3, _ := time.Parse(time.RFC3339, "2015-08-22T17:19:53Z")
	pnt1 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t1)
	pnt2 := getTrackPointWithTime(50.11495750, 8.684874770, 108.0, t2)
	pnt3 := getTrackPointWithTime(50.11484790, 8.684885500, 109.0, t3)
	points := []TrackPoint{pnt1, pnt2, pnt3}

	FillDistancesTrackPoint(&points[0], TrackPoint{}, points[1])
	FillDistancesTrackPoint(&points[1], points[0], points[2])
	FillDistancesTrackPoint(&points[2], points[1], TrackPoint{})
	FillValuesTrackPointArray(points, "none", 0.3, 10.0)
	seg := TrackSegment{}
	seg.TrackPoints = points
	ret := Track{}
	FillTrackSegmentValues(&seg)
	ret.TrackSegments = []TrackSegment{seg}
	FillTrackValues(&ret)
	ret.NumberOfSegments = 1

	return ret
}

func getTrackFileTwoTracksWithThreeSegmentsWithTime() TrackFile {
	trackFile := getTrackFileTwoTracksWithTime()
	trackFile.Tracks[0].TrackSegments = append(trackFile.Tracks[0].TrackSegments, getSimpleTrackFileWithTime().Tracks[0].TrackSegments[0])
	FillTrackValues(&trackFile.Tracks[0])

	return trackFile
}

func getTrackFileTwoTracksWithTime() TrackFile {
	trackFile := getSimpleTrackFileWithTime()
	trackFile.Tracks = append(trackFile.Tracks, getSimpleTrackFileWithTime().Tracks...)
	FillTrackFileValues(&trackFile)

	return trackFile
}

func getTrackFileTwoTracksWithThreeSegments() TrackFile {
	trackFile := getTrackFileTwoTracks()
	trackFile.Tracks[0].TrackSegments = append(trackFile.Tracks[0].TrackSegments, getSimpleTrackFile().Tracks[0].TrackSegments[0])
	FillTrackValues(&trackFile.Tracks[0])

	return trackFile
}

func getTrackFileTwoTracks() TrackFile {
	trackFile := getSimpleTrackFile()
	trackFile.Tracks = append(trackFile.Tracks, getSimpleTrackFile().Tracks...)
	FillTrackFileValues(&trackFile)

	return trackFile
}
