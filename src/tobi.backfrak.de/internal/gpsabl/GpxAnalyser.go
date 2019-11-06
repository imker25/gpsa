package gpsabl

// TrackInfo - Struct that contains information about a track
type TrackInfo struct {
	Track               Trk
	Name                string
	Description         string
	NumberOfSegments    int
	NumberOfTrackPoints int
	Distance            int64
	AtituteRange        int32
}

// GetTrackInfo - Creats a TrackInfo struct out of a Trk struct
func GetTrackInfo(track Trk) TrackInfo {
	info := TrackInfo{}
	info.Track = track
	info.Name = track.Name
	info.Description = track.Description
	info.NumberOfSegments = len(track.TrackSegments)
	info.NumberOfTrackPoints = getNumberOfTrackPoints(track)

	return info
}

func getNumberOfTrackPoints(track Trk) int {
	count := 0
	for _, seg := range track.TrackSegments {
		count = count + len(seg.TrackPoints)
	}

	return count
}
