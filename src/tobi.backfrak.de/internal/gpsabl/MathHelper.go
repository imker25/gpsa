package gpsabl

import "math"

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// CompareFloat64With4Digits - Compare two float64 to 4 digits after decimal
func CompareFloat64With4Digits(in1, in2 float64) bool {
	return RoundFloat64To4Digits(in1) == RoundFloat64To4Digits(in2)
}

// RoundFloat64To4Digits - Rounds a float64 to 4 digits after decimal
func RoundFloat64To4Digits(in float64) float64 {
	return math.Round(in*10000) / 10000
}

// RoundFloat64To2Digits - Rounds a float64 to 2 digits after decimal
func RoundFloat64To2Digits(in float64) float64 {
	return math.Round(in*100) / 100
}
