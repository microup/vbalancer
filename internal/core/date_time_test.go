package core_test

import (
	"fmt"
	"testing"
	"time"
	"vbalancer/internal/core"

	"github.com/stretchr/testify/assert"
)

// TestGetDateTimeStr tests the GetDateTimeStr function.
func TestGetDateTimeStr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		recordTime      time.Time
		expectedDateStr string
		expectedTimeStr string
	}{
		{time.Date(2021, time.April, 11, 15, 20, 5, 0, time.UTC), "2021-04-11", "15:20:05"},
		{time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), "2000-01-01", "00:00:00"},
		{time.Date(2022, time.December, 31, 23, 59, 59, 0, time.UTC), "2022-12-31", "23:59:59"},
	}

	for _, test := range testCases {
		dateStr, timeStr := core.GetDateTimeStr(test.recordTime)
		assert.Equal(t, test.expectedDateStr, dateStr, "unexpected date string for recordTime: %v", test.recordTime)

		assert.Equal(t, fmt.Sprintf("%s %s", test.expectedDateStr, test.expectedTimeStr),
			fmt.Sprintf("%s %s", dateStr, timeStr),
			"unexpected combined date and time string for recordTime: %v", test.recordTime)
	}
}
