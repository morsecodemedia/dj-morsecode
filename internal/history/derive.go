package history

import (
	"sort"

	"github.com/morsecodemedia/dj-morsecode/internal/observation"
)

func Derive(
	stations []observation.StationObservation,
	media []observation.MediaObservation,
) []Entry {

	entries := make(
		[]Entry,
		0,
		len(stations)+len(media),
	)

	for _, observed := range stations {

		if observed.Kind !=
			observation.StationTuneConfirmed {

			continue
		}

		entries = append(
			entries,
			Entry{
				Kind:      KindStation,
				StationID: observed.StationID,
				PlayedAt:  observed.ObservedAt,
			},
		)

	}

	for _, observed := range media {

		if !observed.Item.Valid() {
			continue
		}

		entries = append(
			entries,
			Entry{
				Kind:     KindMedia,
				Media:    observed.Item,
				PlayedAt: observed.ObservedAt,
			},
		)

	}

	sort.SliceStable(
		entries,
		func(i int, j int) bool {

			return entries[i].PlayedAt.Before(
				entries[j].PlayedAt,
			)

		},
	)

	return entries
}
