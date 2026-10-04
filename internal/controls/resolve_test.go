package controls

import "testing"

func TestResolve(
	t *testing.T,
) {

	tests := []struct {
		name string
		mode Mode
		key  string
		want Action
	}{
		{
			name: "enter controls",
			mode: ModeNone,
			key:  "c",
			want: ActionEnterControls,
		},
		{
			name: "enter tune",
			mode: ModeNone,
			key:  "t",
			want: ActionEnterTune,
		},
		{
			name: "enter enhancements",
			mode: ModeNone,
			key:  "e",
			want: ActionEnterEnhancements,
		},
		{
			name: "enter info",
			mode: ModeNone,
			key:  "i",
			want: ActionEnterInfo,
		},
		{
			name: "quit",
			mode: ModeNone,
			key:  "q",
			want: ActionQuit,
		},
		{
			name: "back",
			mode: ModeControls,
			key:  "b",
			want: ActionBack,
		},
		{
			name: "next",
			mode: ModeControls,
			key:  "n",
			want: ActionNext,
		},
		{
			name: "pause",
			mode: ModeControls,
			key:  "p",
			want: ActionTogglePause,
		},
		{
			name: "mute",
			mode: ModeControls,
			key:  "m",
			want: ActionToggleMute,
		},
		{
			name: "enter volume",
			mode: ModeControls,
			key:  "v",
			want: ActionEnterVolume,
		},
		{
			name: "volume up",
			mode: ModeVolume,
			key:  "up",
			want: ActionVolumeUp,
		},
		{
			name: "volume down",
			mode: ModeVolume,
			key:  "down",
			want: ActionVolumeDown,
		},
		{
			name: "stations",
			mode: ModeTune,
			key:  "s",
			want: ActionOpenStations,
		},
		{
			name: "genres",
			mode: ModeTune,
			key:  "g",
			want: ActionOpenGenres,
		},
		{
			name: "moods",
			mode: ModeTune,
			key:  "m",
			want: ActionOpenMoods,
		},
		{
			name: "vibes",
			mode: ModeTune,
			key:  "v",
			want: ActionOpenVibes,
		},
		{
			name: "lyrics",
			mode: ModeEnhancements,
			key:  "l",
			want: ActionToggleLyrics,
		},
		{
			name: "observations",
			mode: ModeInfo,
			key:  "o",
			want: ActionOpenObservations,
		},
		{
			name: "history",
			mode: ModeInfo,
			key:  "h",
			want: ActionOpenHistory,
		},
		{
			name: "cancel controls",
			mode: ModeControls,
			key:  "esc",
			want: ActionCancel,
		},
		{
			name: "cancel volume",
			mode: ModeVolume,
			key:  "esc",
			want: ActionCancel,
		},
		{
			name: "unknown command",
			mode: ModeControls,
			key:  "x",
			want: ActionNone,
		},
	}

	for _, test := range tests {

		t.Run(
			test.name,
			func(t *testing.T) {

				got := Resolve(
					test.mode,
					test.key,
				)

				if got != test.want {

					t.Fatalf(
						"expected %v, got %v",
						test.want,
						got,
					)

				}

			},
		)

	}

}
