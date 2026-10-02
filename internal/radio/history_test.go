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

func TestHistoryZeroValueHasNoCurrent(
	t *testing.T,
) {

	var history History

	if _, ok := history.Current(); ok {
		t.Fatal("expected no current tune")
	}

	if history.CanBack() {
		t.Fatal("expected no back history")
	}

	if history.CanForward() {
		t.Fatal("expected no forward history")
	}

}

func TestHistoryNavigatesBackwardAndForward(
	t *testing.T,
) {

	var history History

	now := time.Now()

	history.Add("a", now)
	history.Add("b", now.Add(time.Second))
	history.Add("c", now.Add(2*time.Second))

	current, ok := history.Current()
	if !ok || current.StationID != "c" {
		t.Fatalf(
			"expected current station c, got %+v",
			current,
		)
	}

	back, ok := history.Back()
	if !ok || back.StationID != "b" {
		t.Fatalf(
			"expected back station b, got %+v",
			back,
		)
	}

	back, ok = history.Back()
	if !ok || back.StationID != "a" {
		t.Fatalf(
			"expected back station a, got %+v",
			back,
		)
	}

	if _, ok := history.Back(); ok {
		t.Fatal(
			"expected beginning of history",
		)
	}

	forward, ok := history.Forward()
	if !ok || forward.StationID != "b" {
		t.Fatalf(
			"expected forward station b, got %+v",
			forward,
		)
	}

	forward, ok = history.Forward()
	if !ok || forward.StationID != "c" {
		t.Fatalf(
			"expected forward station c, got %+v",
			forward,
		)
	}

	if _, ok := history.Forward(); ok {
		t.Fatal(
			"expected end of history",
		)
	}

}

func TestHistoryAddTruncatesForwardHistory(
	t *testing.T,
) {

	var history History

	now := time.Now()

	history.Add("a", now)
	history.Add("b", now.Add(time.Second))
	history.Add("c", now.Add(2*time.Second))
	history.Add("d", now.Add(3*time.Second))

	_, _ = history.Back()
	_, _ = history.Back()

	history.Add(
		"e",
		now.Add(4*time.Second),
	)

	if len(history.Tunes) != 3 {
		t.Fatalf(
			"expected 3 tunes, got %d",
			len(history.Tunes),
		)
	}

	expected := []string{
		"a",
		"b",
		"e",
	}

	for i, stationID := range expected {

		if history.Tunes[i].StationID !=
			stationID {

			t.Errorf(
				"expected tune %d to be %q, got %q",
				i,
				stationID,
				history.Tunes[i].StationID,
			)

		}

	}

	if history.CanForward() {
		t.Fatal(
			"expected forward history to be discarded",
		)
	}

}

func TestHistoryAddCurrentAfterBackDoesNotTruncateForward(
	t *testing.T,
) {

	var history History

	now := time.Now()

	history.Add("a", now)
	history.Add("b", now.Add(time.Second))
	history.Add("c", now.Add(2*time.Second))

	_, _ = history.Back()

	history.Add(
		"b",
		now.Add(3*time.Second),
	)

	if len(history.Tunes) != 3 {
		t.Fatalf(
			"expected forward history preserved, got %d tunes",
			len(history.Tunes),
		)
	}

	if !history.CanForward() {
		t.Fatal(
			"expected forward history to remain",
		)
	}

}

func TestHistoryIdentifiesCurrentEntryAfterBack(
	t *testing.T,
) {

	var history History

	now := time.Now()

	history.Add(
		"z100",
		now,
	)

	history.Add(
		"skaworld",
		now.Add(time.Second),
	)

	_, ok := history.Back()
	if !ok {
		t.Fatal(
			"expected back navigation",
		)
	}

	if !history.IsCurrent(0) {
		t.Fatal(
			"expected first entry to be current",
		)
	}

	if history.IsCurrent(1) {
		t.Fatal(
			"expected forward entry not to be current",
		)
	}

}
