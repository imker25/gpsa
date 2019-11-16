package gpxbl

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

// Gpx - Represents the content of a GPX file
type Gpx struct {
	Name        string `xml:"name"`
	Description string `xml:"desc"`
	Tracks      []Trk  `xml:"trk"`
}

// Trk - Represents the content of a GPX track
type Trk struct {
	Name          string   `xml:"name"`
	Number        int      `xml:"number"`
	Description   string   `xml:"desc"`
	TrackSegments []Trkseg `xml:"trkseg"`
}

// Trkseg - Represents a track segement, basicaly a arry of Trkpt
type Trkseg struct {
	TrackPoints []Trkpt `xml:"trkpt"`
}

// Trkpt - Represents a track point
type Trkpt struct {
	Elevation float32 `xml:"ele"`
	Latitude  float32 `xml:"lat,attr"`
	Longitude float32 `xml:"lon,attr"`
}

// GpxFileError - Error when trying to load not a gpx file
type GpxFileError struct {
	err string
	// File - The path to the file that caused this error
	File string
}

func (e *GpxFileError) Error() string { // Implement the Error Interface for the GpxFileError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// newGpxFileError - Get a new GpxFileError struct
func newGpxFileError(fileName string) *GpxFileError {
	return &GpxFileError{fmt.Sprintf("The file \"%s\" was not a gpx file", fileName), fileName}
}

// ReadGPX - Read a GPX file
func ReadGPX(fileName string) (Gpx, error) {
	xmlfile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return Gpx{}, err
	}
	return readGPXBuffer(xmlfile, fileName)
}

func readGPXBuffer(fileBuffer []byte, fileName string) (Gpx, error) {
	gpx := Gpx{}
	err := xml.Unmarshal([]byte(fileBuffer), &gpx)

	if len(gpx.Tracks) > 0 || err != nil {
		return gpx, err
	}
	return gpx, newGpxFileError(fileName)

}
