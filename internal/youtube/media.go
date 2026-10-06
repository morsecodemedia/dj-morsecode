package youtube

import "github.com/morsecodemedia/dj-morsecode/internal/source"

func MediaItem(
	ref Reference,
) (source.MediaItem, bool) {

	var kind source.MediaKind
	var id string

	switch ref.Kind {

	case ReferenceVideo:

		kind = source.MediaVideo
		id = ref.VideoID

	case ReferencePlaylist:

		kind = source.MediaPlaylist
		id = ref.PlaylistID

	default:

		return source.MediaItem{}, false

	}

	item := source.MediaItem{
		Kind: kind,
		Ref: source.ItemRef{
			Source: source.Source{
				Kind: source.KindYouTube,
			},
			ID:  id,
			URI: ref.URI,
		},
	}

	if !item.Valid() {
		return source.MediaItem{}, false
	}

	return item, true

}
