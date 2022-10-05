//go:build mage
// +build mage

package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
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

	hash, errHash := getGitHash()
	if errHash != nil {
		return errHash
	}
	// fmt.Println(fmt.Sprintf("Git Hash: %s", hash))
	gpsaBuildContext.GitHash = hash

	hight, errHight := getGitHight()
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

	errFindBuild := filepath.Walk(filepath.Join(gpsaBuildContext.SourceDir, PROJECT_NAME_SPACE, "cmd"), func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return nil
		}

		packToBuild := filepath.Dir(path)
		if !info.IsDir() && filepath.Base(path) == "go.mod" && !listContains(gpsaBuildContext.PackagesToBuild, packToBuild) {
			gpsaBuildContext.PackagesToBuild = append(gpsaBuildContext.PackagesToBuild, packToBuild)
		}

		return nil
	})
	if errFindBuild != nil {
		return errFindBuild
	}

	errFindTest := filepath.Walk(filepath.Join(gpsaBuildContext.SourceDir, PROJECT_NAME_SPACE), func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ".go" {
			packToTest := filepath.Dir(path)
			if strings.HasSuffix(path, "_test.go") && !listContains(gpsaBuildContext.PackagesToTest, packToTest) {
				gpsaBuildContext.PackagesToTest = append(gpsaBuildContext.PackagesToTest, packToTest)
			}
		}

		return nil
	})
	if errFindTest != nil {
		return errFindTest
	}

	return nil
}

func GetBuildName() error {
	mg.Deps(getEnvironment, Clean)
	fmt.Println(fmt.Sprintf("Create gpsa Version files..."))

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

	return nil
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(getEnvironment, Clean, GetBuildName)
	fmt.Println(fmt.Sprintf("Building gpsa V%s ...", gpsaBuildContext.ProgramVersion))

	if _, err := os.Stat(gpsaBuildContext.BinDir); os.IsNotExist(err) {
		errCreate := os.Mkdir(gpsaBuildContext.BinDir, 0755)
		if errCreate != nil {
			return errCreate
		}
	}

	for _, packToBuild := range gpsaBuildContext.PackagesToBuild {
		outPutPath := filepath.Join(gpsaBuildContext.BinDir, filepath.Base(packToBuild))
		fmt.Println(fmt.Sprintf("Compile package '%s' to '%s'", packToBuild, outPutPath))

		ldfFlags := fmt.Sprintf("-ldflags=\"-X main.version=%s\"", gpsaBuildContext.ProgramVersion)
		fmt.Println(fmt.Sprintf("Run in %s: %s %s %s %s %s %s", packToBuild, "go", "build", "-o", gpsaBuildContext.BinDir, "-v", ldfFlags))
		cmd := exec.Command("go", "build", "-o", outPutPath, "-v", ldfFlags)
		cmd.Dir = packToBuild
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		errBuild := cmd.Run()
		buildOutStr := out.String()
		if buildOutStr != "" {
			fmt.Println(buildOutStr)
		}
		if errBuild != nil {
			fmt.Println(stderr.String())
			fmt.Println(errBuild.Error())
			return errBuild
		}
	}

	return nil
}

func Test() error {
	mg.Deps(getEnvironment, Clean, GetBuildName, InstallTestDeps)
	fmt.Println(fmt.Sprintf("Testing gpsa... "))

	if _, err := os.Stat(gpsaBuildContext.LogDir); os.IsNotExist(err) {
		errCreate := os.Mkdir(gpsaBuildContext.LogDir, 0755)
		if errCreate != nil {
			return errCreate
		}
	}

	logPath := filepath.Join(gpsaBuildContext.LogDir, "TestsRun.log")
	xmlResult := filepath.Join(gpsaBuildContext.LogDir, "TestsResult.xml")
	logFile, errOpen := os.Create(logPath)
	if errOpen != nil {
		return errOpen
	}

	for _, packToTest := range gpsaBuildContext.PackagesToTest {

		fmt.Println(fmt.Sprintf("Test package '%s', logging to '%s'", packToTest, logPath))
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			fmt.Println(fmt.Sprintf("Run in %s: %s %s %s >> %s", packToTest, "go", "test", "-v", logPath))
			cmd = exec.Command("go", "test", "-v")
		} else {
			fmt.Println(fmt.Sprintf("Run in %s: %s %s %s %s >> %s", packToTest, "go", "test", "-v", "-race", logPath))
			cmd = exec.Command("go", "test", "-v", "-race")
		}
		cmd.Dir = packToTest
		cmd.Stderr = logFile
		cmd.Stdout = logFile
		errTest := cmd.Run()
		if errTest != nil {
			logFile.Close()
			fmt.Println(errTest.Error())
			return errTest
		}
	}
	logFile.Close()

	fmt.Println(fmt.Sprintf("Convert the test results %s to %s", logPath, xmlResult))
	cmd := exec.Command("go", "run", "github.com/tebeka/go2xunit", "-input", logPath, "-output", xmlResult)
	cmd.Dir = filepath.Join(gpsaBuildContext.WorkDir, "build")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	errConvert := cmd.Run()
	buildOutStr := out.String()
	if buildOutStr != "" {
		fmt.Println(buildOutStr)
	}
	if errConvert != nil {
		fmt.Println(stderr.String())
		fmt.Println(errConvert.Error())
		return errConvert
	}

	return nil
}

func Cover() error {
	mg.Deps(getEnvironment, Clean, GetBuildName, InstallTestDeps)
	fmt.Println(fmt.Sprintf("Testing gpsa... "))

	if _, err := os.Stat(gpsaBuildContext.LogDir); os.IsNotExist(err) {
		errCreate := os.Mkdir(gpsaBuildContext.LogDir, 0755)
		if errCreate != nil {
			return errCreate
		}
	}

	logPath := filepath.Join(gpsaBuildContext.LogDir, "TestsCoverRun.log")
	logFile, errOpen := os.Create(logPath)
	if errOpen != nil {
		return errOpen
	}

	for _, packToTest := range gpsaBuildContext.PackagesToTest {

		fmt.Println(fmt.Sprintf("Test package '%s', logging to '%s'", packToTest, logPath))
		fmt.Println(fmt.Sprintf("Run in %s: %s %s %s %s >> %s", packToTest, "go", "test", "-v", "-cover", logPath))
		cmd := exec.Command("go", "test", "-v", "-cover")

		cmd.Dir = packToTest
		cmd.Stderr = logFile
		cmd.Stdout = logFile
		errTest := cmd.Run()
		if errTest != nil {
			logFile.Close()
			fmt.Println(errTest.Error())
			return errTest
		}
	}
	logFile.Close()

	return nil
}

// Manage your deps, or running package managers.
func InstallTestDeps() error {
	mg.Deps(Clean)
	fmt.Println("Installing Test Dependencies...")
	cmd := exec.Command("go", "install", "-v", "github.com/tebeka/go2xunit@v1.4.10")
	cmd.Dir = filepath.Join(gpsaBuildContext.WorkDir, "build")
	return cmd.Run()
}

// Clean up after yourself
func Clean() error {
	mg.Deps(getEnvironment)
	fmt.Println("Cleaning...")

	return removePathes([]string{
		gpsaBuildContext.BinDir,
		gpsaBuildContext.PackageDir,
		gpsaBuildContext.LogDir,
		gpsaBuildContext.BinZipPath,
		gpsaBuildContext.BuildZipPath,
	})
}

func CreateBuildZip() error {
	mg.Deps(getEnvironment)
	fmt.Println("Zipping...")

	errBuildZip := zipSources([]string{gpsaBuildContext.BinDir, gpsaBuildContext.PackageDir, gpsaBuildContext.LogDir}, gpsaBuildContext.BuildZipPath)
	if errBuildZip != nil {
		return errBuildZip
	}

	errBinZip := zipSources([]string{gpsaBuildContext.BinDir}, gpsaBuildContext.BinZipPath)
	if errBinZip != nil {
		return errBinZip
	}

	return nil
}

func getGitHash() (string, error) {
	cmd := exec.Command("git", "describe", "--always", "--long", "--dirty")
	cmd.Dir = gpsaBuildContext.WorkDir
	hash, err := cmd.Output()
	if err != nil {
		return "", err
	}
	hashStr := strings.TrimSpace(string(hash))
	return hashStr, nil
}

func getGitHight() (int, error) {
	cmd := exec.Command("git", "log", "--pretty=format:\"%H\"", "-n 1", "--follow", VERSION_FILE)
	cmd.Dir = gpsaBuildContext.WorkDir
	lastChange, errLast := cmd.Output()
	if errLast != nil {
		return -1, errLast
	}
	lastChangeStr := strings.ReplaceAll(strings.TrimSpace(string(lastChange)), "\"", "")

	cmd = exec.Command("git", "log", "--pretty=format:\"%H\"", "-n 1")
	cmd.Dir = gpsaBuildContext.WorkDir
	head, errHead := cmd.Output()
	if errHead != nil {
		return -1, errHead
	}

	headStr := strings.ReplaceAll(strings.TrimSpace(string(head)), "\"", "")

	cmd = exec.Command("git", "rev-list", "--count", lastChangeStr+".."+headStr)
	cmd.Dir = gpsaBuildContext.WorkDir
	hight, hightErr := cmd.Output()
	if hightErr != nil {
		return -1, hightErr
	}

	hightStr := strings.TrimSpace(string(hight))
	hightInt, errCon := strconv.Atoi(hightStr)
	if errCon != nil {
		return -1, nil
	}

	return hightInt, nil
}

func readVersionMaster() (string, error) {
	content, err := ioutil.ReadFile(gpsaBuildContext.VersionFilePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

func listContains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}

	return false
}

func zipSources(sources []string, target string) error {
	fmt.Println(fmt.Sprintf("Zip %s into %s", sources, target))
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	for _, source := range sources {

		if _, err := os.Stat(source); os.IsNotExist(err) {
			continue
		}
		// 2. Go through all the files of the source
		packSourceErr := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 3. Create a local file header
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// set compression
			header.Method = zip.Deflate

			// 4. Set relative path of a file as the header name
			header.Name, err = filepath.Rel(filepath.Dir(source), path)
			if err != nil {
				return err
			}
			if info.IsDir() {
				header.Name += "/"
			}

			// 5. Create writer for the file header and save content of the file
			headerWriter, err := writer.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(headerWriter, f)
			if err != nil {
				return err
			}

			return nil
		})

		if packSourceErr != nil {
			return packSourceErr
		}
	}

	return nil
}

func removePathes(pathes []string) error {
	for _, path := range pathes {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}
	return nil
}
