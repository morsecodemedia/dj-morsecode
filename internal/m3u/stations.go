package m3u

import (
	"strings"

	"github.com/morsecodemedia/dj-morsecode/internal/radio"
	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

func StationProposal(
	entry Entry,
) (radio.StationProposal, bool) {

	streamURL := strings.TrimSpace(
		entry.Ref.URI,
	)

	if streamURL == "" {
		return radio.StationProposal{}, false
	}

	proposal := radio.StationProposal{
		Name: strings.TrimSpace(
			entry.Ref.Name,
		),
		StreamURL: streamURL,
		Source:    entry.Ref,
	}

	if !proposal.Valid() {
		return radio.StationProposal{}, false
	}

	return proposal, true

}

func StationProposals(
	entries []Entry,
) []radio.StationProposal {

	proposals := make(
		[]radio.StationProposal,
		0,
		len(entries),
	)

	for _, entry := range entries {

		proposal, ok :=
			StationProposal(
				entry,
			)

		if !ok {
			continue
		}

		proposals = append(
			proposals,
			proposal,
		)

	}

	return proposals

}

func StationMediaItem(
	entry Entry,
) (source.MediaItem, bool) {

	if !entry.Ref.Valid() {
		return source.MediaItem{}, false
	}

	item := source.MediaItem{
		Kind: source.MediaStation,
		Ref:  entry.Ref,
		Title: strings.TrimSpace(
			entry.Ref.Name,
		),
	}

	if !item.Valid() {
		return source.MediaItem{}, false
	}

	return item, true

}
