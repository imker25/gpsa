package gpsabl

import (
	"testing"

	"tobi.backfrak.de/internal/testhelper"
)

func TestAnalyseSimpeTrack(t *testing.T) {
	gpx, err := ReadGPX(testhelper.GetValideGPX("01.gpx"))
	if err != nil {
		t.Errorf("Something wrong when reading a valide gpx file %s: %s", testhelper.GetValideGPX("01.gpx"), err.Error())
	}

	info := GetTrackInfo(gpx.Tracks[0])

	if info.Track.Name != gpx.Tracks[0].Name {
		t.Errorf("TrackInfo.Track is not the expected one")
	}

	if info.NumberOfSegments != 1 {
		t.Errorf("Number of segments was not 1 as expecpected, it was %d", info.NumberOfSegments)
	}

	if info.Name != gpx.Tracks[0].Name {
		t.Errorf("The TrackInfo.Name was not %s as expected, it was %s", gpx.Tracks[0].Name, info.Name)
	}

	if info.Description != gpx.Tracks[0].Description {
		t.Errorf("The TrackInfo.Description was not %s as expected, it was %s", gpx.Tracks[0].Description, info.Description)
	}

	if info.NumberOfTrackPoints != 637 {
		t.Errorf("The TrackInfo.NumberOfTrackPoints was not not %d as expected, it was %d", 637, info.NumberOfTrackPoints)
	}

}
