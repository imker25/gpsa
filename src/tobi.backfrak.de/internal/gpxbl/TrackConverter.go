package gpxbl

import (
	"tobi.backfrak.de/internal/gpsabl"
)

// ConvertTrk - Convert a gpxbl.Trk to a gpsabl.Track
func ConvertTrk(track Trk) gpsabl.Track {

	res := GetTrackInfo(track)

	return ConvertTrackInfo(res)
}

// ConvertTrackInfo - Convert a gpxbl.TrackInfo to a gpsabl.Track
func ConvertTrackInfo(track TrackInfo) gpsabl.Track {

	res := gpsabl.Track{}

	res.Name = track.Name
	res.NumberOfSegments = track.NumberOfSegments
	res.Description = track.Description
	res.AtituteRange = track.AtituteRange
	res.MaximumAtitute = track.MaximumAtitute
	res.MinimumAtitute = track.MinimumAtitute

	// res.Distance = // Calc the overall distance
	// res.TrackSegments // Get the needed data

	return res
}
