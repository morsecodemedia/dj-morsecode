package metadata

type TrackContext struct {
	Release ReleaseContext
	Genres  []string
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

type TagEvidence struct {
	Name string

	Weight int

	Provider string
}
