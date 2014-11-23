// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import "time"

type NullableTime struct {
	data    time.Time
	notNull bool
}

func (n *NullableTime) Null() bool {
	return !n.notNull
}

func (n *NullableTime) NotNull() bool {
	return n.notNull
}

func (n *NullableTime) Value() time.Time {
	return n.data
}

func (n *NullableTime) SetValue(data time.Time) {
	n.data = data
	n.notNull = true
}

func (n *NullableTime) SetNull() {
	var defaultValue time.Time
	n.data = defaultValue
	n.notNull = false
}

func NewNullableTime(data time.Time) *NullableTime {
	result := new(NullableTime)
	result.data = data
	result.notNull = true
	return result
}
