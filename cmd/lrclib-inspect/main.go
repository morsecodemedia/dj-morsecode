package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/library"
	"github.com/morsecodemedia/dj-morsecode/internal/lrclib"
)

func main() {

	durationSeconds := flag.Float64(
		"duration",
		0,
		"playback duration in seconds",
	)

	writeCache := flag.Bool(
		"write-cache",
		false,
		"write selected lyrics to the DJ MorseCode cache",
	)

	flag.Parse()

	if flag.NArg() != 2 {

		fmt.Fprintf(
			os.Stderr,
			"usage: %s [-duration seconds] [-write-cache] <artist> <title>\n",
			os.Args[0],
		)

		os.Exit(2)
	}

	artist := flag.Arg(0)
	title := flag.Arg(1)

	duration := time.Duration(
		*durationSeconds *
			float64(time.Second),
	)

	client := lrclib.NewClient()

	results, err := client.Search(
		artist,
		title,
	)
	if err != nil {

		fmt.Fprintf(
			os.Stderr,
			"search failed: %v\n",
			err,
		)

		os.Exit(1)
	}

	fmt.Printf(
		"query: %s - %s\n",
		artist,
		title,
	)

	if duration > 0 {

		fmt.Printf(
			"playback duration: %.3fs\n",
			duration.Seconds(),
		)

	}

	fmt.Printf(
		"results: %d\n\n",
		len(results),
	)

	for i, result := range results {

		resultDuration :=
			time.Duration(
				result.Duration *
					float64(time.Second),
			)

		fmt.Printf(
			"[%d]\n",
			i,
		)

		fmt.Printf(
			"  artist:   %s\n",
			result.ArtistName,
		)

		fmt.Printf(
			"  track:    %s\n",
			result.TrackName,
		)

		fmt.Printf(
			"  album:    %s\n",
			result.AlbumName,
		)

		fmt.Printf(
			"  duration: %.3fs\n",
			result.Duration,
		)

		if duration > 0 {

			delta :=
				resultDuration - duration

			if delta < 0 {
				delta = -delta
			}

			fmt.Printf(
				"  delta:    %.3fs\n",
				delta.Seconds(),
			)

		}

		fmt.Printf(
			"  synced:   %t\n",
			result.SyncedLyrics != "",
		)

		fmt.Printf(
			"  plain:    %t\n\n",
			result.PlainLyrics != "",
		)

	}

	match, ok :=
		lrclib.BestMatch(
			results,
			duration,
		)

	fmt.Printf(
		"BestMatch: %t\n",
		ok,
	)

	if !ok {
		return
	}

	fmt.Printf(
		"Selected: %s - %s (%.3fs)\n",
		match.ArtistName,
		match.TrackName,
		match.Duration,
	)

	if !*writeCache {
		return
	}

	content := match.SyncedLyrics

	if content == "" {
		content = match.PlainLyrics
	}

	path, err := library.Store(
		artist,
		title,
		content,
	)
	if err != nil {

		fmt.Printf(
			"Cache write: FAILED: %v\n",
			err,
		)

		return
	}

	fmt.Println(
		"Cache write: OK",
	)

	fmt.Printf(
		"Cache path: %s\n",
		path,
	)

}
