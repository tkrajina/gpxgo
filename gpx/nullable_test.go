package gpx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullable(t *testing.T) {
	t.Parallel()

	number := NewNullableInt(1)

	assert.False(t, number.Null())
	assert.True(t, number.NotNull())
	assert.Equal(t, 1, number.Value())

	{
		val, hasVal, err := unmarshall("17", int(0))
		assert.Nil(t, err)
		assert.True(t, hasVal)
		assert.Equal(t, 17, val)
	}
	{
		val, hasVal, err := unmarshall("", int(0))
		assert.Nil(t, err)
		assert.False(t, hasVal)
		assert.Equal(t, 0, val)
	}
}
