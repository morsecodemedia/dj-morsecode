package radio

import (
	"testing"
	"time"
)

func TestSelectFirstCandidate(t *testing.T) {

	candidates := []Station{
		{
			ID: "wxpn",
		},
		{
			ID: "radioparadise",
		},
	}

	station, ok := Select(
		candidates,
		History{},
		SelectionOptions{},
	)
	if !ok {
		t.Fatal("expected station selection")
	}

	if station.ID != "wxpn" {
		t.Errorf(
			"expected station %q, got %q",
			"wxpn",
			station.ID,
		)
	}

}

func TestSelectAvoidsRecentStations(t *testing.T) {

	var history History

	history.Add(
		"wxpn",
		time.Now(),
	)

	candidates := []Station{
		{
			ID: "wxpn",
		},
		{
			ID: "radioparadise",
		},
	}

	station, ok := Select(
		candidates,
		history,
		SelectionOptions{
			RecentLimit: 1,
		},
	)
	if !ok {
		t.Fatal("expected station selection")
	}

	if station.ID != "radioparadise" {
		t.Errorf(
			"expected station %q, got %q",
			"radioparadise",
			station.ID,
		)
	}

}

func TestSelectUsesRecentLimit(t *testing.T) {

	var history History

	history.Add(
		"wxpn",
		time.Now(),
	)

	history.Add(
		"wrti",
		time.Now(),
	)

	candidates := []Station{
		{
			ID: "wxpn",
		},
		{
			ID: "wrti",
		},
		{
			ID: "radioparadise",
		},
	}

	station, ok := Select(
		candidates,
		history,
		SelectionOptions{
			RecentLimit: 2,
		},
	)
	if !ok {
		t.Fatal("expected station selection")
	}

	if station.ID != "radioparadise" {
		t.Errorf(
			"expected station %q, got %q",
			"radioparadise",
			station.ID,
		)
	}

}

func TestSelectAllowsOlderStation(t *testing.T) {

	var history History

	history.Add(
		"wxpn",
		time.Now(),
	)

	history.Add(
		"wrti",
		time.Now(),
	)

	candidates := []Station{
		{
			ID: "wxpn",
		},
		{
			ID: "wrti",
		},
	}

	station, ok := Select(
		candidates,
		history,
		SelectionOptions{
			RecentLimit: 1,
		},
	)
	if !ok {
		t.Fatal("expected station selection")
	}

	if station.ID != "wxpn" {
		t.Errorf(
			"expected station %q, got %q",
			"wxpn",
			station.ID,
		)
	}

}

func TestSelectNoCandidates(t *testing.T) {

	_, ok := Select(
		nil,
		History{},
		SelectionOptions{},
	)

	if ok {
		t.Fatal("expected no station selection")
	}

}

func TestSelectAllCandidatesRecent(t *testing.T) {

	var history History

	history.Add(
		"wxpn",
		time.Now(),
	)

	history.Add(
		"radioparadise",
		time.Now(),
	)

	candidates := []Station{
		{
			ID: "wxpn",
		},
		{
			ID: "radioparadise",
		},
	}

	_, ok := Select(
		candidates,
		history,
		SelectionOptions{
			RecentLimit: 2,
		},
	)

	if ok {
		t.Fatal("expected no station selection")
	}

}

func TestSelectWithChooser(t *testing.T) {

	candidates := []Station{
		{
			ID: "wxpn",
		},
		{
			ID: "radioparadise",
		},
		{
			ID: "wrti",
		},
	}

	chooseLast := func(
		eligible []Station,
	) (Station, bool) {

		return eligible[len(eligible)-1], true

	}

	station, ok := SelectWithChooser(
		candidates,
		History{},
		SelectionOptions{},
		chooseLast,
	)
	if !ok {
		t.Fatal("expected station selection")
	}

	if station.ID != "wrti" {
		t.Errorf(
			"expected station %q, got %q",
			"wrti",
			station.ID,
		)
	}

}

func TestSelectWithChooserOnlyReceivesEligibleCandidates(t *testing.T) {

	var history History

	history.Add(
		"wxpn",
		time.Now(),
	)

	candidates := []Station{
		{
			ID: "wxpn",
		},
		{
			ID: "radioparadise",
		},
		{
			ID: "wrti",
		},
	}

	chooser := func(
		eligible []Station,
	) (Station, bool) {

		for _, station := range eligible {

			if station.ID == "wxpn" {
				t.Fatal(
					"chooser received recent station wxpn",
				)
			}

		}

		return eligible[0], true

	}

	_, ok := SelectWithChooser(
		candidates,
		history,
		SelectionOptions{
			RecentLimit: 1,
		},
		chooser,
	)
	if !ok {
		t.Fatal("expected station selection")
	}

}
func TestRandomCandidateEmpty(t *testing.T) {

	_, ok := RandomCandidate(
		nil,
	)

	if ok {
		t.Fatal("expected no random candidate")
	}

}
