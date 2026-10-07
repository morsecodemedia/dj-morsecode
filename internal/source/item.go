package source

type MediaKind string

const (
	MediaUnknown  MediaKind = ""
	MediaStation  MediaKind = "station"
	MediaTrack    MediaKind = "track"
	MediaVideo    MediaKind = "video"
	MediaPlaylist MediaKind = "playlist"
	MediaAlbum    MediaKind = "album"
	MediaArtist   MediaKind = "artist"
)

type MediaItem struct {
	Kind MediaKind

	Ref ItemRef

	Artist string
	Title  string
}

func (i MediaItem) Valid() bool {

	return i.Kind != MediaUnknown &&
		i.Ref.Valid()

}
