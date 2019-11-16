package main

import "testing"

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestHandleComandlineOptions(t *testing.T) {
	handleComandlineOptions()

	if HelpFlag == true {
		t.Errorf("The HelpFlag is true, but should not")
	}
}

func TestHandleError(t *testing.T) {
	HandleError(nil, "my/path")

}
