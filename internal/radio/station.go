package radio

type Station struct {
	ID   string
	Name string

	StreamURL string
	Homepage  string

	Genre string
	Tags  []string

	Moods    []string
	Contexts []string

	Energy int
}

/**
 *
 * GENRE
 * Primary musical family
 *
**/

/**
 *
 * TAGS
 * Musical/style characteristics
 *
**/

/**
 *
 * ENERGY SCALE
 * 1 = extremely low / atmospheric
 * 2 = mellow / relaxed
 * 3 = moderate / variable
 * 4 = energetic
 * 5 = intense / peak-energy
 *
**/

/**
 *
 * MOODS TAXONOMY - What do I want this to fee like?
 * adventurous
 * aggressive
 * calm
 * chilled
 * cozy
 * curious
 * dreamy
 * energetic
 * euphoric
 * familiar
 * groovy
 * independent
 * intense
 * meditative
 * mellow
 * nostalgic
 * otherworldly
 * peaceful
 * playful
 * positive
 * raw
 * rebellious
 * sensual
 * smooth
 * spacious
 * sunny
 * thoughtful
 * upbeat
 * uplifting
 * varied
 *
**/

/**
 *
 * CONTEXTS TAXONOMY - What am I doing?
 * active
 * background
 * casual
 * coding
 * discovery
 * driving
 * focus
 * gaming
 * late-night
 * meditation
 * morning
 * party
 * relaxing
 * sleep
 * social
 * workout
 *
**/

var Stations = []Station{
	{
		ID:        "z100",
		Name:      "Z100",
		StreamURL: "https://stream.revma.ihrhls.com/zc1469",
		Homepage:  "https://z100.iheart.com/",
		Genre:     "Pop",
		Tags: []string{
			"pop",
			"contemporary hits",
			"mainstream",
		},
		Moods: []string{
			"upbeat",
			"energetic",
			"familiar",
		},
		Contexts: []string{
			"casual",
			"social",
			"driving",
		},
		Energy: 4,
	},
	{
		ID:        "wxpn",
		Name:      "WXPN",
		StreamURL: "https://wxpn.xpn.org/xpnmp3hi",
		Homepage:  "",
		Genre:     "Eclectic",
		Tags: []string{
			"adult alternative",
			"indie",
			"rock",
			"roots",
			"folk",
			"discovery",
		},
		Moods: []string{
			"curious",
			"varied",
			"thoughtful",
		},
		Contexts: []string{
			"discovery",
			"casual",
			"background",
		},
		Energy: 3,
	},
	{
		ID:        "wxpn2",
		Name:      "WXPN2",
		StreamURL: "https://wxpn.xpn.org/xpn2mp3hi",
		Homepage:  "",
		Genre:     "Eclectic",
		Tags: []string{
			"indie",
			"rock",
			"soul",
			"folk",
			"americana",
			"deep cuts",
		},
		Moods: []string{
			"curious",
			"varied",
			"adventurous",
		},
		Contexts: []string{
			"discovery",
			"background",
			"casual",
		},
		Energy: 3,
	},
	{
		ID:        "hardrockradiofm",
		Name:      "Hard Rock Radio FM",
		StreamURL: "http://67.249.184.45:8015/listen.pls",
		Homepage:  "",
		Genre:     "Hard Rock",
		Tags: []string{
			"hard rock",
			"classic hard rock",
			"metal",
		},
		Moods: []string{
			"aggressive",
			"energetic",
			"intense",
		},
		Contexts: []string{
			"active",
			"driving",
			"workout",
		},
		Energy: 5,
	},
	{
		ID:        "lofi247",
		Name:      "Lofi 24/7",
		StreamURL: "http://usa9.fastcast4u.com/proxy/jamz?mp=/1",
		Homepage:  "",
		Genre:     "Lo-Fi",
		Tags: []string{
			"lo-fi",
			"beats",
			"chillhop",
		},
		Moods: []string{
			"calm",
			"mellow",
			"cozy",
		},
		Contexts: []string{
			"focus",
			"coding",
			"background",
			"late-night",
		},
		Energy: 1,
	},
	{
		ID:        "realpunkradio",
		Name:      "Real Punk Radio",
		StreamURL: "http://149.56.155.73:8080/stream",
		Homepage:  "https://realpunkradio.com/",
		Genre:     "Punk",
		Tags: []string{
			"punk rock",
			"ska",
			"rockabilly",
			"garage rock",
			"hardcore",
		},
		Moods: []string{
			"rebellious",
			"raw",
			"energetic",
		},
		Contexts: []string{
			"active",
			"driving",
			"discovery",
		},
		Energy: 5,
	},
	{
		ID:        "radioparadise",
		Name:      "Radio Paradise",
		StreamURL: "http://stream-dc1.radioparadise.com/rp_192m.ogg",
		Homepage:  "",
		Genre:     "Eclectic",
		Tags: []string{
			"rock",
			"world",
			"electronic",
			"jazz",
			"curated",
		},
		Moods: []string{
			"adventurous",
			"varied",
			"thoughtful",
		},
		Contexts: []string{
			"discovery",
			"background",
			"casual",
		},
		Energy: 3,
	},
	{
		ID:        "yachtrockmiami",
		Name:      "Yacht Rock Miami",
		StreamURL: "https://usa20.fastcast4u.com:4100/1753014835",
		Homepage:  "",
		Genre:     "Yacht Rock",
		Tags: []string{
			"soft rock",
			"yacht rock",
			"70s",
			"80s",
		},
		Moods: []string{
			"smooth",
			"nostalgic",
			"sunny",
		},
		Contexts: []string{
			"relaxing",
			"casual",
			"background",
		},
		Energy: 2,
	},
	{
		ID:        "1mixradio",
		Name:      "1Mix Radio",
		StreamURL: "http://fr5.1mix.co.uk:8060/256",
		Homepage:  "",
		Genre:     "Trance",
		Tags: []string{
			"trance",
			"edm",
			"progressive",
			"dj sets",
		},
		Moods: []string{
			"uplifting",
			"energetic",
			"euphoric",
		},
		Contexts: []string{
			"active",
			"coding",
			"driving",
			"late-night",
		},
		Energy: 4,
	},
	{
		ID:        "radiocaroline",
		Name:      "Radio Caroline",
		StreamURL: "http://78.129.202.10:8030/",
		Homepage:  "",
		Genre:     "Rock",
		Tags: []string{
			"album rock",
			"classic rock",
			"alternative",
			"deep cuts",
		},
		Moods: []string{
			"nostalgic",
			"varied",
			"independent",
		},
		Contexts: []string{
			"discovery",
			"casual",
			"background",
		},
		Energy: 3,
	},
	{
		ID:        "slayradio",
		Name:      "SLAY Radio",
		StreamURL: "http://relay1.slayradio.org:8000/",
		Homepage:  "https://www.slayradio.org/news",
		Genre:     "Chiptune",
		Tags: []string{
			"c64",
			"amiga",
			"game music",
			"remixes",
			"retro",
		},
		Moods: []string{
			"playful",
			"nostalgic",
			"energetic",
		},
		Contexts: []string{
			"coding",
			"gaming",
			"discovery",
		},
		Energy: 4,
	},
	{
		ID:        "somafm-lush",
		Name:      "SomaFM - Lush",
		StreamURL: "http://ice2.somafm.com/lush-128-aac",
		Homepage:  "",
		Genre:     "Downtempo",
		Tags: []string{
			"downtempo",
			"electronica",
			"female vocals",
			"chillout",
		},
		Moods: []string{
			"dreamy",
			"sensual",
			"mellow",
		},
		Contexts: []string{
			"relaxing",
			"background",
			"late-night",
		},
		Energy: 2,
	},
	{
		ID:        "somafm-deepspaceone",
		Name:      "SomaFM - Deep Space One",
		StreamURL: "http://ice2.somafm.com/deepspaceone-128-aac",
		Homepage:  "",
		Genre:     "Ambient",
		Tags: []string{
			"deep ambient",
			"experimental",
			"space music",
			"electronic",
		},
		Moods: []string{
			"spacious",
			"meditative",
			"otherworldly",
		},
		Contexts: []string{
			"focus",
			"background",
			"meditation",
			"late-night",
		},
		Energy: 1,
	},
	{
		ID:        "somafm-dronezone",
		Name:      "SomaFM - Drone Zone",
		StreamURL: "http://ice2.somafm.com/dronezone-128-aac",
		Homepage:  "",
		Genre:     "Ambient",
		Tags: []string{
			"drone",
			"atmospheric",
			"minimal",
			"instrumental",
		},
		Moods: []string{
			"calm",
			"spacious",
			"meditative",
		},
		Contexts: []string{
			"focus",
			"background",
			"meditation",
			"sleep",
		},
		Energy: 1,
	},
	{
		ID:        "somafm-groovesalad",
		Name:      "SomaFM - Groove Salad",
		StreamURL: "http://ice2.somafm.com/groovesalad-128-aac",
		Homepage:  "",
		Genre:     "Downtempo",
		Tags: []string{
			"downtempo",
			"ambient",
			"electronic",
			"beats",
		},
		Moods: []string{
			"chilled",
			"mellow",
			"groovy",
		},
		Contexts: []string{
			"focus",
			"coding",
			"background",
			"relaxing",
		},
		Energy: 2,
	},
	{
		ID:        "amambient",
		Name:      "A.M. Ambient",
		StreamURL: "http://radio.stereoscenic.com/ama-h",
		Homepage:  "https://amambient.com/",
		Genre:     "Ambient",
		Tags: []string{
			"ambient",
			"instrumental",
			"bright ambient",
		},
		Moods: []string{
			"calm",
			"positive",
			"peaceful",
		},
		Contexts: []string{
			"focus",
			"morning",
			"background",
			"relaxing",
		},
		Energy: 1,
	},
	{
		ID:        "afterhoursfm",
		Name:      "After Hours FM",
		StreamURL: "http://nl.ah.fm:8000/live",
		Homepage:  "https://ah.fm/",
		Genre:     "Trance",
		Tags: []string{
			"trance",
			"progressive trance",
			"edm",
			"dj sets",
		},
		Moods: []string{
			"euphoric",
			"energetic",
		},
		Contexts: []string{
			"active",
			"driving",
			"workout",
			"late-night",
		},
		Energy: 5,
	},
	{
		ID:        "skafari",
		Name:      "SKAfari",
		StreamURL: "http://skafari.stream.laut.fm/skafari",
		Homepage:  "https://laut.fm/skafari",
		Genre:     "Ska",
		Tags: []string{
			"ska",
			"ska-punk",
			"punk",
			"reggae",
		},
		Moods: []string{
			"upbeat",
			"playful",
			"energetic",
		},
		Contexts: []string{
			"active",
			"party",
			"driving",
		},
		Energy: 5,
	},
	{
		ID:        "skaworld",
		Name:      "SKA World",
		StreamURL: "http://skaworld.stream.laut.fm/skaworld",
		Homepage:  "https://laut.fm/skaworld",
		Genre:     "Ska",
		Tags: []string{
			"ska",
			"ska-punk",
			"reggae",
			"rocksteady",
		},
		Moods: []string{
			"upbeat",
			"playful",
			"sunny",
		},
		Contexts: []string{
			"casual",
			"party",
			"driving",
		},
		Energy: 4,
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
