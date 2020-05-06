package tcxbl

import (
	"testing"

	"tobi.backfrak.de/internal/testhelper"
)

func TestConvertValidTcx01(t *testing.T) {
	tcx, readErr := ReadTcx(testhelper.GetValidTcx("01.tcx"))
	if readErr != nil {
		t.Errorf("ReadError, but none expected")
	}
	trackFile, convertErr := ConvertTcx(tcx, testhelper.GetValidTcx("01.tcx"), "none", 0.3, 10)
	if convertErr != nil {
		t.Errorf("ConvertError, but none expected")
	}
	if trackFile.GetDistance() != 6216.201383825188 {
		t.Errorf("The Distance is %f, but should be %f", trackFile.GetDistance(), 6216.201383825188)
	}
}
