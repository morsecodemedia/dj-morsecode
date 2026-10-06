package youtube

import (
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

func TestMediaItemVideo(
	t *testing.T,
) {

	ref := Reference{
		Kind:    ReferenceVideo,
		VideoID: "abc123",
		URI:     "https://www.youtube.com/watch?v=abc123",
	}

	item, ok := MediaItem(ref)

	if !ok {
		t.Fatal("expected media item")
	}

	if item.Kind != source.MediaVideo {
		t.Errorf(
			"expected video kind, got %q",
			item.Kind,
		)
	}

	if item.Ref.Source.Kind !=
		source.KindYouTube {

		t.Errorf(
			"expected YouTube provenance, got %q",
			item.Ref.Source.Kind,
		)
	}

	if item.Ref.ID != "abc123" {
		t.Errorf(
			"unexpected video ID %q",
			item.Ref.ID,
		)
	}

}

func TestMediaItemPlaylist(
	t *testing.T,
) {

	ref := Reference{
		Kind:       ReferencePlaylist,
		PlaylistID: "PL123",
		URI:        "https://www.youtube.com/playlist?list=PL123",
	}

	item, ok := MediaItem(ref)

	if !ok {
		t.Fatal("expected media item")
	}

	if item.Kind != source.MediaPlaylist {
		t.Errorf(
			"expected playlist kind, got %q",
			item.Kind,
		)
	}

	if item.Ref.ID != "PL123" {
		t.Errorf(
			"unexpected playlist ID %q",
			item.Ref.ID,
		)
	}

}

func TestMediaItemRejectsUnknownReference(
	t *testing.T,
) {

	_, ok := MediaItem(
		Reference{},
	)

	if ok {
		t.Fatal(
			"expected unknown reference to be rejected",
		)
	}

}
