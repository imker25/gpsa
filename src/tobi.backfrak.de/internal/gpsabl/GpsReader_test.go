package gpsabl

import (
	"testing"

	"tobi.backfrak.de/internal/testhelper"
)

func TestReadGPX(t *testing.T) {
	file := testhelper.GetValideGPX("01.gpx")

	if file == "" {
		t.Errorf("Test failed, expected not to get an empty string")
	}

	ReadGPX(file)

}
