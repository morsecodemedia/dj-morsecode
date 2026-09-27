package musicbrainz

import (
	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

func ReleaseEvidence(
	release Release,
) metadata.ReleaseEvidence {

	return metadata.ReleaseEvidence{
		ReleaseTitle: release.Title,
		ReleaseDate:  release.Date,

		GroupID: release.ReleaseGroup.ID,

		GroupTitle: release.ReleaseGroup.Title,

		GroupFirstReleaseDate: release.ReleaseGroup.FirstReleaseDate,

		PrimaryType: release.ReleaseGroup.PrimaryType,

		SecondaryTypes: append(
			[]string(nil),
			release.ReleaseGroup.SecondaryTypes...,
		),

		Provider: providerName,
	}

}

func RecordingReleaseEvidence(
	recording Recording,
) []metadata.ReleaseEvidence {

	evidence := make(
		[]metadata.ReleaseEvidence,
		0,
		len(recording.Releases),
	)

	for _, release := range recording.Releases {

		evidence = append(
			evidence,
			ReleaseEvidence(
				release,
			),
		)

	}

	return evidence

}
