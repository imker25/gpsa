package main

// Runs the mage build script build/workflow/magefiles/magefile.go

import (
	"os"

	"github.com/magefile/mage/mage"
)

func main() { os.Exit(mage.Main()) }
