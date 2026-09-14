package radio

type Station struct {
	ID   string
	Name string

	StreamURL string

	Homepage string

	Genre string

	Tags []string
}

var Stations = []Station{
	{
		ID:        "anonradio",
		Name:      "AnonRadio",
		StreamURL: "http://anonradio.net:8000/anonradio",
	},
	{
		ID:        "classicfm",
		Name:      "Classic FM",
		StreamURL: "https://www.globalplayer.com/live/classicfm/uk",
	},
	{
		ID:        "evergreen",
		Name:      "Evergreen",
		StreamURL: "https://emg.streamguys1.com/evergreen-website",
	},
	{
		ID:        "sleepbot",
		Name:      "Sleepbot Ambient",
		StreamURL: "http://www.sleepbot.com/ambience/cgi/listen.cgi/listen.pls",
	},
	{
		ID:        "thechristmasstation",
		Name:      "The Christmas Station",
		StreamURL: "https://stream.radio.co/s63541ce9b/listen",
	},
	{
		ID:        "tilderadio",
		Name:      "Tilde Radio",
		StreamURL: "https://azuracast.tilderadio.org/radio/8000/320k.ogg",
	},
	{
		ID:        "wrti",
		Name:      "WRTI Classical",
		StreamURL: "https://wrti-live.streamguys1.com/classical-mp3",
	},
	{
		ID:        "wxpn",
		Name:      "WXPN",
		StreamURL: "https://wxpn.xpn.org/xpnmp3hi",
	},
	{
		ID:        "wxpn2",
		Name:      "WXPN2",
		StreamURL: "https://wxpn.xpn.org/xpn2mp3hi",
	},
}

func Find(id string) (*Station, bool) {

	for i := range Stations {

		if Stations[i].ID == id {
			return &Stations[i], true
		}

	}

	return nil, false

}
