package library

import (
	"fmt"

	"github.com/morsecodemedia/dj-morsecode/internal/lyrics"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
)

var demoLibrary = map[string]string{
	"Blast":             "assets/blast.lrc",
	"Snap Your Fingers": "assets/snap-your-fingers.lrc",
}

func Load(artist, title string) (music.Song, bool) {

	path, ok := demoLibrary[title]

	if ok {

		song, err := lyrics.LoadSong(path)
		if err != nil {
			fmt.Println(err)
			return music.Song{}, false
		}

		return song, true

	}

	return LoadCached(
		artist,
		title,
	)

}
