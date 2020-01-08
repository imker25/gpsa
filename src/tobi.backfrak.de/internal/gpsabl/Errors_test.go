package gpsabl

import (
	"strings"
	"testing"
)

func TestDepthParameterNotKnownErrorStruct(t *testing.T) {
	val := "asdgfg"
	err := NewDepthParameterNotKnownError(val)

	if err.GivenValue != val {
		t.Errorf("The GivenValue was %s, but %s was expected", err.GivenValue, val)
	}

	if strings.Contains(err.Error(), val) == false {
		t.Errorf("The error message of DepthParameterNotKnownError does not contain the expected GivenValue")
	}
}

func TestCorrectionParameterNotKnownError(t *testing.T) {
	val := "asdgfg"
	err := NewCorrectionParameterNotKnownError(val)

	if err.GivenValue != val {
		t.Errorf("The GivenValue was %s, but %s was expected", err.GivenValue, val)
	}

	if strings.Contains(err.Error(), val) == false {
		t.Errorf("The error messaage of CorrectionParameterNotKnownError does not contain the expected GivenValue")
	}
}
