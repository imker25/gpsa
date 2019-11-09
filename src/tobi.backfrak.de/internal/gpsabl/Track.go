package gpsabl

type Track struct {
	Name string
	Description string
	NumberOfSegments    int
	Distance            float32
	AtituteRange        float32
	MinimumAtitute      float32
	MaximumAtitute      float32
	
	TrackSegments []TrackSegment
}

type TrackSegment struct {
	TrackPoints []TrackPoint
}

type TrackPoint struct {
	Elevation float32 
	Latitude  float32 
	Longitude float32 
	HorizontalDistanceBefore float32
	HorizontalDistanceNext float32
	VerticalDistanceBefore float32
	VerticalDistanceNext float32
}