package tcxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"encoding/xml"
	"io/ioutil"
)

// Tcx - Represents the content of a txc file
type Tcx struct {
	ActivityArray []Activities `xml:"Activities"`
}

// Activities - Represents the Activities of a txc file
type Activities struct {
	Activities []Activity `xml:"Activity"`
}

// Activity - Represents one Activity in a TCX file
type Activity struct {
	ID   string `xml:"Id"`
	Laps []Lap  `xml:"Lap"`
}

// Lap - Represents one Lap in a TCX file
type Lap struct {
	DistanceMeters   string  `xml:"DistanceMeters"`
	TotalTimeSeconds string  `xml:"TotalTimeSeconds"`
	Tracks           []Track `xml:"Track"`

	// ToDo: Read TotalTimeSeconds, DistanceMeters and StartTime
}

// Track - Represents one Track in a TCX file
type Track struct {
	Trackpoints []Trackpoint `xml:"Trackpoint"`
}

// Trackpoint - Represents one Trackpoint in a TCX file
type Trackpoint struct {
	Time           string          `xml:"Time"`
	AltitudeMeters float32         `xml:"AltitudeMeters"`
	Position       PositionWrapper `xml:"Position"`
}

// PositionWrapper - Represents the Position in a TCX file
type PositionWrapper struct {
	LatitudeDegrees  float32 `xml:"LatitudeDegrees"`
	LongitudeDegrees float32 `xml:"LongitudeDegrees"`
}

// ReadTcx - Read a Tcx file
func ReadTcx(fileName string) (Tcx, error) {
	xmlfile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return Tcx{}, err
	}
	return readGPXBuffer(xmlfile, fileName)
}

func readGPXBuffer(fileBuffer []byte, fileName string) (Tcx, error) {
	tcx := Tcx{}
	err := xml.Unmarshal([]byte(fileBuffer), &tcx)

	if err != nil {
		return Tcx{}, err
	}

	if len(tcx.ActivityArray) <= 0 {
		return Tcx{}, newTcxFileError(fileName)
	}
	if len(tcx.ActivityArray[0].Activities) <= 0 {
		return Tcx{}, newEmptyTcxFileError(fileName)
	}
	if len(tcx.ActivityArray[0].Activities[0].Laps) <= 0 {
		return Tcx{}, newEmptyTcxFileError(fileName)
	}

	if len(tcx.ActivityArray[0].Activities[0].Laps[0].Tracks) <= 0 {
		return Tcx{}, newEmptyTcxFileError(fileName)
	}

	if tcx.ActivityArray[0].Activities[0].ID == "" {
		return Tcx{}, newTcxFileError(fileName)
	}

	return tcx, nil
}
