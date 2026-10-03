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

func (p *Player) Load(path string) error {

	return p.client.Load(path)

}

func (p *Player) Filename() string {

	filename, err := p.client.Filename()
	if err != nil {
		return ""
	}

	return filename

}

func (p *Player) Path() string {

	path, err := p.client.Path()
	if err != nil {
		return ""
	}

	return path

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

func (p *Player) TrackID() string {

	path := p.Path()
	title := p.Title()

	if path == "" {
		return title
	}

	return path + "\x00" + title

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

func (p *Player) Metadata() map[string]string {

	metadata, err := p.client.Metadata()
	if err != nil {
		return nil
	}

	return metadata

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

func (p *Player) IsNetwork() bool {

	network, err := p.client.DemuxerViaNetwork()
	if err != nil {
		return false
	}

	return network

}

func (p *Player) IsIdle() bool {

	idle, err := p.client.CoreIdle()
	if err != nil {
		return false
	}

	return idle

}

func (p *Player) Paused() bool {

	paused, err := p.client.Paused()
	if err != nil {
		return false
	}

	return paused

}

func (p *Player) SetPaused(
	paused bool,
) error {

	return p.client.SetPaused(
		paused,
	)

}

func (p *Player) Muted() bool {

	muted, err := p.client.Muted()
	if err != nil {
		return false
	}

	return muted

}

func (p *Player) SetMuted(
	muted bool,
) error {

	return p.client.SetMuted(
		muted,
	)

}

func (p *Player) Volume() float64 {

	volume, err := p.client.Volume()
	if err != nil {
		return 0
	}

	return volume

}

func (p *Player) SetVolume(
	volume float64,
) error {

	return p.client.SetVolume(
		volume,
	)

}
