package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/music"
	"github.com/morsecodemedia/dj-morsecode/internal/player"
	"github.com/morsecodemedia/dj-morsecode/internal/radio"
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

func RenderStationHistory(
	history radio.History,
	width int,
) string {

	if width == 0 {
		width = 72
	}

	var s strings.Builder

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Header.Render("STATION HISTORY"),
		width,
	))

	s.WriteString("\n\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	if len(history.Tunes) == 0 {

		s.WriteString(
			Artist.Render("No stations tuned yet."),
		)

	} else {

		for i := len(history.Tunes) - 1; i >= 0; i-- {

			tune := history.Tunes[i]

			station, ok := radio.Find(
				tune.StationID,
			)

			name := tune.StationID

			if ok {
				name = station.Name
			}

			prefix := "  "

			if i == len(history.Tunes)-1 {
				prefix = "▶ "
			}

			s.WriteString(
				Artist.Render(prefix + name),
			)
			s.WriteString("\n")

		}

	}

	s.WriteString("\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Footer.Render("h or esc to return"),
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

func Render(
	song music.Song,
	width int,
	onAir bool,
	currentCue int,
	elapsed time.Duration,
	nowPlaying string,
	lyricsStatus string,
	stationName string,
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

	if song.Year != 0 {
		albumLine += fmt.Sprintf(" • %d", song.Year)
	}

	s.WriteString(Album.Render(albumLine))
	s.WriteString("\n")

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

	s.WriteString(
		Album.Render("LYRICS • " + lyricsStatus),
	)
	s.WriteString("\n\n")

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	if len(song.Timeline) == 0 {

		s.WriteString(Lyric.Render("No timeline loaded."))

	} else {

		viewport := BuildViewport(
			song.Timeline,
			currentCue,
		)

		preRoll := false

		if len(viewport) > 0 {
			preRoll = elapsed < viewport[0].Time
		}

		for i, cue := range viewport {

			if i == 0 {

				switch {

				case preRoll:

					CurrentLyric.Render(
						cueMarker(cue, true) + " " + cue.Text,
					)
					continue

				case cue.Type == music.CueLyric:

					s.WriteString(
						CurrentLyric.Render(cueMarker(cue, preRoll) + " " + cue.Text),
					)
					s.WriteString("\n")
					continue

				case cue.Type == music.CueBreak:

					s.WriteString(
						Cue.Render(cueMarker(cue, preRoll)),
					)
					s.WriteString("\n")
					continue

				}
			}

			switch cue.Type {

			case music.CueLyric:

				prefix := "  "

				if i > 0 {
					prefix = " " + upcomingMarker() + " "
				}

				s.WriteString(
					Lyric.Render(prefix + cue.Text),
				)

			case music.CueBreak:

				s.WriteString("")

			}

			s.WriteString("\n")

		}

	}

	s.WriteString("\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		// Footer.Render("Thanks for tuning in. • Press q to sign off."),
		Footer.Render("s stations • g genre • m vibes • h history • q sign off"),
		width,
	))

	return s.String()
}
