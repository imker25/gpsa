package gpsabl

import (
	"testing"

	"tobi.backfrak.de/internal/testhelper"
)

func TestGetTrackAtituteInfoJustOnePointGPX11(t *testing.T) {
	gpx, err := ReadGPX(testhelper.GetValideGPX("09.gpx"))
	if err != nil {
		t.Errorf("Something wrong when reading a valide gpx file %s: %s", testhelper.GetValideGPX("01.gpx"), err.Error())
	}

	info := GetTrackInfo(gpx.Tracks[0])

	if info.MinimumAtitute != 11.1 {
		t.Errorf("The TrackInfo.MinimumAtitute was not not %f as expected, it was %f", 11.1, info.MinimumAtitute)
	}

	if info.MaximumAtitute != 11.1 {
		t.Errorf("The TrackInfo.MaximumAtitute was not not %f as expected, it was %f", 11.1, info.MaximumAtitute)
	}

	if info.AtituteRange != 0.0 {
		t.Errorf("The TrackInfo.AtituteRange was not not %f as expected, it was %f", 0.0, info.AtituteRange)
	}
}

func TestGetTrackAtituteInfoJustOnePointGPX10(t *testing.T) {
	gpx, err := ReadGPX(testhelper.GetValideGPX("08.gpx"))
	if err != nil {
		t.Errorf("Something wrong when reading a valide gpx file %s: %s", testhelper.GetValideGPX("01.gpx"), err.Error())
	}

	info := GetTrackInfo(gpx.Tracks[0])

	if info.MinimumAtitute != 11.1 {
		t.Errorf("The TrackInfo.MinimumAtitute was not not %f as expected, it was %f", 11.1, info.MinimumAtitute)
	}

	if info.MaximumAtitute != 11.1 {
		t.Errorf("The TrackInfo.MaximumAtitute was not not %f as expected, it was %f", 11.1, info.MaximumAtitute)
	}

	if info.AtituteRange != 0.0 {
		t.Errorf("The TrackInfo.AtituteRange was not not %f as expected, it was %f", 0.0, info.AtituteRange)
	}
}

func TestGetTrackAtituteInfo(t *testing.T) {
	gpx, err := ReadGPX(testhelper.GetValideGPX("01.gpx"))
	if err != nil {
		t.Errorf("Something wrong when reading a valide gpx file %s: %s", testhelper.GetValideGPX("01.gpx"), err.Error())
	}

	info := GetTrackInfo(gpx.Tracks[0])

	if info.MinimumAtitute != 298.00 {
		t.Errorf("The TrackInfo.MinimumAtitute was not not %f as expected, it was %f", 298.00, info.MinimumAtitute)
	}

	if info.MaximumAtitute != 402.00 {
		t.Errorf("The TrackInfo.MaximumAtitute was not not %f as expected, it was %f", 402.00, info.MaximumAtitute)
	}

	if info.AtituteRange != 104.00 {
		t.Errorf("The TrackInfo.AtituteRange was not not %f as expected, it was %f", 104.00, info.AtituteRange)
	}

}

func TestGetTrackPointInfo(t *testing.T) {
	gpx, err := ReadGPX(testhelper.GetValideGPX("01.gpx"))
	if err != nil {
		t.Errorf("Something wrong when reading a valide gpx file %s: %s", testhelper.GetValideGPX("01.gpx"), err.Error())
	}

	info := GetTrackInfo(gpx.Tracks[0])

	if info.NumberOfTrackPoints != 637 {
		t.Errorf("The TrackInfo.NumberOfTrackPoints was not not %d as expected, it was %d", 637, info.NumberOfTrackPoints)
	}

	if len(info.GetAllTrackPoints()) != info.NumberOfTrackPoints {
		t.Errorf("The number of elements in info.GetAllTrackPoints() is not the same as the info.NumberOfTrackPoints")
	}
}

func TestGetBasicTrackInfo(t *testing.T) {
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
}
