package gpsabl

// InputFileType - a string type to implement the enum pattern
type InputFileType string

const (
	// FilePath - The input file is given as file path
	FilePath InputFileType = "FilePath"
)

// ValidInputFileTypes - Array of inputFileType that contains all valid inputFileTypes
var ValidInputFileTypes = []InputFileType{FilePath}

// inputFile - A struct to tell what kind of input we have, so we can proccess it right
type InputFile struct {
	// Type - The type of the input file
	Type InputFileType
	// Name - The Files path in case of Type=FilePath, the name of the buffer in other cases
	Name string
	// Buffer - nil in case of Type=FilePath, the files content in other cases
	Buffer []byte
}

// NewInputFileWithPath - Get a new inputFile struct from a file path
func NewInputFileWithPath(filePath string) *InputFile {
	file := InputFile{}
	file.Name = filePath
	file.Type = FilePath
	file.Buffer = nil

	return &file
}

// Check if this inputFile has a valid inputFileType as Type
func (file InputFile) InputFileTypeValid() bool {
	return inputFileTypeValid(file.Type)
}

// Check if the given inputFileType is valid
func inputFileTypeValid(fileType InputFileType) bool {
	for _, t := range ValidInputFileTypes {
		if t == fileType {
			return true
		}
	}

	return false
}
