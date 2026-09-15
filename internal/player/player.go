package player

import (
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/player/mpv"
)

type Player struct {
	client *mpv.Client
}

func New(socket string) (*Player, error) {

	client, err := mpv.Connect(socket)
	if err != nil {
		return nil, err
	}

	return &Player{
		client: client,
	}, nil

}

func (p *Player) Title() string {

	title, err := p.client.MediaTitle()
	if err != nil {
		return ""
	}

	return title

}

func (p *Player) Artist() string {

	artist, err := p.client.Artist()
	if err != nil {
		return ""
	}

	return artist

}

func (p *Player) TrackTitle() string {

	title, err := p.client.Title()
	if err != nil {
		return ""
	}

	return title

}

func (p *Player) Album() string {

	album, err := p.client.Album()
	if err != nil {
		return ""
	}

	return album

}

func (p *Player) Position() time.Duration {

	position, err := p.client.PlaybackTime()
	if err != nil {
		return 0
	}

	return position

}

func (p *Player) Duration() time.Duration {

	duration, err := p.client.Duration()
	if err != nil {
		return 0
	}

	return duration

}
