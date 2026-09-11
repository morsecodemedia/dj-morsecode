package ui

import (
	"fmt"

	"strings"

	"github.com/morsecodemedia/dj-morsecode/internal/music"
)

func Render(song music.Song, width int) string {

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

	s.WriteString(Section.Render("● ON AIR"))

	s.WriteString("\n\n")

	// s.WriteString(Title.Render("♫ INTERSTATE LOVE SONG"))
	s.WriteString(Title.Render("♫ " + song.Title))
	s.WriteString("\n\n")

	s.WriteString(Artist.Render(song.Artist))
	s.WriteString("\n\n")

	s.WriteString(
		Album.Render(
			fmt.Sprintf("%s • %d", song.Album, song.Year),
		),
	)
	s.WriteString("\n\n")

	s.WriteString(Divider(width))

	s.WriteString("\n\n")

	s.WriteString(Lyric.Render("      " + song.Lyrics[0].Text))
	s.WriteString("\n\n")

	s.WriteString(CurrentLyric.Render("▶ " + song.Lyrics[1].Text))
	s.WriteString("\n\n")

	s.WriteString(Lyric.Render("      " + song.Lyrics[2].Text))

	s.WriteString("\n\n")

	s.WriteString(Divider(width))

	s.WriteString("\n\n")

	s.WriteString(Center(
		Footer.Render("Thanks for tuning in. • Press q to sign off."),
		width,
	))

	return s.String()

}
