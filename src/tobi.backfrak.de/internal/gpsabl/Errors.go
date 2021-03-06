package gpsabl

import "fmt"

// DepthParameterNotKnownError - Error when the given depth parameter is not known
type DepthParameterNotKnownError struct {
	err string
	// File - The path to the dir that caused this error
	GivenValue DepthArg
}

func (e *DepthParameterNotKnownError) Error() string { // Implement the Error Interface for the DepthParameterNotKnownError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// NewDepthParameterNotKnownError - Get a new DepthParameterNotKnownError struct
func NewDepthParameterNotKnownError(givenValue DepthArg) *DepthParameterNotKnownError {
	return &DepthParameterNotKnownError{fmt.Sprintf("The given -depth \"%s\" is not known.", givenValue), givenValue}
}

// CorrectionParameterNotKnownError - Error when the given correction parameter is not known
type CorrectionParameterNotKnownError struct {
	err string
	// File - The path to the dir that caused this error
	GivenValue CorrectionParameter
}

func (e *CorrectionParameterNotKnownError) Error() string { // Implement the Error Interface for the CorrectionParameterNotKnownError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// NewCorrectionParameterNotKnownError - Get a new CorrectionParameterNotKnownError struct
func NewCorrectionParameterNotKnownError(givenValue CorrectionParameter) *CorrectionParameterNotKnownError {
	return &CorrectionParameterNotKnownError{fmt.Sprintf("The given -correction \"%s\" is not known.", givenValue), givenValue}
}

// MinimalMovingSpeedLessThenZero - Error when the given -minimal-moving-speed is less then zero
type MinimalMovingSpeedLessThenZero struct {
	err string
	// File - The path to the dir that caused this error
	GivenValue float64
}

func (e *MinimalMovingSpeedLessThenZero) Error() string { // Implement the Error Interface for the MinimalMovingSpeedLessThenZero struct
	return fmt.Sprintf("%s", e.err)
}

// NewMinimalMovingSpeedLessThenZero - Get a new MinimalMovingSpeedLessThenZero struct
func NewMinimalMovingSpeedLessThenZero(givenValue float64) *MinimalMovingSpeedLessThenZero {
	return &MinimalMovingSpeedLessThenZero{fmt.Sprintf("The given -minimal-moving-speed \"%f\" is less then 0.0.", givenValue), givenValue}
}

// MinimalStepHightLessThenZero - Error when the given -minimal-step-hight is less then zero
type MinimalStepHightLessThenZero struct {
	err string
	// File - The path to the dir that caused this error
	GivenValue float64
}

func (e *MinimalStepHightLessThenZero) Error() string { // Implement the Error Interface for the MinimalMovingSpeedLessThenZero struct
	return fmt.Sprintf("%s", e.err)
}

// NewMinimalStepHightLessThenZero - Get a new MinimalMovingSpeedLessThenZero struct
func NewMinimalStepHightLessThenZero(givenValue float64) *MinimalStepHightLessThenZero {
	return &MinimalStepHightLessThenZero{fmt.Sprintf("The given -minimal-step-hight \"%f\" is less then 0.0.", givenValue), givenValue}
}

// SummaryParamaterNotKnown - Error when the given -summary is not known
type SummaryParamaterNotKnown struct {
	err string
	// File - The path to the dir that caused this error
	GivenValue SummaryArg
}

func (e *SummaryParamaterNotKnown) Error() string { // Implement the Error Interface for the SummaryParamaterNotKnown struct
	return fmt.Sprintf("%s", e.err)
}

// NewSummaryParamaterNotKnown - Get a new SummaryParamaterNotKnown struct
func NewSummaryParamaterNotKnown(givenValue SummaryArg) *SummaryParamaterNotKnown {
	return &SummaryParamaterNotKnown{fmt.Sprintf("The given -summary \"%s\" is not known.", givenValue), givenValue}
}

// UnKnowninputFileTypeError - Error when getting a unknown inputFileType
type UnKnownInputFileTypeError struct {
	err string
	// File - The path to the file that caused this error
	inputFileType string
}

func (e *UnKnownInputFileTypeError) Error() string { // Implement the Error Interface for the UnKnowninputFileTypeError struct
	return fmt.Sprintf("%s", e.err)
}

// newUnKnownInputStreamError - Get a new UnKnownFileTypeError struct
func NewUnKnownInputFileTypeError(inputFileType string) *UnKnownInputFileTypeError {
	return &UnKnownInputFileTypeError{fmt.Sprintf("Can not process inputFileType \"%s\".", inputFileType), inputFileType}
}

// TimeFormatNotKnown - Error when the given -summary is not known
type TimeFormatNotKnown struct {
	err string
	// File - The path to the dir that caused this error
	GivenValue TimeFormat
}

func (e *TimeFormatNotKnown) Error() string { // Implement the Error Interface for the TimeFormatNotKnown struct
	return fmt.Sprintf("%s", e.err)
}

// NewTimeFormatNotKnown - Get a new TimeFormatNotKnown struct
func NewTimeFormatNotKnown(givenValue TimeFormat) *TimeFormatNotKnown {
	return &TimeFormatNotKnown{fmt.Sprintf("The given -summary \"%s\" is not known.", givenValue), givenValue}
}
