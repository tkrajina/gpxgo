// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

type NullableString struct {
	data    string
	notNull bool
}

func (n *NullableString) Null() bool {
	return !n.notNull
}

func (n *NullableString) NotNull() bool {
	return n.notNull
}

func (n *NullableString) Value() string {
	return n.data
}

func (n *NullableString) SetValue(data string) {
	n.data = data
	n.notNull = true
}

func (n *NullableString) SetNull() {
	var defaultValue string
	n.data = defaultValue
	n.notNull = false
}

func NewNullableString(data string) *NullableString {
	result := new(NullableString)
	result.data = data
	result.notNull = true
	return result
}
