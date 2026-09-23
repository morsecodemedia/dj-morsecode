package metadata

import "strings"

func normalizeLautFM(
	fields map[string]string,
) PlaybackItem {

	rawTitle := metadataTitle(
		fields,
	)

	item := PlaybackItem{
		Type:         PlaybackUnknown,
		RawTitle:     rawTitle,
		SourceFields: fields,
	}

	if rawTitle == "" {
		return item
	}

	if strings.Contains(
		strings.ToLower(rawTitle),
		"verbraucherinformationen",
	) {

		return item
	}

	artist, title, ok := splitArtistTitle(
		rawTitle,
	)
	if !ok {
		return item
	}

	item.Type = PlaybackTrack
	item.Artist = artist
	item.Title = title

	return item

}
