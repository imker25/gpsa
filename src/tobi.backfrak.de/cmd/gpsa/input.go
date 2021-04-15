package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"tobi.backfrak.de/internal/csvbl"
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

var stdOutFormatParameterValues = []string{"CSV", "JSON"}

func ReadInputStreamBuffer(reader *bufio.Reader) ([]string, error) {

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
		gpxBuffers := getGpxBuffers(buffers)
		if len(gpxBuffers) != 0 {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("Got %d gpx files as stream", len(gpxBuffers)))
		}

		tcxBuffers := getTcxBuffers(buffers)
		if len(tcxBuffers) != 0 {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("Got %d tcx files as stream", len(tcxBuffers)))
		}
		return nil, nil
	}

	fileArgs, errProcFileName := getFilePathFromInputStream(inputBytes)
	if errProcFileName != nil {
		return nil, errProcFileName
	}

	return fileArgs, nil
}

func getGpxBuffers(buffers [][]byte) [][]byte {
	var retVal [][]byte
	for _, buffer := range buffers {
		for i, _ := range buffer {
			section := buffer[i : i+4]
			if string(section) == "<gpx" {
				retVal = append(retVal, buffer)
				break
			}
		}
	}
	return retVal
}

func getTcxBuffers(buffers [][]byte) [][]byte {
	var retVal [][]byte
	for _, buffer := range buffers {
		for i, _ := range buffer {
			section := buffer[i : i+23]
			if string(section) == "<TrainingCenterDatabase" {
				retVal = append(retVal, buffer)
				break
			}
		}
	}
	return retVal
}

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
			oneFile = inputBytes[index : startBytes[i+1]-1]
		} else {
			oneFile = inputBytes[index:size]
		}
		retVal = append(retVal, oneFile)
	}
	return retVal
}

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
	flag.StringVar(&OutFileParameter, "out-file", "", "Decide where to write the output. StdOut is used when not explicitly set. *.csv and *.json are supported file endings, the format will be set according the given ending.")
	flag.BoolVar(&DontPanicFlag, "dont-panic", true, "Decide if the program will exit with panic or with negative exit code in error cases. Possible values are [true false]")
	flag.StringVar(&DepthParameter, "depth", string(gpsabl.TRACK),
		fmt.Sprintf("Define the way the program should analyse the files. Possible values are [%s]", gpsabl.GetValidDepthArgsString()))
	flag.StringVar(&CorrectionParameter, "correction", string(gpsabl.STEPS),
		fmt.Sprintf("Define how to correct the elevation data read in from the track. Possible values are [%s]", gpsabl.GetValidCorrectionParametersString()))
	flag.BoolVar(&PrintElevationOverDistanceFlag, "print-elevation-over-distance", false, "Tell if \"ElevationOverDistance.csv\" should be created for each track. The files will be locate in tmp dir.")
	flag.StringVar(&StdOutFormatParameter, "std-out-format", "CSV",
		fmt.Sprintf("The output format when stdout is the used output. Ignored when out-file is given. Possible values are [%s]", getStdOutFormatParameterValuesStr()))
	flag.StringVar(&SummaryParameter, "summary", string(gpsabl.NONE),
		fmt.Sprintf("Tell if you want to get a summary report. Possible values are [%s]", gpsabl.GetValidSummaryArgsString()))
	flag.StringVar(&TimeFormatParameter, "time-format", string(csvbl.RFC850),
		fmt.Sprintf("Tell how the csv output formater should format times. Possible values are [%s]", csvbl.GetValidTimeFormatsString()))
	// Overwrite the std Usage function with some custom stuff
	flag.Usage = customHelpMessage

	// Read the given flags
	flag.Parse()
}

func customHelpMessage() {
	fmt.Fprintln(os.Stdout, fmt.Sprintf("%s: Reads in GPS track files, and writes out basic statistic data found in the track as a CSV or JSON style report", os.Args[0]))
	fmt.Fprintln(os.Stdout, fmt.Sprintf("Program %s", getVersion()))
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, fmt.Sprintf("Usage: %s [options] [files]", os.Args[0]))
	fmt.Fprintln(os.Stdout, "  files")
	fmt.Fprintln(os.Stdout, "        One or more track files (only *.gpx and *.tcx supported at the moment)")
	fmt.Fprintln(os.Stdout, "Options:")
	flag.PrintDefaults()
}

func getStdOutFormatParameterValuesStr() string {
	ret := ""
	for _, arg := range stdOutFormatParameterValues {
		ret = fmt.Sprintf("%s %s", arg, ret)
	}
	return ret
}

func checkStdOutFormatParameterValue(val string) bool {
	for _, arg := range stdOutFormatParameterValues {
		if strings.ToLower(arg) == strings.ToLower(val) {
			return true
		}
	}

	return false
}
