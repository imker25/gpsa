package gpsabl

import "os"

// OutputFormater - Interface for classes that can format a track output into a file format and write this file
type OutputFormater interface {

	// AddOutPut - Add the output values of a TrackFile to the out file buffer
	AddOutPut(trackFile TrackFile, depth string)

	// AddHeader- Add a AddHeader line to the out file buffer
	AddHeader()

	// WriteOutput - Write the output to the output file
	WriteOutput(outFile *os.File) error
}
