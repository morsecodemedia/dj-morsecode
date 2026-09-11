package ui

import "strings"

func Render() string {

	const width = 72

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

	s.WriteString(Title.Render("♫ INTERSTATE LOVE SONG"))
	s.WriteString("\n\n")

	s.WriteString(Artist.Render("Stone Temple Pilots"))
	s.WriteString("\n\n")

	s.WriteString(Album.Render("Purple • 1994"))

	s.WriteString("\n\n")

	s.WriteString(Divider(width))

	s.WriteString("\n\n")

	s.WriteString(Lyric.Render("      Leaving on a southern train..."))
	s.WriteString("\n\n")

	s.WriteString(CurrentLyric.Render("▶ Only yesterday you lied..."))
	s.WriteString("\n\n")

	s.WriteString(Lyric.Render("      Promises of what I seemed..."))

	s.WriteString("\n\n")

	s.WriteString(Divider(width))

	s.WriteString("\n\n")

	s.WriteString(Center(
		Footer.Render("Thanks for tuning in. • Press q to sign off."),
		width,
	))

	return s.String()

}
