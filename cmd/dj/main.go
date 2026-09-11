package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/morsecodemedia/dj-morsecode/internal/music"
	"github.com/morsecodemedia/dj-morsecode/internal/ui"
)

type model struct {
	Width  int
	Height int
	Song   music.Song
}

func (m model) Init() tea.Cmd {
	return nil
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

	}

	return m, nil
}

func (m model) View() string {

	return ui.Render(m.Song, m.Width)
}

func main() {

	p := tea.NewProgram(model{
		Song: music.DemoSong,
	})

	if _, err := p.Run(); err != nil {

		fmt.Println(err)

		os.Exit(1)

	}

}
