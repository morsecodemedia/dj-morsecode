package main

import (
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/observation"
	"github.com/morsecodemedia/dj-morsecode/internal/radio"
)

func (m model) reconcileStationIdentity(
	path string,
) (model, *radio.Station, bool) {

	station, found :=
		radio.FindByStreamURL(
			path,
		)

	if found {

		changed :=
			m.CurrentStationID != "" &&
				m.CurrentStationID !=
					station.ID

		if changed {

			m.PlaybackItem =
				metadata.PlaybackItem{}

			m = m.clearTrackState()

		}

		m.CurrentStationID =
			station.ID

		return m, station, true
	}

	if m.PendingStationID == "" {

		if m.CurrentStationID != "" {

			m.PlaybackItem =
				metadata.PlaybackItem{}

			m = m.clearTrackState()

		}

		m.CurrentStationID = ""

	}

	return m, nil, false

}

func (m model) reconcileStationLifecycle(
	station *radio.Station,
	stationFound bool,
	idle bool,
	now time.Time,
) (model, bool) {

	if stationFound {

		if !idle {

			m.StationHistory.Add(
				station.ID,
				now,
			)

		}

		if station.ID ==
			m.PendingStationID &&
			!idle {

			m.observeStation(
				observation.StationTuneConfirmed,
				station.ID,
			)

			m.PendingStationID = ""
			m.PendingSince = time.Time{}

		}

	}

	if m.PendingStationID != "" &&
		!m.PendingSince.IsZero() &&
		now.Sub(m.PendingSince) >=
			stationTuneGracePeriod &&
		idle {

		failedStationID :=
			m.PendingStationID

		m.observeStation(
			observation.StationTuneFailed,
			failedStationID,
		)

		m.PendingStationID = ""
		m.PendingSince = time.Time{}

		m = m.failStation(
			failedStationID,
		)

		m = m.nextStation()

		return m, true
	}

	if m.shouldRotateStation(
		now,
	) {

		m = m.nextStation()

		return m, true
	}

	return m, false

}
