package lyrics

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/music"
)

var (
	metadataPattern  = regexp.MustCompile(`^\[[a-zA-Z]+:.*\]$`)
	lyricPattern     = regexp.MustCompile(`^\[\d{2}:\d{2}\.\d{2}\]\s*.+`)
	timestampPattern = regexp.MustCompile(`^\[\d{2}:\d{2}(?::\d{2}|\.\d{2})\]$`)
)

func IsMetadata(line string) bool {
	return metadataPattern.MatchString(line)
}

func IsLyric(line string) bool {
	return lyricPattern.MatchString(line)
}

func IsLyricBreak(line string) bool {
	return timestampPattern.MatchString(line)
}

func IsBlank(line string) bool {
	return strings.TrimSpace(line) == ""
}

func ParseSong(lines []string) music.Song {

	metadata := ParseMetadata(lines)

	return music.Song{
		Artist: metadata["ar"],
		Title:  metadata["ti"],
		Album:  metadata["al"],
		Length: metadata["length"],

		Lyrics: ParseLyrics(lines),
	}

}

func ParseMetadata(lines []string) map[string]string {

	metadata := make(map[string]string)

	for _, line := range lines {

		if !IsMetadata(line) {
			continue
		}

		line = strings.TrimPrefix(line, "[")
		line = strings.TrimSuffix(line, "]")

		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		metadata[key] = value

	}

	return metadata

}

func ParseTimestamp(value string) (time.Duration, error) {

	parts := strings.Split(value, ":")

	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid timestamp: %s", value)
	}

	minutePart := parts[0]
	secondPart := parts[1]

	secondParts := strings.Split(secondPart, ".")

	minutes, err := strconv.Atoi(minutePart)
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(secondParts[0])
	if err != nil {
		return 0, err
	}

	hundredths, err := strconv.Atoi(secondParts[1])
	if err != nil {
		return 0, err
	}

	if len(secondParts) != 2 {
		return 0, fmt.Errorf("invalid timestamp: %s", value)
	}

	duration :=
		time.Duration(minutes)*time.Minute +
			time.Duration(seconds)*time.Second +
			time.Duration(hundredths)*10*time.Millisecond

	return duration, nil

}

func ParseLyrics(lines []string) []music.Lyric {

	lyrics := []music.Lyric{}

	for _, line := range lines {

		if !IsLyric(line) {
			continue
		}

		parts := strings.SplitN(line, "]", 2)

		timestamp := strings.TrimPrefix(parts[0], "[")

		lyricText := strings.TrimSpace(parts[1])

		duration, err := ParseTimestamp(timestamp)

		if len(parts) != 2 {
			continue
		}

		if err != nil {
			continue
		}

		lyrics = append(lyrics, music.Lyric{
			Time: duration,
			Text: lyricText,
		})

	}

	return lyrics

}

func Load(path string) ([]string, error) {

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), "\n"), nil

}
