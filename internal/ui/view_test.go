package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
)

func renderTestTrack(
	song music.Song,
	trackContext metadata.TrackContext,
) string {

	return Render(
		song,
		trackContext,
		72,
		true,
		0,
		30*time.Second,
		"",
		"UNAVAILABLE",
		"",
		"",
		"",
		"c controls • t tune • e enhancements • i info • q sign off",
	)

}

func TestRenderUsesTrackContextRelease(
	t *testing.T,
) {

	song := music.Song{
		Artist: "Robert Kraft",
		Title:  "Out With My Ex",
		Album:  "Fallback Album",
		Year:   1999,
	}

	trackContext := metadata.TrackContext{
		Release: metadata.ReleaseContext{
			Title: "Retro Active",
			Date:  "1982",
		},
	}

	view := renderTestTrack(
		song,
		trackContext,
	)

	if !strings.Contains(
		view,
		"Retro Active • 1982",
	) {

		t.Fatalf(
			"expected canonical release context in view",
		)

	}

	if strings.Contains(
		view,
		"Fallback Album",
	) {

		t.Fatal(
			"expected canonical release context to override song album",
		)

	}

	if strings.Contains(
		view,
		"1999",
	) {

		t.Fatal(
			"expected canonical release year to override song year",
		)

	}

}

func TestRenderFallsBackToSongRelease(
	t *testing.T,
) {

	song := music.Song{
		Artist: "Artist",
		Title:  "Track",
		Album:  "Local Album",
		Year:   1997,
	}

	view := renderTestTrack(
		song,
		metadata.TrackContext{},
	)

	if !strings.Contains(
		view,
		"Local Album • 1997",
	) {

		t.Fatalf(
			"expected song release fallback in view",
		)

	}

}

func TestRenderTrackContextDisplaysReleaseYear(
	t *testing.T,
) {

	song := music.Song{
		Artist: "Artist",
		Title:  "Track",
	}

	trackContext := metadata.TrackContext{
		Release: metadata.ReleaseContext{
			Title: "Release",
			Date:  "1982-04-12",
		},
	}

	view := renderTestTrack(
		song,
		trackContext,
	)

	if !strings.Contains(
		view,
		"Release • 1982",
	) {

		t.Fatalf(
			"expected release year in view",
		)

	}

	if strings.Contains(
		view,
		"1982-04-12",
	) {

		t.Fatal(
			"expected UI to display year rather than full release date",
		)

	}

}
