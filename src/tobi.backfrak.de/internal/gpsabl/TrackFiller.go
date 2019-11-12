package gpsabl

// FillDistances - Adds the distance values to the basePoint
func FillDistances(basePoint, beforePoint, nextPoint TrackPoint) TrackPoint {
	retPoint := TrackPoint{}
	retPoint.Elevation = basePoint.Elevation
	retPoint.Latitude = basePoint.Latitude
	retPoint.Longitude = basePoint.Longitude

	if (beforePoint != TrackPoint{}) {
		retPoint.HorizontalDistanceBefore = Distance(basePoint, beforePoint)
		retPoint.VerticalDistanceBefore = basePoint.Elevation - beforePoint.Elevation
	}

	if (nextPoint != TrackPoint{}) {
		retPoint.HorizontalDistanceNext = Distance(basePoint, nextPoint)
		retPoint.VerticalDistanceNext = nextPoint.Elevation - basePoint.Elevation
	}

	return retPoint
}
