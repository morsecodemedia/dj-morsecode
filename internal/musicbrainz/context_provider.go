package musicbrainz

import (
	"context"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

type ContextProvider struct {
	client *Client
}

func NewContextProvider(
	client *Client,
) *ContextProvider {

	return &ContextProvider{
		client: client,
	}

}

func (p *ContextProvider) Contextualize(
	ctx context.Context,
	track metadata.CanonicalTrack,
) (
	metadata.TrackContext,
	bool,
	error,
) {

	recordingID, ok := identifierValue(
		track.Identifiers,
		metadata.IdentifierMusicBrainz,
	)
	if !ok {

		return metadata.TrackContext{},
			false,
			nil

	}

	recording, err := p.client.LookupRecording(
		ctx,
		recordingID,
	)
	if err != nil {

		return metadata.TrackContext{},
			false,
			err

	}

	evidence := RecordingReleaseEvidence(
		recording,
	)

	release, ok := metadata.SelectReleaseContext(
		evidence,
	)
	if !ok {

		return metadata.TrackContext{},
			false,
			nil

	}

	return metadata.TrackContext{
		Release: release,
	}, true, nil

}

func identifierValue(
	identifiers []metadata.Identifier,
	scheme string,
) (string, bool) {

	for _, identifier := range identifiers {

		if identifier.Scheme != scheme {
			continue
		}

		if identifier.Value == "" {
			continue
		}

		return identifier.Value, true

	}

	return "", false

}
