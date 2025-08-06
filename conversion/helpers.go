package conversion

import "math"

func roundToTwoDecimals(val float64) float64 {
	return math.Round(val*100) / 100
}
