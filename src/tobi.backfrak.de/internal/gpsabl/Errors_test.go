package gpsabl

import (
	"strings"
	"testing"
)

func TestDepthParametrNotKnownErrorStruct(t *testing.T) {
	val := "asdgfg"
	err := NewDepthParametrNotKnownError(val)

	if err.GivenValue != val {
		t.Errorf("The GivenValue was %s, but %s was expected", err.GivenValue, val)
	}

	if strings.Contains(err.Error(), val) == false {
		t.Errorf("The error messaage of DepthParametrNotKnownError does not contain the expected GivenValue")
	}
}

func TestCorectionParamterNotKnownError(t *testing.T) {
	val := "asdgfg"
	err := NewCorectionParamterNotKnownError(val)

	if err.GivenValue != val {
		t.Errorf("The GivenValue was %s, but %s was expected", err.GivenValue, val)
	}

	if strings.Contains(err.Error(), val) == false {
		t.Errorf("The error messaage of CorectionParamterNotKnownError does not contain the expected GivenValue")
	}
}
