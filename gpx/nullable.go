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

// NullableGeneric implements a nullable generic
type NullableGeneric[T any] struct {
	data    T
	notNull bool
}

// Null checks if value is null
func (n *NullableGeneric[T]) Null() bool {
	return !n.notNull
}

// NotNull checks if value is not null
func (n *NullableGeneric[T]) NotNull() bool {
	return n.notNull
}

// Value returns the value
func (n *NullableGeneric[T]) Value() T {
	return n.data
}

// SetValue sets the value
func (n *NullableGeneric[T]) SetValue(data T) {
	n.data = data
	n.notNull = true
}

// SetNull sets the value to null
func (n *NullableGeneric[T]) SetNull() {
	var defaultValue T
	n.data = defaultValue
	n.notNull = false
}

func (n NullableGeneric[T]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !n.notNull {
		return nil
	}
	xmlName := xml.Name{Local: start.Name.Local}
	if err := e.EncodeToken(xml.StartElement{Name: xmlName}); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.CharData([]byte(fmt.Sprintf("%v", n.data)))); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.EndElement{Name: xmlName}); err != nil {
		return err
	}
	return nil
}
func (n NullableGeneric[T]) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var result xml.Attr
	if !n.notNull {
		return result, nil
	}
	return xml.Attr{
			Name:  xml.Name{Local: name.Local},
			Value: fmt.Sprintf("%v", n.data)},
		nil

}

func NewNullableInt(i int) NullableInt {
	var result NullableGeneric[int]
	result.data = i
	result.notNull = true
	return NullableInt{
		NullableGeneric: result,
	}
}

func NewNullableFloat(f float64) NullableFloat {
	var result NullableGeneric[float64]
	result.data = f
	result.notNull = true
	return NullableFloat{
		NullableGeneric: result,
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Common utilities for custom unmarshallers
////////////////////////////////////////////////////////////////////////////////////////////////////

func unmarshall(str string, defaultVal any) (val any, hasVal bool, e error) {
	val = defaultVal
	str = strings.TrimSpace(str)
	if str == "" {
		hasVal = false
		return
	}
	var d any = defaultVal
	switch d.(type) {
	case int, int64:
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			e = err
			hasVal = false
			return
		}
		hasVal = true
		var v any = int(value)
		val = v
		return
	case float64:
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			e = err
			hasVal = false
			return
		}
		hasVal = true
		var v any = value
		val = v
		return
	}
	e = fmt.Errorf("invalid type %T", d)
	return
}

// UnmarshalXML implements xml unmarshalling
func UnmarshalXML(d *xml.Decoder, start xml.StartElement, defaultVal any) (val any, notNil bool, e error) {
	t, err := d.Token()
	if err != nil {
		notNil = false
		return
	}
	if charData, ok := t.(xml.CharData); ok {
		strData := strings.Trim(string(charData), " ")
		val, notNil, e = unmarshall(strData, defaultVal)
	}
	d.Skip()
	return
}

// UnmarshalXMLAttr implements xml attribute unmarshalling
func UnmarshalXMLAttr(attr xml.Attr, typ any) (val any, notNil bool, e error) {
	strData := strings.TrimSpace(string(attr.Value))
	val, notNil, e = unmarshall(strData, typ)
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////

type NullableInt struct {
	NullableGeneric[int]
}

var _ xml.Marshaler = new(NullableInt)
var _ xml.MarshalerAttr = new(NullableInt)
var _ xml.Unmarshaler = new(NullableInt)
var _ xml.UnmarshalerAttr = new(NullableInt)

func (n *NullableInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	val, notNil, err := UnmarshalXML(d, start, int(0))
	n.data = val.(int)
	n.notNull = notNil
	return err
}
func (n *NullableInt) UnmarshalXMLAttr(attr xml.Attr) error {
	val, notNil, err := UnmarshalXMLAttr(attr, int(0))
	n.data = val.(int)
	n.notNull = notNil
	return err
}

////////////////////////////////////////////////////////////////////////////////////////////////////

type NullableFloat struct {
	NullableGeneric[float64]
}

var _ xml.Marshaler = new(NullableFloat)
var _ xml.MarshalerAttr = new(NullableFloat)
var _ xml.Unmarshaler = new(NullableFloat)
var _ xml.UnmarshalerAttr = new(NullableFloat)

func (n *NullableFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	val, notNil, err := UnmarshalXML(d, start, float64(0))
	n.data = val.(float64)
	n.notNull = notNil
	return err
}
func (n *NullableFloat) UnmarshalXMLAttr(attr xml.Attr) error {
	val, notNil, err := UnmarshalXMLAttr(attr, float64(0))
	n.data = val.(float64)
	n.notNull = notNil
	return err
}
