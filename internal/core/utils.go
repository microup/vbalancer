package core

import (
	"math"
	"strconv"
)

const (
	place      = 2
	prec       = -1
)

// HumanFileSize returns a human-readable file size.
func HumanFileSize(size float64) string {
	// Declare a string array with 5 elements.
	var suffixes [5]string

	// Assign the second element of the array to "KB".
	suffixes[0] = "B"

	// Assign the third element of the array to "MB".
	suffixes[1] = "KB"

	// Assign the fourth element of the array to "GB".
	suffixes[2] = "MB"

	// Assign the fifth element of the array to "TB".
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	// Calculate the size of the file.
	base := math.Log(size) / math.Log(LengthByte)
	getSize := Round(math.Pow(LengthByte, base-math.Floor(base)), RoundOne, place)

	// If the base is greater than 0, assign the suffix to the corresponding element of the array.
	var suffix string
	if base > 0 {
		suffix = suffixes[int(math.Floor(base))]
	}

	// Return the size of the file with the suffix.
	return strconv.FormatFloat(getSize, 'f', prec, BitSize) + " " + suffix
}

// Round returns a rounded value.
func Round(val float64, roundOn float64, places int) float64 {
	// Calculate the power of the value.
	pow := math.Pow(PowX, float64(places))

	// Return the rounded value.
	return math.Round(val*pow) / pow
}
