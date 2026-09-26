package metadata

type ContextTag struct {
	Name string

	TrackCount  int
	ArtistCount int

	Providers []string
}

type TrackContext struct {
	Release ReleaseContext
	Genres  []string
	Tags    []ContextTag
}

type ReleaseContext struct {
	Title string
	Date  string
}

type ReleaseEvidence struct {
	ReleaseTitle string
	ReleaseDate  string

	GroupID               string
	GroupTitle            string
	GroupFirstReleaseDate string

	PrimaryType    string
	SecondaryTypes []string

	Provider string
}

type TagScope string

const (
	TagScopeTrack  TagScope = "track"
	TagScopeArtist TagScope = "artist"
)

type TagEvidence struct {
	Name  string
	Count int
	Scope TagScope

	Provider string
}
