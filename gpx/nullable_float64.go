package gpx

type NullableFloat64 struct {
	data float64
	null bool
}

func (n *NullableFloat64) Null() bool {
	return n.null
}

func (n *NullableFloat64) NotNull() bool {
	return !n.null
}

func (n *NullableFloat64) Value() float64 {
	return n.data
}

func (n *NullableFloat64) SetValue(data float64) {
	n.data = data
}

func NewNullableFloat64(data float64) *NullableFloat64 {
	result := new(NullableFloat64)
	result.data = data
	return result
}
