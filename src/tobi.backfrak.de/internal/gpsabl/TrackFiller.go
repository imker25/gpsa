package gpsabl

import (
	"math"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// MinimalStepHight - The minimal evelation difference for the setp algorythm
const MinimalStepHight = 10.0

// FillDistancesTrackPoint - Adds the distance values to the basePoint.
// The Values of Elevation, Latitude and Longitude had to be set to all points before!
func FillDistancesTrackPoint(basePoint *TrackPoint, beforePoint TrackPoint, nextPoint TrackPoint) {

	if (beforePoint != TrackPoint{}) {
		basePoint.HorizontalDistanceBefore = HaversineDistance(*basePoint, beforePoint)
		basePoint.DistanceBefore = DistanceFromHaversine(basePoint.HorizontalDistanceBefore, *basePoint, beforePoint)
	}

	if (nextPoint != TrackPoint{}) {
		basePoint.HorizontalDistanceNext = HaversineDistance(*basePoint, nextPoint)
		basePoint.DistanceNext = DistanceFromHaversine(basePoint.HorizontalDistanceNext, *basePoint, nextPoint)
	}

}

// FillValuesTrackPointArray - Fills all the values of all in points in the array, but not distances.
// You may use FillDistancesTrackPoint to get the distance values
// The Array must be soreted by the points Number!
func FillValuesTrackPointArray(pnts []TrackPoint, correction string) error {
	fillDistanceToThisPoint(pnts)
	err := fillCorectedElevationTrackPoint(pnts, correction)
	if err != nil {
		return err
	}
	fillElevationGainLoseTrackPoint(pnts)
	fillCountUpDownWards(pnts, correction)

	return nil
}

// FillTrackSegmentValues - Fills the distance and atitute fields of a tack segment by adding up all TrackPoint distances
// All TrackPoint values has to be set before. See FillDistancesTrackPoint and FillValuesTrackPointArray
func FillTrackSegmentValues(segment *TrackSegment) {
	iPnts := []TrackSummaryProvider{}
	for i := range segment.TrackPoints {
		iPnt := TrackSummaryProvider(&segment.TrackPoints[i])
		iPnts = append(iPnts, iPnt)
	}

	fillTrackSummaryValues(segment, iPnts)
}

// FillTrackValues - Fills the distance and atitute fields of a tack  by adding up all TrackSegments distances
// All TrackSegment has to be set before. See FillTrackSegmentValues
func FillTrackValues(track *Track) {
	iSegs := []TrackSummaryProvider{}
	for i := range track.TrackSegments {
		iSeg := TrackSummaryProvider(&track.TrackSegments[i])
		iSegs = append(iSegs, iSeg)
	}

	fillTrackSummaryValues(track, iSegs)
}

// FillTrackFileValues - Fills the distance and atitute fields of a tack  by adding up all TrackSegments distances
// All Track values has to be set before. See FillTrackValues
func FillTrackFileValues(file *TrackFile) {
	iTrks := []TrackSummaryProvider{}
	for i := range file.Tracks {
		itrk := TrackSummaryProvider(&file.Tracks[i])
		iTrks = append(iTrks, itrk)
	}

	fillTrackSummaryValues(file, iTrks)
}

// GetValideCorectionParamters - Get the valide paramters for FillCorectedElevationTrackPoint corection paramter
func GetValideCorectionParamters() []string {
	return []string{"none", "linear", "steps"}
}

// GetValideCorectionParamtersString - Get the valide paramters for FillCorectedElevationTrackPoint corection paramter as one string
func GetValideCorectionParamtersString() string {
	ret := ""
	for _, str := range GetValideCorectionParamters() {
		ret = str + " " + ret
	}

	return ret
}

// CheckValideCorectionParamters - Check if a string is a valide paramter for FillCorectedElevationTrackPoint corection paramter
func CheckValideCorectionParamters(given string) bool {
	for _, str := range GetValideCorectionParamters() {
		if str == given {
			return true
		}
	}

	return false
}

// fillCorectedElevationTrackPoint - Set the CorectedElevation value in a list of TrackPoints
// Basicaly this will run a somthing algorythm over the Elevation
func fillCorectedElevationTrackPoint(pnts []TrackPoint, corection string) error {

	switch corection {
	case GetValideCorectionParamters()[1]:
		fillCorectedElevationTrackPointLinear(pnts)
	case GetValideCorectionParamters()[0]:
		fillCorectedElevationTrackPointNone(pnts)
	case GetValideCorectionParamters()[2]:
		fillCorectedElevationTrackPointSteps(pnts)
	default:
		return NewCorectionParamterNotKnownError(corection)
	}

	return nil
}

func fillCorectedElevationTrackPointNone(pnts []TrackPoint) {
	for i := range pnts {
		pnts[i].CorectedElevation = pnts[i].Elevation
	}
}

func fillCorectedElevationTrackPointSteps(pnts []TrackPoint) {
	numPnts := len(pnts)
	for i := range pnts {
		if i > 0 && i < (numPnts-1) {
			pnts[i].CorectedElevation = getCorrectedElevationSteps(pnts[i], pnts[i-1], pnts[i+1])
		} else {
			pnts[i].CorectedElevation = pnts[i].Elevation
		}
		//fmt.Println(fmt.Sprintf("%f;%f;%f;", pnts[i].DistanceToThisPoint, pnts[i].Elevation, pnts[i].CorectedElevation))
	}
}

func fillCorectedElevationTrackPointLinear(pnts []TrackPoint) {
	numPnts := len(pnts)
	for i := range pnts {
		if i > 0 && i < (numPnts-1) {
			pnts[i].CorectedElevation = getCorrectedElevationLinear(pnts[i], pnts[i-1], pnts[i+1])
		} else {
			pnts[i].CorectedElevation = pnts[i].Elevation
		}
		// fmt.Println(fmt.Sprintf("%f;%f;%f;", pnts[i].DistanceToThisPoint, pnts[i].Elevation, pnts[i].CorectedElevation))
	}
}

// fillElevationGainLoseTrackPoint - Set the VerticalDistanceBefore and VerticalDistanceNext values
func fillElevationGainLoseTrackPoint(pnts []TrackPoint) {
	numPnts := len(pnts)
	for i := range pnts {
		if i > 0 { // Evaluation of the first point don't count
			pnts[i].VerticalDistanceBefore = pnts[i].CorectedElevation - pnts[i-1].CorectedElevation
		} else {
			pnts[i].VerticalDistanceBefore = 0.0
		}
		if i < (numPnts - 1) { // Evaluation of the last point don't count
			pnts[i].VerticalDistanceNext = pnts[i+1].CorectedElevation - pnts[i].CorectedElevation
		} else {
			pnts[i].VerticalDistanceNext = 0.0
		}
	}
}

// fillDistanceToThisPoint - Fills the DistanceToThisPoint value
func fillDistanceToThisPoint(pnts []TrackPoint) {
	disToHere := float64(0.0)
	for i := range pnts {
		disToHere += pnts[i].DistanceBefore
		pnts[i].DistanceToThisPoint = disToHere
	}
}

// FillCountUpDownWards - Fills the CountUpwards and CountDownwards value
func fillCountUpDownWards(pnts []TrackPoint, correction string) {
	numPnts := len(pnts)
	if correction == GetValideCorectionParamters()[2] { // In case we do steps correction, CorectedElevation will make no sense
		for i := range pnts {
			if i < (numPnts - 1) { // Evaluation of the last point don't count
				eveDiff := pnts[i+1].Elevation - pnts[i].Elevation
				if eveDiff > 0 {
					pnts[i].CountUpwards = true
				}

				if eveDiff < 0 {
					pnts[i].CountDownwards = true
				}

			}
		}
	} else {
		for i := range pnts {
			if i < (numPnts - 1) { // Evaluation of the last point don't count
				eveDiff := pnts[i+1].CorectedElevation - pnts[i].CorectedElevation
				if eveDiff > 0 {
					pnts[i].CountUpwards = true
				}

				if eveDiff < 0 {
					pnts[i].CountDownwards = true
				}

			}
		}
	}
}

func fillTrackSummaryValues(target TrackSummarySetter, input []TrackSummaryProvider) {
	var dist float64
	var minimumAtitute float32
	var maximumAtitute float32
	var elevationGain float32
	var elevationLose float32
	var upwardsDistance float64
	var downwardsDistance float64

	for i, sum := range input {
		dist = dist + sum.GetDistance()
		elevationGain = elevationGain + sum.GetElevationGain()
		elevationLose = elevationLose + sum.GetElevationLose()
		upwardsDistance = upwardsDistance + sum.GetUpwardsDistance()
		downwardsDistance = downwardsDistance + sum.GetDownwardsDistance()

		if i == 0 || sum.GetMaximumAtitute() > maximumAtitute {
			maximumAtitute = sum.GetMaximumAtitute()
		}

		if i == 0 || sum.GetMinimumAtitute() < minimumAtitute {
			minimumAtitute = sum.GetMinimumAtitute()
		}
	}

	target.SetValues(dist, minimumAtitute, maximumAtitute, elevationGain, elevationLose, upwardsDistance, downwardsDistance)
}

func getCorrectedElevationLinear(basePoint TrackPoint, beforePoint TrackPoint, nextPoint TrackPoint) float32 {

	if beforePoint.Elevation != 0 && nextPoint.Elevation != 0 {
		dEve := nextPoint.Elevation - beforePoint.Elevation
		dx := basePoint.HorizontalDistanceBefore + basePoint.HorizontalDistanceNext
		a := dEve / float32(dx)

		return beforePoint.Elevation + (a * float32(basePoint.HorizontalDistanceBefore))
	}

	return basePoint.Elevation
}

func getCorrectedElevationSteps(basePoint TrackPoint, beforePoint TrackPoint, nextPoint TrackPoint) float32 {

	eveDiffBefore := basePoint.Elevation - beforePoint.CorectedElevation
	eveDiffAfter := nextPoint.Elevation - basePoint.Elevation
	sameDirection := eveDiffBefore * eveDiffAfter

	if math.Abs(float64(eveDiffBefore)) >= MinimalStepHight && (sameDirection > 0 || math.Abs(float64(eveDiffAfter)) < MinimalStepHight) {
		return basePoint.Elevation
	}

	return beforePoint.CorectedElevation
}
