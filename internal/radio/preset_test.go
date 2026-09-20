package radio

import "testing"

func TestPresets(t *testing.T) {

	if len(Presets) == 0 {
		t.Fatal("expected presets")
	}

	ids := make(map[string]bool)

	for _, preset := range Presets {

		t.Run(preset.ID, func(t *testing.T) {

			if preset.ID == "" {
				t.Fatal("expected preset ID")
			}

			if ids[preset.ID] {
				t.Fatalf(
					"duplicate preset ID %q",
					preset.ID,
				)
			}

			ids[preset.ID] = true

			if preset.Name == "" {
				t.Fatal("expected preset name")
			}

			matches := Match(
				preset.Criteria,
			)

			if len(matches) == 0 {
				t.Fatalf(
					"preset %q produced no station matches",
					preset.ID,
				)
			}

		})

	}

}

func TestFindPreset(t *testing.T) {

	preset, ok := FindPreset(
		"focus",
	)
	if !ok {
		t.Fatal("expected focus preset")
	}

	if preset.Name != "Focus" {
		t.Errorf(
			"expected preset name %q, got %q",
			"Focus",
			preset.Name,
		)
	}

}

func TestFindPresetMiss(t *testing.T) {

	_, ok := FindPreset(
		"definitely-not-a-preset",
	)

	if ok {
		t.Fatal("expected preset lookup to miss")
	}

}

func TestPresetCanChooseStation(t *testing.T) {

	preset, ok := FindPreset(
		"focus",
	)
	if !ok {
		t.Fatal("expected focus preset")
	}

	station, ok := Choose(
		preset.Criteria,
		History{},
		ChooseOptions{
			RecentLimit: 3,
		},
	)
	if !ok {
		t.Fatal("expected focus preset to choose a station")
	}

	if station.ID == "" {
		t.Fatal("expected selected station ID")
	}

}
