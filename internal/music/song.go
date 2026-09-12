package music

import "time"

type Lyric struct {
	Time time.Duration
	Text string
}

type Song struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length string

	Lyrics []Lyric
}
