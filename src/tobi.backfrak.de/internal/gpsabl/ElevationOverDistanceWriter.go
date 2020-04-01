package gpsabl

import (
	"fmt"
	"os"
)

type printElevationPoint struct {
	DistanceToThisPoint float64
	Elevation           float32
	CorectedElevation   float32
}

// WriteElevationOverDistance - Write out a Elevation over Distance table
func WriteElevationOverDistance(trackFile TrackFile, outFile *os.File, outputSeperator string) error {

	pnts := getTrackPoints(trackFile)
	lines := getOutPutLines(pnts, outputSeperator)

	for _, line := range lines {
		_, err := outFile.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}

func getOutPutLines(trackPoints []printElevationPoint, outputSeperator string) []string {
	lines := []string{fmt.Sprintf("Distance [km];Elevation [m];CorrectedElevation [m];%s", GetNewLine())}

	for _, pnt := range trackPoints {

		line := fmt.Sprintf("%f%s%f%s%f%s%s",
			RoundFloat64To2Digits(pnt.DistanceToThisPoint/1000), outputSeperator,
			RoundFloat64To2Digits(float64(pnt.Elevation)), outputSeperator,
			RoundFloat64To2Digits(float64(pnt.CorectedElevation)), outputSeperator,
			GetNewLine())

		lines = append(lines, line)
	}

	return lines
}

func getTrackPoints(trackFile TrackFile) []printElevationPoint {

	pnts := []printElevationPoint{}
	startDist := 0.0

	for _, track := range trackFile.Tracks {
		for _, seg := range track.TrackSegments {
			pntCount := len(seg.TrackPoints)
			for i, tPnt := range seg.TrackPoints {
				pnt := printElevationPoint{}
				pnt.DistanceToThisPoint = startDist + tPnt.DistanceToThisPoint
				pnt.Elevation = tPnt.Elevation
				pnt.CorectedElevation = tPnt.CorectedElevation
				pnts = append(pnts, pnt)

				if i == pntCount-1 {
					startDist = startDist + tPnt.DistanceToThisPoint
				}
			}
		}
	}

	return pnts
}
