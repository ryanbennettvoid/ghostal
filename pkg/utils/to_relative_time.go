package utils

import (
	"fmt"
	"time"
)

func ToRelativeTime(t time.Time, start time.Time) string {
	diff := start.Sub(t)

	seconds := int(diff.Seconds())
	minutes := int(diff.Minutes())
	hours := int(diff.Hours())
	days := hours / 24

	switch {
	case seconds < 60:
		return fmt.Sprintf("%ds", seconds)
	case minutes < 60:
		return fmt.Sprintf("%dm", minutes)
	case hours < 24:
		return fmt.Sprintf("%dh", hours)
	default:
		return fmt.Sprintf("%dd", days)
	}
}
