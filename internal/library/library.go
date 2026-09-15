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

func Load(rawTitle string) (music.Song, bool) {

	path, ok := demoLibrary[rawTitle]
	if !ok {
		return music.Song{}, false
	}

	song, err := lyrics.LoadSong(path)
	if err != nil {
		fmt.Println(err)
		return music.Song{}, false
	}

	return song, true

}
