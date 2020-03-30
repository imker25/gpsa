package gpsabl

import "testing"

func TestGetOutPutLines(t *testing.T) {
	file := getTrackFileTwoTracksWithThreeSegments()
	pnts := getTrackPoints(file)

	lines := getOutPutLines(pnts)

	if len(lines) != 10 {
		t.Errorf("The number of lines %d is not the expected value 10", len(lines))
	}
}
