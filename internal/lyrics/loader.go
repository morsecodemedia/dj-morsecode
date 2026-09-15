package lyrics

import (
	"strings"

	"github.com/morsecodemedia/dj-morsecode/internal/music"
)

func LoadSong(path string) (music.Song, error) {

	lines, err := Load(path)
	if err != nil {
		return music.Song{}, err
	}

	return FromLines(lines), nil

}

func FromString(content string) music.Song {

	lines := strings.Split(
		strings.ReplaceAll(content, "\r\n", "\n"),
		"\n",
	)

	return FromLines(lines)

}

func FromLines(lines []string) music.Song {

	metadata := ParseMetadata(lines)

	timeline := ParseTimeline(lines)

	return music.Song{
		Title:    metadata.Title,
		Artist:   metadata.Artist,
		Album:    metadata.Album,
		Duration: metadata.Duration,
		Timeline: timeline,
	}

}
