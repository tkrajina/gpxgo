package gpx

import (
	"fmt"
	"math"
	"strings"
)

func formatNumber(f float64) string {
	if f == math.Round(f) {
		return fmt.Sprintf("%d", int(f))
	}
	res := fmt.Sprintf("%.10f", f)
	if strings.Contains(res, ".") {
		return strings.TrimRight(res, "0")
	}
	return res
}
