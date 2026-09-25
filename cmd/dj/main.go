package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/morsecodemedia/dj-morsecode/internal/enrichment"
	"github.com/morsecodemedia/dj-morsecode/internal/library"
	"github.com/morsecodemedia/dj-morsecode/internal/lrclib"
	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
	"github.com/morsecodemedia/dj-morsecode/internal/musicbrainz"
	"github.com/morsecodemedia/dj-morsecode/internal/player"
	"github.com/morsecodemedia/dj-morsecode/internal/player/mpv"
	"github.com/morsecodemedia/dj-morsecode/internal/radio"
	"github.com/morsecodemedia/dj-morsecode/internal/ui"
)

type lyricsState int

const mpvSocketPath = "/tmp/dj-morsecode.sock"
const stationRotationInterval = 30 * time.Minute
const stationTuneGracePeriod = 15 * time.Second

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
	PlaybackItem       metadata.PlaybackItem
	LyricsState        lyricsState
	EnrichmentService  *enrichment.Service
	EnrichmentMatch    metadata.EnrichmentMatch
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
	PendingStationID   string
	PendingSince       time.Time
	FailedStationIDs   []string
}

func (m model) CurrentStation() (*radio.Station, bool) {

	if m.CurrentStationID == "" {
		return nil, false
	}

	return radio.Find(
		m.CurrentStationID,
	)

}

func (m model) failStation(
	stationID string,
) model {

	if stationID == "" {
		return m
	}

	for _, id := range m.FailedStationIDs {

		if id == stationID {
			return m
		}

	}

	m.FailedStationIDs = append(
		m.FailedStationIDs,
		stationID,
	)

	return m

}

func (m model) clearTrackState() model {

	m.NowPlaying = ""
	m.LastTrack = ""

	m.Song = music.Song{}
	m.CurrentCue = 0
	m.LyricsState = lyricsUnavailable

	return m

}

func (m model) nextStation() model {

	if !m.ActiveIntent.Active() {
		return m
	}

	station, ok := radio.Choose(
		m.ActiveIntent.Criteria,
		m.StationHistory,
		radio.ChooseOptions{
			RecentLimit: 3,
			Chooser:     radio.RandomCandidate,
			ExcludeIDs:  m.FailedStationIDs,
		},
	)
	if !ok {
		return m
	}

	err := m.Player.Load(
		station.StreamURL,
	)
	if err != nil {
		return m
	}

	m.PendingStationID = station.ID
	m.PendingSince = time.Now()

	return m

}

func (m model) shouldRotateStation(
	now time.Time,
) bool {

	if !m.ActiveIntent.Active() {
		return false
	}

	current, ok := m.StationHistory.Current()
	if !ok {
		return false
	}

	return now.Sub(
		current.TunedAt,
	) >= stationRotationInterval

}

type tickMsg time.Time

type lrclibSongMsg struct {
	TrackID      string
	Song         music.Song
	SyncedLyrics string
	Err          error
}

type enrichmentMsg struct {
	TrackID string

	Match  metadata.EnrichmentMatch
	Status metadata.MatchStatus
	Err    error
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

func enrichTrack(
	service *enrichment.Service,
	trackID string,
	item metadata.PlaybackItem,
) tea.Cmd {

	if service == nil {
		return nil
	}

	return func() tea.Msg {

		match, status, err := service.Enrich(
			context.Background(),
			item,
		)

		return enrichmentMsg{
			TrackID: trackID,
			Match:   match,
			Status:  status,
			Err:     err,
		}

	}

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
				m.FailedStationIDs = nil
				station, ok := radio.Choose(
					intent.Criteria,
					m.StationHistory,
					radio.ChooseOptions{
						RecentLimit: 3,
						Chooser:     radio.RandomCandidate,
						ExcludeIDs:  m.FailedStationIDs,
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

				m.PendingStationID = station.ID
				m.PendingSince = time.Now()

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
				m.FailedStationIDs = nil
				station, ok := radio.Choose(
					intent.Criteria,
					m.StationHistory,
					radio.ChooseOptions{
						RecentLimit: 3,
						Chooser:     radio.RandomCandidate,
						ExcludeIDs:  m.FailedStationIDs,
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

				m.PendingStationID = station.ID
				m.PendingSince = time.Now()

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
				m.FailedStationIDs = nil
				station, ok := radio.Choose(
					intent.Criteria,
					m.StationHistory,
					radio.ChooseOptions{
						RecentLimit: 3,
						Chooser:     radio.RandomCandidate,
						ExcludeIDs:  m.FailedStationIDs,
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

				m.PendingStationID = station.ID
				m.PendingSince = time.Now()

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

				m.ActiveIntent = radio.Intent{}
				m.PendingStationID = ""
				m.PendingSince = time.Time{}
				m.StationPickerOpen = false
				m.FailedStationIDs = nil

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

		case "n":
			m = m.nextStation()
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

		station, stationFound := radio.FindByStreamURL(
			path,
		)

		stationChanged := false

		if stationFound {

			stationChanged =
				m.CurrentStationID != "" &&
					m.CurrentStationID != station.ID

			if stationChanged {

				m.PlaybackItem = metadata.PlaybackItem{}
				m = m.clearTrackState()

			}

			m.CurrentStationID = station.ID

		} else {

			if m.CurrentStationID != "" {

				m.PlaybackItem = metadata.PlaybackItem{}
				m = m.clearTrackState()

			}

			m.CurrentStationID = ""

		}

		track := metadata.Resolve(
			rawTitle,
		)

		if isNetwork {

			observedItem := metadata.Normalize(
				m.Player.Metadata(),
			)

			if observedItem.Observed() {

				m.PlaybackItem = observedItem

			}

			if m.PlaybackItem.IsTrack() {

				track.RawTitle = m.PlaybackItem.RawTitle
				track.Artist = m.PlaybackItem.Artist
				track.Title = m.PlaybackItem.Title
				track.Valid = true

				album = m.PlaybackItem.Album

				if m.PlaybackItem.Duration > 0 {
					duration = m.PlaybackItem.Duration
				}

				trackID = path +
					"\x00" +
					m.PlaybackItem.Artist +
					"\x00" +
					m.PlaybackItem.Title

			} else {

				track.Valid = false

			}

		} else {

			m.PlaybackItem = metadata.PlaybackItem{}

		}

		if stationFound {

			idle := m.Player.IsIdle()

			if !idle {

				m.StationHistory.Add(
					station.ID,
					time.Now(),
				)

			}

			if station.ID == m.PendingStationID &&
				!idle {

				m.PendingStationID = ""
				m.PendingSince = time.Time{}

			}

		}

		if m.PendingStationID != "" &&
			!m.PendingSince.IsZero() &&
			time.Since(m.PendingSince) >= stationTuneGracePeriod &&
			m.Player.IsIdle() {

			failedStationID := m.PendingStationID

			m.PendingStationID = ""
			m.PendingSince = time.Time{}

			m = m.failStation(
				failedStationID,
			)

			m = m.nextStation()

			return m, tick()

		}

		if m.shouldRotateStation(
			time.Now(),
		) {

			m = m.nextStation()

			return m, tick()

		}

		if !isNetwork {
			if artist != "" {
				track.Artist = artist
			}

			if title != "" {
				track.Title = title
			}
		}

		if isNetwork &&
			!m.PlaybackItem.IsTrack() {

			m = m.clearTrackState()

			return m, tick()

		}

		if track.RawTitle == "" {
			return m, tick()
		}

		if track.Valid {

			if isNetwork {

				m.NowPlaying = m.PlaybackItem.DisplayTitle()

			} else {

				m.NowPlaying = track.RawTitle

			}

			if trackID != m.LastTrack {

				m.Song = music.Song{
					Title:    track.Title,
					Artist:   track.Artist,
					Album:    album,
					Duration: duration,
				}

				m.CurrentCue = 0
				m.EnrichmentMatch =
					metadata.EnrichmentMatch{}

				enrichmentDuration := time.Duration(0)

				if isNetwork &&
					m.PlaybackItem.Duration > 0 {

					enrichmentDuration =
						m.PlaybackItem.Duration

				} else if !isNetwork {

					enrichmentDuration = duration

				}

				enrichmentItem := metadata.PlaybackItem{
					Type:     metadata.PlaybackTrack,
					Artist:   track.Artist,
					Title:    track.Title,
					Duration: enrichmentDuration,
				}

				enrichmentCmd := enrichTrack(
					m.EnrichmentService,
					trackID,
					enrichmentItem,
				)

				song, ok := library.Load(
					track.Artist,
					track.Title,
				)

				if ok {

					m.Song.Timeline = song.Timeline
					m.LyricsState = lyricsLocal
					m.LastTrack = trackID

					return m, tea.Batch(
						tick(),
						enrichmentCmd,
					)

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
					enrichmentCmd,
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

	case enrichmentMsg:

		if msg.Err != nil {
			return m, nil
		}

		if msg.TrackID != m.LastTrack {
			return m, nil
		}

		if msg.Status != metadata.MatchAccepted {
			return m, nil
		}

		m.EnrichmentMatch = msg.Match

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

	if m.NowPlaying == "" &&
		m.CurrentStationID == "" {

		return ui.RenderIdle(
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

	view := ui.Render(
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

	if m.EnrichmentMatch.Track.Title != "" {

		view += fmt.Sprintf(
			"\nENRICHED • %s • %.0f%%",
			m.EnrichmentMatch.Provider,
			m.EnrichmentMatch.Confidence*100,
		)

		for _, identifier := range m.EnrichmentMatch.Track.Identifiers {

			if identifier.Scheme ==
				metadata.IdentifierMusicBrainz {

				view += fmt.Sprintf(
					"\nMBID • %s",
					identifier.Value,
				)

				break
			}

		}

	}

	return view

}

func main() {

	musicBrainzClient := musicbrainz.NewClient(
		"DJ MorseCode/1.0.0 (https://github.com/morsecodemedia/dj-morsecode)",
	)

	musicBrainzEnricher := musicbrainz.NewEnricher(
		musicBrainzClient,
	)

	enrichmentStore := metadata.NewEnrichmentStore()

	enrichmentService := enrichment.NewService(
		enrichmentStore,
		musicBrainzEnricher,
	)

	mpvProcess := mpv.NewProcess(
		mpvSocketPath,
	)

	if err := mpvProcess.Start(); err != nil {

		fmt.Println(
			"Unable to summon MPV:",
			err,
		)

		os.Exit(1)

	}

	if err := mpvProcess.WaitReady(); err != nil {

		_ = mpvProcess.Stop()

		fmt.Println(
			"MPV failed to answer the call:",
			err,
		)

		os.Exit(1)

	}

	defer func() {

		if err := mpvProcess.Stop(); err != nil {

			fmt.Println(
				"Unable to stop MPV:",
				err,
			)

		}

	}()

	playback, err := player.New(
		mpvSocketPath,
	)
	if err != nil {

		fmt.Println(
			"Unable to connect to MPV:",
			err,
		)

		return

	}

	p := tea.NewProgram(model{
		Player:            playback,
		EnrichmentService: enrichmentService,
	})

	if _, err := p.Run(); err != nil {

		fmt.Println(err)
		os.Exit(1)

	}

}
