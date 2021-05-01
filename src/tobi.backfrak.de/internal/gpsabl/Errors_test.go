package gpsabl

import (
	"fmt"
	"strings"
	"testing"
)

func TestDepthParameterNotKnownErrorStruct(t *testing.T) {
	val := "asdgfg"
	err := NewDepthParameterNotKnownError(DepthArg(val))

	if err.GivenValue != DepthArg(val) {
		t.Errorf("The GivenValue was %s, but %s was expected", err.GivenValue, val)
	}

	if strings.Contains(err.Error(), val) == false {
		t.Errorf("The error message of DepthParameterNotKnownError does not contain the expected GivenValue")
	}
}

func TestSummaryParameterNotKnownErrorStruct(t *testing.T) {
	val := "asdgfg"
	err := NewSummaryParamaterNotKnown(SummaryArg(val))

	if err.GivenValue != SummaryArg(val) {
		t.Errorf("The GivenValue was %s, but %s was expected", err.GivenValue, val)
	}

	if strings.Contains(err.Error(), val) == false {
		t.Errorf("The error message of SummaryParameterNotKnownError does not contain the expected GivenValue")
	}
}

func TestCorrectionParameterNotKnownError(t *testing.T) {
	val := "asdgfg"
	err := NewCorrectionParameterNotKnownError(CorrectionParameter(val))

	if err.GivenValue != CorrectionParameter(val) {
		t.Errorf("The GivenValue was %s, but %s was expected", err.GivenValue, val)
	}

	if strings.Contains(err.Error(), val) == false {
		t.Errorf("The error message of CorrectionParameterNotKnownError does not contain the expected GivenValue")
	}
}

func TestMinimalMovingSpeedLessThenZeroError(t *testing.T) {

	val := -1.0
	err := NewMinimalMovingSpeedLessThenZero(val)

	if err.GivenValue != val {
		t.Errorf("The GivenValue was %f, but %f was expected", err.GivenValue, val)
	}

	if strings.Contains(err.Error(), fmt.Sprintf("%f", val)) == false {
		t.Errorf("The error message of MinimalMovingSpeedLessThenZero does not contain the expected GivenValue")
	}
}

func TestMinimalStepHightLessThenZero(t *testing.T) {

	val := -1.0
	err := NewMinimalStepHightLessThenZero(val)

	if err.GivenValue != val {
		t.Errorf("The GivenValue was %f, but %f was expected", err.GivenValue, val)
	}

	if strings.Contains(err.Error(), fmt.Sprintf("%f", val)) == false {
		t.Errorf("The error message of MinimalStepHightLessThenZero does not contain the expected GivenValue")
	}
}

func TestUnKnownInputFileTypeError(t *testing.T) {
	path := "myType"
	err := NewUnKnownInputFileTypeError(path)

	if err.inputFileType != path {
		t.Errorf("The Type was %s, but %s was expected", err.inputFileType, path)
	}

	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error message of UnKnownInputFileTypeError does not contain the expected type")
	}
}
