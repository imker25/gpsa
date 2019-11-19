package gpsabl

import "testing"

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestRoundFloat64To2Digits(t *testing.T) {
	if RoundFloat64To2Digits(23.123) != 23.12 {
		t.Errorf("Round down does not work")
	}

	if RoundFloat64To2Digits(23.127) != 23.13 {
		t.Errorf("Round up does not work")
	}

	if RoundFloat64To2Digits(23.1200) != 23.12 {
		t.Errorf("Round 0 does not work")
	}
}

func TestRoundFloat64To4Digits(t *testing.T) {
	if RoundFloat64To4Digits(23.12344) != 23.1234 {
		t.Errorf("Round down does not work")
	}

	if RoundFloat64To4Digits(23.12717) != 23.1272 {
		t.Errorf("Round up does not work")
	}

	if RoundFloat64To4Digits(23.123300) != 23.1233 {
		t.Errorf("Round 0 does not work")
	}
}

func TestCompareFloat64With4Digits(t *testing.T) {

	if !CompareFloat64With4Digits(23.12344, 23.1234) {
		t.Errorf("Round down does not work")
	}

	if !CompareFloat64With4Digits(23.12348, 23.1235) {
		t.Errorf("Round up does not work")
	}

	if !CompareFloat64With4Digits(23.123700, 23.1237) {
		t.Errorf("Round 0 does not work")
	}
}
