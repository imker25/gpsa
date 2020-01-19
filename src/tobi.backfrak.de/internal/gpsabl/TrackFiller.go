package gpsabl

import (
	"math"
	"time"
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
	err := fillCorrectedElevationTrackPoint(pnts, correction)
	if err != nil {
		return err
	}
	fillElevationGainLoseTrackPoint(pnts)
	fillCountUpDownWards(pnts, correction)
	fillSpeedValues(pnts)

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

	fillTrackSummaryValues(segment, iPnts, true)
}

// FillTrackValues - Fills the distance and atitute fields of a tack  by adding up all TrackSegments distances
// All TrackSegment has to be set before. See FillTrackSegmentValues
func FillTrackValues(track *Track) {
	iSegs := []TrackSummaryProvider{}
	for i := range track.TrackSegments {
		iSeg := TrackSummaryProvider(&track.TrackSegments[i])
		iSegs = append(iSegs, iSeg)
	}

	fillTrackSummaryValues(track, iSegs, false)
}

// FillTrackFileValues - Fills the distance and atitute fields of a tack  by adding up all TrackSegments distances
// All Track values has to be set before. See FillTrackValues
func FillTrackFileValues(file *TrackFile) {
	iTrks := []TrackSummaryProvider{}
	for i := range file.Tracks {
		itrk := TrackSummaryProvider(&file.Tracks[i])
		iTrks = append(iTrks, itrk)
	}

	fillTrackSummaryValues(file, iTrks, false)
}

// GetValidCorrectionParameters - Get the valid parameters for fillCorrectedElevationTrackPoint correction parameter
func GetValidCorrectionParameters() []string {
	return []string{"none", "linear", "steps"}
}

// GetValidCorrectionParametersString - Get the valid parameters for fillCorrectedElevationTrackPoint correction parameter as one string
func GetValidCorrectionParametersString() string {
	ret := ""
	for _, str := range GetValidCorrectionParameters() {
		ret = str + " " + ret
	}

	return ret
}

// CheckValidCorrectionParameters - Check if a string is a valid parameter for fillCorrectedElevationTrackPoint correction parameter
func CheckValidCorrectionParameters(given string) bool {
	for _, str := range GetValidCorrectionParameters() {
		if str == given {
			return true
		}
	}

	return false
}

// fillCorrectedElevationTrackPoint - Set the CorectedElevation value in a list of TrackPoints
// Basicaly this will run a somthing algorythm over the Elevation
func fillCorrectedElevationTrackPoint(pnts []TrackPoint, correction string) error {

	switch correction {
	case GetValidCorrectionParameters()[1]:
		fillCorrectedElevationTrackPointLinear(pnts)
	case GetValidCorrectionParameters()[0]:
		fillCorrectedElevationTrackPointNone(pnts)
	case GetValidCorrectionParameters()[2]:
		fillCorrectedElevationTrackPointSteps(pnts)
	default:
		return NewCorrectionParameterNotKnownError(correction)
	}

	return nil
}

func fillCorrectedElevationTrackPointNone(pnts []TrackPoint) {
	for i := range pnts {
		pnts[i].CorectedElevation = pnts[i].Elevation
	}
}

func fillCorrectedElevationTrackPointSteps(pnts []TrackPoint) {
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

func fillCorrectedElevationTrackPointLinear(pnts []TrackPoint) {
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

func fillSpeedValues(pnts []TrackPoint) {
	startTime := pnts[0].Time
	for i, pnt := range pnts {
		if pnt.TimeValid {
			pnts[i].MovingTime = pnt.Time.Sub(startTime)
			pnts[i].AvarageSpeed = pnt.DistanceToThisPoint / float64((pnt.MovingTime / 1000000000))
		}
	}

}

// FillCountUpDownWards - Fills the CountUpwards and CountDownwards value
func fillCountUpDownWards(pnts []TrackPoint, correction string) {
	numPnts := len(pnts)
	if correction == GetValidCorrectionParameters()[2] { // In case we do steps correction, CorectedElevation will make no sense
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

func fillTrackSummaryValues(target TrackSummarySetter, input []TrackSummaryProvider, inputIsPoints bool) {
	var dist float64
	var minimumAltitude float32
	var maximumAltitude float32
	var elevationGain float32
	var elevationLose float32
	var upwardsDistance float64
	var downwardsDistance float64
	timeDataValid := true
	var startTime time.Time
	var endTime time.Time
	var movingTime time.Duration
	var movingTimeSum time.Duration

	for i, sum := range input {
		dist = dist + sum.GetDistance()
		elevationGain = elevationGain + sum.GetElevationGain()
		elevationLose = elevationLose + sum.GetElevationLose()
		upwardsDistance = upwardsDistance + sum.GetUpwardsDistance()
		downwardsDistance = downwardsDistance + sum.GetDownwardsDistance()
		movingTimeSum = movingTimeSum + sum.GetMovingTime()

		if i == 0 || sum.GetMaximumAltitude() > maximumAltitude {
			maximumAltitude = sum.GetMaximumAltitude()
		}

		if i == 0 || sum.GetMinimumAltitude() < minimumAltitude {
			minimumAltitude = sum.GetMinimumAltitude()
		}

		if sum.GetTimeDataValid() == false {
			timeDataValid = false
		}
	}

	// If the input has no elements, there can not be valide time data
	if len(input) <= 0 {
		timeDataValid = false
	}
	if timeDataValid {
		startTime = input[0].GetStartTime()
		endTime = input[len(input)-1].GetEndTime()
		if inputIsPoints {
			movingTime = input[len(input)-1].GetMovingTime()
		} else {
			movingTime = movingTimeSum
		}
	}
	target.SetValues(dist, minimumAltitude, maximumAltitude, elevationGain, elevationLose, upwardsDistance, downwardsDistance,
		timeDataValid, startTime, endTime, movingTime)
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
