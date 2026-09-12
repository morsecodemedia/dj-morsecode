package lyrics

import (
	"os"
	"regexp"
	"strings"
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

func Load(path string) ([]string, error) {

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), "\n"), nil

}
