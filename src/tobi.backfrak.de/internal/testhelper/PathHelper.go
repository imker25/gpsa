package testhelper

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetProjectRoot - Get the root folder of the gpsa project
func GetProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	for !strings.HasSuffix(wd, "gpsa") {
		wd = filepath.Dir(wd)
	}

	return wd
}

// GetValideGPX - Get the file path to a valide gpx file with the given name
func GetValideGPX(name string) string {
	rootDir := GetProjectRoot()

	return filepath.Join(rootDir, "testdata", "valide-gpx", name)
}

// GetUnValideGPX - Get the file path to a valide gpx file with the given name
func GetUnValideGPX(name string) string {
	rootDir := GetProjectRoot()

	return filepath.Join(rootDir, "testdata", "unvalide-gpx", name)
}
