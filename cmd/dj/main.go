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
	"github.com/morsecodemedia/dj-morsecode/internal/radio"
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
	Width              int
	Height             int
	Song               music.Song
	OnAir              bool
	CurrentCue         int
	Player             *player.Player
	NowPlaying         string
	LastTrack          string
	LyricsState        lyricsState
	StationPickerOpen  bool
	StationHistoryOpen bool
	VibePickerOpen     bool
	MoodPickerOpen     bool
	GenrePickerOpen    bool
	MoodIndex          int
	GenreIndex         int
	VibeIndex          int
	StationIndex       int
	CurrentStationID   string
	StationHistory     radio.History
	ActiveIntent       radio.Intent
}

func (m model) CurrentStation() (*radio.Station, bool) {

	if m.CurrentStationID == "" {
		return nil, false
	}

	return radio.Find(
		m.CurrentStationID,
	)

}

type tickMsg time.Time

type lrclibSongMsg struct {
	TrackID      string
	Song         music.Song
	SyncedLyrics string
	Err          error
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
	trackID string,
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
				TrackID: trackID,
				Err:     err,
			}
		}

		result, ok := lrclib.BestMatch(
			results,
			duration,
		)
		if !ok {
			return lrclibSongMsg{
				TrackID: trackID,
				Err: fmt.Errorf(
					"no LRCLIB match for %s",
					title,
				),
			}
		}

		return lrclibSongMsg{
			TrackID:      trackID,
			Song:         lrclib.Song(result),
			SyncedLyrics: result.SyncedLyrics,
		}

	}

}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		key := msg.String()

		if m.MoodPickerOpen {

			switch key {

			case "esc":
				m.MoodPickerOpen = false
				return m, nil

			case "up", "k":

				if m.MoodIndex > 0 {
					m.MoodIndex--
				}

				return m, nil

			case "down", "j":

				if m.MoodIndex < len(radio.MoodFamilies)-1 {
					m.MoodIndex++
				}

				return m, nil

			case "enter":

				if len(radio.MoodFamilies) == 0 {
					return m, nil
				}

				family := radio.MoodFamilies[m.MoodIndex]
				intent := radio.MoodIntent(
					family,
				)

				station, ok := radio.Choose(
					intent.Criteria,
					m.StationHistory,
					radio.ChooseOptions{
						RecentLimit: 3,
					},
				)
				if !ok {
					return m, nil
				}

				err := m.Player.Load(
					station.StreamURL,
				)
				if err != nil {
					return m, nil
				}

				m.ActiveIntent = intent
				m.MoodPickerOpen = false

				return m, nil

			}

			return m, nil

		}

		if m.GenrePickerOpen {

			genres := radio.Genres()

			switch key {

			case "esc":
				m.GenrePickerOpen = false
				return m, nil

			case "up", "k":

				if m.GenreIndex > 0 {
					m.GenreIndex--
				}

				return m, nil

			case "down", "j":

				if m.GenreIndex < len(genres)-1 {
					m.GenreIndex++
				}

				return m, nil

			case "enter":

				if len(genres) == 0 {
					return m, nil
				}

				genre := genres[m.GenreIndex]
				intent := radio.GenreIntent(
					genre,
				)

				station, ok := radio.Choose(
					intent.Criteria,
					m.StationHistory,
					radio.ChooseOptions{
						RecentLimit: 3,
					},
				)
				if !ok {
					return m, nil
				}

				err := m.Player.Load(
					station.StreamURL,
				)
				if err != nil {
					return m, nil
				}

				m.ActiveIntent = intent
				m.GenrePickerOpen = false

				return m, nil

			}

			return m, nil

		}

		if m.VibePickerOpen {

			switch key {

			case "esc":
				m.VibePickerOpen = false
				return m, nil

			case "up", "k":

				if m.VibeIndex > 0 {
					m.VibeIndex--
				}

				return m, nil

			case "down", "j":

				if m.VibeIndex < len(radio.Presets)-1 {
					m.VibeIndex++
				}

				return m, nil

			case "enter":

				if len(radio.Presets) == 0 {
					return m, nil
				}

				preset := radio.Presets[m.VibeIndex]
				intent := radio.VibeIntent(
					preset,
				)

				station, ok := radio.Choose(
					intent.Criteria,
					m.StationHistory,
					radio.ChooseOptions{
						RecentLimit: 3,
					},
				)
				if !ok {
					return m, nil
				}

				err := m.Player.Load(
					station.StreamURL,
				)
				if err != nil {
					return m, nil
				}

				m.ActiveIntent = intent
				m.VibePickerOpen = false

				return m, nil

			}

			return m, nil

		}

		if m.StationHistoryOpen {

			switch key {

			case "esc", "h":
				m.StationHistoryOpen = false
				return m, nil

			}

			return m, nil

		}

		if m.StationPickerOpen {

			switch key {

			case "enter":

				if len(radio.Stations) == 0 {
					return m, nil
				}

				station := radio.Stations[m.StationIndex]

				err := m.Player.Load(
					station.StreamURL,
				)
				if err != nil {
					return m, nil
				}

				m.CurrentStationID = station.ID
				m.StationPickerOpen = false

				return m, nil

			case "esc":
				m.StationPickerOpen = false
				return m, nil

			case "up", "k":

				if m.StationIndex > 0 {
					m.StationIndex--
				}

				return m, nil

			case "down", "j":

				if m.StationIndex < len(radio.Stations)-1 {
					m.StationIndex++
				}

				return m, nil

			}

			return m, nil

		}

		switch key {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "s":
			m.StationPickerOpen = true
			m.StationIndex = 0

			return m, nil

		case "h":
			m.StationHistoryOpen = true
			return m, nil

		case "v":
			m.VibePickerOpen = true
			m.VibeIndex = 0

			return m, nil

		case "m":
			m.MoodPickerOpen = true
			m.MoodIndex = 0

			return m, nil

		case "g":
			m.GenrePickerOpen = true
			m.GenreIndex = 0

			return m, nil

		}

	case tea.WindowSizeMsg:

		m.Width = msg.Width
		m.Height = msg.Height

		return m, nil

	case tickMsg:

		position := m.Player.Position()
		duration := m.Player.Duration()

		if duration > 0 {
			m.Song.Duration = duration
		}

		rawTitle := m.Player.Title()
		artist := m.Player.Artist()
		title := m.Player.TrackTitle()
		album := m.Player.Album()
		trackID := m.Player.TrackID()
		path := m.Player.Path()
		isNetwork := m.Player.IsNetwork()
		track := metadata.Resolve(rawTitle)
		station, ok := radio.FindByStreamURL(path)

		if ok {

			m.CurrentStationID = station.ID

			m.StationHistory.Add(
				station.ID,
				time.Now(),
			)

		} else {

			m.CurrentStationID = ""

		}

		if !isNetwork {
			if artist != "" {
				track.Artist = artist
			}

			if title != "" {
				track.Title = title
			}
		}

		if track.RawTitle == "" {
			return m, tick()
		}

		if track.Valid {

			m.NowPlaying = track.RawTitle

			if trackID != m.LastTrack {

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

					m.Song.Timeline = song.Timeline
					m.LyricsState = lyricsLocal
					m.LastTrack = trackID

					return m, tick()

				}

				m.LastTrack = trackID
				m.LyricsState = lyricsSearching

				return m, tea.Batch(
					tick(),
					loadLRCLIBSong(
						trackID,
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

		if msg.TrackID != m.LastTrack {
			return m, nil
		}

		if msg.Err != nil {
			m.LyricsState = lyricsUnavailable
			return m, nil
		}

		_, err := library.Store(
			m.Song.Artist,
			m.Song.Title,
			msg.SyncedLyrics,
		)
		if err != nil {
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

	if m.MoodPickerOpen {

		return ui.RenderMoodPicker(
			radio.MoodFamilies,
			m.MoodIndex,
			m.Width,
		)

	}

	if m.GenrePickerOpen {

		return ui.RenderGenrePicker(
			radio.Genres(),
			m.GenreIndex,
			m.Width,
		)

	}

	if m.VibePickerOpen {

		return ui.RenderVibePicker(
			radio.Presets,
			m.VibeIndex,
			m.Width,
		)

	}

	if m.StationHistoryOpen {

		return ui.RenderStationHistory(
			m.StationHistory,
			m.Width,
		)

	}

	if m.StationPickerOpen {

		return ui.RenderStationPicker(
			radio.Stations,
			m.StationIndex,
			m.Width,
		)

	}

	stationName := ""
	intentType := ""
	intentName := ""

	if m.ActiveIntent.Active() {
		intentType = string(m.ActiveIntent.Type)
		intentName = m.ActiveIntent.Name
	}
	if station, ok := m.CurrentStation(); ok {
		stationName = station.Name
	}

	return ui.Render(
		m.Song,
		m.Width,
		m.OnAir,
		m.CurrentCue,
		m.Player.Position(),
		m.NowPlaying,
		m.LyricsState.String(),
		stationName,
		intentType,
		intentName,
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
