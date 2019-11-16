package gpsabl

import "testing"

func TestNewTrackFile(t *testing.T) {
	file := NewTrackFile("myFile")

	if file.FilePath != "myFile" {
		t.Errorf("The FilePath is %s, but %s was expected", file.FilePath, "myFile")
	}
}

func TestTrackFileIsTrackSummary(t *testing.T) {
	file := TrackFile{}
	sum := TrackInfoProvider(file)

	if sum == nil {
		t.Errorf("The Track struct does not implement the TrackInfoProvider interface as expected")
	}

	if sum.GetDistance() != 0.0 {
		t.Errorf("The GetDistance() does return %f, but %f was expected", sum.GetDistance(), 0.0)
	}

	if sum.GetAtituteRange() != 0.0 {
		t.Errorf("The GetAtituteRange() does return %f, but %f was expected", sum.GetAtituteRange(), 0.0)
	}

	if sum.GetMaximumAtitute() != 0.0 {
		t.Errorf("The GetMaximumAtitute() does return %f, but %f was expected", sum.GetMaximumAtitute(), 0.0)
	}

	if sum.GetMinimumAtitute() != 0.0 {
		t.Errorf("The GetMinimumAtitute() does return %f, but %f was expected", sum.GetMinimumAtitute(), 0.0)
	}
}

func TestTrackIsTrackSummary(t *testing.T) {
	trk := Track{}
	sum := TrackInfoProvider(trk)

	if sum == nil {
		t.Errorf("The Track struct does not implement the TrackInfoProvider interface as expected")
	}

	if sum.GetDistance() != 0.0 {
		t.Errorf("The GetDistance() does return %f, but %f was expected", sum.GetDistance(), 0.0)
	}

	if sum.GetAtituteRange() != 0.0 {
		t.Errorf("The GetAtituteRange() does return %f, but %f was expected", sum.GetAtituteRange(), 0.0)
	}

	if sum.GetMaximumAtitute() != 0.0 {
		t.Errorf("The GetMaximumAtitute() does return %f, but %f was expected", sum.GetMaximumAtitute(), 0.0)
	}

	if sum.GetMinimumAtitute() != 0.0 {
		t.Errorf("The GetMinimumAtitute() does return %f, but %f was expected", sum.GetMinimumAtitute(), 0.0)
	}
}

func TestTrackSegmentIsTrackSummary(t *testing.T) {
	seg := TrackSegment{}
	sum := TrackInfoProvider(seg)

	if sum == nil {
		t.Errorf("The Track struct does not implement the TrackInfoProvider interface as expected")
	}

	if sum.GetDistance() != 0.0 {
		t.Errorf("The GetDistance() does return %f, but %f was expected", sum.GetDistance(), 0.0)
	}

	if sum.GetAtituteRange() != 0.0 {
		t.Errorf("The GetAtituteRange() does return %f, but %f was expected", sum.GetAtituteRange(), 0.0)
	}

	if sum.GetMaximumAtitute() != 0.0 {
		t.Errorf("The GetMaximumAtitute() does return %f, but %f was expected", sum.GetMaximumAtitute(), 0.0)
	}

	if sum.GetMinimumAtitute() != 0.0 {
		t.Errorf("The GetMinimumAtitute() does return %f, but %f was expected", sum.GetMinimumAtitute(), 0.0)
	}
}
