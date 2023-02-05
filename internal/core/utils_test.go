package core_test

import (
	"math"
	"testing"
	"vbalancer/internal/core"
)

func TestHumanFileSize(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		size float64
		want string
	}{
		{1023, "1023 B"},
		{1024, "1 KB"},
		{1048575, "1024 KB"},
		{1048576, "1 MB"},
		{1073741823, "1024 MB"},
		{1073741824, "1 GB"},
		{1099511627775, "1024 GB"},
		{1099511627776, "1 TB"},
	}

	for _, tc := range testCases {
		got := core.HumanFileSize(tc.size)
		if got != tc.want {
			t.Errorf("HumanFileSize(%f) = %s; want %s", tc.size, got, tc.want)
		}
	}
}

func TestRound(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		val      float64
		roundOn  float64
		places   int
		expected float64
	}{
		{3.1415, 0.5, 2, 3.14},
		{1.2345, 0.5, 3, 1.235},
		{2.6666, 0.5, 2, 2.67},
		{0.0, 0.5, 2, 0.0},
		{-1.234, 0.5, 2, -1.23},
	}

	for _, tc := range testCases {
		actual := core.Round(tc.val, tc.roundOn, tc.places)
		if math.Abs(actual-tc.expected) > 0.0001 {
			t.Errorf("For val %v, roundOn %v and places %v, expected %v but got %v",
				tc.val, tc.roundOn, tc.places, tc.expected, actual)
		}
	}
}
