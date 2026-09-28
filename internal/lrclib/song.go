package lrclib

import (
	"strings"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/lyrics"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
)

func Song(result Result) music.Song {

	content := result.SyncedLyrics

	if strings.TrimSpace(
		content,
	) == "" {

		content = result.PlainLyrics
	}

	song := lyrics.FromString(
		content,
	)

	song.Title = result.TrackName
	song.Artist = result.ArtistName
	song.Album = result.AlbumName
	song.Duration = time.Duration(
		result.Duration * float64(time.Second),
	)

	return song

}
