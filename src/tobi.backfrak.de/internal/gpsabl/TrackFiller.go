package gpsabl

import (
	"math"
	"time"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

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

// FillValuesTrackPointArray - Fills all the values of all in points in the array, but not distances and basic info
// like Elevation, Latitude and Longitude and Time (including TimeValid)
// You may use FillDistancesTrackPoint to get the distance values
// The Array must be soreted by the points Number!
func FillValuesTrackPointArray(pnts []TrackPoint, correction string, minimalMovingSpeed float64, minimalStepHight float64) error {

	if minimalMovingSpeed < 0.0 {
		return NewMinimalMovingSpeedLessThenZero(minimalMovingSpeed)
	}

	if minimalStepHight < 0.0 {
		return NewMinimalStepHightLessThenZero(minimalStepHight)
	}

	fillDistanceTimeAndSpeedValues(pnts, minimalMovingSpeed)
	err := fillCorrectedElevationTrackPoint(pnts, correction, minimalStepHight)
	if err != nil {
		return err
	}
	fillElevationGainLoseTrackPoint(pnts)
	fillCountUpDownWards(pnts, correction)

	return nil
}

// FillTrackSegmentValues - Fills the distance and attitude fields of a tack segment by adding up all TrackPoint distances
// All TrackPoint values has to be set before. See FillDistancesTrackPoint and FillValuesTrackPointArray
func FillTrackSegmentValues(segment *TrackSegment) {
	iPnts := []TrackSummaryProvider{}
	for i := range segment.TrackPoints {
		iPnt := TrackSummaryProvider(&segment.TrackPoints[i])
		iPnts = append(iPnts, iPnt)
	}

	fillTrackSummaryValues(segment, iPnts, true)
}

// FillTrackValues - Fills the distance and attitude fields of a tack  by adding up all TrackSegments distances
// All TrackSegment has to be set before. See FillTrackSegmentValues
func FillTrackValues(track *Track) {
	iSegs := []TrackSummaryProvider{}
	for i := range track.TrackSegments {
		iSeg := TrackSummaryProvider(&track.TrackSegments[i])
		iSegs = append(iSegs, iSeg)
	}

	fillTrackSummaryValues(track, iSegs, false)
}

// FillTrackFileValues - Fills the distance and attitude fields of a tack  by adding up all TrackSegments distances
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

// fillCorrectedElevationTrackPoint - Set the CorrectedElevation value in a list of TrackPoints
// Basically this will run a somthing algorithm over the Elevation
func fillCorrectedElevationTrackPoint(pnts []TrackPoint, correction string, minimalStepHight float64) error {

	switch correction {
	case GetValidCorrectionParameters()[1]:
		fillCorrectedElevationTrackPointLinear(pnts)
	case GetValidCorrectionParameters()[0]:
		fillCorrectedElevationTrackPointNone(pnts)
	case GetValidCorrectionParameters()[2]:
		fillCorrectedElevationTrackPointSteps(pnts, minimalStepHight)
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

func fillCorrectedElevationTrackPointSteps(pnts []TrackPoint, minimalStepHight float64) {
	numPnts := len(pnts)
	for i := range pnts {
		if i > 0 && i < (numPnts-1) {
			pnts[i].CorectedElevation = getCorrectedElevationSteps(pnts[i], pnts[i-1], pnts[i+1], minimalStepHight)
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

func fillDistanceTimeAndSpeedValues(pnts []TrackPoint, minimalMovingSpeed float64) {
	disToHere := float64(0.0)
	var movingTime time.Duration
	for i, pnt := range pnts {
		if pnt.TimeValid {

			if i > 0 {
				pnts[i].TimeDurationBefore = pnt.Time.Sub(pnts[i-1].Time)
				pnts[i].SpeedBefore = pnt.DistanceBefore / float64((pnts[i].TimeDurationBefore / time.Second))
				if pnts[i].SpeedBefore >= minimalMovingSpeed {
					movingTime = movingTime + pnts[i].TimeDurationBefore
					pnts[i].CountMoving = true
					disToHere += pnts[i].DistanceBefore
				} else {
					pnts[i].CountMoving = false
				}
			} else {
				// Make sure the first point counts as moving
				pnts[i].CountMoving = true
			}

			if i < (len(pnts) - 1) {
				pnts[i].TimeDurationNext = pnts[i+1].Time.Sub(pnt.Time)
				pnts[i].SpeedNext = pnt.DistanceNext / float64((pnts[i].TimeDurationNext / time.Second))
			}

			pnts[i].DistanceToThisPoint = disToHere
			pnts[i].MovingTime = movingTime
			if pnts[i].MovingTime > 0 {
				pnts[i].AvarageSpeed = pnts[i].DistanceToThisPoint / float64((pnts[i].MovingTime / time.Second))
			}
		} else {
			// If we can not calc the speed because of missing time info, all points count
			disToHere += pnts[i].DistanceBefore
			pnts[i].DistanceToThisPoint = disToHere
			pnts[i].CountMoving = true
		}
	}

}

// FillCountUpDownWards - Fills the CountUpwards and CountDownwards value
func fillCountUpDownWards(pnts []TrackPoint, correction string) {
	var upwardsTime time.Duration
	var downwardsTime time.Duration
	//numPnts := len(pnts)

	for i := range pnts {
		if correction == GetValidCorrectionParameters()[2] { // In case we do steps correction, CorectedElevation will make no sense
			if i > 0 { // Evaluation of the last point don't count
				eveDiff := pnts[i].Elevation - pnts[i-1].Elevation
				if eveDiff > 0 {
					pnts[i].CountUpwards = true
				}

				if eveDiff < 0 {
					pnts[i].CountDownwards = true
				}
			}
		} else {

			if i > 0 { // Evaluation of the first point don't count
				eveDiff := pnts[i].CorectedElevation - pnts[i-1].CorectedElevation
				if eveDiff > 0 {
					pnts[i].CountUpwards = true
				}

				if eveDiff < 0 {
					pnts[i].CountDownwards = true
				}
			}
		}
		if i > 0 && pnts[i].CountMoving && pnts[i].TimeValid { // Time of the first point don't count
			if pnts[i].CountDownwards {
				downwardsTime = downwardsTime + pnts[i].TimeDurationBefore
			}

			if pnts[i].CountUpwards {
				upwardsTime = upwardsTime + pnts[i].TimeDurationBefore
			}
		}

		pnts[i].UpwardsTime = upwardsTime
		pnts[i].DownwardsTime = downwardsTime
	}

}

func fillTrackSummaryValues(target TrackSummarySetter, input []TrackSummaryProvider, inputIsPoints bool) {
	var dist float64
	var horizontalDist float64
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
	var upwardsTimeSum time.Duration
	var downwarsTimeSum time.Duration
	var upwardsTime time.Duration
	var downwarsTime time.Duration

	for i, sum := range input {
		dist = dist + sum.GetDistance()
		horizontalDist = horizontalDist + sum.GetHorizontalDistance()
		elevationGain = elevationGain + sum.GetElevationGain()
		elevationLose = elevationLose + sum.GetElevationLose()
		upwardsDistance = upwardsDistance + sum.GetUpwardsDistance()
		downwardsDistance = downwardsDistance + sum.GetDownwardsDistance()
		movingTimeSum = movingTimeSum + sum.GetMovingTime()
		downwarsTimeSum = downwarsTimeSum + sum.GetDownwardsTime()
		upwardsTimeSum = upwardsTimeSum + sum.GetUpwardsTime()

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
			upwardsTime = input[len(input)-1].GetUpwardsTime()
			downwarsTime = input[len(input)-1].GetDownwardsTime()
		} else {
			movingTime = movingTimeSum
			downwarsTime = downwarsTimeSum
			upwardsTime = upwardsTimeSum
		}
	}
	target.SetValues(dist, horizontalDist, minimumAltitude, maximumAltitude, elevationGain, elevationLose, upwardsDistance, downwardsDistance,
		timeDataValid, startTime, endTime, movingTime, upwardsTime, downwarsTime)
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

func getCorrectedElevationSteps(basePoint TrackPoint, beforePoint TrackPoint, nextPoint TrackPoint, minimalStepHight float64) float32 {

	eveDiffBefore := basePoint.Elevation - beforePoint.CorectedElevation
	eveDiffAfter := nextPoint.Elevation - basePoint.Elevation
	sameDirection := eveDiffBefore * eveDiffAfter

	if math.Abs(float64(eveDiffBefore)) >= minimalStepHight && (sameDirection > 0 || math.Abs(float64(eveDiffAfter)) < minimalStepHight) {
		return basePoint.Elevation
	} else if minimalStepHight == 0.0 {
		return basePoint.Elevation
	}

	return beforePoint.CorectedElevation
}
