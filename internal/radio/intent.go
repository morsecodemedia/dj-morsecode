package radio

type IntentType string

const (
	IntentNone  IntentType = ""
	IntentVibe  IntentType = "vibe"
	IntentMood  IntentType = "mood"
	IntentGenre IntentType = "genre"
)

type Intent struct {
	Type     IntentType
	ID       string
	Name     string
	Criteria Criteria
}

func VibeIntent(
	preset Preset,
) Intent {

	return Intent{
		Type:     IntentVibe,
		ID:       preset.ID,
		Name:     preset.Name,
		Criteria: preset.Criteria,
	}

}

func MoodIntent(
	family MoodFamily,
) Intent {

	return Intent{
		Type:     IntentMood,
		ID:       family.ID,
		Name:     family.Name,
		Criteria: family.Criteria(),
	}

}

func GenreIntent(
	genre string,
) Intent {

	return Intent{
		Type: IntentGenre,
		ID:   genre,
		Name: genre,
		Criteria: Criteria{
			Genres: []string{
				genre,
			},
		},
	}

}

func (i Intent) Active() bool {

	return i.Type != IntentNone

}
