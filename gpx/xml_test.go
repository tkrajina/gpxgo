// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"encoding/xml"
	"testing"
)

func TestParseTime(t *testing.T) {
	time, err := parseGPXTime("")
	if time != nil {
		t.Errorf("Empty string should not return a nonnil time")
	}
	if err == nil {
		t.Errorf("Empty string should result in error")
	}
}

type testXml struct {
	XMLName   xml.Name        `xml:"gpx"`
	Float     NullableFloat64 `xml:"float"`
	Int       NullableInt     `xml:"int"`
	FloatAttr NullableFloat64 `xml:"floatattr,attr"`
	IntAttr   NullableInt     `xml:"intattr,attr"`
}

func TestInvalidFloat(t *testing.T) {
	xmlStr := `<gpx floatattr="1"><float>...a</float></gpx>`
	testXmlDoc := testXml{}
	xml.Unmarshal([]byte(xmlStr), &testXmlDoc)
	if testXmlDoc.Float.NotNull() {
		t.Error("Float is invalid in ", xmlStr)
	}
}

func TestValidFloat(t *testing.T) {
	xmlStr := `<gpx floatattr="13"><float>12</float><aaa  /></gpx>`
	testFloat(xmlStr, 12, 13, t)
}

func TestValidFloat2(t *testing.T) {
	xmlStr := `<gpx floatattr=" 13.4"><float> 12.3</float></gpx>`
	testFloat(xmlStr, 12.3, 13.4, t)
}

func TestValidFloat3(t *testing.T) {
	xmlStr := `<gpx floatattr="13.5   " ><float>12.3    </float></gpx>`
	testFloat(xmlStr, 12.3, 13.5, t)
}

func testFloat(xmlStr string, expectedFloat float64, expectedFloatAttribute float64, t *testing.T) {
	testXmlDoc := testXml{}
	xml.Unmarshal([]byte(xmlStr), &testXmlDoc)
	if testXmlDoc.Float.Null() || testXmlDoc.Float.Value() != expectedFloat {
		t.Error("Float is valid in ", xmlStr)
	}
	if testXmlDoc.FloatAttr.Null() || testXmlDoc.FloatAttr.Value() != expectedFloatAttribute {
		t.Error("Float attribute is valid in ", xmlStr)
	}
}

func TestValidInt(t *testing.T) {
	xmlStr := `<gpx intattr="15"><int>12</int></gpx>`
	testInt(xmlStr, 12, 15, t)
}

func TestValidInt2(t *testing.T) {
	xmlStr := `<gpx intattr="  17.2"><int> 12.3</int></gpx>`
	testInt(xmlStr, 12, 17, t)
}

func TestValidInt3(t *testing.T) {
	xmlStr := `<gpx intattr="18   "><int>12.3    </int></gpx>`
	testInt(xmlStr, 12, 18, t)
}

func testInt(xmlStr string, expectedInt int, expectedIntAttribute int, t *testing.T) {
	testXmlDoc := testXml{}
	xml.Unmarshal([]byte(xmlStr), &testXmlDoc)
	if testXmlDoc.Int.Null() || testXmlDoc.Int.Value() != expectedInt {
		t.Error("Int is valid in ", xmlStr)
	}
	if testXmlDoc.IntAttr.Null() || testXmlDoc.IntAttr.Value() != expectedIntAttribute {
		t.Error("Int attribute is valid in ", xmlStr)
	}
}
