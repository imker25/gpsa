package gpxbl

// TrackInfo - Struct that contains information about a track
type TrackInfo struct {
	Track               Trk
	Name                string
	Description         string
	NumberOfSegments    int
	NumberOfTrackPoints int
	Distance            int64
	AtituteRange        float32
	MinimumAtitute      float32
	MaximumAtitute      float32
}

// GetTrackInfo - Creats a TrackInfo struct out of a Trk struct
func GetTrackInfo(track Trk) TrackInfo {
	info := TrackInfo{}
	info.Track = track
	info.Name = track.Name
	info.Description = track.Description
	info.NumberOfSegments = len(track.TrackSegments)
	info.NumberOfTrackPoints = getNumberOfTrackPoints(track)
	info.MinimumAtitute = getMinimumAtitute(track)
	info.MaximumAtitute = getMaximumAtitute(track)
	info.AtituteRange = info.MaximumAtitute - info.MinimumAtitute

	return info
}

// GetAllTrackPoints - Get all points of this track as array
func (i TrackInfo) GetAllTrackPoints() []Trkpt {
	return getAllTrackPoints(i.Track)
}

func getAllTrackPoints(track Trk) []Trkpt {
	var pntList []Trkpt
	for _, seg := range track.TrackSegments {
		pntList = append(pntList, seg.TrackPoints...) // The "..." operator will expand the array see https://stackoverflow.com/questions/16248241/concatenate-two-slices-in-go
	}

	return pntList
}

func getNumberOfTrackPoints(track Trk) int {
	count := 0
	for _, seg := range track.TrackSegments {
		count = count + len(seg.TrackPoints)
	}

	return count
}

func getMinimumAtitute(track Trk) float32 {

	var min float32
	for i, pnt := range getAllTrackPoints(track) {
		if i == 0 || min > pnt.Elevation {
			min = pnt.Elevation
		}
	}

	return min
}

func getMaximumAtitute(track Trk) float32 {

	var max float32
	for i, pnt := range getAllTrackPoints(track) {
		if i == 0 || max < pnt.Elevation {
			max = pnt.Elevation
		}
	}

	return max
}
