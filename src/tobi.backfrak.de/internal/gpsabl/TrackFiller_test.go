package gpsabl

import "testing"

func TestFillDistancesThreePoints(t *testing.T) {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	pnt2 = FillDistancesTrackPoint(pnt2, pnt1, pnt3)

	if pnt2.VerticalDistanceBefore != -1.0 {
		t.Errorf("The VerticalDistanceBefore is %f, but %f was expected", pnt2.VerticalDistanceBefore, -1.0)
	}

	if pnt2.HorizontalDistanceBefore != pnt2.HorizontalDistanceNext {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was expected", pnt2.HorizontalDistanceBefore, pnt2.HorizontalDistanceNext)
	}

	if pnt2.VerticalDistanceNext != 1.0 {
		t.Errorf("The VerticalDistanceNext is %f, but %f was expected", pnt2.VerticalDistanceNext, 1.0)
	}
}

func TestFillDistancesTwoPointBefore(t *testing.T) {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := TrackPoint{}

	pnt2 = FillDistancesTrackPoint(pnt2, pnt1, pnt3)

	if pnt2.VerticalDistanceBefore != -1.0 {
		t.Errorf("The VerticalDistanceBefore is %f, but %f was expected", pnt2.VerticalDistanceBefore, -1.0)
	}

	if pnt2.HorizontalDistanceBefore == 0.0 {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was not expected", pnt2.HorizontalDistanceBefore, 0.0)
	}

	if pnt2.HorizontalDistanceNext != 0.0 {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was expected", pnt2.HorizontalDistanceNext, 0.0)
	}

	if pnt2.VerticalDistanceNext != 0.0 {
		t.Errorf("The VerticalDistanceNext is %f, but %f was expected", pnt2.VerticalDistanceNext, 0.0)
	}
}

func TestFillDistancesTwoPointNext(t *testing.T) {
	pnt1 := TrackPoint{}
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)
	pnt3 := getTrackPoint(50.11484790, 8.684885500, 109.0)

	pnt2 = FillDistancesTrackPoint(pnt2, pnt1, pnt3)

	if pnt2.VerticalDistanceBefore != 0.0 {
		t.Errorf("The VerticalDistanceBefore is %f, but %f was expected", pnt2.VerticalDistanceBefore, -1.0)
	}

	if pnt2.HorizontalDistanceBefore != 0.0 {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was expected", pnt2.HorizontalDistanceBefore, 0.0)
	}

	if pnt2.HorizontalDistanceNext == 0.0 {
		t.Errorf("The HorizontalDistanceBefore is %f, but %f was not expected", pnt2.HorizontalDistanceNext, 0.0)
	}

	if pnt2.VerticalDistanceNext != 1.0 {
		t.Errorf("The VerticalDistanceNext is %f, but %f was expected", pnt2.VerticalDistanceNext, 0.0)
	}
}
