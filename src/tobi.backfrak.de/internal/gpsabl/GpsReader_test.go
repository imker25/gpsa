package gpsabl

import (
	"fmt"
	"testing"

	"tobi.backfrak.de/internal/testhelper"
)

func TestReadValideGPX(t *testing.T) {
	file := testhelper.GetValideGPX("01.gpx")

	if file == "" {
		t.Errorf("Test failed, expected not to get an empty string")
	}

	gpx, err := ReadGPX(file)
	if err != nil {
		t.Errorf("Something wrong when reading a valide gpx file: %s", err.Error())
	}
	fmt.Println(gpx.Name)
	fmt.Println(gpx.Description)
}
