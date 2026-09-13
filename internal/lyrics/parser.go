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

type Metadata struct {
	Title    string
	Artist   string
	Album    string
	Duration time.Duration
}

var (
	metadataPattern  = regexp.MustCompile(`^\[[a-zA-Z]+:.*\]$`)
	lyricPattern     = regexp.MustCompile(`^\[\d{2}:\d{2}\.\d{2}\]\s*.+`)
	timestampPattern = regexp.MustCompile(`^\[\d{2}:\d{2}\.\d{2}\]$`)
)

func Load(path string) ([]string, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), "\n"), nil
}

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

func ParseMetadata(lines []string) Metadata {

	values := make(map[string]string)

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

		values[parts[0]] = parts[1]
	}

	duration, err := ParseDuration(values["length"])

	if err != nil {
		duration = 0
	}

	return Metadata{
		Artist:   values["ar"],
		Title:    values["ti"],
		Album:    values["al"],
		Duration: duration,
	}
}

func ParseTimeline(lines []string) []music.Cue {

	var timeline []music.Cue

	for _, line := range lines {

		switch {

		case IsLyric(line):

			parts := strings.SplitN(line, "]", 2)

			if len(parts) != 2 {
				continue
			}

			duration, err := ParseTimestamp(
				strings.TrimPrefix(parts[0], "["),
			)

			if err != nil {
				continue
			}

			timeline = append(timeline, music.Cue{
				Time: duration,
				Type: music.CueLyric,
				Text: strings.TrimSpace(parts[1]),
			})

		case IsLyricBreak(line):

			duration, err := ParseTimestamp(
				strings.Trim(line, "[]"),
			)

			if err != nil {
				continue
			}

			timeline = append(timeline, music.Cue{
				Time: duration,
				Type: music.CueBreak,
			})
		}
	}

	return timeline
}

func ParseTimestamp(value string) (time.Duration, error) {

	parts := strings.Split(value, ":")

	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid timestamp: %s", value)
	}

	secondParts := strings.Split(parts[1], ".")

	if len(secondParts) != 2 {
		return 0, fmt.Errorf("invalid timestamp: %s", value)
	}

	minutes, err := strconv.Atoi(parts[0])
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

	duration :=
		time.Duration(minutes)*time.Minute +
			time.Duration(seconds)*time.Second +
			time.Duration(hundredths)*10*time.Millisecond

	return duration, nil
}

func ParseDuration(value string) (time.Duration, error) {

	parts := strings.Split(value, ":")

	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid duration: %s", value)
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	duration :=
		time.Duration(minutes)*time.Minute +
			time.Duration(seconds)*time.Second

	return duration, nil

}
