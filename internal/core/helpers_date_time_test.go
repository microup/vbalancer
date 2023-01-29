package core_test

import (
	"testing"
	"time"
	"vbalancer/internal/core"
)

func TestGetDateTimeStr(t *testing.T) {
	t.Parallel()
	
	testCases := []struct {
		recordTime time.Time
		expectedDateStr string
		expectedTimeStr string
	}{
		{time.Date(2021, time.April, 11, 15, 20, 5, 0, time.UTC), "2021-04-11", "15:20:05"},
		{time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), "2000-01-01", "00:00:00"},
		{time.Date(2022, time.December, 31, 23, 59, 59, 0, time.UTC), "2022-12-31", "23:59:59"},
	}

	for _, tc := range testCases {
		dateStr, timeStr := core.GetDateTimeStr(tc.recordTime)
		if dateStr != tc.expectedDateStr || timeStr != tc.expectedTimeStr {
			t.Errorf("GetDateTimeStr(%v) = (%v, %v), want (%v, %v)", tc.recordTime,
				dateStr, timeStr, tc.expectedDateStr, tc.expectedTimeStr)
		}
	}
}