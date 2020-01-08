package gpsabl

import "fmt"

// DepthParameterNotKnownError - Error when the given depth parameter is not known
type DepthParameterNotKnownError struct {
	err string
	// File - The path to the dir that caused this error
	GivenValue string
}

func (e *DepthParameterNotKnownError) Error() string { // Implement the Error Interface for the DepthParameterNotKnownError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// NewDepthParameterNotKnownError - Get a new DepthParameterNotKnownError struct
func NewDepthParameterNotKnownError(givenValue string) *DepthParameterNotKnownError {
	return &DepthParameterNotKnownError{fmt.Sprintf("The given -depth \"%s\" is not known.", givenValue), givenValue}
}

// CorrectionParameterNotKnownError - Error when the given correction parameter is not known
type CorrectionParameterNotKnownError struct {
	err string
	// File - The path to the dir that caused this error
	GivenValue string
}

func (e *CorrectionParameterNotKnownError) Error() string { // Implement the Error Interface for the DepthParameterNotKnownError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// NewCorrectionParameterNotKnownError - Get a new DepthParameterNotKnownError struct
func NewCorrectionParameterNotKnownError(givenValue string) *CorrectionParameterNotKnownError {
	return &CorrectionParameterNotKnownError{fmt.Sprintf("The given -correction \"%s\" is not known.", givenValue), givenValue}
}
