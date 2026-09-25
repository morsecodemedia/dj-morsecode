package metadata

type TrackContext struct {
	Release ReleaseContext
	Genres  []string
}

type ReleaseContext struct {
	Title string
	Date  string
}
