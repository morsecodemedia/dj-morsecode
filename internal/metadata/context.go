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

func MergeTrackContext(
	current TrackContext,
	incoming TrackContext,
) TrackContext {

	if incoming.Release.Title != "" ||
		incoming.Release.Date != "" {

		current.Release =
			incoming.Release

	}

	if len(incoming.Tags) > 0 {

		current.Tags = mergeContextTags(
			current.Tags,
			incoming.Tags,
		)

	}

	if len(incoming.Genres) > 0 {

		current.Genres = append(
			[]string(nil),
			incoming.Genres...,
		)

	}

	return current

}
