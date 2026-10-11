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

func TestRenderTransportState(
	t *testing.T,
) {

	tests := []struct {
		name string

		paused bool
		muted  bool
		volume float64

		wantMarker string
		wantVolume string
	}{
		{
			name:       "playing",
			volume:     60,
			wantMarker: "▶",
			wantVolume: "VOL 60%",
		},
		{
			name:       "paused",
			paused:     true,
			volume:     60,
			wantMarker: "Ⅱ",
			wantVolume: "VOL 60%",
		},
		{
			name:       "muted",
			muted:      true,
			volume:     60,
			wantMarker: "▶",
			wantVolume: "MUTED",
		},
	}

	for _, test := range tests {

		t.Run(
			test.name,
			func(t *testing.T) {

				view := Render(
					music.Song{
						Title:  "Messy",
						Artist: "Lola Young",
						Duration: 4*time.Minute +
							44*time.Second,
					},
					metadata.TrackContext{},
					72,
					true,
					0,
					21*time.Second,
					"Messy",
					test.paused,
					test.muted,
					test.volume,
					"SYNCED",
					false,
					"Z100",
					"",
					"",
					"",
				)

				if !strings.Contains(
					view,
					test.wantMarker,
				) {

					t.Errorf(
						"expected transport marker %q",
						test.wantMarker,
					)

				}

				if !strings.Contains(
					view,
					test.wantVolume,
				) {

					t.Errorf(
						"expected volume state %q",
						test.wantVolume,
					)

				}

			},
		)

	}

}

func renderTestTrack(
	song music.Song,
	trackContext metadata.TrackContext,
) string {

	return Render(
		song,
		trackContext,
		72,
		false,
		0,
		0,
		"",    // nowPlaying
		false, // paused
		false, // muted
		60,    // volume
		"Unavailable",
		false,
		"",
		"",
		"",
		"",
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

func TestRenderHidesLyricsBody(
	t *testing.T,
) {

	view := Render(
		music.Song{
			Lyrics: "Secret lyric",
		},
		metadata.TrackContext{},

		72,
		false,
		0,
		0,
		"",    // nowPlaying
		false, // paused
		false, // muted
		60,    // volume
		"Unavailable",
		false,
		"",
		"",
		"",
		"",
	)

	if strings.Contains(
		view,
		"Secret lyric",
	) {

		t.Fatal(
			"expected lyrics body to be hidden",
		)
	}

	if !strings.Contains(
		view,
		"LYRICS • AVAILABLE",
	) {

		t.Fatal(
			"expected lyrics availability to remain visible",
		)
	}

}

func TestRenderShowsSyncedLyricsViewport(
	t *testing.T,
) {

	view := Render(
		music.Song{
			Lyrics: strings.Join(
				[]string{
					"First lyric",
					"Current lyric",
					"Next lyric",
					"Another lyric",
					"Distant lyric",
				},
				"\n",
			),

			Timeline: []music.Cue{
				{
					Time: 0,
					Type: music.CueLyric,
					Text: "First lyric",
				},
				{
					Time: 10 * time.Second,
					Type: music.CueLyric,
					Text: "Current lyric",
				},
				{
					Time: 20 * time.Second,
					Type: music.CueLyric,
					Text: "Next lyric",
				},
				{
					Time: 30 * time.Second,
					Type: music.CueLyric,
					Text: "Another lyric",
				},
				{
					Time: 40 * time.Second,
					Type: music.CueLyric,
					Text: "Distant lyric",
				},
			},
		},
		metadata.TrackContext{},

		72,
		true,
		1,
		10*time.Second,
		"",
		false,
		false,
		60,
		"Remote",
		true,
		"",
		"",
		"",
		"",
	)

	if !strings.Contains(
		view,
		"LYRICS • SYNCED",
	) {

		t.Fatal(
			"expected synced lyrics status",
		)
	}

	if !strings.Contains(
		view,
		"♫ Current lyric",
	) {

		t.Fatal(
			"expected current lyric",
		)
	}

	if !strings.Contains(
		view,
		"· Next lyric",
	) {

		t.Fatal(
			"expected next lyric",
		)
	}

	if !strings.Contains(
		view,
		"· Another lyric",
	) {

		t.Fatal(
			"expected second upcoming lyric",
		)
	}

	if strings.Contains(
		view,
		"First lyric",
	) {

		t.Fatal(
			"expected previous lyric to be outside viewport",
		)
	}

	if strings.Contains(
		view,
		"Distant lyric",
	) {

		t.Fatal(
			"expected distant lyric to be outside viewport",
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

func TestRenderTransportPlaying(
	t *testing.T,
) {

	view := renderTransport(
		21*time.Second,
		4*time.Minute+44*time.Second,
		false,
		false,
		60,
		24,
	)

	if !strings.Contains(
		view,
		"▶ 00:21",
	) {

		t.Fatal(
			"expected playing transport",
		)
	}

	if !strings.Contains(
		view,
		"04:44",
	) {

		t.Fatal(
			"expected duration transport",
		)
	}

	if !strings.Contains(
		view,
		"VOL 60%",
	) {

		t.Fatal(
			"expected volume",
		)
	}

}

func TestRenderTransportPausedAndMuted(
	t *testing.T,
) {

	view := renderTransport(
		21*time.Second,
		4*time.Minute+44*time.Second,
		true,
		true,
		60,
		24,
	)

	if !strings.Contains(
		view,
		"Ⅱ 00:21",
	) {

		t.Fatal(
			"expected paused transport",
		)
	}

	if !strings.Contains(
		view,
		"04:44",
	) {

		t.Fatal(
			"expected duration transport",
		)
	}

	if !strings.Contains(
		view,
		"MUTED",
	) {

		t.Fatal(
			"expected muted transport",
		)
	}

}

func TestRenderTransportFallsBackWhenNarrow(
	t *testing.T,
) {

	view := renderTransport(
		21*time.Second,
		4*time.Minute+44*time.Second,
		false,
		false,
		60,
		24,
	)

	if !strings.Contains(
		view,
		"▶ 00:21 / 04:44",
	) {

		t.Fatal(
			"expected compact transport",
		)
	}

}

func TestRenderNowPlaying(
	t *testing.T,
) {

	view := renderNowPlaying(
		music.Song{
			Title:    "Loser",
			Artist:   "Tame Impala",
			Album:    "Currents",
			Year:     2015,
			Duration: 2*time.Minute + 11*time.Second,
		},
		metadata.TrackContext{},
		true,
		48*time.Second,
		"Tame Impala - Loser",
		false,
		false,
		60,
		"SYNCED",
		"Z100",
		"",
		"",
		120,
	)

	for _, expected := range []string{
		"● ON AIR",
		"STATION • Z100",
		"♫ Loser",
		"Tame Impala",
		"Currents • 2015",
		"▶ 00:48",
		"02:11",
		"VOL 60%",
		"LYRICS • SYNCED",
	} {

		if !strings.Contains(
			view,
			expected,
		) {

			t.Errorf(
				"expected now playing to contain %q",
				expected,
			)
		}

	}

	if strings.Contains(
		view,
		"Tame Impala - Loser",
	) {

		t.Fatal(
			"expected normalized title instead of raw now playing",
		)
	}

}

func renderLyricsLayout(
	width int,
) string {

	return Render(
		music.Song{
			Title:    "Loser",
			Artist:   "Tame Impala",
			Duration: 2*time.Minute + 3*time.Second,

			Timeline: []music.Cue{
				{
					Time: 10 * time.Second,
					Type: music.CueLyric,
					Text: "Current lyric",
				},
				{
					Time: 20 * time.Second,
					Type: music.CueLyric,
					Text: "Next lyric",
				},
			},
		},
		metadata.TrackContext{},
		width,
		true,
		0,
		10*time.Second,
		"Tame Impala - Loser",
		false,
		false,
		60,
		"Remote",
		true,
		"Z100",
		"",
		"",
		"",
	)
}

func TestRenderWideLyricsBesideNowPlaying(
	t *testing.T,
) {

	view :=
		renderLyricsLayout(
			120,
		)

	lines :=
		strings.Split(
			view,
			"\n",
		)

	found := false

	for _, line := range lines {

		if strings.Contains(
			line,
			"● ON AIR",
		) &&
			strings.Contains(
				line,
				"LYRICS",
			) {

			found = true
			break
		}

	}

	if !found {

		t.Fatal(
			"expected lyrics beside now playing at wide width",
		)
	}

}

func TestRenderNarrowLyricsBelowNowPlaying(
	t *testing.T,
) {

	view :=
		renderLyricsLayout(
			72,
		)

	lines :=
		strings.Split(
			view,
			"\n",
		)

	onAirLine := -1
	lyricsLine := -1

	for i, line := range lines {

		if onAirLine < 0 &&
			strings.Contains(
				line,
				"● ON AIR",
			) {

			onAirLine = i
			continue
		}

		if onAirLine >= 0 &&
			strings.Contains(
				line,
				"SYNCED • KARAOKE",
			) {

			lyricsLine = i
			break
		}

	}

	if onAirLine < 0 ||
		lyricsLine < 0 ||
		lyricsLine <= onAirLine {

		t.Fatal(
			"expected lyrics below now playing at narrow width",
		)
	}

}

func TestRenderFooterUsesDesignedCompactLayout(
	t *testing.T,
) {

	view := renderFooter(
		"c controls • t tune • e enhancements • i info • q sign off",
		48,
	)

	lines := strings.Split(
		view,
		"\n",
	)

	if len(lines) != 3 {

		t.Fatalf(
			"expected three footer rows, got %d",
			len(lines),
		)
	}

	for _, expected := range []string{
		"c controls • t tune",
		"e enhancements • i info",
		"q sign off",
	} {

		if !strings.Contains(
			view,
			expected,
		) {

			t.Errorf(
				"expected compact footer to contain %q",
				expected,
			)
		}
	}
}

func TestRenderFooterPreservesContextualFooter(
	t *testing.T,
) {

	footer :=
		"CONTROLS • b previous • n next • p resume • m mute • esc cancel"

	view := renderFooter(
		footer,
		120,
	)

	if !strings.Contains(
		view,
		footer,
	) {

		t.Fatal(
			"expected contextual footer to remain unchanged",
		)
	}

	if strings.Contains(
		view,
		"e enhancements • i info",
	) {

		t.Fatal(
			"expected contextual footer not to use root footer layout",
		)
	}

}
