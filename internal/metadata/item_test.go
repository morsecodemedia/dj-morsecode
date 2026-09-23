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

func TestPlaybackItemDisplayTitle(t *testing.T) {

	item := PlaybackItem{
		Type:   PlaybackTrack,
		Artist: "Steve Lacy",
		Title:  "Oh Yeah",
	}

	if got := item.DisplayTitle(); got != "Steve Lacy - Oh Yeah" {
		t.Errorf(
			"expected display title %q, got %q",
			"Steve Lacy - Oh Yeah",
			got,
		)
	}

}

func TestNonTrackPlaybackItemHasNoDisplayTitle(t *testing.T) {

	item := PlaybackItem{
		Type:  PlaybackStationID,
		Title: "Z100",
	}

	if got := item.DisplayTitle(); got != "" {
		t.Errorf(
			"expected empty display title, got %q",
			got,
		)
	}

}

func TestPlaybackItemObserved(t *testing.T) {

	item := PlaybackItem{
		RawTitle: "Dominic Fike - Babydoll",
	}

	if !item.Observed() {
		t.Fatal("expected playback item to be observed")
	}

}

func TestZeroPlaybackItemNotObserved(t *testing.T) {

	var item PlaybackItem

	if item.Observed() {
		t.Fatal("expected zero playback item not to be observed")
	}

}

func TestNormalizeIHeartWithoutPlaybackItem(t *testing.T) {

	fields := loadMetadataFixture(
		t,
		"z100-no-item.json",
	)

	item := Normalize(
		fields,
	)

	if item.Observed() {
		t.Fatalf(
			"expected no playback observation, got %q",
			item.RawTitle,
		)
	}

	if item.Type != PlaybackUnknown {
		t.Fatalf(
			"expected unknown, got %q",
			item.Type,
		)
	}

}
