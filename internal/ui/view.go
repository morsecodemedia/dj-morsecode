package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/music"
)

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
	s.WriteString("\n\n")

	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	if len(song.Timeline) == 0 {

		s.WriteString(Lyric.Render("No timeline loaded."))

	} else {

		preRoll := elapsed < song.Timeline[0].Time

		start := currentCue - 2
		if start < 0 {
			start = 0
		}

		end := currentCue + 3
		if end > len(song.Timeline) {
			end = len(song.Timeline)
		}

		topPadding := 0
		if currentCue < 2 {
			topPadding = 2 - currentCue
		}

		bottomPadding := 0
		remaining := len(song.Timeline) - currentCue - 1
		if remaining < 2 {
			bottomPadding = 2 - remaining
		}

		for i := 0; i < topPadding; i++ {
			s.WriteString("\n\n")
		}

		for i := start; i < end; i++ {

			cue := song.Timeline[i]

			if i == currentCue {

				if preRoll {

					s.WriteString(Cue.Render("▶"))
					s.WriteString("\n\n")

				} else {

					switch cue.Type {

					case music.CueLyric:

						s.WriteString(
							CurrentLyric.Render("♫ " + cue.Text),
						)

					case music.CueBreak:

						s.WriteString(
							Cue.Render("●"),
						)

					}

					s.WriteString("\n\n")
					continue
				}
			}

			switch cue.Type {

			case music.CueLyric:

				s.WriteString(
					Lyric.Render("      " + cue.Text),
				)

			case music.CueBreak:

				s.WriteString("")

			}

			s.WriteString("\n\n")

		}

		for i := 0; i < bottomPadding; i++ {
			s.WriteString("\n\n")
		}
	}

	s.WriteString("\n\n")
	s.WriteString(Divider(width))
	s.WriteString("\n\n")

	s.WriteString(Center(
		Footer.Render("Thanks for tuning in. • Press q to sign off."),
		width,
	))

	return s.String()
}
