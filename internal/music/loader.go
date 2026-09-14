package music

import "github.com/morsecodemedia/dj-morsecode/internal/lyrics"

func Load(path string) (Song, error) {

	lines, err := lyrics.Load(path)
	if err != nil {
		return Song{}, err
	}

	metadata := lyrics.ParseMetadata(lines)

	timeline := lyrics.ParseTimeline(lines)

	return Song{
		Title:    metadata.Title,
		Artist:   metadata.Artist,
		Album:    metadata.Album,
		Duration: metadata.Duration,
		Timeline: timeline,
	}, nil

}
