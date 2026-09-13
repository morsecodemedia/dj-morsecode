package player

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {

	totalSeconds := int(d.Seconds())

	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf(
		"%02d:%02d",
		minutes,
		seconds,
	)

}
