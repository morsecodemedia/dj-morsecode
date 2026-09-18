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
		ID:        "z100",
		Name:      "Z100",
		StreamURL: "https://stream.revma.ihrhls.com/zc1469",
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
	{
		ID:        "classicvinyl",
		Name:      "Classic Vinyl",
		StreamURL: "https://icecast.walmradio.com:8443/classic_opus",
	},
	{
		ID:        "hardrockradiofm",
		Name:      "Hard Rock Radio FM",
		StreamURL: "http://67.249.184.45:8015/listen.pls",
	},
	{
		ID:        "lofi247",
		Name:      "Lofi 24/7",
		StreamURL: "http://usa9.fastcast4u.com/proxy/jamz?mp=/1",
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

func FindByStreamURL(streamURL string) (*Station, bool) {

	for i := range Stations {

		if Stations[i].StreamURL == streamURL {
			return &Stations[i], true
		}

	}

	return nil, false

}
