package controls

type Mode int

const (
	ModeNone Mode = iota
	ModeControls
	ModeTune
	ModeEnhancements
	ModeInfo
	ModeVolume
)
