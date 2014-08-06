package gpx

import (
    "github.com/joeshaw/gengen/generic"
)

type NullableGeneric struct {
    data generic.T
    null bool
}

func (n *NullableGeneric) Null() bool {
    return n.null
}

func (n *NullableGeneric) NotNull() bool {
    return !n.null
}

func (n *NullableGeneric) Value() generic.T {
    return n.data
}

func (n *NullableGeneric) SetValue(data generic.T) {
    n.data = data
}

func (n *NullableGeneric) SetNull() {
    var defaultValue generic.T
    n.data = defaultValue
    n.null = true
}

func NewNullableGeneric(data generic.T) (*NullableGeneric) {
    result := new(NullableGeneric)
    result.data = data
    result.null = false
    return result
}
