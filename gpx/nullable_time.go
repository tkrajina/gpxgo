package gpx

import (
    "time"
)

type NullableTime struct {
	data time.Time
	null bool
}

func (n *NullableTime) Null() bool {
	return n.null
}

func (n *NullableTime) NotNull() bool {
	return !n.null
}

func (n *NullableTime) Value() time.Time {
	return n.data
}

func (n *NullableTime) SetValue(data time.Time) {
	n.data = data
}

func (n *NullableTime) SetNull() {
	var defaultValue time.Time
	n.data = defaultValue
	n.null = true
}

func NewNullableTime(data time.Time) *NullableTime {
	result := new(NullableTime)
	result.data = data
	result.null = false
	return result
}
