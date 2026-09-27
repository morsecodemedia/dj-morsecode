package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/morsecodemedia/dj-morsecode/internal/lastfm"
)

func main() {

	if len(os.Args) != 3 {

		fmt.Fprintf(
			os.Stderr,
			"usage: %s <artist> <title>\n",
			os.Args[0],
		)

		os.Exit(2)
	}

	apiKey := os.Getenv(
		"LASTFM_API_KEY",
	)

	if apiKey == "" {

		fmt.Fprintln(
			os.Stderr,
			"LASTFM_API_KEY is not set",
		)

		os.Exit(2)
	}

	artist := os.Args[1]
	title := os.Args[2]

	client := lastfm.NewClient(
		apiKey,
	)

	ctx := context.Background()

	fmt.Println("QUERY")
	fmt.Printf(
		"%s - %s\n\n",
		artist,
		title,
	)

	info, err := client.TrackInfo(
		ctx,
		artist,
		title,
	)

	fmt.Println("TRACK INFO")

	if err != nil {

		fmt.Printf(
			"ERROR: %v\n",
			err,
		)

	} else {

		fmt.Printf(
			"Resolved: %s - %s\n",
			info.Artist,
			info.Title,
		)

		if info.MBID != "" {

			fmt.Printf(
				"MBID: %s\n",
				info.MBID,
			)

		} else {

			fmt.Println(
				"MBID: (none)",
			)

		}

		fmt.Printf(
			"Listeners: %d\n",
			info.Listeners,
		)

		fmt.Printf(
			"Playcount: %d\n",
			info.Playcount,
		)

		fmt.Printf(
			"Info tags (%d):\n",
			len(info.Tags),
		)

		for _, tag := range info.Tags {

			fmt.Printf(
				"  %s\n",
				tag.Name,
			)

		}

	}

	fmt.Println()

	tags, err := client.TrackTopTags(
		ctx,
		artist,
		title,
	)

	fmt.Println("TOP TAGS")

	if err != nil {

		fmt.Printf(
			"ERROR: %v\n",
			err,
		)

		return
	}

	fmt.Printf(
		"Count: %d\n",
		len(tags),
	)

	for _, tag := range tags {

		fmt.Printf(
			"%5d  %s\n",
			tag.Count,
			tag.Name,
		)

	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	artistTags, err := client.ArtistTopTags(
		ctx,
		artist,
	)

	fmt.Println("ARTIST TOP TAGS")

	if err != nil {

		fmt.Printf(
			"ERROR: %v\n",
			err,
		)

		return
	}

	fmt.Printf(
		"Count: %d\n",
		len(artistTags),
	)

	for _, tag := range artistTags {

		fmt.Printf(
			"%5d  %s\n",
			tag.Count,
			tag.Name,
		)

	}

}
