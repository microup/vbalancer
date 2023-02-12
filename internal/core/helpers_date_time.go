package core

import "time"

// GetDateTimeStr returns the date and time strings of the given time.
func GetDateTimeStr(recordTime time.Time) (string, string) {
	// Format the date string.
	var dateStr = recordTime.Format("2006-01-02")

	// Format the time string.
	var timeStr = recordTime.Format("15:04:05")

	// Return the date and time strings.
	return dateStr, timeStr
}
