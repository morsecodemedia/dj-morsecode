package lastfm

type TrackInfo struct {
	Artist string
	Title  string

	MBID string
	URL  string

	Listeners int
	Playcount int

	Tags []Tag
}
