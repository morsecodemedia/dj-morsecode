package lastfm

import (
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

func TestMediaItem(
	t *testing.T,
) {

	track := TrackInfo{
		Artist: "Stone Temple Pilots",
		Title:  "Interstate Love Song",
		MBID:   "recording-123",
		URL: "https://www.last.fm/music/" +
			"Stone+Temple+Pilots/_/" +
			"Interstate+Love+Song",
	}

	item, ok := MediaItem(
		track,
	)

	if !ok {
		t.Fatal(
			"expected media item",
		)
	}

	if item.Kind != source.MediaTrack {

		t.Errorf(
			"expected track kind, got %q",
			item.Kind,
		)

	}

	if item.Artist !=
		"Stone Temple Pilots" {

		t.Errorf(
			"unexpected artist %q",
			item.Artist,
		)

	}

	if item.Title !=
		"Interstate Love Song" {

		t.Errorf(
			"unexpected title %q",
			item.Title,
		)

	}

	if item.Ref.Source.Kind !=
		source.KindLastFM {

		t.Errorf(
			"expected Last.fm provenance, got %q",
			item.Ref.Source.Kind,
		)

	}

	if item.Ref.ID != "recording-123" {

		t.Errorf(
			"unexpected provider ID %q",
			item.Ref.ID,
		)

	}

	if item.Ref.URI != track.URL {

		t.Errorf(
			"unexpected provider URI %q",
			item.Ref.URI,
		)

	}

}

func TestMediaItemAcceptsMBIDWithoutURL(
	t *testing.T,
) {

	track := TrackInfo{
		Artist: "Stone Temple Pilots",
		Title:  "Interstate Love Song",
		MBID:   "recording-123",
	}

	item, ok := MediaItem(
		track,
	)

	if !ok {
		t.Fatal(
			"expected MBID to provide source identity",
		)
	}

	if item.Ref.ID != "recording-123" {
		t.Errorf(
			"unexpected ID %q",
			item.Ref.ID,
		)
	}

}

func TestMediaItemAcceptsURLWithoutMBID(
	t *testing.T,
) {

	track := TrackInfo{
		Artist: "Stone Temple Pilots",
		Title:  "Interstate Love Song",
		URL:    "https://www.last.fm/example",
	}

	item, ok := MediaItem(
		track,
	)

	if !ok {
		t.Fatal(
			"expected URL to provide source identity",
		)
	}

	if item.Ref.URI != track.URL {
		t.Errorf(
			"unexpected URI %q",
			item.Ref.URI,
		)
	}

}

func TestMediaItemRejectsTrackWithoutProviderIdentity(
	t *testing.T,
) {

	track := TrackInfo{
		Artist: "Stone Temple Pilots",
		Title:  "Interstate Love Song",
	}

	_, ok := MediaItem(
		track,
	)

	if ok {
		t.Fatal(
			"expected track without Last.fm identity to be rejected",
		)
	}

}

func TestMediaItemTrimsIdentity(
	t *testing.T,
) {

	track := TrackInfo{
		Artist: "  Stone Temple Pilots  ",
		Title:  "  Interstate Love Song  ",
		MBID:   "  recording-123  ",
		URL:    "  https://www.last.fm/example  ",
	}

	item, ok := MediaItem(
		track,
	)

	if !ok {
		t.Fatal(
			"expected media item",
		)
	}

	if item.Artist !=
		"Stone Temple Pilots" {

		t.Errorf(
			"unexpected artist %q",
			item.Artist,
		)

	}

	if item.Title !=
		"Interstate Love Song" {

		t.Errorf(
			"unexpected title %q",
			item.Title,
		)

	}

	if item.Ref.ID !=
		"recording-123" {

		t.Errorf(
			"unexpected ID %q",
			item.Ref.ID,
		)

	}

	if item.Ref.URI !=
		"https://www.last.fm/example" {

		t.Errorf(
			"unexpected URI %q",
			item.Ref.URI,
		)

	}

}
