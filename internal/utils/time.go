package utils

// original project : https://github.com/andanhm/go-prettytime

import (
	"strconv"
	"strings"
	"time"
)

// Unix epoch (or Unix time or POSIX time or Unix timestamp)  1 year (365.24 days)
const infinity float64 = 31556926 * 1000

// Handler function which determines the time difference based on defined time spams
func handler(timeIntervalThreshold float64, timePeriod, message string) func(float64) string {
	return func(difference float64) string {
		var str strings.Builder
		n := difference / timeIntervalThreshold
		nStr := strconv.FormatFloat(n, 'f', 0, 64)
		str.WriteString(nStr)
		str.WriteString(" ")
		str.WriteString(timePeriod)
		if int(n) > 1 {
			str.WriteString("s ")
			str.WriteString(message)
			return str.String()
		}
		str.WriteString(" ")
		str.WriteString(message)
		return str.String()
	}
}

// timeLapse condition struct
type timeLapse struct {
	// Time stamp threshold to handle the time lap condition
	Threshold float64
	// Handler function which determines the time lapse based on the condition
	Handler func(float64) string
}

var timeLapses = []timeLapse{
	{Threshold: -31535999, Handler: handler(-31536000, "year", "from now")},
	{Threshold: -2591999, Handler: handler(-2592000, "month", "from now")},
	{Threshold: -604799, Handler: handler(-604800, "week", "from now")},
	{Threshold: -172799, Handler: handler(-86400, "day", "from now")},
	{Threshold: -86399, Handler: func(diff float64) string {
		return "tomorrow"
	}},
	{Threshold: -3599, Handler: handler(-3600, "hour", "from now")},
	{Threshold: -59, Handler: handler(-60, "minute", "from now")},
	{Threshold: -0.9999, Handler: handler(-1, "second", "from now")},
	{Threshold: 1, Handler: func(diff float64) string {
		return "just now"
	}},
	{Threshold: 60, Handler: handler(1, "second", "ago")},
	{Threshold: 3600, Handler: handler(60, "minute", "ago")},
	{Threshold: 86400, Handler: handler(3600, "hour", "ago")},
	{Threshold: 172800, Handler: func(diff float64) string {
		return "yesterday"
	}},
	{Threshold: 604800, Handler: handler(86400, "day", "ago")},
	{Threshold: 2592000, Handler: handler(604800, "week", "ago")},
	{Threshold: 31536000, Handler: handler(2592000, "month", "ago")},
	{Threshold: infinity, Handler: handler(31536000, "year", "ago")},
}

//Format returns a string describing how long it has been since the time argument passed int
func RelativeTimeDiff(t time.Time) (timeSince string) {
	timestamp := t.Unix()
	now := time.Now().Unix()

	if timestamp > now || timestamp <= 0 {
		timeSince = ""
	}

	timeElapsed := float64(now - timestamp)
	for _, formatter := range timeLapses {
		if timeElapsed < formatter.Threshold {
			timeSince = formatter.Handler(timeElapsed)
			break
		}
	}
	return timeSince
}