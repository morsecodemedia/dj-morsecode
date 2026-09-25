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
	Title string
	Date  string

	ReleaseGroupID string

	PrimaryType    string
	SecondaryTypes []string

	Provider string
}
