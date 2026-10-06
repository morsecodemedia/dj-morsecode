package youtube

import "testing"

func TestParseURL(
	t *testing.T,
) {

	tests := []struct {
		name string
		raw  string

		kind       ReferenceKind
		videoID    string
		playlistID string

		ok bool
	}{
		{
			name: "standard video",
			raw:  "https://www.youtube.com/watch?v=abc123",

			kind:    ReferenceVideo,
			videoID: "abc123",
			ok:      true,
		},
		{
			name: "short video",
			raw:  "https://youtu.be/abc123",

			kind:    ReferenceVideo,
			videoID: "abc123",
			ok:      true,
		},
		{
			name: "shorts video",
			raw:  "https://www.youtube.com/shorts/abc123",

			kind:    ReferenceVideo,
			videoID: "abc123",
			ok:      true,
		},
		{
			name: "mobile video",
			raw:  "https://m.youtube.com/watch?v=abc123",

			kind:    ReferenceVideo,
			videoID: "abc123",
			ok:      true,
		},
		{
			name: "playlist",
			raw:  "https://www.youtube.com/playlist?list=PL123",

			kind:       ReferencePlaylist,
			playlistID: "PL123",
			ok:         true,
		},
		{
			name: "video in playlist",
			raw:  "https://www.youtube.com/watch?v=abc123&list=PL123",

			kind:       ReferenceVideo,
			videoID:    "abc123",
			playlistID: "PL123",
			ok:         true,
		},
		{
			name: "unsupported host",
			raw:  "https://example.com/watch?v=abc123",

			ok: false,
		},
		{
			name: "missing identity",
			raw:  "https://www.youtube.com/",

			ok: false,
		},
		{
			name: "shorts missing video ID",
			raw:  "https://www.youtube.com/shorts/",

			ok: false,
		},
		{
			name: "empty",
			raw:  "",

			ok: false,
		},
	}

	for _, test := range tests {

		t.Run(
			test.name,
			func(t *testing.T) {

				ref, ok := ParseURL(
					test.raw,
				)

				if ok != test.ok {

					t.Fatalf(
						"expected ok %v, got %v",
						test.ok,
						ok,
					)

				}

				if !test.ok {
					return
				}

				if ref.Kind != test.kind {

					t.Errorf(
						"expected kind %v, got %v",
						test.kind,
						ref.Kind,
					)

				}

				if ref.VideoID !=
					test.videoID {

					t.Errorf(
						"expected video ID %q, got %q",
						test.videoID,
						ref.VideoID,
					)

				}

				if ref.PlaylistID !=
					test.playlistID {

					t.Errorf(
						"expected playlist ID %q, got %q",
						test.playlistID,
						ref.PlaylistID,
					)

				}

				if ref.URI != test.raw {

					t.Errorf(
						"expected URI %q, got %q",
						test.raw,
						ref.URI,
					)

				}

			},
		)

	}

}
