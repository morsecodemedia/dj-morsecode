package metadata

import "testing"

func TestResolveContextTagsMergesScopedEvidence(
	t *testing.T,
) {

	evidence := []TagEvidence{
		{
			Name:     "Soul",
			Count:    16,
			Scope:    TagScopeTrack,
			Provider: "lastfm",
		},
		{
			Name:     " soul ",
			Count:    56,
			Scope:    TagScopeArtist,
			Provider: "lastfm",
		},
	}

	tags := ResolveContextTags(
		evidence,
	)

	if len(tags) != 1 {

		t.Fatalf(
			"expected one context tag, got %d",
			len(tags),
		)

	}

	tag := tags[0]

	if tag.Name != "soul" {
		t.Errorf(
			"expected name %q, got %q",
			"soul",
			tag.Name,
		)
	}

	if tag.TrackCount != 16 {
		t.Errorf(
			"expected track count %d, got %d",
			16,
			tag.TrackCount,
		)
	}

	if tag.ArtistCount != 56 {
		t.Errorf(
			"expected artist count %d, got %d",
			56,
			tag.ArtistCount,
		)
	}

	if len(tag.Providers) != 1 ||
		tag.Providers[0] != "lastfm" {

		t.Errorf(
			"unexpected providers %v",
			tag.Providers,
		)

	}

}

func TestResolveContextTagsPreservesSemanticEvidence(
	t *testing.T,
) {

	evidence := []TagEvidence{
		{
			Name:     "pop soul",
			Count:    100,
			Scope:    TagScopeTrack,
			Provider: "lastfm",
		},
		{
			Name:     "2024",
			Count:    100,
			Scope:    TagScopeTrack,
			Provider: "lastfm",
		},
		{
			Name:     "Love",
			Count:    16,
			Scope:    TagScopeTrack,
			Provider: "lastfm",
		},
		{
			Name:     "irish",
			Count:    21,
			Scope:    TagScopeArtist,
			Provider: "lastfm",
		},
	}

	tags := ResolveContextTags(
		evidence,
	)

	expected := []string{
		"2024",
		"irish",
		"love",
		"pop soul",
	}

	if len(tags) != len(expected) {

		t.Fatalf(
			"expected %d tags, got %d",
			len(expected),
			len(tags),
		)

	}

	for index, name := range expected {

		if tags[index].Name != name {

			t.Errorf(
				"expected tag %d to be %q, got %q",
				index,
				name,
				tags[index].Name,
			)

		}

	}

}

func TestResolveContextTagsKeepsStrongestCountPerScope(
	t *testing.T,
) {

	evidence := []TagEvidence{
		{
			Name:     "ambient",
			Count:    78,
			Scope:    TagScopeArtist,
			Provider: "provider-a",
		},
		{
			Name:     "Ambient",
			Count:    100,
			Scope:    TagScopeArtist,
			Provider: "provider-b",
		},
		{
			Name:     "ambient",
			Count:    42,
			Scope:    TagScopeTrack,
			Provider: "provider-a",
		},
	}

	tags := ResolveContextTags(
		evidence,
	)

	if len(tags) != 1 {

		t.Fatalf(
			"expected one tag, got %d",
			len(tags),
		)

	}

	if tags[0].ArtistCount != 100 {

		t.Errorf(
			"expected strongest artist count %d, got %d",
			100,
			tags[0].ArtistCount,
		)

	}

	if tags[0].TrackCount != 42 {

		t.Errorf(
			"expected track count %d, got %d",
			42,
			tags[0].TrackCount,
		)

	}

}

func TestResolveContextTagsDeduplicatesProviders(
	t *testing.T,
) {

	evidence := []TagEvidence{
		{
			Name:     "jazz",
			Count:    100,
			Scope:    TagScopeArtist,
			Provider: "lastfm",
		},
		{
			Name:     "Jazz",
			Count:    90,
			Scope:    TagScopeTrack,
			Provider: "LASTFM",
		},
		{
			Name:     "jazz",
			Count:    80,
			Scope:    TagScopeArtist,
			Provider: "musicbrainz",
		},
	}

	tags := ResolveContextTags(
		evidence,
	)

	if len(tags) != 1 {
		t.Fatalf(
			"expected one tag, got %d",
			len(tags),
		)
	}

	if len(tags[0].Providers) != 2 {

		t.Fatalf(
			"expected two providers, got %v",
			tags[0].Providers,
		)

	}

}

func TestResolveContextTagsRejectsEmptyNames(
	t *testing.T,
) {

	evidence := []TagEvidence{
		{
			Name:     "",
			Count:    100,
			Scope:    TagScopeTrack,
			Provider: "lastfm",
		},
		{
			Name:     "   ",
			Count:    100,
			Scope:    TagScopeArtist,
			Provider: "lastfm",
		},
	}

	tags := ResolveContextTags(
		evidence,
	)

	if len(tags) != 0 {

		t.Fatalf(
			"expected empty evidence to be rejected, got %v",
			tags,
		)

	}

}
