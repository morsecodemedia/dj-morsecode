package m3u

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

type Entry struct {
	Ref source.ItemRef

	Duration int
}

func Parse(
	reader io.Reader,
	origin source.Source,
) ([]Entry, error) {

	scanner := bufio.NewScanner(
		reader,
	)

	var entries []Entry
	var pendingName string
	var pendingDuration int

	for scanner.Scan() {

		line := strings.TrimSpace(
			scanner.Text(),
		)

		if line == "" {
			continue
		}

		if strings.HasPrefix(
			line,
			"#EXTINF:",
		) {

			pendingDuration,
				pendingName =
				parseEXTINF(line)

			continue
		}

		if strings.HasPrefix(
			line,
			"#",
		) {

			continue
		}

		entries = append(
			entries,
			Entry{
				Ref: source.ItemRef{
					Source: origin,
					URI:    line,
					Name:   pendingName,
				},
				Duration: pendingDuration,
			},
		)

		pendingName = ""
		pendingDuration = 0

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil

}

func parseEXTINF(
	line string,
) (int, string) {

	value := strings.TrimPrefix(
		line,
		"#EXTINF:",
	)

	parts := strings.SplitN(
		value,
		",",
		2,
	)

	duration, err := strconv.Atoi(
		strings.TrimSpace(
			parts[0],
		),
	)
	if err != nil {
		duration = 0
	}

	name := ""

	if len(parts) == 2 {
		name = strings.TrimSpace(
			parts[1],
		)
	}

	return duration, name

}
