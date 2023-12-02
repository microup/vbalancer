package core_test

import (
	"fmt"
	"math"
	"testing"
	"vbalancer/internal/core"

	"github.com/stretchr/testify/assert"
)

// TestHumanFileSize tests the HumanFileSize function.
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

	for _, test := range testCases {
		got := core.HumanFileSize(test.size)

		assert.Equalf(t, test.want, got, fmt.Sprintf("size %f want %s", test.size, test.want))
	}
}

// TestRound tests the Round function.
func TestRound(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		val      float64
		roundOn  float64
		places   int
		expected float64
	}{
		{3.1415, 0.4, 2, 3.14},
		{1.2345, 0.5, 3, 1.235},
		{2.6666, 0.5, 2, 2.67},
		{0.0, 0.5, 2, 0.0},
		{-1.234, 0.5, 2, -1.23},
	}

	for _, testCase := range testCases {
		actual := core.Round(testCase.val, testCase.roundOn, testCase.places)
		if math.Abs(actual-testCase.expected) > 0.0001 {
			assert.InDelta(
				t,
				testCase.expected,
				actual,
				0.0001,
				"for val %v, roundOn %v and places %v, expected %v but got %v",
				testCase.val, testCase.roundOn, testCase.places, actual, testCase.expected,
			)
		}
	}
}
