package main

import (
	"os"

	"tobi.backfrak.de/internal/gpsabl"
)

func main() {
	gpsabl.ReadGPX("")

	os.Exit(0)
}
