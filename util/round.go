package util

import "math"

func Round(number float64, precision int) float64 {
	if precision == 0 {
		return math.Round(number)
	}
	p := math.Pow10(precision)
	return math.Round(number*10*p/10) / p
}
