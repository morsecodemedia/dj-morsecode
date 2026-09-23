package metadata

import "testing"

func TestPlaybackItemIsTrack(t *testing.T) {

	item := PlaybackItem{
		Type: PlaybackTrack,
	}

	if !item.IsTrack() {
		t.Fatal("expected track item")
	}

}

func TestZeroPlaybackItemIsNotTrack(t *testing.T) {

	var item PlaybackItem

	if item.IsTrack() {
		t.Fatal("expected zero playback item not to be track")
	}

}
