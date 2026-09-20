package radio

type MoodFamily struct {
	ID    string
	Name  string
	Moods []string
}

var MoodFamilies = []MoodFamily{
	{
		ID:   "calm",
		Name: "Calm",
		Moods: []string{
			"calm",
			"chilled",
			"cozy",
			"mellow",
			"peaceful",
			"smooth",
		},
	},
	{
		ID:   "dreamy",
		Name: "Dreamy",
		Moods: []string{
			"dreamy",
			"meditative",
			"otherworldly",
			"sensual",
			"spacious",
		},
	},
	{
		ID:   "energetic",
		Name: "Energetic",
		Moods: []string{
			"energetic",
			"intense",
			"upbeat",
		},
	},
	{
		ID:   "uplifting",
		Name: "Uplifting",
		Moods: []string{
			"euphoric",
			"positive",
			"sunny",
			"uplifting",
		},
	},
	{
		ID:   "adventurous",
		Name: "Adventurous",
		Moods: []string{
			"adventurous",
			"curious",
			"varied",
		},
	},
	{
		ID:   "edgy",
		Name: "Edgy",
		Moods: []string{
			"aggressive",
			"raw",
			"rebellious",
		},
	},
	{
		ID:   "nostalgic",
		Name: "Nostalgic",
		Moods: []string{
			"familiar",
			"nostalgic",
		},
	},
	{
		ID:   "playful",
		Name: "Playful",
		Moods: []string{
			"groovy",
			"playful",
		},
	},
	{
		ID:   "thoughtful",
		Name: "Thoughtful",
		Moods: []string{
			"independent",
			"thoughtful",
		},
	},
}

func FindMoodFamily(id string) (*MoodFamily, bool) {

	for i := range MoodFamilies {

		if MoodFamilies[i].ID == id {
			return &MoodFamilies[i], true
		}

	}

	return nil, false

}

func (m MoodFamily) Criteria() Criteria {

	return Criteria{
		Moods: m.Moods,
	}

}
