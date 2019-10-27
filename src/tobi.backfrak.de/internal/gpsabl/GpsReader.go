package gpsabl

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
