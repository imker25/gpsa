package gpsabl

import "fmt"

// DepthParametrNotKnownError - Error when the given depth paramter is not known
type DepthParametrNotKnownError struct {
	err string
	// File - The path to the dir that caused this error
	GivenValue string
}

func (e *DepthParametrNotKnownError) Error() string { // Implement the Error Interface for the DepthParametrNotKnownError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// NewDepthParametrNotKnownError - Get a new DepthParametrNotKnownError struct
func NewDepthParametrNotKnownError(givenValue string) *DepthParametrNotKnownError {
	return &DepthParametrNotKnownError{fmt.Sprintf("The given -depth \"%s\" is not known.", givenValue), givenValue}
}

// CorectionParamterNotKnownError - Error when the given corection paramter is not known
type CorectionParamterNotKnownError struct {
	err string
	// File - The path to the dir that caused this error
	GivenValue string
}

func (e *CorectionParamterNotKnownError) Error() string { // Implement the Error Interface for the DepthParametrNotKnownError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// NewCorectionParamterNotKnownError - Get a new DepthParametrNotKnownError struct
func NewCorectionParamterNotKnownError(givenValue string) *CorectionParamterNotKnownError {
	return &CorectionParamterNotKnownError{fmt.Sprintf("The given -corection \"%s\" is not known.", givenValue), givenValue}
}
