package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/morsecodemedia/dj-morsecode/internal/lyrics"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
	"github.com/morsecodemedia/dj-morsecode/internal/ui"
)

type model struct {
	Width       int
	Height      int
	Song        music.Song
	OnAir       bool
	CurrentLine int
}

type tickMsg time.Time

func tick() tea.Cmd {

	return tea.Tick(
		time.Second,
		func(t time.Time) tea.Msg {
			return tickMsg(t)
		},
	)

}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		}

	case tea.WindowSizeMsg:

		m.Width = msg.Width
		m.Height = msg.Height

		return m, nil

	case tickMsg:

		m.OnAir = !m.OnAir
		m.CurrentLine++

		if m.CurrentLine >= len(m.Song.Lyrics) {
			m.CurrentLine = 0
		}

		return m, tick()

	}

	return m, nil
}

func (m model) View() string {

	return ui.Render(m.Song, m.Width, m.OnAir, m.CurrentLine)
}

func main() {
	lines, err := lyrics.Load("assets/interstate-love-song.lrc")

	metadata := lyrics.ParseMetadata(lines)

	fmt.Println()

	fmt.Println("Metadata")

	fmt.Println("--------")

	for key, value := range metadata {

		fmt.Printf(
			"%s = %s\n",
			key,
			value,
		)

	}

	if err != nil {

		fmt.Println(err)

		os.Exit(1)

	}

	fmt.Printf(
		"Loaded %d lines\n",
		len(lines),
	)

	fmt.Println()

	for i, line := range lines {

		switch {

		case lyrics.IsMetadata(line):

			fmt.Printf("%02d | META  | %s\n", i, line)

		case lyrics.IsLyric(line):

			fmt.Printf("%02d | LYRIC | %s\n", i, line)

		case lyrics.IsLyricBreak(line):

			fmt.Printf("%02d | BREAK | %s\n", i, line)

		case lyrics.IsBlank(line):

			fmt.Printf("%02d | BLANK | %s\n", i, line)

		default:

			fmt.Printf("%02d | UNKNOWN | %s\n", i, line)

		}

	}

	duration, err := lyrics.ParseTimestamp("00:35.18")
	if err != nil {
		panic(err)
	}
	fmt.Println(duration)

	song := lyrics.ParseSong(lines)
	fmt.Println()

	fmt.Println("Lyrics")

	fmt.Println("------")

	for i, lyric := range song.Lyrics {

		fmt.Printf(
			"%02d | %s\n",
			i,
			lyric.Text,
		)

	}

	p := tea.NewProgram(model{
		Song:        song,
		CurrentLine: 1,
	})

	if _, err := p.Run(); err != nil {

		fmt.Println(err)

		os.Exit(1)

	}

}
