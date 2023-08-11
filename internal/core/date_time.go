package core

import "time"

// GetDateTimeStr returns the date and time string in the format of YYYY-MM-DD HH:MM:SS.
func GetDateTimeStr(recordTime time.Time) (string, string) {
	var dateStr = recordTime.Format("2006-01-02")
	
	var timeStr = recordTime.Format("15:04:05")

	return dateStr, timeStr
}
