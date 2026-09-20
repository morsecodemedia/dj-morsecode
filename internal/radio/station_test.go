package radio

import (
	"strings"
	"testing"
)

var validMoods = map[string]bool{
	"adventurous":   true,
	"aggressive":    true,
	"calm":          true,
	"chilled":       true,
	"contemplative": true,
	"cozy":          true,
	"curious":       true,
	"dreamy":        true,
	"energetic":     true,
	"euphoric":      true,
	"familiar":      true,
	"groovy":        true,
	"independent":   true,
	"intense":       true,
	"meditative":    true,
	"mellow":        true,
	"nostalgic":     true,
	"otherworldly":  true,
	"peaceful":      true,
	"playful":       true,
	"positive":      true,
	"raw":           true,
	"rebellious":    true,
	"refined":       true,
	"sensual":       true,
	"smooth":        true,
	"spacious":      true,
	"sunny":         true,
	"thoughtful":    true,
	"upbeat":        true,
	"uplifting":     true,
	"varied":        true,
}

var validContexts = map[string]bool{
	"active":     true,
	"background": true,
	"casual":     true,
	"coding":     true,
	"discovery":  true,
	"driving":    true,
	"focus":      true,
	"gaming":     true,
	"late-night": true,
	"meditation": true,
	"morning":    true,
	"party":      true,
	"relaxing":   true,
	"sleep":      true,
	"social":     true,
	"workout":    true,
}

func TestStations(t *testing.T) {

	if len(Stations) == 0 {
		t.Fatal("expected station catalog")
	}

	ids := make(map[string]bool)

	for _, station := range Stations {

		t.Run(station.ID, func(t *testing.T) {

			if strings.TrimSpace(station.ID) == "" {
				t.Error("expected station ID")
			}

			if ids[station.ID] {
				t.Errorf(
					"duplicate station ID %q",
					station.ID,
				)
			}

			ids[station.ID] = true

			if strings.TrimSpace(station.Name) == "" {
				t.Error("expected station name")
			}

			if strings.TrimSpace(station.StreamURL) == "" {
				t.Error("expected station stream URL")
			}

			if strings.TrimSpace(station.Genre) == "" {
				t.Error("expected station genre")
			}

			if station.Energy < 1 || station.Energy > 5 {
				t.Errorf(
					"expected energy between 1 and 5, got %d",
					station.Energy,
				)
			}

			validateValues(
				t,
				"tag",
				station.Tags,
				nil,
			)

			validateValues(
				t,
				"mood",
				station.Moods,
				validMoods,
			)

			validateValues(
				t,
				"context",
				station.Contexts,
				validContexts,
			)

		})

	}

}

func TestFindByStreamURL(t *testing.T) {

	station, ok := FindByStreamURL(
		"https://wxpn.xpn.org/xpnmp3hi",
	)
	if !ok {
		t.Fatal("expected station to be found")
	}

	if station.ID != "wxpn" {
		t.Errorf(
			"expected station ID %q, got %q",
			"wxpn",
			station.ID,
		)
	}

}

func TestFindByStreamURLMiss(t *testing.T) {

	_, ok := FindByStreamURL(
		"https://example.com/not-a-station",
	)

	if ok {
		t.Fatal("expected station lookup to miss")
	}

}

func validateValues(
	t *testing.T,
	name string,
	values []string,
	allowed map[string]bool,
) {

	t.Helper()

	seen := make(map[string]bool)

	for _, value := range values {

		if strings.TrimSpace(value) == "" {
			t.Errorf(
				"%s must not be empty",
				name,
			)
			continue
		}

		if value != strings.ToLower(value) {
			t.Errorf(
				"%s %q must be lowercase",
				name,
				value,
			)
		}

		if seen[value] {
			t.Errorf(
				"duplicate %s %q",
				name,
				value,
			)
		}

		seen[value] = true

		if allowed != nil && !allowed[value] {
			t.Errorf(
				"unknown %s %q",
				name,
				value,
			)
		}

	}

}
