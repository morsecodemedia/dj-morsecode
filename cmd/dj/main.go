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
	Width      int
	Height     int
	Song       music.Song
	OnAir      bool
	CurrentCue int
	Elapsed    time.Duration
}

type tickMsg time.Time

const TickRate = time.Second

func tick() tea.Cmd {

	return tea.Tick(
		TickRate,
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
		m.Elapsed += TickRate

		for i := len(m.Song.Timeline) - 1; i >= 0; i-- {

			if m.Elapsed >= m.Song.Timeline[i].Time {

				m.CurrentCue = i

				break

			}

		}

		return m, tick()

	}

	return m, nil
}

func (m model) View() string {

	return ui.Render(
		m.Song,
		m.Width,
		m.OnAir,
		m.CurrentCue,
		m.Elapsed,
	)

}

func main() {

	lines, err := lyrics.Load("assets/interstate-love-song.lrc")
	if err != nil {

		fmt.Println(err)
		os.Exit(1)

	}

	metadata := lyrics.ParseMetadata(lines)

	timeline := lyrics.ParseTimeline(lines)

	song := music.Song{
		Title:    metadata.Title,
		Artist:   metadata.Artist,
		Album:    metadata.Album,
		Length:   metadata.Length,
		Timeline: timeline,
	}

	p := tea.NewProgram(model{
		Song: song,
	})

	if _, err := p.Run(); err != nil {

		fmt.Println(err)
		os.Exit(1)

	}

}
