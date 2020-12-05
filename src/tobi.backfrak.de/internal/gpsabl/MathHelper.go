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

func sumTimeDurationArray(data []time.Duration) time.Duration {
	var ret time.Duration
	for _, value := range data {
		ret += value
	}

	return ret
}
