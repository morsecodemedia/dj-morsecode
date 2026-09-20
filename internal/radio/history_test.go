package radio

import (
	"testing"
	"time"
)

func TestHistoryAddAndCurrent(t *testing.T) {

	var history History

	tunedAt := time.Date(
		2026,
		time.September,
		20,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	history.Add(
		"wxpn",
		tunedAt,
	)

	current, ok := history.Current()
	if !ok {
		t.Fatal("expected current tune")
	}

	if current.StationID != "wxpn" {
		t.Errorf(
			"expected station %q, got %q",
			"wxpn",
			current.StationID,
		)
	}

	if !current.TunedAt.Equal(tunedAt) {
		t.Errorf(
			"expected tuned time %s, got %s",
			tunedAt,
			current.TunedAt,
		)
	}

}

func TestHistoryPrevious(t *testing.T) {

	var history History

	history.Add(
		"wxpn",
		time.Now(),
	)

	history.Add(
		"wrti",
		time.Now(),
	)

	previous, ok := history.Previous()
	if !ok {
		t.Fatal("expected previous tune")
	}

	if previous.StationID != "wxpn" {
		t.Errorf(
			"expected previous station %q, got %q",
			"wxpn",
			previous.StationID,
		)
	}

}

func TestHistoryPreviousRequiresTwoTunes(t *testing.T) {

	var history History

	history.Add(
		"wxpn",
		time.Now(),
	)

	_, ok := history.Previous()

	if ok {
		t.Fatal("expected no previous tune")
	}

}

func TestHistoryContainsRecent(t *testing.T) {

	var history History

	history.Add("wxpn", time.Now())
	history.Add("wrti", time.Now())
	history.Add("radioparadise", time.Now())
	history.Add("lofi247", time.Now())

	if !history.ContainsRecent("wrti", 3) {
		t.Fatal("expected wrti in last 3 tunes")
	}

	if history.ContainsRecent("wxpn", 3) {
		t.Fatal("did not expect wxpn in last 3 tunes")
	}

}

func TestHistoryIgnoresEmptyStation(t *testing.T) {

	var history History

	history.Add(
		"",
		time.Now(),
	)

	if len(history.Tunes) != 0 {
		t.Fatalf(
			"expected no tunes, got %d",
			len(history.Tunes),
		)
	}

}

func TestHistoryIgnoresConsecutiveDuplicateStation(t *testing.T) {

	var history History

	first := time.Date(
		2026,
		time.September,
		20,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	second := first.Add(
		time.Minute,
	)

	history.Add(
		"wxpn",
		first,
	)

	history.Add(
		"wxpn",
		second,
	)

	if len(history.Tunes) != 1 {
		t.Fatalf(
			"expected 1 tune, got %d",
			len(history.Tunes),
		)
	}

	current, ok := history.Current()
	if !ok {
		t.Fatal("expected current tune")
	}

	if !current.TunedAt.Equal(first) {
		t.Errorf(
			"expected original tuned time %s, got %s",
			first,
			current.TunedAt,
		)
	}

}
