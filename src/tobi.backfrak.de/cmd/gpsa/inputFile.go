package main

// inputFileType - a string type to implement the enum pattern
type inputFileType string

const (
	// FilePath - The input file is given as file path
	FilePath inputFileType = "FilePath"
	// GpxBuffer - The input file is given as buffer containing a gpx files content
	GpxBuffer inputFileType = "GpxBuffer"
	// TcxBuffer -  The input file is given as buffer containing a tcx files content
	TcxBuffer inputFileType = "TcxBuffer"
)

// validInputFileTypes - Array of inputFileType that contains all valid inputFileTypes
var validInputFileTypes = []inputFileType{FilePath, GpxBuffer, TcxBuffer}

// inputFile - A struct to tell what kind of input we have, so we can proccess it right
type inputFile struct {
	// Type - The type of the input file
	Type inputFileType
	// Name - The Files path in case of Type=FilePath, the name of the buffer in other cases
	Name string
	// Buffer - nil in case of Type=FilePath, the files content in other cases
	Buffer []byte
}

// newInputFileWithPath - Get a new inputFile struct from a file path
func newInputFileWithPath(filePath string) *inputFile {
	file := inputFile{}
	file.Name = filePath
	file.Type = FilePath
	file.Buffer = nil

	return &file
}

// newInputFileGpxBuffer - Get a new inputFilw from a buffer containing a gpx files content
func newInputFileGpxBuffer(buffer []byte, name string) *inputFile {
	file := inputFile{}
	file.Name = name
	file.Type = GpxBuffer
	file.Buffer = buffer

	return &file
}

// newInputFileGpxBuffer - Get a new inputFilw from a buffer containing a tcx files content
func newInputFileTcxBuffer(buffer []byte, name string) *inputFile {
	file := inputFile{}
	file.Name = name
	file.Type = TcxBuffer
	file.Buffer = buffer

	return &file
}

// Check if this inputFile has a valid inputFileType as Type
func (file inputFile) inputFileTypeValid() bool {
	return inputFileTypeValid(file.Type)
}

// Check if the given inputFileType is valid
func inputFileTypeValid(fileType inputFileType) bool {
	for _, t := range validInputFileTypes {
		if t == fileType {
			return true
		}
	}

	return false
}
