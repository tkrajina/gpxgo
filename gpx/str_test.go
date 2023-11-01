package gpx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFloat(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "7", formatNumber(7.))
	assert.Equal(t, "7.1", formatNumber(7.1))
	assert.Equal(t, "7.01", formatNumber(7.01))
	assert.Equal(t, "0.01", formatNumber(.01))
}
