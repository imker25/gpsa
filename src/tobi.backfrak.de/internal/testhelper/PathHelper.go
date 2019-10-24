package testhelper

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
