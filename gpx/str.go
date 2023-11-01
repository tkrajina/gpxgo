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
	return strings.TrimRight(fmt.Sprintf("%.10f", f), "0")
}
