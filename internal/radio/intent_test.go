package radio

import "testing"

func TestVibeIntent(t *testing.T) {

	preset, ok := FindPreset("focus")
	if !ok {
		t.Fatal("expected focus preset")
	}

	intent := VibeIntent(*preset)

	if intent.Type != IntentVibe {
		t.Errorf(
			"expected vibe intent, got %q",
			intent.Type,
		)
	}

	if intent.ID != "focus" {
		t.Errorf(
			"expected intent ID %q, got %q",
			"focus",
			intent.ID,
		)
	}

	if !intent.Active() {
		t.Fatal("expected active intent")
	}

	if len(intent.Criteria.Contexts) == 0 {
		t.Fatal("expected vibe criteria")
	}

}

func TestMoodIntent(t *testing.T) {

	family, ok := FindMoodFamily("dreamy")
	if !ok {
		t.Fatal("expected dreamy mood family")
	}

	intent := MoodIntent(*family)

	if intent.Type != IntentMood {
		t.Errorf(
			"expected mood intent, got %q",
			intent.Type,
		)
	}

	if intent.Name != "Dreamy" {
		t.Errorf(
			"expected intent name %q, got %q",
			"Dreamy",
			intent.Name,
		)
	}

	if len(intent.Criteria.Moods) == 0 {
		t.Fatal("expected mood criteria")
	}

}

func TestGenreIntent(t *testing.T) {

	intent := GenreIntent("Ska")

	if intent.Type != IntentGenre {
		t.Errorf(
			"expected genre intent, got %q",
			intent.Type,
		)
	}

	if intent.Name != "Ska" {
		t.Errorf(
			"expected intent name %q, got %q",
			"Ska",
			intent.Name,
		)
	}

	if len(intent.Criteria.Genres) != 1 ||
		intent.Criteria.Genres[0] != "Ska" {

		t.Fatalf(
			"unexpected genre criteria: %v",
			intent.Criteria.Genres,
		)
	}

}

func TestZeroIntentIsInactive(t *testing.T) {

	var intent Intent

	if intent.Active() {
		t.Fatal("expected zero intent to be inactive")
	}

}
