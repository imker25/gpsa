package main

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"fmt"
	"testing"
)

func TestMainNoArgs(t *testing.T) {
	main()

	fmt.Println("OK")
}
