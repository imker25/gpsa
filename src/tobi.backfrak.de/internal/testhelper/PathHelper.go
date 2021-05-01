package testhelper

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"bufio"
	"fmt"
	"io"
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

func GetValidGpxBuffer(name string) ([]byte, error) {
	filePath1 := GetValidGPX(name)
	file1, _ := os.Open(filePath1)

	var inputBytes []byte
	reader1 := bufio.NewReader(file1)
	for {
		input, errRead1 := reader1.ReadByte()
		if errRead1 != nil {
			if errRead1 == io.EOF {
				break
			} else {
				return nil, errRead1
			}
		}

		inputBytes = append(inputBytes, input)
	}

	return inputBytes, nil
}

func GetInvalidGpxBuffer(name string) ([]byte, error) {
	filePath1 := GetInvalidGPX(name)
	file1, _ := os.Open(filePath1)

	var inputBytes []byte
	reader1 := bufio.NewReader(file1)
	for {
		input, errRead1 := reader1.ReadByte()
		if errRead1 != nil {
			if errRead1 == io.EOF {
				break
			} else {
				return nil, errRead1
			}
		}

		inputBytes = append(inputBytes, input)
	}

	return inputBytes, nil
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

func GetValidTcxBuffer(name string) ([]byte, error) {
	filePath1 := GetValidTcx(name)
	file1, _ := os.Open(filePath1)

	var inputBytes []byte
	reader1 := bufio.NewReader(file1)
	for {
		input, errRead1 := reader1.ReadByte()
		if errRead1 != nil {
			if errRead1 == io.EOF {
				break
			} else {
				return nil, errRead1
			}
		}

		inputBytes = append(inputBytes, input)
	}

	return inputBytes, nil
}

func GetInvalidTcxBuffer(name string) ([]byte, error) {
	filePath1 := GetInvalidTcx(name)
	file1, _ := os.Open(filePath1)

	var inputBytes []byte
	reader1 := bufio.NewReader(file1)
	for {
		input, errRead1 := reader1.ReadByte()
		if errRead1 != nil {
			if errRead1 == io.EOF {
				break
			} else {
				return nil, errRead1
			}
		}

		inputBytes = append(inputBytes, input)
	}

	return inputBytes, nil
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
