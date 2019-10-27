package gpsabl

// Copyright 2019 by Tobias Zellner. All
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
}

// ReadGPX - Read a GPX file
func ReadGPX(filename string) (Gpx, error) {

	if filename == "" {
		return *(&Gpx{}), nil
	}

	fmt.Println("Read file: " + filename)
	xmlfile, err := ioutil.ReadFile(filename)
	if err != nil {
		HandleError(err)
	}
	return readGPXBuffer(xmlfile)
}

func readGPXBuffer(fileBuffer []byte) (Gpx, error) {
	gpx := &Gpx{}
	err := xml.Unmarshal([]byte(fileBuffer), &gpx)

	return *gpx, err
}
