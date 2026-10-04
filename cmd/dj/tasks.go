package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	trackcontext "github.com/morsecodemedia/dj-morsecode/internal/context"
	"github.com/morsecodemedia/dj-morsecode/internal/enrichment"
	"github.com/morsecodemedia/dj-morsecode/internal/lrclib"
	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

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
