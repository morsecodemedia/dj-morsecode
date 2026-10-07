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

type SimilarTrack struct {
	Track TrackInfo
	Match float64
}
