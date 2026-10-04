package radio

import (
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

func TestStationProposalValid(
	t *testing.T,
) {

	proposal := StationProposal{
		Name:      "Example FM",
		StreamURL: "https://example.com/live",
		Source: source.ItemRef{
			Source: source.Source{
				Kind: source.KindM3U,
				URI:  "stations.m3u",
			},
			URI: "https://example.com/live",
		},
	}

	if !proposal.Valid() {
		t.Fatal(
			"expected station proposal to be valid",
		)
	}

}

func TestStationProposalRejectsMissingStream(
	t *testing.T,
) {

	proposal := StationProposal{
		Source: source.ItemRef{
			Source: source.Source{
				Kind: source.KindM3U,
			},
			ID: "station-1",
		},
	}

	if proposal.Valid() {
		t.Fatal(
			"expected proposal without stream URL to be invalid",
		)
	}

}

func TestStationProposalRejectsMissingSource(
	t *testing.T,
) {

	proposal := StationProposal{
		StreamURL: "https://example.com/live",
	}

	if proposal.Valid() {
		t.Fatal(
			"expected proposal without source provenance to be invalid",
		)
	}

}
