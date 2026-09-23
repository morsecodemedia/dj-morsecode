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

	if isLautFMAdvertisement(
		rawTitle,
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

func isLautFMAdvertisement(
	rawTitle string,
) bool {

	value := strings.ToLower(
		rawTitle,
	)

	if strings.Contains(
		value,
		"verbraucherinformationen",
	) {

		return true
	}

	artist, _, ok := splitArtistTitle(
		rawTitle,
	)
	if !ok {
		return false
	}

	return isDomainLike(
		artist,
	)

}

func isDomainLike(
	value string,
) bool {

	value = strings.TrimSpace(
		strings.ToLower(value),
	)

	return strings.Contains(
		value,
		".",
	) &&
		!strings.Contains(
			value,
			" ",
		)

}
