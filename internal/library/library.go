package library

import "github.com/morsecodemedia/dj-morsecode/internal/music"

func Load(
	artist string,
	title string,
) (music.Song, bool) {

	return LoadCached(
		artist,
		title,
	)

}
