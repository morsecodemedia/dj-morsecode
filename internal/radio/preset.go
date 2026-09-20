package radio

type Preset struct {
	ID       string
	Name     string
	Criteria Criteria
}

var Presets = []Preset{
	{
		ID:   "focus",
		Name: "Focus",
		Criteria: Criteria{
			Contexts: []string{
				"focus",
				"coding",
			},
			MaxEnergy: 2,
		},
	},
	{
		ID:   "discovery",
		Name: "Discovery",
		Criteria: Criteria{
			Contexts: []string{
				"discovery",
			},
			MinEnergy: 2,
			MaxEnergy: 4,
		},
	},
	{
		ID:   "energy",
		Name: "Energy",
		Criteria: Criteria{
			Contexts: []string{
				"active",
			},
			MinEnergy: 4,
		},
	},
	{
		ID:   "wind-down",
		Name: "Wind Down",
		Criteria: Criteria{
			Moods: []string{
				"calm",
				"mellow",
				"dreamy",
			},
			MaxEnergy: 2,
		},
	},
}

func FindPreset(id string) (*Preset, bool) {

	for i := range Presets {

		if Presets[i].ID == id {
			return &Presets[i], true
		}

	}

	return nil, false

}
