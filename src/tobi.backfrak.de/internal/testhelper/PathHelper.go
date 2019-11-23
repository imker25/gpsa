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

	file := filepath.Join(wd, "testdata", "valide-gpx", "01.gpx")
	if fileExists(file) {
		return true
	}

	return false
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
