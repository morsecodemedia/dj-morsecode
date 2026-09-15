package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/morsecodemedia/dj-morsecode/internal/library"
	"github.com/morsecodemedia/dj-morsecode/internal/lyrics"
	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
	"github.com/morsecodemedia/dj-morsecode/internal/player"
	"github.com/morsecodemedia/dj-morsecode/internal/ui"
)

type model struct {
	Width      int
	Height     int
	Song       music.Song
	OnAir      bool
	CurrentCue int
	Player     *player.Player
	NowPlaying string
	LastTitle  string
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

		position := m.Player.Position()

		track := metadata.Resolve(
			m.Player.Title(),
		)
		if track.RawTitle == "" {
			return m, tick()
		}

		if track.Valid {

			m.NowPlaying = track.RawTitle

			if track.RawTitle != m.LastTitle {

				fmt.Printf("Loaded: %s\n", track.RawTitle)
				song, ok := library.Load(track.RawTitle)

				if ok {

					m.Song = song
					m.CurrentCue = 0

				}

				m.LastTitle = track.RawTitle

			}

		}

		m.CurrentCue = player.CurrentCue(
			m.Song.Timeline,
			position,
		)

		m.OnAir = !m.OnAir
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
		m.Player.Position(),
		m.NowPlaying,
	)

}

func main() {

	song, err := lyrics.LoadSong(
		"assets/interstate-love-song.lrc",
	)
	if err != nil {

		fmt.Println(err)
		os.Exit(1)

	}

	playback, err := player.New("/tmp/dj-morsecode.sock")
	if err != nil {
		panic(err)
	}

	p := tea.NewProgram(model{
		Song:   song,
		Player: playback,
	})

	if _, err := p.Run(); err != nil {

		fmt.Println(err)
		os.Exit(1)

	}

}
