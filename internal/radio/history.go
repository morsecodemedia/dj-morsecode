package radio

import "time"

type Tune struct {
	StationID string
	TunedAt   time.Time
}

type History struct {
	Tunes []Tune
}

func (h *History) Add(
	stationID string,
	tunedAt time.Time,
) {

	if stationID == "" {
		return
	}

	current, ok := h.Current()

	if ok && current.StationID == stationID {
		return
	}

	h.Tunes = append(
		h.Tunes,
		Tune{
			StationID: stationID,
			TunedAt:   tunedAt,
		},
	)

}

func (h History) Current() (Tune, bool) {

	if len(h.Tunes) == 0 {
		return Tune{}, false
	}

	return h.Tunes[len(h.Tunes)-1], true

}

func (h History) Previous() (Tune, bool) {

	if len(h.Tunes) < 2 {
		return Tune{}, false
	}

	return h.Tunes[len(h.Tunes)-2], true

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
