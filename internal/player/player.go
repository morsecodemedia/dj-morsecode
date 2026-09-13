package player

import "time"

type Player struct {
	started time.Time
}

func New() *Player {

	return &Player{
		started: time.Now(),
	}

}

func (p *Player) Elapsed() time.Duration {

	return time.Since(p.started)

}
