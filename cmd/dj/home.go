package main

import "github.com/morsecodemedia/dj-morsecode/internal/controls"

type homeItem struct {
	Label  string
	Action controls.Action
}

var homeItems = []homeItem{
	{
		Label:  "Pick a vibe",
		Action: controls.ActionOpenVibes,
	},
	{
		Label:  "Set a mood",
		Action: controls.ActionOpenMoods,
	},
	{
		Label:  "Choose a genre",
		Action: controls.ActionOpenGenres,
	},
	{
		Label:  "Tune a station",
		Action: controls.ActionOpenStations,
	},
	{
		Label:  "Play something",
		Action: controls.ActionOpenYouTube,
	},
}

func (m model) homeActive() bool {

	return m.NowPlaying == "" &&
		m.CurrentStationID == "" &&
		m.PendingStationID == ""
}

func homeLabels() []string {

	labels := make(
		[]string,
		0,
		len(homeItems),
	)

	for _, item := range homeItems {

		labels = append(
			labels,
			item.Label,
		)

	}

	return labels
}
