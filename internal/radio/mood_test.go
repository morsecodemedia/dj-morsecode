package radio

import "testing"

func TestMoodFamilies(t *testing.T) {

	if len(MoodFamilies) == 0 {
		t.Fatal("expected mood families")
	}

	ids := make(map[string]bool)

	for _, family := range MoodFamilies {

		t.Run(family.ID, func(t *testing.T) {

			if family.ID == "" {
				t.Fatal("expected mood family ID")
			}

			if ids[family.ID] {
				t.Fatalf(
					"duplicate mood family ID %q",
					family.ID,
				)
			}

			ids[family.ID] = true

			if family.Name == "" {
				t.Fatal("expected mood family name")
			}

			if len(family.Moods) == 0 {
				t.Fatal("expected mood family moods")
			}

			matches := Match(
				family.Criteria(),
			)

			if len(matches) == 0 {
				t.Fatalf(
					"mood family %q produced no station matches",
					family.ID,
				)
			}

		})

	}

}

func TestMoodFamiliesCoverCatalogMoods(t *testing.T) {

	covered := make(map[string]bool)

	for _, family := range MoodFamilies {

		for _, mood := range family.Moods {
			covered[mood] = true
		}

	}

	for _, mood := range Moods() {

		if !covered[mood] {
			t.Errorf(
				"catalog mood %q is not assigned to a mood family",
				mood,
			)
		}

	}

}

func TestMoodFamilyMoodsExistInCatalog(t *testing.T) {

	catalogMoods := make(map[string]bool)

	for _, mood := range Moods() {
		catalogMoods[mood] = true
	}

	for _, family := range MoodFamilies {

		for _, mood := range family.Moods {

			if !catalogMoods[mood] {
				t.Errorf(
					"mood family %q references unknown catalog mood %q",
					family.ID,
					mood,
				)
			}

		}

	}

}

func TestFindMoodFamily(t *testing.T) {

	family, ok := FindMoodFamily(
		"calm",
	)
	if !ok {
		t.Fatal("expected calm mood family")
	}

	if family.Name != "Calm" {
		t.Errorf(
			"expected mood family name %q, got %q",
			"Calm",
			family.Name,
		)
	}

}

func TestFindMoodFamilyMiss(t *testing.T) {

	_, ok := FindMoodFamily(
		"definitely-not-a-mood",
	)

	if ok {
		t.Fatal("expected mood family lookup to miss")
	}

}

func TestMoodFamilyCanChooseStation(t *testing.T) {

	family, ok := FindMoodFamily(
		"dreamy",
	)
	if !ok {
		t.Fatal("expected dreamy mood family")
	}

	station, ok := Choose(
		family.Criteria(),
		History{},
		ChooseOptions{
			RecentLimit: 3,
		},
	)
	if !ok {
		t.Fatal("expected dreamy mood family to choose a station")
	}

	if station.ID == "" {
		t.Fatal("expected selected station ID")
	}

}
