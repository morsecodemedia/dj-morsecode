package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/morsecodemedia/dj-morsecode/internal/library"
	"github.com/morsecodemedia/dj-morsecode/internal/lrclib"
	"github.com/morsecodemedia/dj-morsecode/internal/lyrics"
	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
	"github.com/morsecodemedia/dj-morsecode/internal/player"
	"github.com/morsecodemedia/dj-morsecode/internal/ui"
)

type lyricsState int

const (
	lyricsUnavailable lyricsState = iota
	lyricsLocal
	lyricsSearching
	lyricsRemote
)

func (s lyricsState) String() string {

	switch s {

	case lyricsLocal:
		return "LOCAL"

	case lyricsSearching:
		return "SEARCHING"

	case lyricsRemote:
		return "LRCLIB"

	default:
		return "UNAVAILABLE"

	}

}

type model struct {
	Width       int
	Height      int
	Song        music.Song
	OnAir       bool
	CurrentCue  int
	Player      *player.Player
	NowPlaying  string
	LastTitle   string
	LyricsState lyricsState
}

type tickMsg time.Time

type lrclibSongMsg struct {
	RawTitle string
	Song     music.Song
	Err      error
}

const TickRate = time.Second

func tick() tea.Cmd {

	return tea.Tick(
		TickRate,
		func(t time.Time) tea.Msg {
			return tickMsg(t)
		},
	)

}

func loadLRCLIBSong(
	rawTitle string,
	artist string,
	title string,
	duration time.Duration,
) tea.Cmd {

	return func() tea.Msg {

		client := lrclib.NewClient()

		results, err := client.Search(
			artist,
			title,
		)
		if err != nil {
			return lrclibSongMsg{
				RawTitle: rawTitle,
				Err:      err,
			}
		}

		result, ok := lrclib.BestMatch(
			results,
			duration,
		)
		if !ok {
			return lrclibSongMsg{
				RawTitle: rawTitle,
				Err: fmt.Errorf(
					"no LRCLIB match for %s",
					rawTitle,
				),
			}
		}

		return lrclibSongMsg{
			RawTitle: rawTitle,
			Song:     lrclib.Song(result),
		}

	}

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
		duration := m.Player.Duration()
		rawTitle := m.Player.Title()
		artist := m.Player.Artist()
		title := m.Player.TrackTitle()
		album := m.Player.Album()
		track := metadata.Resolve(rawTitle)

		if artist != "" {
			track.Artist = artist
		}

		if title != "" {
			track.Title = title
		}

		if track.RawTitle == "" {
			return m, tick()
		}

		if track.Valid {

			m.NowPlaying = track.RawTitle

			if track.RawTitle != m.LastTitle {

				m.Song = music.Song{
					Title:    track.Title,
					Artist:   track.Artist,
					Album:    album,
					Duration: duration,
				}

				m.CurrentCue = 0

				song, ok := library.Load(
					track.Artist,
					track.Title,
				)

				if ok {

					m.Song = song
					m.LyricsState = lyricsLocal
					m.LastTitle = track.RawTitle

					return m, tick()

				}

				m.LastTitle = track.RawTitle
				m.LyricsState = lyricsSearching

				return m, tea.Batch(
					tick(),
					loadLRCLIBSong(
						track.RawTitle,
						track.Artist,
						track.Title,
						duration,
					),
				)

			}

		}

		m.CurrentCue = player.CurrentCue(
			m.Song.Timeline,
			position,
		)

		m.OnAir = !m.OnAir
		return m, tick()

	case lrclibSongMsg:

		if msg.RawTitle != m.LastTitle {
			return m, nil
		}

		if msg.Err != nil {
			m.LyricsState = lyricsUnavailable
			return m, nil
		}

		m.Song.Timeline = msg.Song.Timeline
		m.LyricsState = lyricsRemote
		m.CurrentCue = player.CurrentCue(
			m.Song.Timeline,
			m.Player.Position(),
		)

		return m, nil

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
		m.LyricsState.String(),
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
