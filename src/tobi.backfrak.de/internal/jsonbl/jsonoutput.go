package jsonbl

import (
	"encoding/json"
	"os"
	"strings"
	"sync"

	"tobi.backfrak.de/internal/gpsabl"
)

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// JSONOutputFormatertype - The gpsabl.OutputFormaterType this formater is responsible for
const JSONOutputFormatertype gpsabl.OutputFormaterType = "JSON"
const FileExtension = ".json"

// JSONOutput - Structure of the json file
type JSONOutput struct {
	Statistics []gpsabl.OutputLine
	Summary    []gpsabl.OutputLine
}

// JSONOutputFormater - type that formats TrackSummary into json style
type JSONOutputFormater struct {
	lineBuffer []gpsabl.OutputLine
	mux        sync.Mutex
}

// NewJSONOutputFormater - Get a new instance of the JSONOutputFormater
func NewJSONOutputFormater() *JSONOutputFormater {
	ret := JSONOutputFormater{}
	ret.lineBuffer = []gpsabl.OutputLine{}

	return &ret
}

// NewOutputFormater -  Get a new gpsabl.OutputFormater of this type
func (formater *JSONOutputFormater) NewOutputFormater() gpsabl.OutputFormater {
	ret := NewJSONOutputFormater()

	return gpsabl.OutputFormater(ret)
}

// AddOutPut - Add the output values of a TrackFile to the out file buffer. Implements the gpsabl.OutputFormater interface
func (formater *JSONOutputFormater) AddOutPut(trackFile gpsabl.TrackFile, depth gpsabl.DepthArg, filterDuplicate bool) error {
	var lines []gpsabl.OutputLine
	linesFromFile, err := gpsabl.GetOutlines(trackFile, depth)
	if err != nil {
		return err
	}
	if filterDuplicate {
		for _, line := range gpsabl.StripOutlines(linesFromFile) {
			if gpsabl.OutputContainsLineByTimeStamps(lines, line) == false && gpsabl.OutputContainsLineByTimeStamps(formater.lineBuffer, line) == false {
				lines = append(lines, line)
			}
		}
	} else {
		lines = gpsabl.StripOutlines(linesFromFile)
	}

	if len(lines) > 0 {
		formater.mux.Lock()
		defer formater.mux.Unlock()
		formater.lineBuffer = append(formater.lineBuffer, lines...)
	}

	return nil
}

// WriteOutput - Write the output to the output file
func (formater *JSONOutputFormater) WriteOutput(outFile *os.File, summary gpsabl.SummaryArg) error {

	output, errGet := formater.GetOutput(summary)
	if errGet != nil {
		return errGet
	}

	errWrite := writeJSON(outFile, output)
	if errWrite != nil {
		return errWrite
	}

	return nil
}

// GetOutput - Get the output that will be written to the file
func (formater *JSONOutputFormater) GetOutput(summary gpsabl.SummaryArg) (JSONOutput, error) {

	if !gpsabl.CheckValidSummaryArg(string(summary)) {
		return JSONOutput{}, gpsabl.NewSummaryParamaterNotKnown(summary)
	}

	formater.mux.Lock()
	defer formater.mux.Unlock()
	var ret JSONOutput
	switch summary {
	case gpsabl.NONE:
		ret.Statistics = formater.lineBuffer
	case gpsabl.ONLY:
		ret.Summary = formater.getSummaryEntires()
	case gpsabl.ADDITIONAL:
		ret.Statistics = formater.lineBuffer
		ret.Summary = formater.getSummaryEntires()
	default:
		return JSONOutput{}, gpsabl.NewSummaryParamaterNotKnown(summary)
	}

	return ret, nil
}

// CheckOutputFormaterType - Check if this OutputFormater is responsible for the given gpsabl.OutputFormaterType
func (formater *JSONOutputFormater) CheckOutputFormaterType(formaterType gpsabl.OutputFormaterType) bool {
	if formaterType == JSONOutputFormatertype {
		return true
	}

	return false
}

// GetOutputFormaterTypes - Get the list of gpsabl.OutputFormaterType this formater can write
func (formater *JSONOutputFormater) GetOutputFormaterTypes() []gpsabl.OutputFormaterType {
	return []gpsabl.OutputFormaterType{JSONOutputFormatertype}
}

// CheckFileExtension - Check if this OutputFormater can write the given output file
func (formater *JSONOutputFormater) CheckFileExtension(filePath string) bool {
	if strings.HasSuffix(strings.ToLower(filePath), FileExtension) {
		return true
	}

	return false
}

// GetFileExtensions - Get the list of file extensions this formater can write
func (formater *JSONOutputFormater) GetFileExtensions() []string {
	return []string{FileExtension}
}

func (formater *JSONOutputFormater) getSummaryEntires() []gpsabl.OutputLine {
	stats := gpsabl.GetStatisticSummaryData(formater.lineBuffer)

	ret := []gpsabl.OutputLine{}
	sumLine := gpsabl.OutputLine{}
	sumLine.Name = "Sum"
	sumLine.Data = stats.Sum
	ret = append(ret, sumLine)

	avgLine := gpsabl.OutputLine{}
	avgLine.Name = "Average"
	avgLine.Data = stats.Average
	ret = append(ret, avgLine)

	minLine := gpsabl.OutputLine{}
	minLine.Name = "Minimum"
	minLine.Data = stats.Minimum
	ret = append(ret, minLine)

	maxLine := gpsabl.OutputLine{}
	maxLine.Name = "Maximum"
	maxLine.Data = stats.Maximum
	ret = append(ret, maxLine)

	return ret
}

func writeJSON(outFile *os.File, output JSONOutput) error {
	file, errConv := json.MarshalIndent(output, "", " ")
	if errConv != nil {
		return errConv
	}

	count, errWrite := outFile.Write(file)
	if errWrite != nil {
		return errWrite
	}

	if count != len(file) {
		return os.ErrClosed
	}

	return nil
}
