package global

import (
	"time"
)

// NanoToTime takes in a nano unix time and converts it to time.Time object
func NanoToTime(t int64) time.Time {
	return time.Unix(0, t)
}
