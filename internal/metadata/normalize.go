package metadata

import "strings"

func Normalize(
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

func metadataTitle(
	fields map[string]string,
) string {

	if value := strings.TrimSpace(
		fields["icy-title"],
	); value != "" {

		return value

	}

	return strings.TrimSpace(
		fields["title"],
	)

}

func splitArtistTitle(
	value string,
) (string, string, bool) {

	parts := strings.SplitN(
		value,
		" - ",
		2,
	)

	if len(parts) != 2 {
		return "", "", false
	}

	artist := strings.TrimSpace(
		parts[0],
	)

	title := strings.TrimSpace(
		parts[1],
	)

	if artist == "" ||
		title == "" {

		return "", "", false
	}

	return artist, title, true

}
