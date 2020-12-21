package gpsabl

import "fmt"

// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// GetOutlines - Get the OutputLine entries from a TrackFile depending on the analisis depths
func GetOutlines(trackFile TrackFile, depth DepthArg) ([]OutputLine, error) {
	ret := []OutputLine{}

	switch depth {
	case FILE:
		ret = append(ret, getOutlineFromTrackFile(trackFile))
	case TRACK:
		ret = append(ret, getOutlinesFromTracks(trackFile)...)
	case SEGMENT:
		ret = append(ret, getOutlinesFromTrackSegments(trackFile)...)
	default:
		return nil, NewDepthParameterNotKnownError(depth)
	}

	return ret, nil
}

// StripOutlines - Get the input outlines stripped of from inner data, so serialization will work fine
func StripOutlines(lines []OutputLine) []OutputLine {
	ret := []OutputLine{}
	// segs := []gpsabl.TrackSegment{}

	for _, line := range lines {
		data := ExtendedTrackSummary{}
		data.Distance = line.Data.GetDistance()
		data.HorizontalDistance = line.Data.GetHorizontalDistance()
		data.MaximumAltitude = line.Data.GetMaximumAltitude()
		data.MinimumAltitude = line.Data.GetMinimumAltitude()
		data.ElevationGain = line.Data.GetElevationGain()
		data.ElevationLose = line.Data.GetElevationLose()
		data.AltitudeRange = float64(line.Data.GetAltitudeRange())
		data.UpwardsDistance = line.Data.GetUpwardsDistance()
		data.DownwardsDistance = line.Data.GetDownwardsDistance()

		data.TimeDataValid = line.Data.GetTimeDataValid()
		if data.TimeDataValid {
			data.StartTime = line.Data.GetStartTime()
			data.EndTime = line.Data.GetEndTime()
			data.Duration = data.EndTime.Sub(data.StartTime)
			data.MovingTime = line.Data.GetMovingTime()
			data.UpwardsTime = line.Data.GetUpwardsTime()
			data.UpwardsSpeed = line.Data.GetUpwardsSpeed()
			data.DownwardsTime = line.Data.GetDownwardsTime()
			data.DownwardsSpeed = line.Data.GetDownwardsSpeed()
			data.AverageSpeed = line.Data.GetAvarageSpeed()
		}

		newLine := OutputLine{}
		newLine.Name = line.Name
		newLine.Data = data

		ret = append(ret, newLine)
	}

	return ret
}

// getOutlineFromTrackFile - Get the Outline for File depth analisis
func getOutlineFromTrackFile(trackFile TrackFile) OutputLine {
	return *NewOutputLine(getLineNameFromTrackFile(trackFile), TrackSummaryProvider(trackFile))
}

// getOutlinesFromTrackSegments - Get the Outlines for Segment depth analisis
func getOutlinesFromTrackSegments(trackFile TrackFile) []OutputLine {
	ret := []OutputLine{}
	for iTrack, track := range trackFile.Tracks {
		for iSeg, seg := range track.TrackSegments {
			info := TrackSummaryProvider(seg)
			name := fmt.Sprintf("%s: Segment #%d", getLineNameFromTrack(track, trackFile, iTrack), iSeg+1)
			entry := NewOutputLine(name, info)
			ret = append(ret, *entry)
		}
	}

	return ret
}

// getOutlinesFromTracks - Get the Outlines for Track depth analisis
func getOutlinesFromTracks(trackFile TrackFile) []OutputLine {
	ret := []OutputLine{}
	for i, track := range trackFile.Tracks {
		info := TrackSummaryProvider(track)
		name := getLineNameFromTrack(track, trackFile, i)
		entry := NewOutputLine(name, info)
		ret = append(ret, *entry)
	}

	return ret
}

func getLineNameFromTrack(track Track, parent TrackFile, index int) string {
	if track.Name != "" {
		return fmt.Sprintf("%s: %s", getLineNameFromTrackFile(parent), track.Name)
	}

	return fmt.Sprintf("%s: Track #%d", getLineNameFromTrackFile(parent), index+1)

}

func getLineNameFromTrackFile(trackFile TrackFile) string {
	if trackFile.Name != "" {
		return trackFile.Name
	}

	return trackFile.FilePath
}
