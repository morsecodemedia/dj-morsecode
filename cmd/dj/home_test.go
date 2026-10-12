package main

import (
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/controls"
)

func TestHomeItems(
	t *testing.T,
) {

	expected := []struct {
		label  string
		action controls.Action
	}{
		{
			label:  "Pick a vibe",
			action: controls.ActionOpenVibes,
		},
		{
			label:  "Set a mood",
			action: controls.ActionOpenMoods,
		},
		{
			label:  "Choose a genre",
			action: controls.ActionOpenGenres,
		},
		{
			label:  "Tune a station",
			action: controls.ActionOpenStations,
		},
		{
			label:  "Play something",
			action: controls.ActionOpenYouTube,
		},
	}

	if len(homeItems) !=
		len(expected) {

		t.Fatalf(
			"expected %d home items, got %d",
			len(expected),
			len(homeItems),
		)
	}

	for i, item := range expected {

		if homeItems[i].Label !=
			item.label {

			t.Errorf(
				"item %d: expected label %q, got %q",
				i,
				item.label,
				homeItems[i].Label,
			)
		}

		if homeItems[i].Action !=
			item.action {

			t.Errorf(
				"item %d: unexpected action %v, got %v",
				i,
				item.action,
				homeItems[i].Action,
			)
		}

	}

}

func TestHomeActiveWhenNothingIsPlaying(
	t *testing.T,
) {

	m := model{}

	if !m.homeActive() {

		t.Fatal(
			"expected home to be active",
		)
	}

}

func TestHomeInactiveWithNowPlaying(
	t *testing.T,
) {

	m := model{
		NowPlaying: "Tame Impala - Loser",
	}

	if m.homeActive() {

		t.Fatal(
			"expected home to be inactive",
		)
	}

}
