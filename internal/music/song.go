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
		{Text: "Waitin' on a Sunday afternoon"},
		{Text: "For what I've read between the lines"},
		{Text: "Your lies"},
		{Text: "Feelin' like a hand in rusted shame"},
		{Text: "So do you laugh at those who cry?"},
		{Text: "Reply?"},
		{Text: "Leavin' on a southern train, only yesterday you lied"},
		{Text: "Promises of what I seemed to be, only watched the time go by"},
		{Text: "All of these things you said to me"},
		{Text: "Breathin' is the hardest thing to do"},
		{Text: "With all I've said and all that's dead for you"},
		{Text: "You lied, goodbye"},
		{Text: "Leavin' on a southern train, only yesterday you lied"},
		{Text: "Promises of what I seemed to be, only watched the time go by"},
		{Text: "All of these things I said to you"},
	},
}
