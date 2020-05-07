package tcxbl

import (
	"testing"
	"time"

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

	if trackFile.GetMovingTime() != 1103000000000 {
		t.Errorf("The MovingTime is %d, but should be %d", trackFile.GetMovingTime(), 1103000000000)
	}

	if trackFile.GetStartTime().Format(time.RFC3339) != "2016-06-05T10:45:59Z" {
		t.Errorf("The StartTime is %s, but should be %s", trackFile.GetStartTime().Format(time.RFC3339), "2016-06-05T10:45:59Z")
	}

	if trackFile.GetUpwardsSpeed() != 4.797203586378565 {
		t.Errorf("The UpwardsSpeed is %f, but should be %f", trackFile.GetUpwardsSpeed(), 4.797204)
	}

	if trackFile.GetDownwardsSpeed() != 7.683785245071919 {
		t.Errorf("The UpwardsSpeed is %f, but should be %f", trackFile.GetDownwardsSpeed(), 7.683785)
	}
}

func TestConvertValidTcx03(t *testing.T) {
	tcx, readErr := ReadTcx(testhelper.GetValidTcx("03.tcx"))
	if readErr != nil {
		t.Errorf("ReadError, but none expected")
	}
	trackFile, convertErr := ConvertTcx(tcx, testhelper.GetValidTcx("03.tcx"), "none", 0.3, 10)
	if convertErr != nil {
		t.Errorf("ConvertError, but none expected")
	}
	if trackFile.GetDistance() != 1250.000000 {
		t.Errorf("The Distance is %f, but should be %f", trackFile.GetDistance(), 1250.000000)
	}
	if trackFile.GetMovingTime() != 3000*time.Second {
		t.Errorf("The MovingTime is %s, but should be %s", trackFile.GetMovingTime(), 3000*time.Second)
	}

	if trackFile.GetStartTime().Format(time.RFC3339) != "2019-02-05T09:17:02Z" {
		t.Errorf("The StartTime is %s, but should be %s", trackFile.GetStartTime().Format(time.RFC3339), "2019-02-05T09:17:02Z")
	}

	if trackFile.GetEndTime().Format(time.RFC3339) != "2019-02-05T10:07:02Z" {
		t.Errorf("The StartTime is %s, but should be %s", trackFile.GetEndTime().Format(time.RFC3339), "2019-02-05T10:07:02Z")
	}

	if trackFile.GetUpwardsSpeed() != 0.0 {
		t.Errorf("The UpwardsSpeed is %f, but should be %f", trackFile.GetUpwardsSpeed(), 0.0)
	}

	if trackFile.GetDownwardsSpeed() != 0.0 {
		t.Errorf("The UpwardsSpeed is %f, but should be %f", trackFile.GetDownwardsSpeed(), 0.0)
	}
}
