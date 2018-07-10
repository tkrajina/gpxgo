// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"github.com/joeshaw/gengen/generic"
)

//NullableGeneric implements a nullable generic
type NullableGeneric struct {
	data    generic.T
	notNull bool
}

//Null checks if value is null
func (n *NullableGeneric) Null() bool {
	return !n.notNull
}

//NotNull checks if value is not null
func (n *NullableGeneric) NotNull() bool {
	return n.notNull
}

//Value returns the value
func (n *NullableGeneric) Value() generic.T {
	return n.data
}

//SetValue sets the value
func (n *NullableGeneric) SetValue(data generic.T) {
	n.data = data
	n.notNull = true
}

//SetNull sets the value to null
func (n *NullableGeneric) SetNull() {
	var defaultValue generic.T
	n.data = defaultValue
	n.notNull = false
}

//NewNullableGeneric creates a new NullableGeneric
func NewNullableGeneric(data generic.T) *NullableGeneric {
	result := new(NullableGeneric)
	result.data = data
	result.notNull = true
	return result
}
