package library

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/morsecodemedia/dj-morsecode/internal/lyrics"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
)

func cachePath(
	root string,
	artist string,
	title string,
) string {

	key := cacheKey(
		artist,
		title,
	)

	return filepath.Join(
		root,
		"lyrics",
		key+".lrc",
	)

}

func cacheKey(
	artist string,
	title string,
) string {

	artist = normalizeCacheValue(artist)
	title = normalizeCacheValue(title)

	hash := sha256.Sum256(
		[]byte(artist + "\x00" + title),
	)

	return hex.EncodeToString(hash[:])

}

func normalizeCacheValue(value string) string {

	return strings.ToLower(
		strings.TrimSpace(value),
	)

}

func CacheRoot() (string, error) {

	root, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf(
			"resolve user cache directory: %w",
			err,
		)
	}

	return filepath.Join(
		root,
		"dj-morsecode",
	), nil

}

func StoreAt(
	root string,
	artist string,
	title string,
	content string,
) (string, error) {

	path := cachePath(
		root,
		artist,
		title,
	)

	err := os.MkdirAll(
		filepath.Dir(path),
		0o755,
	)
	if err != nil {
		return "", fmt.Errorf(
			"create lyrics cache directory: %w",
			err,
		)
	}

	err = os.WriteFile(
		path,
		[]byte(content),
		0o644,
	)
	if err != nil {
		return "", fmt.Errorf(
			"write cached lyrics: %w",
			err,
		)
	}

	return path, nil

}

func LoadCachedAt(
	root string,
	artist string,
	title string,
) (music.Song, bool) {

	path := cachePath(
		root,
		artist,
		title,
	)

	song, err := lyrics.LoadSong(path)
	if err != nil {
		return music.Song{}, false
	}

	return song, true

}

func Store(
	artist string,
	title string,
	content string,
) (string, error) {

	root, err := CacheRoot()
	if err != nil {
		return "", err
	}

	return StoreAt(
		root,
		artist,
		title,
		content,
	)

}

func LoadCached(
	artist string,
	title string,
) (music.Song, bool) {

	root, err := CacheRoot()
	if err != nil {
		return music.Song{}, false
	}

	return LoadCachedAt(
		root,
		artist,
		title,
	)

}
