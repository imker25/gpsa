package gpsabl

import "testing"

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestHaversineDistanceZeroDegree(t *testing.T) {
	pnt1 := getTrackPoint(30.0, 30.0, 100.0)
	pnt2 := getTrackPoint(30.0, 30.0, 100.0)

	dist := HaversineDistance(pnt1, pnt2)
	if dist != 0.0 {
		t.Errorf("The distance was calculated with %f but %f was expected", dist, 0.0)
	}
}

func TestHaversineDistanceOneDegreeLat(t *testing.T) {
	pnt1 := getTrackPoint(31.0, 30.0, 100.0)
	pnt2 := getTrackPoint(30.0, 30.0, 100.0)

	dist := HaversineDistance(pnt1, pnt2)
	if dist != 111196.67197381073 {
		t.Errorf("The distance was calculated with %f but %f was expected", dist, 111196.67197381073)
	}
}

func TestHaversineDistanceOneDegreeLon(t *testing.T) {
	pnt1 := getTrackPoint(30.0, 30.0, 100.0)
	pnt2 := getTrackPoint(30.0, 31.0, 100.0)

	dist := HaversineDistance(pnt1, pnt2)
	if dist != 96298.83717228581 {
		t.Errorf("The distance was calculated with %f but %f was expected", dist, 96298.83717228581)
	}
}

func TestHaversineDistanceMunichColone(t *testing.T) {
	pnt1 := getTrackPoint(48.13992070, 11.56654350, 529.0)
	pnt2 := getTrackPoint(50.94169280, 6.959409710, 56.0)

	dist := HaversineDistance(pnt1, pnt2)
	if dist != 455454.9517128867 {
		t.Errorf("The distance was calculated with %f but %f was expected", dist, 455454.9517128867)
	}
}

func TestHaversineShortDistance1(t *testing.T) {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)

	dist := HaversineDistance(pnt1, pnt2)
	if dist != 11.900633553301818 {
		t.Errorf("The distance was calculated with %f but %f was expected", dist, 11.900633553301818)
	}
}
func TestHaversineShortDistance2(t *testing.T) {
	pnt2 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt1 := getTrackPoint(50.11495750, 8.684874770, 108.0)

	dist := HaversineDistance(pnt1, pnt2)
	if dist != 11.900633553301818 {
		t.Errorf("The distance was calculated with %f but %f was expected", dist, 11.900633553301818)
	}
}

func TestHaversineDistanceHalfEquator(t *testing.T) {
	pnt1 := getTrackPoint(0.0, 0.0, 0.0)
	pnt2 := getTrackPoint(180.0, 0.0, 0.0)

	dist := HaversineDistance(pnt1, pnt2)
	if dist != 20015086.79602057 {
		t.Errorf("The distance was calculated with %f but %f was expected", dist, 11.900633553301818)
	}
}

func TestHaversineDistanceQuaterMeridian(t *testing.T) {
	pnt1 := getTrackPoint(0.0, 90.0, 0.0)
	pnt2 := getTrackPoint(0.0, 0.0, 0.0)

	dist := HaversineDistance(pnt1, pnt2)
	if dist != 10007543.398010286 {
		t.Errorf("The distance was calculated with %f but %f was expected", dist, 11.900633553301818)
	}
}

func TestDistanceMunichColone(t *testing.T) {
	pnt1 := getTrackPoint(48.13992070, 11.56654350, 529.0)
	pnt2 := getTrackPoint(50.94169280, 6.959409710, 56.0)

	dist1 := HaversineDistance(pnt1, pnt2)
	dist2 := Distance(pnt1, pnt2)
	// As the distance is bigger than 33km elevation gain will be ignored
	if dist2 != dist1 {
		t.Errorf("The distance was calculated with %f but %f was expected", dist2, dist1)
	}
}

func TestShortDistance1(t *testing.T) {
	pnt1 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt2 := getTrackPoint(50.11495750, 8.684874770, 108.0)

	dist1 := HaversineDistance(pnt1, pnt2)
	dist2 := Distance(pnt1, pnt2)
	// As the distance is smaller than 33km elevation gain will not be ignored
	if dist2 == dist1 {
		t.Errorf("The distance was calculated with %f but %f was not expected", dist2, dist1)
	}
}

func TestShortDistance2(t *testing.T) {
	pnt2 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt1 := getTrackPoint(50.11495750, 8.684874770, 108.0)

	dist1 := HaversineDistance(pnt1, pnt2)
	dist2 := Distance(pnt1, pnt2)
	// As the distance is smaller than 33km elevation gain will not be ignored
	if dist2 == dist1 {
		t.Errorf("The distance was calculated with %f but %f was not expected", dist2, dist1)
	}
}

func TestShortDistance3(t *testing.T) {
	pnt2 := getTrackPoint(50.11484790, 8.684885500, 109.0)
	pnt1 := getTrackPoint(50.11495750, 8.684874770, 108.0)

	dist1 := Distance(pnt2, pnt1)
	dist2 := Distance(pnt1, pnt2)
	// As the distance is smaller than 33km elevation gain will not be ignored
	if dist2 != dist1 {
		t.Errorf("The distance was calculated with %f but %f was expected", dist2, dist1)
	}
}

func getTrackPoint(lat, lon, ele float32) TrackPoint {
	pnt := TrackPoint{}
	pnt.Latitude = lat
	pnt.Longitude = lon
	pnt.Elevation = ele

	return pnt
}
