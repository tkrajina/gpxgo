// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"github.com/joeshaw/gengen/generic"
)

type NullableGeneric struct {
	data    generic.T
	notNull bool
}

func (n *NullableGeneric) Null() bool {
	return !n.notNull
}

func (n *NullableGeneric) NotNull() bool {
	return n.notNull
}

func (n *NullableGeneric) Value() generic.T {
	return n.data
}

func (n *NullableGeneric) SetValue(data generic.T) {
	n.data = data
	n.notNull = true
}

func (n *NullableGeneric) SetNull() {
	var defaultValue generic.T
	n.data = defaultValue
	n.notNull = false
}

func NewNullableGeneric(data generic.T) *NullableGeneric {
	result := new(NullableGeneric)
	result.data = data
	result.notNull = true
	return result
}
