package metadata

import (
	"sort"
	"strings"
)

func ResolveContextTags(
	evidence []TagEvidence,
) []ContextTag {

	tags := make(
		map[string]*ContextTag,
	)

	for _, item := range evidence {

		name := normalizeContextTagName(
			item.Name,
		)

		if name == "" {
			continue
		}

		tag, ok := tags[name]
		if !ok {

			tag = &ContextTag{
				Name: name,
			}

			tags[name] = tag

		}

		switch item.Scope {

		case TagScopeTrack:

			if item.Count > tag.TrackCount {
				tag.TrackCount = item.Count
			}

		case TagScopeArtist:

			if item.Count > tag.ArtistCount {
				tag.ArtistCount = item.Count
			}

		}

		tag.Providers = appendProvider(
			tag.Providers,
			item.Provider,
		)

	}

	result := make(
		[]ContextTag,
		0,
		len(tags),
	)

	for _, tag := range tags {

		sort.Strings(
			tag.Providers,
		)

		result = append(
			result,
			*tag,
		)

	}

	sort.Slice(
		result,
		func(i int, j int) bool {

			return result[i].Name <
				result[j].Name

		},
	)

	return result

}

func normalizeContextTagName(
	value string,
) string {

	return strings.ToLower(
		strings.TrimSpace(value),
	)

}

func appendProvider(
	providers []string,
	provider string,
) []string {

	provider = strings.TrimSpace(
		provider,
	)

	if provider == "" {
		return providers
	}

	for _, existing := range providers {

		if strings.EqualFold(
			existing,
			provider,
		) {

			return providers
		}

	}

	return append(
		providers,
		provider,
	)

}

func mergeContextTags(
	current []ContextTag,
	incoming []ContextTag,
) []ContextTag {

	evidence := make(
		[]TagEvidence,
		0,
		len(current)+len(incoming),
	)

	for _, tag := range current {

		for _, provider := range tag.Providers {

			if tag.TrackCount > 0 {

				evidence = append(
					evidence,
					TagEvidence{
						Name:     tag.Name,
						Count:    tag.TrackCount,
						Scope:    TagScopeTrack,
						Provider: provider,
					},
				)

			}

			if tag.ArtistCount > 0 {

				evidence = append(
					evidence,
					TagEvidence{
						Name:     tag.Name,
						Count:    tag.ArtistCount,
						Scope:    TagScopeArtist,
						Provider: provider,
					},
				)

			}

		}

	}

	for _, tag := range incoming {

		for _, provider := range tag.Providers {

			if tag.TrackCount > 0 {

				evidence = append(
					evidence,
					TagEvidence{
						Name:     tag.Name,
						Count:    tag.TrackCount,
						Scope:    TagScopeTrack,
						Provider: provider,
					},
				)

			}

			if tag.ArtistCount > 0 {

				evidence = append(
					evidence,
					TagEvidence{
						Name:     tag.Name,
						Count:    tag.ArtistCount,
						Scope:    TagScopeArtist,
						Provider: provider,
					},
				)

			}

		}

	}

	return ResolveContextTags(
		evidence,
	)

}
