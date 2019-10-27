package gpsabl

import (
	"fmt"
	"os"
)

// HandleError - Handles an error
func HandleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
