package gpsabl

import (
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
	formater := NewCsvOutputFormater(";", false)
	formater.AddOutPut(file, "segment", false)

	arrays := GetTrackDataArrays(formater.lineBuffer)

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
	formater := NewCsvOutputFormater(";", false)
	formater.AddOutPut(file, "segment", false)

	arrays := GetTrackDataArrays(formater.lineBuffer)

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
	formater := NewCsvOutputFormater(";", false)
	formater.AddOutPut(file1, "segment", false)
	formater.AddOutPut(file2, "segment", false)

	arrays := GetTrackDataArrays(formater.lineBuffer)

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
	formater := NewCsvOutputFormater(";", false)
	formater.AddOutPut(file, "segment", false)

	summaries := GetStatisticSummaryData(formater.lineBuffer)

	if summaries.AllTimeDataValid == true {
		t.Errorf("All time data is valid, but should not")
	}
}

func TestGetStatisticSummaryDataWithTime(t *testing.T) {
	file := getTrackFileTwoTracksWithThreeSegmentsWithTime()
	formater := NewCsvOutputFormater(";", false)
	formater.AddOutPut(file, "segment", false)

	summaries := GetStatisticSummaryData(formater.lineBuffer)

	if summaries.AllTimeDataValid == false {
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
