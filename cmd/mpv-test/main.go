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

	title, err := client.MediaTitle()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected!")
	fmt.Println(title)

}
