package lyrics

import (
	"os"
	"regexp"
	"strings"

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

func Load(path string) ([]string, error) {

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), "\n"), nil

}
