package radio

import (
	"testing"
	"time"
)

func TestChooseMatchesCriteria(t *testing.T) {

	station, ok := Choose(
		Criteria{
			Contexts: []string{
				"coding",
			},
			MaxEnergy: 2,
		},
		History{},
		ChooseOptions{
			RecentLimit: 3,
		},
	)
	if !ok {
		t.Fatal("expected station choice")
	}

	if station.ID != "lofi247" {
		t.Errorf(
			"expected station %q, got %q",
			"lofi247",
			station.ID,
		)
	}

}

func TestChooseAvoidsCurrentStation(t *testing.T) {

	var history History

	history.Add(
		"lofi247",
		time.Now(),
	)

	station, ok := Choose(
		Criteria{
			Contexts: []string{
				"coding",
			},
			MaxEnergy: 2,
		},
		history,
		ChooseOptions{
			RecentLimit: 3,
		},
	)
	if !ok {
		t.Fatal("expected station choice")
	}

	if station.ID == "lofi247" {
		t.Fatal("expected current station to be excluded")
	}

}

func TestChooseRelaxesRecentHistory(t *testing.T) {

	var history History

	history.Add(
		"somafm-groovesalad",
		time.Now(),
	)

	history.Add(
		"lofi247",
		time.Now(),
	)

	station, ok := Choose(
		Criteria{
			Contexts: []string{
				"coding",
			},
			MaxEnergy: 2,
		},
		history,
		ChooseOptions{
			RecentLimit: 2,
		},
	)
	if !ok {
		t.Fatal("expected station choice after history relaxation")
	}

	if station.ID == "lofi247" {
		t.Fatal("expected current station to remain excluded")
	}

}

func TestChooseNeverFallsBackToCurrentStation(t *testing.T) {

	var history History

	history.Add(
		"somafm-dronezone",
		time.Now(),
	)

	_, ok := Choose(
		Criteria{
			Genres: []string{
				"Ambient",
			},
			Contexts: []string{
				"sleep",
			},
			MaxEnergy: 1,
		},
		history,
		ChooseOptions{
			RecentLimit: 3,
		},
	)

	if ok {
		t.Fatal("expected no choice when current station is only match")
	}

}

func TestChooseNoMatchingStations(t *testing.T) {

	_, ok := Choose(
		Criteria{
			Genres: []string{
				"Definitely Not A Genre",
			},
		},
		History{},
		ChooseOptions{
			RecentLimit: 3,
		},
	)

	if ok {
		t.Fatal("expected no station choice")
	}

}

func TestChooseUsesCandidateChooser(t *testing.T) {

	chooseLast := func(
		candidates []Station,
	) (Station, bool) {

		return candidates[len(candidates)-1], true

	}

	station, ok := Choose(
		Criteria{
			Contexts: []string{
				"coding",
			},
			MaxEnergy: 2,
		},
		History{},
		ChooseOptions{
			RecentLimit: 3,
			Chooser:     chooseLast,
		},
	)
	if !ok {
		t.Fatal("expected station choice")
	}

	matches := Match(
		Criteria{
			Contexts: []string{
				"coding",
			},
			MaxEnergy: 2,
		},
	)

	expected := matches[len(matches)-1]

	if station.ID != expected.ID {
		t.Errorf(
			"expected station %q, got %q",
			expected.ID,
			station.ID,
		)
	}

}

func TestChooseExcludesStations(t *testing.T) {

	station, ok := Choose(
		Criteria{
			Genres: []string{
				"Trance",
			},
		},
		History{},
		ChooseOptions{
			RecentLimit: 3,
			ExcludeIDs: []string{
				"afterhoursfm",
			},
		},
	)
	if !ok {
		t.Fatal("expected station choice")
	}

	if station.ID == "afterhoursfm" {
		t.Fatal("expected excluded station not to be chosen")
	}

}
