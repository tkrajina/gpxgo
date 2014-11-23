// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type NullableInt struct {
	data    int
	notNull bool
}

func (n *NullableInt) Null() bool {
	return !n.notNull
}

func (n *NullableInt) NotNull() bool {
	return n.notNull
}

func (n *NullableInt) Value() int {
	return n.data
}

func (n *NullableInt) SetValue(data int) {
	n.data = data
	n.notNull = true
}

func (n *NullableInt) SetNull() {
	var defaultValue int
	n.data = defaultValue
	n.notNull = false
}

func NewNullableInt(data int) *NullableInt {
	result := new(NullableInt)
	result.data = data
	result.notNull = true
	return result
}

func (n *NullableInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	t, err := d.Token()
	if err != nil {
		n.SetNull()
		return nil
	}
	if charData, ok := t.(xml.CharData); ok {
		strData := strings.Trim(string(charData), " ")
		value, err := strconv.ParseFloat(strData, 64)
		if err != nil {
			n.SetNull()
			return nil
		}
		n.SetValue(int(value))
	}
	d.Skip()
	return nil
}

func (n *NullableInt) UnmarshalXMLAttr(attr xml.Attr) error {
	strData := strings.Trim(string(attr.Value), " ")
	value, err := strconv.ParseFloat(strData, 64)
	if err != nil {
		n.SetNull()
		return nil
	}
	n.SetValue(int(value))
	return nil
}

func (n NullableInt) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.Null() {
		return nil
	}
	xmlName := xml.Name{Local: start.Name.Local}
	e.EncodeToken(xml.StartElement{Name: xmlName})
	e.EncodeToken(xml.CharData([]byte(fmt.Sprintf("%d", n.Value()))))
	e.EncodeToken(xml.EndElement{Name: xmlName})
	return nil
}

func (n NullableInt) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var result xml.Attr
	if n.Null() {
		return result, nil
	}
	return xml.Attr{
			Name:  xml.Name{Local: name.Local},
			Value: fmt.Sprintf("%d", n.Value())},
		nil
}
