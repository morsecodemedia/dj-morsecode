package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/musicbrainz"
)

func main() {

	verbose := flag.Bool(
		"verbose",
		false,
		"show MusicBrainz lookup details",
	)

	flag.Parse()

	if flag.NArg() != 2 {

		fmt.Fprintf(
			os.Stderr,
			"usage: %s [-verbose] <artist> <title>\n",
			os.Args[0],
		)

		os.Exit(2)
	}

	artist := flag.Arg(0)
	title := flag.Arg(1)

	client := musicbrainz.NewClient(
		"DJ MorseCode/1.0.0 (https://github.com/morsecodemedia/dj-morsecode)",
	)

	recordings, err := client.SearchRecordings(
		context.Background(),
		artist,
		title,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("QUERY")
	fmt.Printf(
		"%s - %s\n\n",
		artist,
		title,
	)

	fmt.Printf(
		"CANDIDATES (%d)\n",
		len(recordings),
	)

	candidates := make(
		[]metadata.EnrichmentCandidate,
		0,
		len(recordings),
	)

	for _, recording := range recordings {

		candidate := musicbrainz.Candidate(
			recording,
		)

		candidates = append(
			candidates,
			candidate,
		)

		fmt.Printf(
			"%3.0f%%  %-10s  %s - %s\n",
			candidate.ProviderScore*100,
			candidate.Duration,
			candidate.Artist,
			candidate.Title,
		)

		fmt.Printf(
			"      MBID: %s\n",
			recording.ID,
		)

		if candidate.Variant != "" {

			fmt.Printf(
				"      Variant: %s\n",
				candidate.Variant,
			)

		}

		if !*verbose {
			continue
		}

		details, err := client.LookupRecording(
			context.Background(),
			recording.ID,
		)
		if err != nil {

			fmt.Printf(
				"      Lookup failed: %v\n",
				err,
			)

			continue
		}

		for _, isrc := range details.ISRCs {

			fmt.Printf(
				"      ISRC: %s\n",
				isrc,
			)

		}

		for _, release := range details.Releases {

			fmt.Printf(
				"      Release: %s | %s | %s | %s\n",
				release.Title,
				release.Date,
				release.Country,
				release.ID,
			)

		}

	}

	observed := metadata.PlaybackItem{
		Type:   metadata.PlaybackTrack,
		Artist: artist,
		Title:  title,
	}

	match, status := metadata.MatchEnrichment(
		observed,
		candidates,
	)

	fmt.Println()
	fmt.Println("RESULT")

	switch status {

	case metadata.MatchAccepted:

		fmt.Println("ACCEPTED")

		fmt.Printf(
			"%s - %s\n",
			match.Track.Artist,
			match.Track.Title,
		)

		fmt.Printf(
			"Provider: %s\n",
			match.Provider,
		)

		fmt.Printf(
			"Confidence: %.0f%%\n",
			match.Confidence*100,
		)

		for _, identifier := range match.Track.Identifiers {

			fmt.Printf(
				"%s: %s\n",
				identifier.Scheme,
				identifier.Value,
			)

		}

	case metadata.MatchAmbiguous:

		fmt.Println("AMBIGUOUS")

	default:

		fmt.Println("NO MATCH")

	}

}
