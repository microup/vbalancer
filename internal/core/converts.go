package core

import (
	"math"
	"strconv"
	"vbalancer/internal/types"
)

const (
	place      = 2
	prec       = -1
)

// HumanFileSize returns a human-readable file size.
func HumanFileSize(size float64) string {
	var suffixes [5]string

	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(size) / math.Log(types.LengthByte)
	getSize := Round(math.Pow(types.LengthByte, base-math.Floor(base)), types.RoundOne, place)

	var suffix string
	if base > 0 {
		suffix = suffixes[int(math.Floor(base))]
	}

	return strconv.FormatFloat(getSize, 'f', prec, types.BitSize) + " " + suffix
}

// Round rounds a float to a given number of decimal places.
func Round(val float64, roundOn float64, places int) float64 {
	pow := math.Pow(types.PowX, float64(places))
	
	return math.Round(val*pow) / pow
}
