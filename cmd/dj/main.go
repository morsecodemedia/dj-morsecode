package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	trackcontext "github.com/morsecodemedia/dj-morsecode/internal/context"
	"github.com/morsecodemedia/dj-morsecode/internal/enrichment"
	"github.com/morsecodemedia/dj-morsecode/internal/lastfm"
	"github.com/morsecodemedia/dj-morsecode/internal/library"
	"github.com/morsecodemedia/dj-morsecode/internal/lrclib"
	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
	"github.com/morsecodemedia/dj-morsecode/internal/musicbrainz"
	"github.com/morsecodemedia/dj-morsecode/internal/observation"
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

type commandMode int

const (
	commandModeNone commandMode = iota
	commandModeControls
	commandModeTune
	commandModeEnhancements
	commandModeInfo
	commandModeVolume
)

type model struct {
	Width                  int
	Height                 int
	Song                   music.Song
	OnAir                  bool
	CurrentCue             int
	Player                 *player.Player
	NowPlaying             string
	LastTrack              string
	PlaybackItem           metadata.PlaybackItem
	LyricsState            lyricsState
	LyricsLookupDuration   time.Duration
	LyricsVisible          bool
	EnrichmentService      *enrichment.Service
	EnrichmentMatch        metadata.EnrichmentMatch
	TrackContext           metadata.TrackContext
	ContextService         *trackcontext.Service
	ObservationService     *observation.Service
	ObservationHistory     *observation.MemorySink
	CommandMode            commandMode
	StationPickerOpen      bool
	StationHistoryOpen     bool
	ObservationHistoryOpen bool
	VibePickerOpen         bool
	MoodPickerOpen         bool
	GenrePickerOpen        bool
	MoodIndex              int
	GenreIndex             int
	VibeIndex              int
	StationIndex           int
	CurrentStationID       string
	StationHistory         radio.History
	ActiveIntent           radio.Intent
	PendingStationID       string
	PendingSince           time.Time
	FailedStationIDs       []string
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
	m.LyricsLookupDuration = 0

	m.TrackContext = metadata.TrackContext{}

	return m

}

func (m model) observeStation(
	kind observation.StationKind,
	stationID string,
) {

	if m.ObservationService == nil {
		return
	}

	_, _, _ = m.ObservationService.ObserveStation(
		kind,
		stationID,
	)

}

func (m model) observePlayback(
	item metadata.PlaybackItem,
	stationID string,
	trackID string,
) {

	if m.ObservationService == nil {
		return
	}

	_, _, _ = m.ObservationService.ObservePlayback(
		item,
		stationID,
		trackID,
	)

}

func (m model) tuneStation(
	station *radio.Station,
) model {

	if station == nil {
		return m
	}

	err := m.Player.Load(
		station.StreamURL,
	)
	if err != nil {
		return m
	}

	m.observeStation(
		observation.StationTuneRequested,
		station.ID,
	)

	m.PendingStationID =
		station.ID

	m.PendingSince =
		time.Now()

	m.PlaybackItem =
		metadata.PlaybackItem{}

	m = m.clearTrackState()

	return m

}

func (m model) nextStation() model {

	if m.ActiveIntent.Active() {

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

		return m.tuneStation(
			&station,
		)
	}

	current, ok :=
		m.CurrentStation()

	if !ok {
		return m
	}

	station, ok := radio.Next(
		current.ID,
	)
	if !ok {
		return m
	}

	return m.tuneStation(
		station,
	)

}

func (m model) previousStation() model {

	if m.ActiveIntent.Active() {
		return m
	}

	current, ok :=
		m.CurrentStation()

	if !ok {
		return m
	}

	station, ok := radio.Previous(
		current.ID,
	)
	if !ok {
		return m
	}

	return m.tuneStation(
		station,
	)

}

func (m model) tuneHistory(
	tune radio.Tune,
) model {

	station, ok := radio.Find(
		tune.StationID,
	)
	if !ok {
		return m
	}

	if err := m.Player.Load(
		station.StreamURL,
	); err != nil {
		return m
	}

	m.observeStation(
		observation.StationTuneRequested,
		station.ID,
	)

	m.PendingStationID = station.ID
	m.PendingSince = time.Now()

	m.PlaybackItem =
		metadata.PlaybackItem{}

	m = m.clearTrackState()

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

func (m model) playbackPosition() time.Duration {

	if m.Player == nil {
		return 0
	}

	if !m.Player.IsNetwork() {
		return m.Player.Position()
	}

	if m.ObservationHistory == nil {
		return 0
	}

	elapsed, ok :=
		m.ObservationHistory.PlaybackElapsed(
			m.LastTrack,
			time.Now(),
		)

	if !ok {
		return 0
	}

	return elapsed

}

func (m model) enterCommandMode(
	mode commandMode,
) model {

	m.CommandMode = mode

	return m
}

func (m model) leaveCommandMode() model {

	m.CommandMode = commandModeNone

	return m
}

func (m model) footerText() string {

	switch m.CommandMode {

	case commandModeControls:

		return fmt.Sprintf(
			"CONTROLS • b previous • n next • %s • %s • %s • esc cancel",
			m.pauseControlLabel(),
			m.muteControlLabel(),
			m.volumeControlLabel(),
		)

	case commandModeVolume:

		if m.Player == nil {
			return "VOLUME • ↑ louder • ↓ quieter • esc back"
		}

		return fmt.Sprintf(
			"VOLUME • %.0f%% • ↑ louder • ↓ quieter • esc back",
			m.Player.Volume(),
		)

	case commandModeTune:

		return "TUNE • s stations • g genres • m moods • v vibes • esc cancel"

	case commandModeEnhancements:

		return "ENHANCEMENTS • l lyrics • t trivia • esc cancel"

	case commandModeInfo:

		return "INFO • o observations • h history • esc cancel"

	default:

		return "c controls • t tune • e enhancements • i info • q sign off"

	}

}

func (m model) pauseControlLabel() string {

	if m.Player != nil &&
		m.Player.Paused() {

		return "p resume"
	}

	return "p pause"

}

func (m model) muteControlLabel() string {

	if m.Player != nil &&
		m.Player.Muted() {

		return "m unmute"
	}

	return "m mute"

}

func (m model) volumeControlLabel() string {

	if m.Player == nil {
		return "v volume"
	}

	volume := m.Player.Volume()

	return fmt.Sprintf(
		"v volume %.0f%%",
		volume,
	)

}

type tickMsg time.Time

type lrclibSongMsg struct {
	TrackID  string
	Song     music.Song
	Content  string
	Err      error
	Duration time.Duration
}

type enrichmentMsg struct {
	TrackID string

	Match  metadata.EnrichmentMatch
	Status metadata.MatchStatus
	Err    error
}

type contextMsg struct {
	TrackID string

	Context metadata.TrackContext
	OK      bool
	Err     error
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

func loadTrackContext(
	task trackcontext.Task,
	trackID string,
	track metadata.CanonicalTrack,
) tea.Cmd {

	if task == nil {
		return nil
	}

	return func() tea.Msg {

		result, ok, err := task(
			context.Background(),
			track,
		)

		return contextMsg{
			TrackID: trackID,
			Context: result,
			OK:      ok,
			Err:     err,
		}

	}

}

func loadTrackContexts(
	service *trackcontext.Service,
	trackID string,
	track metadata.CanonicalTrack,
) tea.Cmd {

	if service == nil {
		return nil
	}

	tasks := service.Tasks()

	commands := make(
		[]tea.Cmd,
		0,
		len(tasks),
	)

	for _, task := range tasks {

		commands = append(
			commands,
			loadTrackContext(
				task,
				trackID,
				track,
			),
		)

	}

	return tea.Batch(
		commands...,
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
				TrackID:  trackID,
				Duration: duration,
				Err:      err,
			}
		}

		result, ok := lrclib.BestMatch(
			results,
			duration,
		)
		if !ok {
			return lrclibSongMsg{
				TrackID:  trackID,
				Duration: duration,
				Err: fmt.Errorf(
					"no LRCLIB match for %s",
					title,
				),
			}
		}

		content := result.SyncedLyrics

		if strings.TrimSpace(
			content,
		) == "" {

			content = result.PlainLyrics
		}

		return lrclibSongMsg{
			TrackID:  trackID,
			Duration: duration,
			Song:     lrclib.Song(result),
			Content:  content,
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

				m.observeStation(
					observation.StationTuneRequested,
					station.ID,
				)

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

				m.observeStation(
					observation.StationTuneRequested,
					station.ID,
				)

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

				m.observeStation(
					observation.StationTuneRequested,
					station.ID,
				)

				m.PendingStationID = station.ID
				m.PendingSince = time.Now()

				m.ActiveIntent = intent
				m.VibePickerOpen = false

				return m, nil

			}

			return m, nil

		}

		if m.ObservationHistoryOpen {

			switch key {

			case "esc", "o":
				m.ObservationHistoryOpen = false
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

				m.observeStation(
					observation.StationTuneRequested,
					station.ID,
				)

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

		if m.CommandMode !=
			commandModeNone {

			switch key {

			case "esc":

				if m.CommandMode ==
					commandModeVolume {

					m = m.enterCommandMode(
						commandModeControls,
					)

					return m, nil
				}

				m = m.leaveCommandMode()

				return m, nil

			}

			switch m.CommandMode {

			case commandModeTune:

				switch key {

				case "s":
					m = m.leaveCommandMode()
					m.StationPickerOpen = true
					m.StationIndex = 0

					return m, nil

				case "g":
					m = m.leaveCommandMode()
					m.GenrePickerOpen = true
					m.GenreIndex = 0

					return m, nil

				case "m":
					m = m.leaveCommandMode()
					m.MoodPickerOpen = true
					m.MoodIndex = 0

					return m, nil

				case "v":
					m = m.leaveCommandMode()
					m.VibePickerOpen = true
					m.VibeIndex = 0

					return m, nil

				}

			case commandModeInfo:

				switch key {

				case "o":
					m = m.leaveCommandMode()
					m.ObservationHistoryOpen = true

					return m, nil

				case "h":
					m = m.leaveCommandMode()
					m.StationHistoryOpen = true

					return m, nil

				}

			case commandModeEnhancements:

				switch key {

				case "l":
					m.LyricsVisible =
						!m.LyricsVisible

					m = m.leaveCommandMode()

					return m, nil

				}

			case commandModeControls:

				switch key {

				case "b":

					if m.StationHistory.CanBack() {

						tune, ok :=
							m.StationHistory.Back()

						if !ok {
							return m, nil
						}

						m = m.tuneHistory(
							tune,
						)

						m = m.leaveCommandMode()

						return m, nil
					}

					m = m.previousStation()
					m = m.leaveCommandMode()

					return m, nil

				case "n":

					if m.StationHistory.CanForward() {

						tune, ok :=
							m.StationHistory.Forward()

						if !ok {
							return m, nil
						}

						m = m.tuneHistory(
							tune,
						)

						m = m.leaveCommandMode()

						return m, nil
					}

					m = m.nextStation()
					m = m.leaveCommandMode()

					return m, nil

				case "p":

					paused := m.Player.Paused()

					if err := m.Player.SetPaused(
						!paused,
					); err != nil {

						return m, nil
					}

					// m = m.leaveCommandMode()

					return m, nil

				case "m":

					muted := m.Player.Muted()

					if err := m.Player.SetMuted(
						!muted,
					); err != nil {

						return m, nil
					}

					// m = m.leaveCommandMode()

					return m, nil

				case "v":

					m = m.enterCommandMode(
						commandModeVolume,
					)

					return m, nil

				}
			case commandModeVolume:

				switch key {

				case "up":

					volume := m.Player.Volume() + 5

					if volume > 100 {
						volume = 100
					}

					if err := m.Player.SetVolume(
						volume,
					); err != nil {

						return m, nil
					}

					return m, nil

				case "down":

					volume := m.Player.Volume() - 5

					if volume < 0 {
						volume = 0
					}

					if err := m.Player.SetVolume(
						volume,
					); err != nil {

						return m, nil
					}

					return m, nil

				}
			}

			return m, nil
		}

		switch key {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "c":
			m = m.enterCommandMode(
				commandModeControls,
			)
			return m, nil

		case "t":
			m = m.enterCommandMode(
				commandModeTune,
			)
			return m, nil

		case "e":
			m = m.enterCommandMode(
				commandModeEnhancements,
			)
			return m, nil

		case "i":
			m = m.enterCommandMode(
				commandModeInfo,
			)
			return m, nil

		}

	case tea.WindowSizeMsg:

		m.Width = msg.Width
		m.Height = msg.Height

		return m, nil

	case tickMsg:

		position := m.playbackPosition()
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

				m.PlaybackItem =
					metadata.PlaybackItem{}

				m = m.clearTrackState()

			}

			m.CurrentStationID = station.ID

		} else if m.PendingStationID == "" {

			if m.CurrentStationID != "" {

				m.PlaybackItem =
					metadata.PlaybackItem{}

				m = m.clearTrackState()

			}

			m.CurrentStationID = ""

		}

		track := metadata.Resolve(
			rawTitle,
		)

		observedPlaybackItem := false

		if isNetwork {

			observedItem := metadata.Normalize(
				m.Player.Metadata(),
			)

			if observedItem.Observed() {

				observedPlaybackItem = true
				m.PlaybackItem = observedItem

			}

			if m.PlaybackItem.IsTrack() {

				track.RawTitle =
					m.PlaybackItem.RawTitle

				track.Artist =
					m.PlaybackItem.Artist

				track.Title =
					m.PlaybackItem.Title

				track.Valid = true

				album = m.PlaybackItem.Album

				if m.PlaybackItem.Duration > 0 {
					duration =
						m.PlaybackItem.Duration
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

			m.PlaybackItem =
				metadata.PlaybackItem{}

		}

		if isNetwork &&
			observedPlaybackItem {

			m.observePlayback(
				m.PlaybackItem,
				m.CurrentStationID,
				trackID,
			)

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

				m.observeStation(
					observation.StationTuneConfirmed,
					station.ID,
				)

				m.PendingStationID = ""
				m.PendingSince = time.Time{}
			}

		}

		if m.PendingStationID != "" &&
			!m.PendingSince.IsZero() &&
			time.Since(m.PendingSince) >= stationTuneGracePeriod &&
			m.Player.IsIdle() {

			failedStationID := m.PendingStationID

			m.observeStation(
				observation.StationTuneFailed,
				failedStationID,
			)

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
			observedPlaybackItem &&
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
				m.TrackContext =
					metadata.TrackContext{}

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
					m.Song.Lyrics = song.Lyrics
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
				m.LyricsLookupDuration = duration

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

		if msg.Duration !=
			m.LyricsLookupDuration {

			return m, nil
		}

		if msg.Err != nil {
			m.LyricsState = lyricsUnavailable
			return m, nil
		}

		_, err := library.Store(
			m.Song.Artist,
			m.Song.Title,
			msg.Content,
		)
		if err != nil {
			m.LyricsState = lyricsUnavailable
			return m, nil
		}

		m.Song.Lyrics = msg.Song.Lyrics
		m.Song.Timeline = msg.Song.Timeline
		m.LyricsState = lyricsRemote
		m.CurrentCue = player.CurrentCue(
			m.Song.Timeline,
			m.playbackPosition(),
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

		contextCmd := loadTrackContexts(
			m.ContextService,
			msg.TrackID,
			msg.Match.Track,
		)

		if msg.Match.Duration <= 0 ||
			msg.Match.Duration == m.LyricsLookupDuration {

			return m, contextCmd
		}

		if m.LyricsState == lyricsLocal ||
			m.LyricsState == lyricsRemote {

			return m, contextCmd
		}

		m.LyricsLookupDuration =
			msg.Match.Duration

		m.LyricsState =
			lyricsSearching

		lyricsCmd := loadLRCLIBSong(
			msg.TrackID,
			msg.Match.Track.Artist,
			msg.Match.Track.Title,
			msg.Match.Duration,
		)

		return m, tea.Batch(
			contextCmd,
			lyricsCmd,
		)

	case contextMsg:

		if msg.TrackID != m.LastTrack {
			return m, nil
		}

		if msg.Err != nil {
			return m, nil
		}

		if !msg.OK {
			return m, nil
		}

		m.TrackContext =
			metadata.MergeTrackContext(
				m.TrackContext,
				msg.Context,
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

	if m.ObservationHistoryOpen {

		if m.ObservationHistory == nil {

			return ui.RenderObservationHistory(
				nil,
				nil,
				m.Width,
			)

		}

		return ui.RenderObservationHistory(
			m.ObservationHistory.Stations(),
			m.ObservationHistory.Playback(),
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
		m.CurrentStationID == "" &&
		m.PendingStationID == "" {

		return ui.RenderIdle(
			m.Width,
			m.footerText(),
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
		m.TrackContext,
		m.Width,
		m.OnAir,
		m.CurrentCue,
		m.playbackPosition(),
		m.NowPlaying,
		m.LyricsState.String(),
		m.LyricsVisible,
		stationName,
		intentType,
		intentName,
		m.footerText(),
	)

	return view

}

func main() {

	lastFMAPIKey := os.Getenv(
		"LASTFM_API_KEY",
	)

	musicBrainzClient := musicbrainz.NewClient(
		"DJ MorseCode/1.0.0 (https://github.com/morsecodemedia/dj-morsecode)",
	)

	musicBrainzEnricher := musicbrainz.NewEnricher(
		musicBrainzClient,
	)

	musicBrainzContextProvider :=
		musicbrainz.NewContextProvider(
			musicBrainzClient,
		)

	enrichmentStore := metadata.NewEnrichmentStore()

	enrichmentService := enrichment.NewService(
		enrichmentStore,
		musicBrainzEnricher,
	)

	var contextProviders []trackcontext.Provider

	contextProviders = append(
		contextProviders,
		musicBrainzContextProvider,
	)

	if lastFMAPIKey != "" {

		lastFMClient := lastfm.NewClient(
			lastFMAPIKey,
		)

		lastFMContextProvider :=
			lastfm.NewContextProvider(
				lastFMClient,
			)

		contextProviders = append(
			contextProviders,
			lastFMContextProvider,
		)

	}

	contextService := trackcontext.NewService(
		contextProviders...,
	)

	observationSink :=
		observation.NewMemorySink()

	observationRecorder :=
		observation.NewRecorder()

	observationService :=
		observation.NewService(
			observationRecorder,
			observationSink,
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
		Player:             playback,
		EnrichmentService:  enrichmentService,
		ContextService:     contextService,
		ObservationService: observationService,
		ObservationHistory: observationSink,
		LyricsVisible:      true,
	})

	if _, err := p.Run(); err != nil {

		fmt.Println(err)
		os.Exit(1)

	}

}
