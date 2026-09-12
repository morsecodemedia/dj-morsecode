package ui

import (
	"fmt"

	"strings"

	"github.com/morsecodemedia/dj-morsecode/internal/music"
)

func Render(song music.Song, width int, onAir bool, currentLine int) string {

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

	s.WriteString(
		Album.Render(albumLine),
	)
	s.WriteString("\n\n")

	s.WriteString(Divider(width))

	s.WriteString("\n\n")

	start := currentLine - 2

	if start < 0 {
		start = 0
	}

	end := currentLine + 3

	if end > len(song.Lyrics) {
		end = len(song.Lyrics)
	}

	topPadding := 0

	if currentLine < 2 {
		topPadding = 2 - currentLine
	}

	bottomPadding := 0

	remaining := len(song.Lyrics) - currentLine - 1

	if remaining < 2 {
		bottomPadding = 2 - remaining
	}

	for i := 0; i < topPadding; i++ {

		s.WriteString("\n\n")

	}

	for i := start; i < end; i++ {

		lyric := song.Lyrics[i]

		if i == currentLine {

			s.WriteString(
				CurrentLyric.Render(
					"▶ " + lyric.Text,
				),
			)

		} else {

			s.WriteString(
				Lyric.Render(
					"      " + lyric.Text,
				),
			)

		}

		s.WriteString("\n\n")

	}

	for i := 0; i < bottomPadding; i++ {

		s.WriteString("\n\n")

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
