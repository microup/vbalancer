package core

import (
	"math"
	"strconv"
)

const (
	lengthByte = 1024
	powX       = 10
	roundOne   = .5
	place      = 2
	prec       = -1
	bitSize    = 64
)

func HumanFileSize(size float64) string {
	var suffixes [5]string

	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(size) / math.Log(lengthByte)
	getSize := Round(math.Pow(lengthByte, base-math.Floor(base)), roundOne, place)

	var suffix string
	if base > 0 {
		suffix = suffixes[int(math.Floor(base))]
	}

	return strconv.FormatFloat(getSize, 'f', prec, bitSize) + " " + suffix
}

func Round(val float64, roundOn float64, places int) float64 {
	pow := math.Pow(powX, float64(places))

	return math.Round(val*pow) / pow
}
