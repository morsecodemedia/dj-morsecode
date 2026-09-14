package metadata

import "strings"

type NowPlaying struct {
	RawTitle string

	Artist string
	Title  string

	Valid bool
}

func Resolve(raw string) NowPlaying {

	nowPlaying := NowPlaying{
		RawTitle: strings.TrimSpace(raw),
	}

	if nowPlaying.RawTitle == "" {
		return nowPlaying
	}

	parts := strings.SplitN(
		nowPlaying.RawTitle,
		" - ",
		2,
	)

	if len(parts) != 2 {
		return nowPlaying
	}

	nowPlaying.Artist = strings.TrimSpace(parts[0])
	nowPlaying.Title = strings.TrimSpace(parts[1])
	nowPlaying.Valid = true

	return nowPlaying

}
