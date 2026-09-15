package lyrics

import "github.com/morsecodemedia/dj-morsecode/internal/music"

func LoadSong(path string) (music.Song, error) {

	lines, err := Load(path)
	if err != nil {
		return music.Song{}, err
	}

	metadata := ParseMetadata(lines)

	timeline := ParseTimeline(lines)

	return music.Song{
		Title:    metadata.Title,
		Artist:   metadata.Artist,
		Album:    metadata.Album,
		Duration: metadata.Duration,
		Timeline: timeline,
	}, nil

}
