package lastfm

import (
	"strings"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

func MediaItem(
	track TrackInfo,
) (source.MediaItem, bool) {

	artist := strings.TrimSpace(
		track.Artist,
	)

	title := strings.TrimSpace(
		track.Title,
	)

	mbid := strings.TrimSpace(
		track.MBID,
	)

	url := strings.TrimSpace(
		track.URL,
	)

	ref := source.ItemRef{
		Source: source.Source{
			Kind: source.KindLastFM,
		},
		ID:  mbid,
		URI: url,
		Name: strings.TrimSpace(
			artist + " - " + title,
		),
	}

	item := source.MediaItem{
		Kind:   source.MediaTrack,
		Ref:    ref,
		Artist: artist,
		Title:  title,
	}

	if !item.Valid() {
		return source.MediaItem{}, false
	}

	return item, true

}
