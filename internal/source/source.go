package source

type Kind string

const (
	KindBuiltin Kind = "builtin"
	KindManual  Kind = "manual"
	KindM3U     Kind = "m3u"
	KindYouTube Kind = "youtube"
	KindLastFM  Kind = "lastfm"
	KindSpotify Kind = "spotify"
)

type Source struct {
	Kind Kind

	ID  string
	URI string
}

func (s Source) Valid() bool {

	return s.Kind != ""

}

type ItemRef struct {
	Source Source

	ID   string
	URI  string
	Name string
}

func (r ItemRef) Valid() bool {

	if !r.Source.Valid() {
		return false
	}

	return r.ID != "" ||
		r.URI != ""

}
