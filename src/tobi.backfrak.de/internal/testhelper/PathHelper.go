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

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// GetProjectRoot - Get the root folder of the gpsa project
func GetProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	orgWD := wd
	for !isRootPath(wd) {
		wd = filepath.Dir(wd)
		if wd == "/" || strings.HasSuffix(wd, ":\\") { // If we reach the root folder (of a drive on windows), we think try to return the orignial value as a gues
			return orgWD
		}
	}

	return wd
}

func isRootPath(wd string) bool {

	file := filepath.Join(wd, "testdata", "valid-gpx", "01.gpx")
	if fileExists(file) {
		return true
	}

	return false
}

// GetValidGPX - Get the file path to a valid gpx file with the given name
func GetValidGPX(name string) string {
	rootDir := GetProjectRoot()

	return filepath.Join(rootDir, "testdata", "valid-gpx", name)
}

// GetInvalidGPX - Get the file path to a valid gpx file with the given name
func GetInvalidGPX(name string) string {
	rootDir := GetProjectRoot()

	return filepath.Join(rootDir, "testdata", "invalid-gpx", name)
}

// GetValidTcx - Get the file path to a valid gpx file with the given name
func GetValidTcx(name string) string {
	rootDir := GetProjectRoot()

	return filepath.Join(rootDir, "testdata", "valid-tcx", name)
}

// GetInvalidTcx - Get the file path to a valid gpx file with the given name
func GetInvalidTcx(name string) string {
	rootDir := GetProjectRoot()

	return filepath.Join(rootDir, "testdata", "invalid-tcx", name)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	if info.IsDir() {
		return false
	}
	return true
}
