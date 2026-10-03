package radio

import "time"

type Tune struct {
	StationID string
	TunedAt   time.Time
}

type History struct {
	Tunes  []Tune
	cursor int
}

func (h *History) Add(stationID string, tunedAt time.Time) {

	if stationID == "" {
		return
	}

	current, ok := h.Current()

	if ok &&
		current.StationID == stationID {

		return
	}

	if h.cursor < len(h.Tunes) {

		h.Tunes = append(
			[]Tune(nil),
			h.Tunes[:h.cursor]...,
		)

	}

	h.Tunes = append(
		h.Tunes,
		Tune{
			StationID: stationID,
			TunedAt:   tunedAt,
		},
	)

	h.cursor = len(h.Tunes)

}

func (h History) Current() (Tune, bool) {

	if h.cursor <= 0 ||
		h.cursor > len(h.Tunes) {

		return Tune{}, false
	}

	return h.Tunes[h.cursor-1], true

}

func (h History) Previous() (Tune, bool) {

	if h.cursor <= 1 ||
		h.cursor > len(h.Tunes) {

		return Tune{}, false
	}

	return h.Tunes[h.cursor-2], true

}

func (h History) ContainsRecent(
	stationID string,
	limit int,
) bool {

	if stationID == "" || limit <= 0 {
		return false
	}

	start := len(h.Tunes) - limit

	if start < 0 {
		start = 0
	}

	for i := len(h.Tunes) - 1; i >= start; i-- {

		if h.Tunes[i].StationID == stationID {
			return true
		}

	}

	return false

}

func (h *History) CanBack() bool {

	return h.cursor > 1

}

func (h *History) CanForward() bool {

	return h.cursor > 0 &&
		h.cursor < len(h.Tunes)

}

func (h *History) Back() (Tune, bool) {

	if !h.CanBack() {
		return Tune{}, false
	}

	h.cursor--

	return h.Current()

}

func (h *History) Forward() (Tune, bool) {

	if !h.CanForward() {
		return Tune{}, false
	}

	h.cursor++

	return h.Current()

}

func (h History) IsCurrent(
	index int,
) bool {

	return h.cursor > 0 &&
		h.cursor <= len(h.Tunes) &&
		index == h.cursor-1

}
