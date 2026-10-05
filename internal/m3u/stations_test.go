package m3u

import (
	"os"
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

func TestStationProposal(
	t *testing.T,
) {

	entry := Entry{
		Ref: source.ItemRef{
			Source: source.Source{
				Kind: source.KindM3U,
				URI:  "stations.m3u",
			},
			URI:  "https://example.com/live",
			Name: "Example FM",
		},
		Duration: -1,
	}

	proposal, ok :=
		StationProposal(
			entry,
		)

	if !ok {
		t.Fatal(
			"expected station proposal",
		)
	}

	if proposal.Name != "Example FM" {
		t.Errorf(
			"expected station name, got %q",
			proposal.Name,
		)
	}

	if proposal.StreamURL !=
		"https://example.com/live" {

		t.Errorf(
			"unexpected stream URL %q",
			proposal.StreamURL,
		)
	}

	if proposal.Source.Source.Kind !=
		source.KindM3U {

		t.Errorf(
			"expected M3U provenance, got %q",
			proposal.Source.Source.Kind,
		)
	}

}

func TestStationProposalDoesNotInventMetadata(
	t *testing.T,
) {

	entry := Entry{
		Ref: source.ItemRef{
			Source: source.Source{
				Kind: source.KindM3U,
			},
			URI:  "https://example.com/live",
			Name: "PARTY VIBE RADIO : POP",
		},
	}

	proposal, ok :=
		StationProposal(
			entry,
		)

	if !ok {
		t.Fatal(
			"expected station proposal",
		)
	}

	if proposal.Genre != "" {
		t.Errorf(
			"expected no inferred genre, got %q",
			proposal.Genre,
		)
	}

	if len(proposal.Tags) != 0 {
		t.Fatal(
			"expected no inferred tags",
		)
	}

	if len(proposal.Moods) != 0 {
		t.Fatal(
			"expected no inferred moods",
		)
	}

	if len(proposal.Contexts) != 0 {
		t.Fatal(
			"expected no inferred contexts",
		)
	}

	if proposal.Energy != 0 {
		t.Errorf(
			"expected unknown energy, got %d",
			proposal.Energy,
		)
	}

}

func TestStationProposalRejectsMissingURI(
	t *testing.T,
) {

	entry := Entry{
		Ref: source.ItemRef{
			Source: source.Source{
				Kind: source.KindM3U,
			},
			Name: "No Stream",
		},
	}

	_, ok := StationProposal(
		entry,
	)

	if ok {
		t.Fatal(
			"expected missing URI to be rejected",
		)
	}

}

func TestStationProposalsSkipsInvalidEntries(
	t *testing.T,
) {

	origin := source.Source{
		Kind: source.KindM3U,
		URI:  "stations.m3u",
	}

	entries := []Entry{
		{
			Ref: source.ItemRef{
				Source: origin,
				URI:    "https://example.com/one",
				Name:   "One",
			},
		},
		{
			Ref: source.ItemRef{
				Source: origin,
				Name:   "Missing URI",
			},
		},
		{
			Ref: source.ItemRef{
				Source: origin,
				URI:    "https://example.com/two",
				Name:   "Two",
			},
		},
	}

	proposals :=
		StationProposals(
			entries,
		)

	if len(proposals) != 2 {
		t.Fatalf(
			"expected 2 proposals, got %d",
			len(proposals),
		)
	}

}
func TestRealExtendedPlaylistProducesStationProposals(
	t *testing.T,
) {

	file, err := os.Open(
		"testdata/extended.m3u",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	origin := source.Source{
		Kind: source.KindM3U,
		URI:  "testdata/extended.m3u",
	}

	entries, err := Parse(
		file,
		origin,
	)
	if err != nil {
		t.Fatal(err)
	}

	proposals :=
		StationProposals(
			entries,
		)

	if len(proposals) == 0 {
		t.Fatal(
			"expected station proposals",
		)
	}

	if len(proposals) != len(entries) {

		t.Fatalf(
			"expected %d proposals, got %d",
			len(entries),
			len(proposals),
		)

	}

	for i, proposal := range proposals {

		if !proposal.Valid() {

			t.Fatalf(
				"expected proposal %d to be valid",
				i,
			)

		}

		if proposal.Source.Source.Kind !=
			source.KindM3U {

			t.Fatalf(
				"proposal %d lost M3U provenance",
				i,
			)

		}

	}

}

func TestStationMediaItem(
	t *testing.T,
) {

	entry := Entry{
		Ref: source.ItemRef{
			Source: source.Source{
				Kind: source.KindM3U,
				URI:  "stations.m3u",
			},
			URI:  "https://example.com/live",
			Name: "Example FM",
		},
	}

	item, ok :=
		StationMediaItem(
			entry,
		)

	if !ok {
		t.Fatal(
			"expected station media item",
		)
	}

	if item.Kind !=
		source.MediaStation {

		t.Errorf(
			"expected station kind, got %q",
			item.Kind,
		)

	}

	if item.Title != "Example FM" {
		t.Errorf(
			"expected station title, got %q",
			item.Title,
		)
	}

	if item.Ref.Source.Kind !=
		source.KindM3U {

		t.Errorf(
			"expected M3U provenance, got %q",
			item.Ref.Source.Kind,
		)

	}

}
