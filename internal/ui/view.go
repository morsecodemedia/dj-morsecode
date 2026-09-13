package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/music"
	"github.com/morsecodemedia/dj-morsecode/internal/player"
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

func Render(
	song music.Song,
	width int,
	onAir bool,
	currentCue int,
	elapsed time.Duration,
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
	s.WriteString("\n\n")

	s.WriteString(Title.Render("♫ " + song.Title))
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
			song.Length

	s.WriteString(
		Artist.Render("▶ " + transport),
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
		Footer.Render("Thanks for tuning in. • Press q to sign off."),
		width,
	))

	return s.String()
}
