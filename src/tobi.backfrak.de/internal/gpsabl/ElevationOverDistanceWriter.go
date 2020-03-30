package gpsabl

import (
	"fmt"
	"os"
)

// WriteElevationOverDistance - Write out a Elevation over Distance table
func WriteElevationOverDistance(trackFile TrackFile, outFile *os.File) error {

	pnts := getTrackPoints(trackFile)
	lines := getOutPutLines(pnts)

	for _, line := range lines {
		_, err := outFile.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}

func getOutPutLines(trackPoints []TrackPoint) []string {
	lines := []string{"Distance [km];Elevation [m];CorrectedElevation [m];"}

	for _, pnt := range trackPoints {

		line := fmt.Sprintf("%f;%f;%f;%s",
			RoundFloat64To2Digits(pnt.DistanceToThisPoint/1000),
			RoundFloat64To2Digits(float64(pnt.Elevation)),
			RoundFloat64To2Digits(float64(pnt.CorectedElevation)),
			GetNewLine())

		lines = append(lines, line)
	}

	return lines
}

func getTrackPoints(trackFile TrackFile) []TrackPoint {

	pnts := []TrackPoint{}

	for _, track := range trackFile.Tracks {
		for _, seg := range track.TrackSegments {
			pnts = append(pnts, seg.TrackPoints...)
		}
	}

	return pnts
}
