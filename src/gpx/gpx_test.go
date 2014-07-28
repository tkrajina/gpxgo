package gpx

import (
	"fmt"
	"testing"
	"time"
)

func assertEquals(t *testing.T, var1 interface{}, var2 interface{}) {
	if var1 != var2 {
		fmt.Println(var1, "not equals to", var2)
		t.Error("Not equals")
	}
}

func assertNil(t *testing.T, var1 interface{}) {
	if var1 != nil {
		fmt.Println(var1)
		t.Error("nil!")
	}
}

func assertNotNil(t *testing.T, var1 interface{}) {
	if var1 == nil {
		fmt.Println(var1)
		t.Error("nil!")
	}
}

func TestParseGPXTimes(t *testing.T) {
	datetimes := []string{
		"2013-01-02T12:07:08Z",
		"2013-01-02 12:07:08Z",
		"2013-01-02T12:07:08",
		"2013-01-02T12:07:08.034Z",
		"2013-01-02 12:07:08.045Z",
		"2013-01-02T12:07:08.123",
	}
	for _, value := range datetimes {
		fmt.Println("datetime:", value)
		parsedTime, err := parseGPXTime(value)
		fmt.Println(parsedTime)
		assertNil(t, err)
		assertNotNil(t, parsedTime)
		assertEquals(t, parsedTime.Year(), 2013)
		assertEquals(t, parsedTime.Month(), time.January)
		assertEquals(t, parsedTime.Day(), 2)
		assertEquals(t, parsedTime.Hour(), 12)
		assertEquals(t, parsedTime.Minute(), 7)
		assertEquals(t, parsedTime.Second(), 8)
	}
}

func TestParseAndReparseGPX11(t *testing.T) {
	gpxDocuments := []*GPX{}

	{
		gpxDoc, err := ParseFile("../../test_files/gpx1.1_with_all_fields.gpx")
		if err != nil || gpxDoc == nil {
			t.Error("Error parsing:" + err.Error())
		}
		gpxDocuments = append(gpxDocuments, gpxDoc)

		// Test after reparsing
		xml, err := gpxDoc.ToXml("1.1")
		//fmt.Println(string(xml))
		if err != nil {
			t.Error("Error serializing to XML:" + err.Error())
		}
		gpxDoc2, err := ParseString(xml)
		if err != nil {
			t.Error("Error parsing XML:" + err.Error())
		}
		gpxDocuments = append(gpxDocuments, gpxDoc2)

		// TODO: ToString 1.0 and check again
	}

	for i := 1; i < len(gpxDocuments); i++ {
		fmt.Println("Testing gpx doc #", i)

		gpxDoc := gpxDocuments[i]

		assertEquals(t, gpxDoc.Version, "1.1")
		assertEquals(t, gpxDoc.Creator, "...")
		assertEquals(t, gpxDoc.Name, "example name")
		assertEquals(t, gpxDoc.AuthorName, "author name")
		assertEquals(t, gpxDoc.AuthorEmail, "aaa@bbb.com")
		assertEquals(t, gpxDoc.Description, "example description")
		assertEquals(t, gpxDoc.AuthorLink, "http://link")
		assertEquals(t, gpxDoc.AuthorLinkText, "link text")
		assertEquals(t, gpxDoc.AuthorLinkType, "link type")
		assertEquals(t, gpxDoc.Copyright, "gpxauth")
		assertEquals(t, gpxDoc.CopyrightYear, "2013")
		assertEquals(t, gpxDoc.CopyrightLicense, "lic")
		assertEquals(t, gpxDoc.Link, "http://link2")
		assertEquals(t, gpxDoc.LinkText, "link text2")
		assertEquals(t, gpxDoc.LinkType, "link type2")
		assertEquals(t, gpxDoc.Time.Format("2013-01-01 12:00:00 +0000 UTC"), time.Date(2013, time.January, 01, 12, 0, 0, 0, time.UTC).Format("2013-01-01 12:00:00 +0000 UTC"))
		assertEquals(t, gpxDoc.Keywords, "example keywords")

	}
}
