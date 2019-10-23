package main

import (
	"os"

	"tobi.backfrak.de/internal/gpsabl"
)

func main() {
	gpsabl.ReadGPX("Let's go :)")

	os.Exit(0)
}
