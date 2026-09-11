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

	Lyrics []Lyric
}

var DemoSong = Song{
	Title:  "Interstate Love Song",
	Artist: "Stone Temple Pilots",
	Album:  "Purple",
	Year:   1994,

	Lyrics: []Lyric{
		{Text: "Leaving on a southern train..."},
		{Text: "Only yesterday you lied..."},
		{Text: "Promises of what I seemed..."},
	},
}
