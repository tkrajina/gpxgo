// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type NilableFloat64 float64

func NewNilableFloat64(val float64) *NilableFloat64 {
	res := NilableFloat64(val)
	return &res
}

func (f *NilableFloat64) Nil() bool {
	return f == nil
}

func (f *NilableFloat64) NotNil() bool {
	return f != nil
}

func (f *NilableFloat64) Value() float64 {
	if f == nil {
		return 0.0
	}
	return float64(*f)
}

var _ xml.Unmarshaler = new(NilableFloat64)

func (f *NilableFloat64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	t, err := d.Token()
	if err != nil {
		f = nil
		return nil
	}
	if charData, ok := t.(xml.CharData); ok {
		strData := strings.Trim(string(charData), " ")
		value, err := strconv.ParseFloat(strData, 64)
		if err != nil {
			f = nil
			return err
		}
		*f = *NewNilableFloat64(value)
	}
	d.Skip()
	return nil
}

type NilableInt float64

func NewNilableint(val int) *NilableInt {
	res := NilableInt(val)
	return &res
}

func (f *NilableInt) Nil() bool {
	return f == nil
}

func (f *NilableInt) NotNil() bool {
	return f != nil
}

func (f *NilableInt) Value() int {
	if f == nil {
		return 0
	}
	return int(*f)
}

var _ xml.Unmarshaler = new(NilableInt)

func (f *NilableInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	t, err := d.Token()
	if err != nil {
		f = nil
		return nil
	}
	if charData, ok := t.(xml.CharData); ok {
		strData := strings.Trim(string(charData), " ")
		value, err := strconv.ParseInt(strData, 10, 64)
		if err != nil {
			f = nil
			return nil
		}
		*f = *NewNilableint(int(value))
	}
	d.Skip()
	return nil
}
