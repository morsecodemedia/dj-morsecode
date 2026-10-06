package youtube

import (
	"net/url"
	"strings"
)

type ReferenceKind int

const (
	ReferenceUnknown ReferenceKind = iota
	ReferenceVideo
	ReferencePlaylist
)

type Reference struct {
	Kind ReferenceKind

	VideoID    string
	PlaylistID string

	URI string
}

func ParseURL(
	raw string,
) (Reference, bool) {

	raw = strings.TrimSpace(
		raw,
	)

	if raw == "" {
		return Reference{}, false
	}

	parsed, err := url.Parse(
		raw,
	)
	if err != nil {
		return Reference{}, false
	}

	host := strings.ToLower(
		parsed.Hostname(),
	)

	switch host {

	case "youtu.be":

		videoID := strings.Trim(
			parsed.Path,
			"/",
		)

		if videoID == "" {
			return Reference{}, false
		}

		return Reference{
			Kind:    ReferenceVideo,
			VideoID: videoID,
			URI:     raw,
		}, true

	case "youtube.com",
		"www.youtube.com",
		"m.youtube.com":

		parts := strings.Split(
			strings.Trim(
				parsed.Path,
				"/",
			),
			"/",
		)

		if len(parts) == 2 &&
			parts[0] == "shorts" &&
			parts[1] != "" {

			return Reference{
				Kind:    ReferenceVideo,
				VideoID: parts[1],
				URI:     raw,
			}, true
		}

		query := parsed.Query()

		videoID := strings.TrimSpace(
			query.Get("v"),
		)

		playlistID := strings.TrimSpace(
			query.Get("list"),
		)

		if videoID != "" {

			return Reference{
				Kind:       ReferenceVideo,
				VideoID:    videoID,
				PlaylistID: playlistID,
				URI:        raw,
			}, true

		}

		if playlistID != "" {

			return Reference{
				Kind:       ReferencePlaylist,
				PlaylistID: playlistID,
				URI:        raw,
			}, true

		}

	}

	return Reference{}, false

}
