package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"tobi.backfrak.de/internal/gpsabl"
)

// HelpFlag - Tells if the program was called with -help
var HelpFlag bool

// VerboseFlag - Tells if the program was called with -verbose
var VerboseFlag bool

// SkipErrorExitFlag - Tells if the program was called with -skip-error-exit
var SkipErrorExitFlag bool

// PrintCsvHeaderFlag - Tells if the program was called with -print-csv-header
var PrintCsvHeaderFlag bool

// OutFileParameter - Tells if and where we should write the output to ( -out-file )
var OutFileParameter string

// DontPanicFlag - Tells if the program was called with -dont-panic
var DontPanicFlag bool

// DepthParameter - Tells for which depth we should perform the analyses ( -depth )
var DepthParameter string

// SummaryParameter - Tells if we should add summary to the output ( -summary )
var SummaryParameter string

// TimeFormatParameter - Tells if we should add summary to the output ( -time-format )
var TimeFormatParameter string

// CorrectionParameter - Tells how we should correct Elevation data from the track ( -correction )
var CorrectionParameter string

// PrintVersionFlag - Tells if the program was called with the -version flag
var PrintVersionFlag bool

// PrintLicenseFlag - Tells if the program was called with the -license flag
var PrintLicenseFlag bool

// SuppressDuplicateOutPutFlag - Tells if the program was called with the -suppressDuplicateOutPut flag
var SuppressDuplicateOutPutFlag bool

// MinimalMovingSpeedParameter - Tells the minimal moving speed for moving time and speed calculation
var MinimalMovingSpeedParameter float64

// MinimalStepHightParameter - Tells the minimal step hight, when "steps" correction is used
var MinimalStepHightParameter float64

// PrintElevationOverDistanceFlag - Tell if the program was called with the -print-elevation-over-distance flag
var PrintElevationOverDistanceFlag bool

// StdOutFormatParameter - Tells the formant when StdOut is the output stream -std-out-format
var StdOutFormatParameter string

// MinStartTime - The minimum StartTime for a track to be added to the output. Formatted in "YYYY-MMM-dd HH:mm:ss", may without seconds or just a date
var MinStartTime string

// MaxStartTime - The maximum StartTime for a track to be added to the output. Formatted in "YYYY-MMM-dd HH:mm:ss", may without seconds or just a date
var MaxStartTime string

// ReadInputStreamBuffer - Read an input stream and figure out what kind of files are given
func ReadInputStreamBuffer(reader *bufio.Reader) ([]gpsabl.InputFile, error) {
	var fileArgs []gpsabl.InputFile
	var inputBytes []byte
	for {
		input, errRead := reader.ReadByte()
		if errRead != nil {
			if errRead == io.EOF {
				break
			} else {
				return nil, errRead
			}
		}

		inputBytes = append(inputBytes, input)
	}

	buffers := getXMlFileBuffersFromInputStream(inputBytes)

	if len(buffers) != 0 {
		// fmt.Fprintln(os.Stdout, fmt.Sprintf("Got %d input files as stream", len(buffers)))

		for i, buffer := range buffers {
			res, input := gpsabl.GetInputFileFromBuffer(ValidReaders, buffer, fmt.Sprintf("Input stream buffer %d", i+1))
			if res == true {
				fileArgs = append(fileArgs, input)
			}
		}
		if len(fileArgs) != 0 && VerboseFlag {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("Got %d files as stream", len(fileArgs)))
		}
		return fileArgs, nil
	}

	fileArgsStr, errProcFileName := getFilePathFromInputStream(inputBytes)
	if errProcFileName != nil {
		return nil, errProcFileName
	}

	for _, fileArgStr := range fileArgsStr {
		res, input := gpsabl.GetInputFileFromPath(ValidReaders, fileArgStr)
		if res == true {
			fileArgs = append(fileArgs, input)
		} else {
			return nil, newUnKnownFileTypeError(fileArgStr)
		}
	}

	return fileArgs, nil
}

// getXMlFileBuffersFromInputStream - Parse the input buffer array and search for xml files
// resturn an array of input buffer arrays, where each array of input buffer array contain exactly one xml file content
func getXMlFileBuffersFromInputStream(inputBytes []byte) [][]byte {
	var startBytes []int
	for i, _ := range inputBytes {
		section := inputBytes[i : i+5]
		if string(section) == "<?xml" {
			startBytes = append(startBytes, i)
		}
	}

	var retVal [][]byte
	size := len(inputBytes)
	numberFiles := len(startBytes)
	for i, index := range startBytes {
		// fmt.Println(fmt.Sprintf("New xml file at index %d", index))
		var oneFile []byte
		if i < numberFiles-1 {
			oneFile = inputBytes[index:startBytes[i+1]]
		} else {
			oneFile = inputBytes[index:size]
		}
		retVal = append(retVal, oneFile)
	}
	return retVal
}

// getFilePathFromInputStream - parse the input bytes array and search for valid file pathes
// resturn the list of valid file pathes
func getFilePathFromInputStream(inputBytes []byte) ([]string, error) {
	var fileArgs []string
	read, write, errCreate := os.Pipe()
	if errCreate != nil {
		return nil, errCreate
	}

	_, errWrite := write.Write(inputBytes)
	if errWrite != nil {
		return nil, errWrite
	}
	write.Close()
	reader := bufio.NewReader(read)
	for {
		input, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		line := string(input)
		if line != "" {
			if strings.Contains(line, string(os.PathSeparator)) && fileExists(line) {
				fileArgs = append(fileArgs, line)
			} else {
				return nil, newUnKnownInputStreamError(line)
			}
		}
	}
	return fileArgs, nil
}

// handleComandlineOptions - Setup and parse the commandline options.
// Defines the usage function as well
func handleComandlineOptions() {

	// Setup the valid comandline flags
	flag.Float64Var(&MinimalStepHightParameter, "minimal-step-hight", 10.0, "The minimal step hight. Only in use when \"steps\"  elevation correction is used. In [m]")
	flag.Float64Var(&MinimalMovingSpeedParameter, "minimal-moving-speed", 0.3, "The minimal speed. Distances traveled with less speed are not counted. In [m/s]")
	flag.BoolVar(&SuppressDuplicateOutPutFlag, "suppress-duplicate-out-put", false, "Suppress the output of duplicate lines. Duplicates are detected by timestamps. Output with non valid time data may still contains duplicates.")
	flag.BoolVar(&HelpFlag, "help", false, "Print help message and exit")
	flag.BoolVar(&PrintVersionFlag, "version", false, "Print version of the program and exit")
	flag.BoolVar(&PrintLicenseFlag, "license", false, "Print license information of the program and exit")
	flag.BoolVar(&VerboseFlag, "verbose", false, "Run the program with verbose output")
	flag.BoolVar(&SkipErrorExitFlag, "skip-error-exit", false, "Don't exit the program on track file processing errors")
	flag.BoolVar(&PrintCsvHeaderFlag, "print-csv-header", true, "Print out a csv header line. Possible values are [true false]")
	flag.StringVar(&OutFileParameter, "out-file", "",
		fmt.Sprintf("Decide where to write the output. StdOut is used when not explicitly set. Supported file endings are: %s. The format will be set according the given ending.", getValidOutputxtensions()))
	flag.BoolVar(&DontPanicFlag, "dont-panic", true, "Decide if the program will exit with panic or with negative exit code in error cases. Possible values are [true false]")
	flag.StringVar(&DepthParameter, "depth", string(gpsabl.TRACK),
		fmt.Sprintf("Define the way the program should analyse the files. Possible values are [%s]", gpsabl.GetValidDepthArgsString()))
	flag.StringVar(&CorrectionParameter, "correction", string(gpsabl.STEPS),
		fmt.Sprintf("Define how to correct the elevation data read in from the track. Possible values are [%s]", gpsabl.GetValidCorrectionParametersString()))
	flag.BoolVar(&PrintElevationOverDistanceFlag, "print-elevation-over-distance", false, "Tell if \"ElevationOverDistance.csv\" should be created for each track. The files will be locate in tmp dir.")
	flag.StringVar(&StdOutFormatParameter, "std-out-format", string(ValidFormaters[0].GetOutputFormaterTypes()[0]),
		fmt.Sprintf("The output format when stdout is the used output. Ignored when out-file is given. Possible values are [%s]", getStdOutFormatParameterValuesStr()))
	flag.StringVar(&SummaryParameter, "summary", string(gpsabl.NONE),
		fmt.Sprintf("Tell if you want to get a summary report. Possible values are [%s]", gpsabl.GetValidSummaryArgsString()))
	flag.StringVar(&TimeFormatParameter, "time-format", string(gpsabl.RFC850),
		fmt.Sprintf("Tell how the csv output formater should format times. Possible values are [%s]", gpsabl.GetValidTimeFormatsString()))
	flag.StringVar(&MinStartTime, "min-start-time", "",
		"The minimum StartTime for a track to be added to the output. Formatted in \"YYYY-MMM-dd HH:mm:ss\", may without seconds or just a date")
	flag.StringVar(&MaxStartTime, "max-start-time", "",
		"The maximum StartTime for a track to be added to the output. Formatted in \"YYYY-MMM-dd HH:mm:ss\", may without seconds or just a date")

	// Overwrite the std Usage function with some custom stuff
	flag.Usage = customHelpMessage

	// Read the given flags
	flag.Parse()
}

// customHelpMessage - Print he customized help message
func customHelpMessage() {
	fmt.Fprintln(os.Stdout, fmt.Sprintf("%s: Reads in GPS track files, and writes out basic statistic data found in the track as a report", os.Args[0]))
	fmt.Fprintln(os.Stdout, fmt.Sprintf("Program %s", getVersion()))
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, fmt.Sprintf("Usage: %s [options] [files]", os.Args[0]))
	fmt.Fprintln(os.Stdout, "  files")
	fmt.Fprintln(os.Stdout, fmt.Sprintf("        One or more track files of the following type: %s", getValidTrackExtensions()))
	fmt.Fprintln(os.Stdout, "Options:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "It is also possible to pipe track file names or track file content into")
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Examples:")
	fmt.Fprintln(os.Stdout, "./gpsa my/test/file.gpx")
	fmt.Fprintln(os.Stdout, "./gpsa -verbose -out-file=gps-statistics.csv my/test/*.gpx")
	fmt.Fprintln(os.Stdout, "find ./testdata/valid-gpx -name \"*.gpx\" | ./bin/gpsa -summary=additional -out-file=./test.json")
	fmt.Fprintln(os.Stdout, "cat  01.gpx 01.tcx 03.tcx 02.gpx | ./bin/gpsa -out-file=./test.json")
}

// getStdOutFormatParameterValuesStr - Get a string that contains all valid stdOutFormatParameterValues
func getStdOutFormatParameterValuesStr() string {
	ret := ""
	for _, formater := range ValidFormaters {
		types := formater.GetOutputFormaterTypes()
		for _, t := range types {
			ret = fmt.Sprintf("%s %s", t, ret)
		}
	}
	return ret
}

// checkStdOutFormatParameterValue - Tell if a given stdOutFormatParameterValue is valid
func checkStdOutFormatParameterValue(val string) bool {
	for _, formater := range ValidFormaters {
		args := formater.GetOutputFormaterTypes()
		for _, arg := range args {
			if strings.ToLower(string(arg)) == strings.ToLower(val) {
				return true
			}
		}
	}

	return false
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getValidTrackExtensions() string {
	var ret string

	for _, reader := range ValidReaders {
		extensions := reader.GetValidFileExtensions()

		for _, extension := range extensions {
			toAdd := fmt.Sprintf("*%s", extension)
			ret = fmt.Sprintf("%s, %s", toAdd, ret)
		}
	}

	return ret
}

func getValidOutputxtensions() string {
	var ret string

	for _, formater := range ValidFormaters {
		extensions := formater.GetFileExtensions()

		for _, extension := range extensions {
			toAdd := fmt.Sprintf("*%s", extension)
			ret = fmt.Sprintf("%s, %s", toAdd, ret)
		}
	}

	return ret
}
