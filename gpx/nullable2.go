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
var _ xml.UnmarshalerAttr = new(NilableFloat64)
var _ xml.Marshaler = new(NilableFloat64)
var _ xml.MarshalerAttr = new(NilableFloat64)

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

func (n *NilableFloat64) UnmarshalXMLAttr(attr xml.Attr) error {
	strData := strings.Trim(string(attr.Value), " ")
	value, err := strconv.ParseFloat(strData, 64)
	if err != nil {
		n = nil
		return err
	}
	*n = *NewNilableFloat64(value)
	return nil
}

// MarshalXML implements xml marshalling
func (n *NilableFloat64) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.Nil() {
		return nil
	}
	xmlName := xml.Name{Local: start.Name.Local}
	if err := e.EncodeToken(xml.StartElement{Name: xmlName}); err != nil {
		return err
	}
	e.EncodeToken(xml.CharData([]byte(formatNumber(n.Value()))))
	e.EncodeToken(xml.EndElement{Name: xmlName})
	return nil
}

// MarshalXMLAttr implements xml attribute marshalling
func (n *NilableFloat64) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var result xml.Attr
	if n.Nil() {
		return result, nil
	}
	return xml.Attr{
		Name:  xml.Name{Local: name.Local},
		Value: formatNumber(n.Value()),
	}, nil
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
var _ xml.UnmarshalerAttr = new(NilableInt)
var _ xml.Marshaler = new(NilableInt)
var _ xml.MarshalerAttr = new(NilableInt)

func (f *NilableInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
		*f = *NewNilableint(int(value))
	}
	d.Skip()
	return nil
}

func (n *NilableInt) UnmarshalXMLAttr(attr xml.Attr) error {
	strData := strings.Trim(string(attr.Value), " ")
	value, err := strconv.ParseFloat(strData, 64)
	if err != nil {
		n = nil
		return err
	}
	*n = *NewNilableint(int(value))
	return nil
}

// MarshalXML implements xml marshalling
func (n *NilableInt) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.Nil() {
		return nil
	}
	xmlName := xml.Name{Local: start.Name.Local}
	if err := e.EncodeToken(xml.StartElement{Name: xmlName}); err != nil {
		return err
	}
	e.EncodeToken(xml.CharData([]byte(fmt.Sprintf("%d", n.Value()))))
	e.EncodeToken(xml.EndElement{Name: xmlName})
	return nil
}

// MarshalXMLAttr implements xml attribute marshalling
func (n *NilableInt) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var result xml.Attr
	if n.Nil() {
		return result, nil
	}
	return xml.Attr{
			Name:  xml.Name{Local: name.Local},
			Value: fmt.Sprintf("%d", n.Value())},
		nil
}
