package music

import "time"

type CueType int

const (
	CueLyric CueType = iota
	CueBreak
)

type Cue struct {
	Time time.Duration
	Type CueType
	Text string
}

type Song struct {
	Title    string
	Artist   string
	Album    string
	Year     int
	Duration time.Duration

	Timeline []Cue
}
