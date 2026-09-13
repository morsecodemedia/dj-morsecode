package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212"))

	Section = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Width(72).
		Align(lipgloss.Center)

	Subtitle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("241"))

	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230"))

	Artist = lipgloss.NewStyle().
		Foreground(lipgloss.Color("250"))

	Album = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	Lyric = lipgloss.NewStyle().
		Foreground(lipgloss.Color("250"))

	CurrentLyric = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229"))

	UpcomingLyric = lipgloss.NewStyle().
			Foreground(lipgloss.Color("242"))

	Cue = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("229"))

	Footer = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
)

func Divider(width int) string {

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("238")).
		Render(strings.Repeat("─", width))

}

func Center(text string, width int) string {

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(text)

}
