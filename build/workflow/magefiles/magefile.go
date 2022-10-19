//go:build mage
// +build mage

package main

// Copyright 2022 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

const VERSION_FILE = "VersionMaster.txt"
const PROJECT_NAME_SPACE = "tobi.backfrak.de"

type buildContext struct {
	GitHight             int
	GitHash              string
	ProgramVersion       string
	ProgramVersionNumber string
	BinDir               string
	PackageDir           string
	LogDir               string
	WorkDir              string
	SourceDir            string
	PackagesToBuild      []string
	PackagesToTest       []string
	BuildZipPath         string
	BinZipPath           string
	VersionFilePath      string
}

var gpsaBuildContext buildContext

func getEnvironment() error {
	workDir, errWorkDir := os.Getwd()
	if errWorkDir != nil {
		return errWorkDir
	}

	if strings.HasSuffix(workDir, "build\\workflow\\") || strings.HasSuffix(workDir, "build\\workflow") || strings.HasSuffix(workDir, "build/workflow/") || strings.HasSuffix(workDir, "build/workflow") {
		workDir = filepath.Join(workDir, "..", "..")
	}

	gpsaBuildContext.WorkDir = workDir
	gpsaBuildContext.BinDir = filepath.Join(workDir, "bin")
	gpsaBuildContext.LogDir = filepath.Join(workDir, "logs")
	gpsaBuildContext.PackageDir = filepath.Join(workDir, "pkg")
	gpsaBuildContext.SourceDir = filepath.Join(workDir, "src")
	gpsaBuildContext.BuildZipPath = filepath.Join(gpsaBuildContext.WorkDir, fmt.Sprintf("%s_Build.zip", runtime.GOOS))
	gpsaBuildContext.BinZipPath = filepath.Join(gpsaBuildContext.WorkDir, fmt.Sprintf("%s_bin.zip", runtime.GOOS))
	gpsaBuildContext.VersionFilePath = filepath.Join(gpsaBuildContext.WorkDir, VERSION_FILE)

	hash, errHash := GetGitHash(gpsaBuildContext.WorkDir)
	if errHash != nil {
		return errHash
	}
	// fmt.Println(fmt.Sprintf("Git Hash: %s", hash))
	gpsaBuildContext.GitHash = hash

	hight, errHight := GetGitHeight(VERSION_FILE, gpsaBuildContext.WorkDir)
	if errHight != nil {
		return errHight
	}
	// fmt.Println(fmt.Sprintf("Git Hight: %d", hight))
	gpsaBuildContext.GitHight = hight

	givenVersion, errVersion := readVersionMaster()
	if errVersion != nil {
		return errVersion
	}

	gpsaBuildContext.ProgramVersionNumber = fmt.Sprintf("%s.%d", givenVersion, hight)
	gpsaBuildContext.ProgramVersion = fmt.Sprintf("%s.%d-%s", givenVersion, hight, hash)

	fmt.Println(fmt.Sprintf("Run gpsa build workflow for V%s", gpsaBuildContext.ProgramVersionNumber))

	var errFinBuild error
	gpsaBuildContext.PackagesToBuild, errFinBuild = FindPackagesToBuild(filepath.Join(gpsaBuildContext.SourceDir, PROJECT_NAME_SPACE, "cmd"))
	if errFinBuild != nil {
		return errFinBuild
	}

	var errFinTest error
	gpsaBuildContext.PackagesToTest, errFinTest = FindPackagesToTest(filepath.Join(gpsaBuildContext.SourceDir, PROJECT_NAME_SPACE))
	if errFinTest != nil {
		return errFinTest
	}

	return nil
}

// Get the build name files
func GetBuildName() error {
	mg.Deps(getEnvironment, Clean)
	fmt.Println(fmt.Sprintf("Create gpsa Version files..."))
	fmt.Println("# ########################################################################################")

	buildNumber := gpsaBuildContext.ProgramVersion
	if os.Getenv("BUILD_NUMBER") != "" {
		buildNumber = fmt.Sprintf("%s-%s", gpsaBuildContext.ProgramVersion, os.Getenv("BUILD_NUMBER"))
	}

	if _, err := os.Stat(gpsaBuildContext.LogDir); os.IsNotExist(err) {
		errCreate := os.Mkdir(gpsaBuildContext.LogDir, 0755)
		if errCreate != nil {
			return errCreate
		}
	}

	errNr := ioutil.WriteFile(filepath.Join(gpsaBuildContext.LogDir, "BuildName.txt"), []byte(buildNumber), 0644)
	if errNr != nil {
		return errNr
	}

	errVersion := ioutil.WriteFile(filepath.Join(gpsaBuildContext.LogDir, "Version.txt"), []byte(gpsaBuildContext.ProgramVersionNumber), 0644)
	if errVersion != nil {
		return errVersion
	}

	errDum := ioutil.WriteFile(filepath.Join(gpsaBuildContext.LogDir, "dummy.json"), []byte("{\"key\":\"value\"}"), 0644)
	if errDum != nil {
		return errDum
	}

	fmt.Println("# ########################################################################################")
	return nil
}

// Compiles the project
func Build() error {
	mg.Deps(getEnvironment, Clean, GetBuildName)
	fmt.Println(fmt.Sprintf("Building gpsa V%s ...", gpsaBuildContext.ProgramVersion))
	fmt.Println("# ########################################################################################")
	ldfFlags := fmt.Sprintf("-X main.version=%s", gpsaBuildContext.ProgramVersion)

	BuildFolders(gpsaBuildContext.PackagesToBuild, gpsaBuildContext.BinDir, ldfFlags)

	fmt.Println("# ########################################################################################")
	return nil
}

// Runs the tests for the project
func Test() error {
	mg.Deps(getEnvironment, Clean, GetBuildName, installTestDeps)
	fmt.Println(fmt.Sprintf("Testing gpsa... "))
	fmt.Println("# ########################################################################################")
	xmlResult := filepath.Join(gpsaBuildContext.LogDir, "TestsResult.xml")
	logFileName := "TestRun.log"

	testErrors := RunTestFolders(gpsaBuildContext.PackagesToTest, gpsaBuildContext.LogDir, logFileName)

	errConv := ConvertTestResults(filepath.Join(gpsaBuildContext.LogDir, logFileName), xmlResult, gpsaBuildContext.WorkDir)
	if errConv != nil {
		return errConv
	}
	if len(testErrors) > 0 {
		return testErrors[0]
	}

	fmt.Println("# ########################################################################################")
	return nil
}

// Runs test coverage for the project
func Cover() error {
	mg.Deps(getEnvironment, Clean, GetBuildName)
	fmt.Println(fmt.Sprintf("Testing with coverage gpsa... "))
	fmt.Println("# ########################################################################################")

	CoverTestFolders(gpsaBuildContext.PackagesToTest, gpsaBuildContext.LogDir, "TestCoverage.log")

	fmt.Println("# ########################################################################################")
	return nil
}

// Remove all build output
func Clean() error {
	mg.Deps(getEnvironment)
	fmt.Println("Cleaning...")
	fmt.Println("# ########################################################################################")

	errClean := RemovePaths([]string{
		gpsaBuildContext.BinDir,
		gpsaBuildContext.PackageDir,
		gpsaBuildContext.LogDir,
		gpsaBuildContext.BinZipPath,
		gpsaBuildContext.BuildZipPath,
	})
	if errClean != nil {
		return errClean
	}

	fmt.Println("# ########################################################################################")
	return nil
}

// Create zip files from the build output and logs
func CreateBuildZip() error {
	mg.Deps(getEnvironment, Build)
	fmt.Println("Zipping...")
	fmt.Println("# ########################################################################################")

	errBuildZip := ZipFolders([]string{gpsaBuildContext.BinDir, gpsaBuildContext.PackageDir, gpsaBuildContext.LogDir}, gpsaBuildContext.BuildZipPath)
	if errBuildZip != nil {
		return errBuildZip
	}

	errBinZip := ZipFolders([]string{gpsaBuildContext.BinDir}, gpsaBuildContext.BinZipPath)
	if errBinZip != nil {
		return errBinZip
	}

	fmt.Println("# ########################################################################################")
	return nil
}

func installTestDeps() error {
	mg.Deps(Clean)
	fmt.Println("Installing Test Dependencies...")
	fmt.Println("# ########################################################################################")

	InstallTestConverter(filepath.Join(gpsaBuildContext.WorkDir, "build"))

	fmt.Println("# ########################################################################################")
	return nil
}

func readVersionMaster() (string, error) {
	content, err := ioutil.ReadFile(gpsaBuildContext.VersionFilePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}
