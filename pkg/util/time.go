package util

import (
	"fmt"
	"time"
)

func RelativeTime(t time.Time) string {
	d := time.Since(t)

	const day = 24
	const week = 7 * day
	const year = 365 * day

	switch {
	case d.Hours() >= year:
		return fmt.Sprintf("%.0fy", d.Hours()/year)
	case d.Hours() >= week:
		return fmt.Sprintf("%.0fw", d.Hours()/week)
	case d.Hours() >= day:
		return fmt.Sprintf("%.0fd", d.Hours()/day)
	case d.Hours() >= 1:
		return fmt.Sprintf("%.0fh", d.Hours())
	case d.Minutes() >= 1:
		return fmt.Sprintf("%.0fm", d.Minutes())
	case d.Seconds() >= 1:
		return fmt.Sprintf("%.0fs", d.Seconds())
	default:
		return "now"
	}
}
