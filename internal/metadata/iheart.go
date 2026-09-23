package metadata

import (
	"regexp"
	"strings"
	"time"
)

var iheartFieldPattern = regexp.MustCompile(
	`([A-Za-z][A-Za-z0-9_]*)="([^"]*)"`,
)

func isIHeartMetadata(
	fields map[string]string,
) bool {

	return strings.Contains(
		strings.ToLower(fields["icy-url"]),
		"iheartradio.com",
	)

}

func normalizeIHeart(
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

	values := parseIHeartFields(
		rawTitle,
	)

	switch values["song_spot"] {

	case "M":

		item.Type = PlaybackTrack
		item.Title = strings.TrimSpace(
			values["text"],
		)

		if index := strings.Index(
			rawTitle,
			" - text=",
		); index > 0 {

			item.Artist = strings.TrimSpace(
				rawTitle[:index],
			)

		}

		if duration, ok := parseIHeartDuration(
			values["length"],
		); ok {

			item.Duration = duration

		}

	case "T":

		item.Type = PlaybackStationID
		item.Title = strings.TrimSpace(
			values["text"],
		)

	}

	return item

}

func parseIHeartFields(
	rawTitle string,
) map[string]string {

	fields := make(
		map[string]string,
	)

	for _, match := range iheartFieldPattern.FindAllStringSubmatch(
		rawTitle,
		-1,
	) {

		if len(match) != 3 {
			continue
		}

		fields[match[1]] = match[2]

	}

	return fields

}

func parseIHeartDuration(
	value string,
) (time.Duration, bool) {

	parts := strings.Split(
		value,
		":",
	)

	if len(parts) != 3 {
		return 0, false
	}

	hours, err := time.ParseDuration(
		parts[0] + "h",
	)
	if err != nil {
		return 0, false
	}

	minutes, err := time.ParseDuration(
		parts[1] + "m",
	)
	if err != nil {
		return 0, false
	}

	seconds, err := time.ParseDuration(
		parts[2] + "s",
	)
	if err != nil {
		return 0, false
	}

	return hours + minutes + seconds, true

}
