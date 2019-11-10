package gpsabl

// Track - the struct to handle track info in gpsa
type Track struct {
	Name             string
	Description      string
	NumberOfSegments int
	Distance         float32
	AtituteRange     float32
	MinimumAtitute   float32
	MaximumAtitute   float32

	TrackSegments []TrackSegment
}

// TrackSegment - the struct to handle track segment info in gpsa
type TrackSegment struct {
	TrackPoints    []TrackPoint
	Distance       float32
	AtituteRange   float32
	MinimumAtitute float32
	MaximumAtitute float32
}

// TrackPoint - the struct to handle track point info in gpsa
type TrackPoint struct {
	Elevation                float32
	Latitude                 float32
	Longitude                float32
	HorizontalDistanceBefore float32
	HorizontalDistanceNext   float32
	VerticalDistanceBefore   float32
	VerticalDistanceNext     float32
}
