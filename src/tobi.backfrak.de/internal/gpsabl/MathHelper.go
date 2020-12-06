package gpsabl

import (
	"math"
	"time"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// CompareFloat64With4Digits - Compare two float64 to 4 digits after decimal
func CompareFloat64With4Digits(in1, in2 float64) bool {
	return RoundFloat64To4Digits(in1) == RoundFloat64To4Digits(in2)
}

// RoundFloat64To4Digits - Rounds a float64 to 4 digits after decimal
func RoundFloat64To4Digits(in float64) float64 {
	return math.Round(in*10000) / 10000
}

// RoundFloat64To2Digits - Rounds a float64 to 2 digits after decimal
func RoundFloat64To2Digits(in float64) float64 {
	return math.Round(in*100) / 100
}

// TrackDataArrays - The tracks data in arrays, sorted by values not the line
type TrackDataArrays struct {
	AllTimeDataValid bool

	Distances           []float64
	HorizontalDistances []float64
	ElevationGains      []float64
	ElevationLoses      []float64
	AltitudeRanges      []float64
	MinimumAltitudes    []float64
	MaximumAltitudes    []float64
	UpwardsDistances    []float64
	DownwardsDistances  []float64
	StartTimes          []time.Time
	EndTimes            []time.Time
	Durations           []time.Duration
	MovingTimes         []time.Duration
	UpwardsTimes        []time.Duration
	DownwardsTimes      []time.Duration
	AverageSpeeds       []float64
	UpwardsSpeeds       []float64
	DownwardsSpeeds     []float64
}

// TrackStatisticSummaryData - Contains statistic data from a bunch of tracks
type TrackStatisticSummaryData struct {
	AllTimeDataValid bool
	InputTackCount   int
	Sum              ExtendedTrackSummary
	Average          ExtendedTrackSummary
	Minimum          ExtendedTrackSummary
	Maximum          ExtendedTrackSummary
}

// ExtendedTrackSummary - The TrackSummary extended by the duration
type ExtendedTrackSummary struct {
	TrackSummary
	Duration       time.Duration
	AverageSpeed   float64
	UpwardsSpeed   float64
	DownwardsSpeed float64
	AltitudeRange  float64
}

// GetTrackDataArrays - Get the tracks data in arrays, sorted by values not the line
func GetTrackDataArrays(lines []OutputLine) TrackDataArrays {
	ret := TrackDataArrays{}
	ret.AllTimeDataValid = allTimeDataValid(lines)

	for _, line := range lines {
		info := line.Data
		ret.Distances = append(ret.Distances, info.GetDistance())
		ret.HorizontalDistances = append(ret.HorizontalDistances, info.GetHorizontalDistance())
		ret.ElevationGains = append(ret.ElevationGains, float64(info.GetElevationGain()))
		ret.ElevationLoses = append(ret.ElevationLoses, float64(info.GetElevationLose()))
		ret.AltitudeRanges = append(ret.AltitudeRanges, float64(info.GetAltitudeRange()))
		ret.MinimumAltitudes = append(ret.MinimumAltitudes, float64(info.GetMinimumAltitude()))
		ret.MaximumAltitudes = append(ret.MaximumAltitudes, float64(info.GetMaximumAltitude()))
		ret.UpwardsDistances = append(ret.UpwardsDistances, info.GetUpwardsDistance())
		ret.DownwardsDistances = append(ret.DownwardsDistances, info.GetDownwardsDistance())
		if ret.AllTimeDataValid {
			ret.Durations = append(ret.Durations, info.GetEndTime().Sub(info.GetStartTime()))
			ret.StartTimes = append(ret.StartTimes, info.GetStartTime())
			ret.EndTimes = append(ret.EndTimes, info.GetEndTime())
			ret.MovingTimes = append(ret.MovingTimes, info.GetMovingTime())
			ret.UpwardsTimes = append(ret.UpwardsTimes, info.GetUpwardsTime())
			ret.DownwardsTimes = append(ret.DownwardsTimes, info.GetDownwardsTime())
			ret.AverageSpeeds = append(ret.AverageSpeeds, info.GetAvarageSpeed())
			ret.UpwardsSpeeds = append(ret.UpwardsSpeeds, info.GetUpwardsSpeed())
			ret.DownwardsSpeeds = append(ret.DownwardsSpeeds, info.GetDownwardsSpeed())
		}
	}

	return ret
}

// GetStatisticSummaryData - Get the TrackStatisticSummaryData of the input tracks
func GetStatisticSummaryData(lines []OutputLine) TrackStatisticSummaryData {
	ret := TrackStatisticSummaryData{}
	arrays := GetTrackDataArrays(lines)
	ret.AllTimeDataValid = arrays.AllTimeDataValid
	ret.InputTackCount = len(lines)

	ret.Sum.Distance = sumFloat64Array(arrays.Distances)
	ret.Average.Distance = ret.Sum.Distance / float64(ret.InputTackCount)
	ret.Minimum.Distance = minFloat64Array(arrays.Distances)
	ret.Maximum.Distance = maxFloat64Array(arrays.Distances)

	ret.Sum.DownwardsDistance = sumFloat64Array(arrays.DownwardsDistances)
	ret.Average.DownwardsDistance = ret.Sum.DownwardsDistance / float64(ret.InputTackCount)
	ret.Minimum.DownwardsDistance = minFloat64Array(arrays.DownwardsDistances)
	ret.Maximum.DownwardsDistance = maxFloat64Array(arrays.DownwardsDistances)

	ret.Sum.ElevationGain = float32(sumFloat64Array(arrays.ElevationGains))
	ret.Average.ElevationGain = ret.Sum.ElevationGain / float32(ret.InputTackCount)
	ret.Minimum.ElevationGain = float32(minFloat64Array(arrays.ElevationGains))
	ret.Maximum.ElevationGain = float32(maxFloat64Array(arrays.ElevationGains))

	ret.Sum.ElevationLose = float32(sumFloat64Array(arrays.ElevationLoses))
	ret.Average.ElevationLose = ret.Sum.ElevationLose / float32(ret.InputTackCount)
	ret.Minimum.ElevationLose = float32(minFloat64Array(arrays.ElevationLoses))
	ret.Maximum.ElevationLose = float32(maxFloat64Array(arrays.ElevationLoses))

	ret.Sum.HorizontalDistance = sumFloat64Array(arrays.HorizontalDistances)
	ret.Average.HorizontalDistance = ret.Sum.HorizontalDistance / float64(ret.InputTackCount)
	ret.Minimum.HorizontalDistance = minFloat64Array(arrays.HorizontalDistances)
	ret.Maximum.HorizontalDistance = maxFloat64Array(arrays.HorizontalDistances)

	ret.Sum.UpwardsDistance = sumFloat64Array(arrays.UpwardsDistances)
	ret.Average.UpwardsDistance = ret.Sum.UpwardsDistance / float64(ret.InputTackCount)
	ret.Minimum.UpwardsDistance = minFloat64Array(arrays.UpwardsDistances)
	ret.Maximum.UpwardsDistance = maxFloat64Array(arrays.UpwardsDistances)

	sumRange := sumFloat64Array(arrays.AltitudeRanges)
	ret.Average.AltitudeRange = sumRange / float64(ret.InputTackCount)
	ret.Minimum.AltitudeRange = minFloat64Array(arrays.AltitudeRanges)
	ret.Maximum.AltitudeRange = maxFloat64Array(arrays.AltitudeRanges)

	ret.Minimum.MinimumAltitude = float32(minFloat64Array(arrays.MinimumAltitudes))
	ret.Minimum.MaximumAltitude = float32(minFloat64Array(arrays.MaximumAltitudes))
	ret.Maximum.MinimumAltitude = float32(maxFloat64Array(arrays.MinimumAltitudes))
	ret.Maximum.MaximumAltitude = float32(maxFloat64Array(arrays.MaximumAltitudes))

	if ret.AllTimeDataValid {
		ret.Sum.DownwardsTime = sumTimeDurationArray(arrays.DownwardsTimes)
		ret.Average.DownwardsTime = averageDuration(ret.Sum.DownwardsTime, ret.InputTackCount)
		ret.Minimum.DownwardsTime = minTimeDurationArray(arrays.DownwardsTimes)
		ret.Maximum.DownwardsTime = maxTimeDurationArray(arrays.DownwardsTimes)

		ret.Sum.Duration = sumTimeDurationArray(arrays.Durations)
		ret.Average.Duration = averageDuration(ret.Sum.Duration, ret.InputTackCount)
		ret.Minimum.Duration = minTimeDurationArray(arrays.Durations)
		ret.Maximum.Duration = maxTimeDurationArray(arrays.Durations)

		ret.Sum.MovingTime = sumTimeDurationArray(arrays.MovingTimes)
		ret.Average.MovingTime = averageDuration(ret.Sum.MovingTime, ret.InputTackCount)
		ret.Minimum.MovingTime = minTimeDurationArray(arrays.MovingTimes)
		ret.Maximum.MovingTime = maxTimeDurationArray(arrays.MovingTimes)

		ret.Sum.UpwardsTime = sumTimeDurationArray(arrays.UpwardsTimes)
		ret.Average.UpwardsTime = averageDuration(ret.Sum.UpwardsTime, ret.InputTackCount)
		ret.Minimum.UpwardsTime = minTimeDurationArray(arrays.UpwardsTimes)
		ret.Maximum.UpwardsTime = maxTimeDurationArray(arrays.UpwardsTimes)

		speedSum := sumFloat64Array(arrays.AverageSpeeds)
		ret.Average.AverageSpeed = speedSum / float64(ret.InputTackCount)
		ret.Maximum.AverageSpeed = maxFloat64Array(arrays.AverageSpeeds)
		ret.Minimum.AverageSpeed = minFloat64Array(arrays.AverageSpeeds)

		speedSum = sumFloat64Array(arrays.UpwardsSpeeds)
		ret.Average.UpwardsSpeed = speedSum / float64(ret.InputTackCount)
		ret.Maximum.UpwardsSpeed = maxFloat64Array(arrays.UpwardsSpeeds)
		ret.Minimum.UpwardsSpeed = minFloat64Array(arrays.UpwardsSpeeds)

		speedSum = sumFloat64Array(arrays.DownwardsSpeeds)
		ret.Average.DownwardsSpeed = speedSum / float64(ret.InputTackCount)
		ret.Minimum.DownwardsSpeed = minFloat64Array(arrays.DownwardsSpeeds)
		ret.Maximum.DownwardsSpeed = maxFloat64Array(arrays.DownwardsSpeeds)

		ret.Maximum.StartTime = maxTimeArray(arrays.StartTimes)
		ret.Maximum.EndTime = maxTimeArray(arrays.EndTimes)

		ret.Minimum.StartTime = minTimeArray(arrays.StartTimes)
		ret.Minimum.EndTime = minTimeArray(arrays.EndTimes)
	}

	return ret
}

func averageDuration(sum time.Duration, count int) time.Duration {
	timeSumNanoSec := int64(sum)
	avrDurationNanoSec := timeSumNanoSec / int64(count)

	return time.Duration(avrDurationNanoSec)
}

func allTimeDataValid(lines []OutputLine) bool {
	for _, line := range lines {
		if line.Data.GetTimeDataValid() == false {
			return false
		}
	}

	return true
}

func sumFloat64Array(data []float64) float64 {
	ret := 0.0
	for _, value := range data {
		ret += value
	}

	return ret
}

func minFloat64Array(data []float64) float64 {
	min := data[0]
	for _, value := range data {
		if value < min {
			min = value
		}
	}

	return min
}

func maxFloat64Array(data []float64) float64 {
	max := data[0]
	for _, value := range data {
		if value > max {
			max = value
		}
	}

	return max
}

func sumTimeDurationArray(data []time.Duration) time.Duration {
	var ret time.Duration
	for _, value := range data {
		ret += value
	}

	return ret
}

func minTimeDurationArray(data []time.Duration) time.Duration {
	min := data[0]
	for _, value := range data {
		if value < min {
			min = value
		}
	}

	return min
}

func maxTimeDurationArray(data []time.Duration) time.Duration {
	max := data[0]
	for _, value := range data {
		if value > max {
			max = value
		}
	}

	return max
}

func minTimeArray(data []time.Time) time.Time {
	min := data[0]
	for _, value := range data {
		if value.Before(min) {
			min = value
		}
	}

	return min
}

func maxTimeArray(data []time.Time) time.Time {
	max := data[0]
	for _, value := range data {
		if max.Before(value) {
			max = value
		}
	}

	return max
}
