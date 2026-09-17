package library

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCachePath(t *testing.T) {

	root := filepath.Join(
		"tmp",
		"dj-morsecode",
	)

	path := cachePath(
		root,
		"Stone Temple Pilots",
		"Interstate Love Song",
	)

	if !strings.HasPrefix(
		path,
		filepath.Join(root, "lyrics"),
	) {
		t.Errorf(
			"expected path under lyrics directory, got %q",
			path,
		)
	}

	if filepath.Ext(path) != ".lrc" {
		t.Errorf(
			"expected .lrc extension, got %q",
			filepath.Ext(path),
		)
	}

	filename := strings.TrimSuffix(
		filepath.Base(path),
		".lrc",
	)

	if len(filename) != 64 {
		t.Errorf(
			"expected 64-character SHA-256 key, got %d",
			len(filename),
		)
	}

}

func TestCacheKeyIsDeterministic(t *testing.T) {

	first := cacheKey(
		"Stone Temple Pilots",
		"Interstate Love Song",
	)

	second := cacheKey(
		"Stone Temple Pilots",
		"Interstate Love Song",
	)

	if first != second {
		t.Errorf(
			"expected identical cache keys, got %q and %q",
			first,
			second,
		)
	}

}

func TestCacheKeyNormalizesCaseAndWhitespace(t *testing.T) {

	first := cacheKey(
		"Stone Temple Pilots",
		"Interstate Love Song",
	)

	second := cacheKey(
		"  stone temple pilots  ",
		"  INTERSTATE LOVE SONG ",
	)

	if first != second {
		t.Errorf(
			"expected normalized values to produce same key, got %q and %q",
			first,
			second,
		)
	}

}

func TestCacheKeySeparatesArtistAndTitle(t *testing.T) {

	first := cacheKey(
		"ab",
		"c",
	)

	second := cacheKey(
		"a",
		"bc",
	)

	if first == second {
		t.Errorf(
			"expected different cache keys, both were %q",
			first,
		)
	}

}

func TestCachePathHandlesFilesystemCharacters(t *testing.T) {

	root := filepath.Join(
		"tmp",
		"dj-morsecode",
	)

	path := cachePath(
		root,
		"AC/DC",
		"It's a Long Way to the Top (If You Wanna Rock 'n' Roll)",
	)

	expectedDir := filepath.Join(
		root,
		"lyrics",
	)

	if filepath.Dir(path) != expectedDir {
		t.Errorf(
			"expected cache file directly under %q, got %q",
			expectedDir,
			path,
		)
	}

}

func TestCacheRoot(t *testing.T) {

	root, err := CacheRoot()
	if err != nil {
		t.Fatalf(
			"CacheRoot returned error: %v",
			err,
		)
	}

	if filepath.Base(root) != "dj-morsecode" {
		t.Errorf(
			"expected cache root to end with dj-morsecode, got %q",
			root,
		)
	}

}

func TestStoreAt(t *testing.T) {

	root := t.TempDir()

	content := `[00:13.00]Waiting on a Sunday afternoon
[00:18.00]For what I read between the lines`

	path, err := StoreAt(
		root,
		"Stone Temple Pilots",
		"Interstate Love Song",
		content,
	)
	if err != nil {
		t.Fatalf(
			"StoreAt returned error: %v",
			err,
		)
	}

	expectedPath := cachePath(
		root,
		"Stone Temple Pilots",
		"Interstate Love Song",
	)

	if path != expectedPath {
		t.Errorf(
			"expected path %q, got %q",
			expectedPath,
			path,
		)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf(
			"read cached lyrics: %v",
			err,
		)
	}

	if string(data) != content {
		t.Errorf(
			"expected cached content %q, got %q",
			content,
			string(data),
		)
	}

}

func TestStoreAndLoadCachedAt(t *testing.T) {

	root := t.TempDir()

	content := `[ti:Interstate Love Song]
[ar:Stone Temple Pilots]
[al:Purple]
[length:03:16]
[00:13.00]Waiting on a Sunday afternoon
[00:18.00]For what I read between the lines`

	_, err := StoreAt(
		root,
		"Stone Temple Pilots",
		"Interstate Love Song",
		content,
	)
	if err != nil {
		t.Fatalf(
			"StoreAt returned error: %v",
			err,
		)
	}

	song, ok := LoadCachedAt(
		root,
		"Stone Temple Pilots",
		"Interstate Love Song",
	)
	if !ok {
		t.Fatal("expected cached song to load")
	}

	if song.Title != "Interstate Love Song" {
		t.Errorf(
			"expected title %q, got %q",
			"Interstate Love Song",
			song.Title,
		)
	}

	if song.Artist != "Stone Temple Pilots" {
		t.Errorf(
			"expected artist %q, got %q",
			"Stone Temple Pilots",
			song.Artist,
		)
	}

	if len(song.Timeline) != 2 {
		t.Fatalf(
			"expected 2 timeline cues, got %d",
			len(song.Timeline),
		)
	}

	if song.Timeline[0].Text != "Waiting on a Sunday afternoon" {
		t.Errorf(
			"unexpected first cue text %q",
			song.Timeline[0].Text,
		)
	}

}

func TestLoadCachedAtMiss(t *testing.T) {

	root := t.TempDir()

	_, ok := LoadCachedAt(
		root,
		"Stone Temple Pilots",
		"Interstate Love Song",
	)

	if ok {
		t.Fatal("expected cache miss")
	}

}
