package main

import (
	"fmt"

	"github.com/morsecodemedia/dj-morsecode/internal/player/mpv"
)

func main() {

	client, err := mpv.Connect("/tmp/dj-morsecode.sock")
	if err != nil {
		panic(err)
	}

	title, _ := client.MediaTitle()

	playback, _ := client.PlaybackTime()

	duration, _ := client.Duration()

	fmt.Println(title)
	fmt.Println(playback)
	fmt.Println(duration)

}
