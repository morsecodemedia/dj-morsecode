package app

import "github.com/morsecodemedia/dj-morsecode/internal/music"

type Model struct {
	Song music.Song

	CurrentLine int

	Width  int
	Height int

	OnAir bool
}
