package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/history"
	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
	"github.com/morsecodemedia/dj-morsecode/internal/source"
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
		true,
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

func TestRenderHidesLyrics(
	t *testing.T,
) {

	view := Render(
		music.Song{
			Lyrics: "Secret lyric",
		},
		metadata.TrackContext{},
		72,
		true,
		0,
		0,
		"Artist - Track",
		"LRCLIB",
		false,
		"Station",
		"",
		"",
		"footer",
	)

	if strings.Contains(
		view,
		"Secret lyric",
	) {

		t.Fatal(
			"expected lyrics to be hidden",
		)
	}

	if strings.Contains(
		view,
		"LYRICS",
	) {

		t.Fatal(
			"expected lyrics section to be hidden",
		)
	}

}

func TestRenderListeningHistoryEmpty(
	t *testing.T,
) {

	view := RenderListeningHistory(
		nil,
		72,
	)

	if !strings.Contains(
		view,
		"HISTORY",
	) {

		t.Fatal(
			"expected history heading",
		)
	}

	if !strings.Contains(
		view,
		"No listening history yet.",
	) {

		t.Fatal(
			"expected empty history message",
		)
	}

}

func TestRenderListeningHistoryStation(
	t *testing.T,
) {

	view := RenderListeningHistory(
		[]history.Entry{
			{
				Kind: history.KindStation,

				StationID: "hardrockradiofm",

				PlayedAt: time.Date(
					2026,
					time.October,
					7,
					22,
					56,
					17,
					0,
					time.UTC,
				),
			},
		},
		72,
	)

	if !strings.Contains(
		view,
		"Hard Rock Radio FM",
	) {

		t.Fatal(
			"expected station name",
		)
	}

	if !strings.Contains(
		view,
		"22:56:17",
	) {

		t.Fatal(
			"expected timestamp",
		)
	}

}

func TestRenderListeningHistoryMedia(
	t *testing.T,
) {

	view := RenderListeningHistory(
		[]history.Entry{
			{
				Kind: history.KindMedia,

				Media: source.MediaItem{
					Kind: source.MediaVideo,

					Ref: source.ItemRef{
						Source: source.Source{
							Kind: source.KindYouTube,
						},

						ID: "qORYO0atB6g",

						Name: "Beastie Boys - Intergalactic",
					},
				},

				PlayedAt: time.Date(
					2026,
					time.October,
					7,
					22,
					55,
					20,
					0,
					time.UTC,
				),
			},
		},
		72,
	)

	if !strings.Contains(
		view,
		"Beastie Boys - Intergalactic",
	) {

		t.Fatal(
			"expected media name",
		)
	}

}
