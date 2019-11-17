package gpsabl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import "math"

// EarthRadius -  The Radius of the earth in meter
const EarthRadius = 6371 * 1000

//ToRad converts angles in ° to radiant
func ToRad(x float64) float64 {
	return x / 180. * math.Pi
}

// HaversineDistance - Calcs the distance between two TrackPoints in Meter.
// Assuming a spherical earth.
// Don't use this function for distance because it will ignore elevation gain
func HaversineDistance(pnt1, pnt2 TrackPoint) float64 {
	dLat := ToRad(float64(pnt1.Latitude - pnt2.Latitude))
	dLon := ToRad(float64(pnt1.Longitude - pnt2.Longitude))
	thisLat1 := ToRad(float64(pnt1.Latitude))
	thisLat2 := ToRad(float64(pnt2.Latitude))

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(thisLat1)*math.Cos(thisLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return float64(EarthRadius+((pnt1.Elevation+pnt2.Elevation)/2)) * c
}

// Distance -  Calcs the distance between two TrackPoints in Meter.
// Assuming a spherical earth.
func Distance(pnt1, pnt2 TrackPoint) float64 {
	dist := HaversineDistance(pnt1, pnt2)

	// When the distance is bigger then 33km, elevation gain will be ignored
	if dist >= 33000.0 {
		return dist
	}

	dEve := math.Abs(float64(pnt1.Elevation - pnt2.Elevation))

	return math.Sqrt(math.Pow(dist, 2) + math.Pow(dEve, 2))
}
