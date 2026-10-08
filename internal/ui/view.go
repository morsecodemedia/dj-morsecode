package ui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/history"
	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
	"github.com/morsecodemedia/dj-morsecode/internal/observation"
	"github.com/morsecodemedia/dj-morsecode/internal/player"
	"github.com/morsecodemedia/dj-morsecode/internal/radio"
	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

// BuildViewport converts the complete song timeline into the
// subset of cues currently visible in the timeline panel.
//
// The renderer should never access Song.Timeline directly.
// From this point forward, it renders only what the viewport
// chooses to expose.
func BuildViewport(
	timeline []music.Cue,
	current int,
) []music.Cue {

	if len(timeline) == 0 {
		return nil
	}

	var viewport []music.Cue

	currentCue := timeline[current]

	viewport = append(viewport, currentCue)

	for i := current + 1; i < len(timeline); i++ {

		cue := timeline[i]

		if cue.Type == music.CueBreak {
			continue
		}

		viewport = append(viewport, cue)

		if len(viewport) == 3 {
			break
		}

	}

	return viewport

}

func cueMarker(
	cue music.Cue,
	preRoll bool,
) string {

	if preRoll {
		return "▶"
	}

	switch cue.Type {

	case music.CueLyric:
		return "♫"

	case music.CueBreak:
		return "●"

	default:
		return "?"
	}
}

func upcomingMarker() string {
	return "·"
}

func releaseYear(
	date string,
) int {

	date = strings.TrimSpace(
		date,
	)

	if len(date) < 4 {
		return 0
	}

	year, err := strconv.Atoi(
		date[:4],
	)
	if err != nil {
		return 0
	}

	return year

}

func RenderMoodPicker(
	families []radio.MoodFamily,
	selected int,
	width int,
) string {

	if width == 0 {
		width = 72
	}

	var s strings.Builder

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Header.Render("SELECT A MOOD"),
		width,
	))

	s.WriteString("\n")

	s.WriteString(Center(
		Subtitle.Render("Let DJ MorseCode set the mood."),
		width,
	))

	s.WriteString("\n\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	for i, family := range families {

		prefix := "  "

		if i == selected {
			prefix = "▶ "
		}

		line := prefix + family.Name

		if i == selected {

			s.WriteString(
				Title.Render(line),
			)

		} else {

			s.WriteString(
				Artist.Render(line),
			)

		}

		s.WriteString("\n")

	}

	s.WriteString("\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Footer.Render("↑/↓ or j/k to navigate • enter to choose • esc to cancel"),
		width,
	))

	return s.String()

}

func RenderGenrePicker(
	genres []string,
	selected int,
	width int,
) string {

	if width == 0 {
		width = 72
	}

	var s strings.Builder

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Header.Render("SELECT A GENRE"),
		width,
	))

	s.WriteString("\n")

	s.WriteString(Center(
		Subtitle.Render("Let DJ MorseCode pick the station."),
		width,
	))

	s.WriteString("\n\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	for i, genre := range genres {

		prefix := "  "

		if i == selected {
			prefix = "▶ "
		}

		line := prefix + genre

		if i == selected {

			s.WriteString(
				Title.Render(line),
			)

		} else {

			s.WriteString(
				Artist.Render(line),
			)

		}

		s.WriteString("\n")

	}

	s.WriteString("\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Footer.Render("↑/↓ or j/k to navigate • enter to choose • esc to cancel"),
		width,
	))

	return s.String()

}

func RenderListeningHistory(
	entries []history.Entry,
	width int,
) string {

	if width == 0 {
		width = 72
	}

	var s strings.Builder

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Header.Render("HISTORY"),
		width,
	))

	s.WriteString("\n\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	if len(entries) == 0 {

		s.WriteString(
			Artist.Render(
				"No listening history yet.",
			),
		)

	} else {

		for _, entry := range entries {

			label := ""

			switch entry.Kind {

			case history.KindStation:

				label = entry.StationID

				if station, ok :=
					radio.Find(
						entry.StationID,
					); ok {

					label = station.Name
				}

			case history.KindMedia:

				label =
					entry.Media.Ref.Name

				if label == "" &&
					entry.Media.Artist != "" &&
					entry.Media.Title != "" {

					label =
						entry.Media.Artist +
							" - " +
							entry.Media.Title
				}

				if label == "" {
					label =
						entry.Media.Title
				}

				if label == "" {
					label =
						entry.Media.Ref.ID
				}

			}

			if label == "" {
				continue
			}

			timestamp :=
				entry.PlayedAt.Format(
					"15:04:05",
				)

			s.WriteString(
				Artist.Render(
					timestamp +
						"  " +
						label,
				),
			)

			s.WriteString("\n")
		}

	}

	s.WriteString("\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Footer.Render(
			"h or esc to return",
		),
		width,
	))

	return s.String()
}

func RenderObservationHistory(
	stations []observation.StationObservation,
	playback []observation.PlaybackObservation,
	media []observation.MediaObservation,
	width int,
) string {

	if width == 0 {
		width = 72
	}

	type historyItem struct {
		observedAt time.Time
		label      string
		detail     string
	}

	items := make(
		[]historyItem,
		0,
		len(stations)+len(playback)+len(media),
	)

	for _, observed := range stations {

		label := strings.ToUpper(
			string(observed.Kind),
		)

		detail := observed.StationID

		if station, ok := radio.Find(
			observed.StationID,
		); ok {

			detail = station.Name
		}

		items = append(
			items,
			historyItem{
				observedAt: observed.ObservedAt,
				label:      label,
				detail:     detail,
			},
		)

	}

	for _, observed := range playback {

		label := strings.ToUpper(
			string(observed.Item.Type),
		)

		detail := observed.Item.DisplayTitle()

		if detail == "" {
			detail = observed.Item.RawTitle
		}

		items = append(
			items,
			historyItem{
				observedAt: observed.ObservedAt,
				label:      label,
				detail:     detail,
			},
		)

	}

	for _, observed := range media {

		label := strings.ToUpper(
			string(observed.Item.Kind),
		)

		detail := observed.Item.Ref.Name

		if detail == "" &&
			observed.Item.Artist != "" &&
			observed.Item.Title != "" {

			detail =
				observed.Item.Artist +
					" - " +
					observed.Item.Title

		}

		if detail == "" {
			detail = observed.Item.Title
		}

		if detail == "" {
			detail = observed.Item.Ref.ID
		}

		items = append(
			items,
			historyItem{
				observedAt: observed.ObservedAt,
				label:      label,
				detail:     detail,
			},
		)

	}

	sort.Slice(
		items,
		func(i int, j int) bool {

			return items[i].observedAt.Before(
				items[j].observedAt,
			)

		},
	)

	var s strings.Builder

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Header.Render("OBSERVATIONS"),
		width,
	))

	s.WriteString("\n")

	s.WriteString(Center(
		Subtitle.Render(
			"What DJ MorseCode has witnessed.",
		),
		width,
	))

	s.WriteString("\n\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	if len(items) == 0 {

		s.WriteString(
			Artist.Render(
				"No observations yet.",
			),
		)

	} else {

		for _, item := range items {

			timestamp := item.observedAt.Format(
				"15:04:05",
			)

			s.WriteString(
				Album.Render(
					timestamp +
						"  " +
						item.label,
				),
			)

			s.WriteString("\n")

			if item.detail != "" {

				s.WriteString(
					Artist.Render(
						"          " +
							item.detail,
					),
				)

				s.WriteString("\n")

			}

			s.WriteString("\n")

		}

	}

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	count := len(items)

	s.WriteString(Center(
		Footer.Render(
			fmt.Sprintf(
				"%d observations • o or esc to return",
				count,
			),
		),
		width,
	))

	return s.String()
}

func RenderVibePicker(
	presets []radio.Preset,
	selected int,
	width int,
) string {

	if width == 0 {
		width = 72
	}

	var s strings.Builder

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Header.Render("SELECT A VIBE"),
		width,
	))

	s.WriteString("\n")

	s.WriteString(Center(
		Subtitle.Render("Let DJ MorseCode pick the station."),
		width,
	))

	s.WriteString("\n\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	for i, preset := range presets {

		prefix := "  "

		if i == selected {
			prefix = "▶ "
		}

		line := prefix + preset.Name

		if i == selected {
			s.WriteString(
				Title.Render(line),
			)
		} else {
			s.WriteString(
				Artist.Render(line),
			)
		}

		s.WriteString("\n")

	}

	s.WriteString("\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Footer.Render("↑/↓ or j/k to navigate • enter to choose • esc to cancel"),
		width,
	))

	return s.String()

}

func RenderStationPicker(
	stations []radio.Station,
	selected int,
	width int,
) string {

	if width == 0 {
		width = 72
	}

	var s strings.Builder

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Header.Render("SELECT A STATION"),
		width,
	))

	s.WriteString("\n")

	s.WriteString(Center(
		Subtitle.Render("Tune in to something good."),
		width,
	))

	s.WriteString("\n\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	for i, station := range stations {

		prefix := "  "

		if i == selected {
			prefix = "▶ "
		}

		line := prefix + station.Name

		if i == selected {
			s.WriteString(Title.Render(line))
		} else {
			s.WriteString(Artist.Render(line))
		}

		s.WriteString("\n")

	}

	s.WriteString("\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Footer.Render("↑/↓ or j/k to navigate • enter to tune • esc to cancel"),
		width,
	))

	return s.String()

}

func RenderYouTubePicker(
	items []source.MediaItem,
	selected int,
	width int,
) string {
	if width == 0 {
		width = 72
	}

	var s strings.Builder

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Header.Render("SELECT FROM YOUTUBE"),
		width,
	))

	s.WriteString("\n")

	s.WriteString(Center(
		Subtitle.Render("Pick something to play."),
		width,
	))

	s.WriteString("\n\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	for i, item := range items {

		prefix := "  "

		if i == selected {
			prefix = "▶ "
		}

		line := prefix + item.Title

		if i == selected {
			s.WriteString(Title.Render(line))
		} else {
			s.WriteString(Artist.Render(line))
		}

		s.WriteString("\n")

	}

	s.WriteString("\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Footer.Render("↑/↓ or j/k to navigate • enter to tune • esc to cancel"),
		width,
	))

	return s.String()
}

func RenderIdle(
	width int,
	footer string,
) string {

	if width == 0 {
		width = 72
	}

	var s strings.Builder

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Header.Render("DJ MORSECODE"),
		width,
	))

	s.WriteString("\n")

	s.WriteString(Center(
		Subtitle.Render("The DJ that quietly codes with you."),
		width,
	))

	s.WriteString("\n\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Section.Render("○ STANDING BY"),
		width,
	))

	s.WriteString("\n\n")

	s.WriteString(Center(
		Title.Render("What are we listening to?"),
		width,
	))

	s.WriteString("\n\n")

	s.WriteString(Center(
		Artist.Render("v  Pick a vibe"),
		width,
	))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Artist.Render("m  Set a mood"),
		width,
	))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Artist.Render("g  Choose a genre"),
		width,
	))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Artist.Render("s  Tune a station"),
		width,
	))

	s.WriteString("\n\n")

	s.WriteString(Center(
		Subtitle.Render("DJ MorseCode is ready when you are."),
		width,
	))

	s.WriteString("\n\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Footer.Render(
			footer,
		),
		width,
	))

	return s.String()

}

func Render(
	song music.Song,
	trackContext metadata.TrackContext,
	width int,
	onAir bool,
	currentCue int,
	elapsed time.Duration,
	nowPlaying string,
	lyricsStatus string,
	lyricsVisible bool,
	stationName string,
	intentType string,
	intentName string,
	footer string,
) string {

	if width == 0 {
		width = 72
	}

	var s strings.Builder

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Header.Render("DJ MORSECODE"),
		width,
	))

	s.WriteString("\n")

	s.WriteString(Center(
		Subtitle.Render("The DJ that quietly codes with you."),
		width,
	))

	s.WriteString("\n\n")

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	status := "○ ON AIR"
	if onAir {
		status = "● ON AIR"
	}

	s.WriteString(Section.Render(status))
	s.WriteString("\n")

	if stationName != "" {

		s.WriteString(
			Album.Render("STATION • " + stationName),
		)
		s.WriteString("\n")

	}

	if intentName != "" {

		label := strings.ToUpper(
			intentType,
		)

		s.WriteString(
			Album.Render(
				label + " • " + intentName,
			),
		)
		s.WriteString("\n")

	}

	s.WriteString("\n")

	title := nowPlaying

	if title == "" {
		title = song.Title
	}

	s.WriteString(Title.Render("♫ " + title))
	s.WriteString("\n\n")

	s.WriteString(Artist.Render(song.Artist))
	s.WriteString("\n\n")

	albumLine := song.Album

	if trackContext.Release.Title != "" {
		albumLine = trackContext.Release.Title
	}

	year := song.Year

	if trackContext.Release.Date != "" {
		year = releaseYear(
			trackContext.Release.Date,
		)
	}

	if year != 0 {

		if albumLine != "" {
			albumLine += " • "
		}

		albumLine += fmt.Sprintf(
			"%d",
			year,
		)

	}

	if albumLine != "" {

		s.WriteString(
			Album.Render(albumLine),
		)

		s.WriteString("\n")

	}

	transport :=
		player.FormatDuration(elapsed) +
			" / " +
			player.FormatDuration(song.Duration)

	progress := player.Progress(
		elapsed.Seconds(),
		song.Duration.Seconds(),
	)

	bar := player.ProgressBar(
		24,
		progress,
	)

	s.WriteString(
		Artist.Render("▶ " + transport),
	)
	s.WriteString("\n")
	s.WriteString(
		Album.Render(bar),
	)
	s.WriteString("\n\n")

	if lyricsVisible {

		s.WriteString(
			Album.Render(
				"LYRICS • " + lyricsStatus,
			),
		)
		s.WriteString("\n\n")

		s.WriteString(Divider(width))
		s.WriteString("\n\n")

		if strings.TrimSpace(
			song.Lyrics,
		) == "" {

			s.WriteString(
				Lyric.Render(
					"No lyrics available.",
				),
			)

		} else {

			lines := strings.Split(
				song.Lyrics,
				"\n",
			)

			for _, line := range lines {

				if strings.TrimSpace(line) == "" {
					s.WriteString("\n")
					continue
				}

				s.WriteString(
					Lyric.Render(line),
				)

				s.WriteString("\n")

			}

		}

		s.WriteString("\n")
		s.WriteString(Divider(width))
		s.WriteString("\n\n")

	}

	s.WriteString("\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Footer.Render(
			footer,
		),
		width,
	))

	return s.String()
}
