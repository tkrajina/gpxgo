package gpx

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const TIME_FORMAT = "2006-01-02T15:04:05Z"

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

func testDetectVersion(t *testing.T, fileName, expectedVersion string) {
	f, _ := os.Open(fileName)
	contents, _ := ioutil.ReadAll(f)
	version, err := guessGPXVersion(contents)
	fmt.Println("Version=", version)
	if err != nil {
		t.Error("Can't detect 1.1 GPX, error=" + err.Error())
	}
	if version != expectedVersion {
		t.Error("Can't detect 1.1 GPX")
	}
}

func TestDetect11GPXVersion(t *testing.T) {
	testDetectVersion(t, "../../test_files/gpx1.1_with_all_fields.gpx", "1.1")
}

func TestDetect10GPXVersion(t *testing.T) {
	testDetectVersion(t, "../../test_files/gpx1.0_with_all_fields.gpx", "1.0")
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
		assertEquals(t, gpxDoc.Time.Format(TIME_FORMAT), time.Date(2013, time.January, 01, 12, 0, 0, 0, time.UTC).Format(TIME_FORMAT))
		assertEquals(t, gpxDoc.Keywords, "example keywords")

		// Waypoints:
		assertEquals(t, len(gpxDoc.Waypoints), 2)
		assertEquals(t, gpxDoc.Waypoints[0].Latitude, 12.3)
		assertEquals(t, gpxDoc.Waypoints[0].Longitue, 45.6)
		assertEquals(t, gpxDoc.Waypoints[0].Elevation, 75.1)
		assertEquals(t, gpxDoc.Waypoints[0].Timestamp.Format(TIME_FORMAT), "2013-01-02T02:03:00Z")
		assertEquals(t, gpxDoc.Waypoints[0].MagneticVariation, "1.1")
		assertEquals(t, gpxDoc.Waypoints[0].GeoidHeight, "2.0")
		assertEquals(t, gpxDoc.Waypoints[0].Name, "example name")
		assertEquals(t, gpxDoc.Waypoints[0].Comment, "example cmt")
		assertEquals(t, gpxDoc.Waypoints[0].Description, "example desc")
		assertEquals(t, gpxDoc.Waypoints[0].Source, "example src")
		// TODO
		// Links       []GpxLink
		assertEquals(t, gpxDoc.Waypoints[0].Symbol, "example sym")
		assertEquals(t, gpxDoc.Waypoints[0].Type, "example type")
		assertEquals(t, gpxDoc.Waypoints[0].TypeOfGpsFix, "2d")
		assertEquals(t, gpxDoc.Waypoints[0].Satellites, 5)
		assertEquals(t, gpxDoc.Waypoints[0].HorizontalDilution, 6.0)
		assertEquals(t, gpxDoc.Waypoints[0].VerticalDilution, 7.0)
		assertEquals(t, gpxDoc.Waypoints[0].PositionalDilution, 8.0)
		assertEquals(t, gpxDoc.Waypoints[0].AgeOfDGpsData, 9.0)
		assertEquals(t, gpxDoc.Waypoints[0].DGpsId, 45)
		// TODO: Extensions

		assertEquals(t, gpxDoc.Waypoints[1].Latitude, 13.4)
		assertEquals(t, gpxDoc.Waypoints[1].Longitue, 46.7)

		// Routes:
		assertEquals(t, len(gpxDoc.Routes), 2)
		assertEquals(t, gpxDoc.Routes[0].Name, "example name")
		assertEquals(t, gpxDoc.Routes[0].Comment, "example cmt")
		assertEquals(t, gpxDoc.Routes[0].Description, "example desc")
		assertEquals(t, gpxDoc.Routes[0].Source, "example src")
		assertEquals(t, gpxDoc.Routes[0].Number, 7)
		assertEquals(t, gpxDoc.Routes[0].Type, "rte type")
		assertEquals(t, len(gpxDoc.Routes[0].Points), 3)
		// TODO: Link
		// TODO: Points
		assertEquals(t, gpxDoc.Routes[0].Points[0].Elevation, 75.1)
		fmt.Println("t=", gpxDoc.Routes[0].Points[0].Timestamp)
		assertEquals(t, gpxDoc.Routes[0].Points[0].Timestamp.Format(TIME_FORMAT), "2013-01-02T02:03:03Z")
		assertEquals(t, gpxDoc.Routes[0].Points[0].MagneticVariation, "1.2")
		assertEquals(t, gpxDoc.Routes[0].Points[0].GeoidHeight, "2.1")
		assertEquals(t, gpxDoc.Routes[0].Points[0].Name, "example name r")
		assertEquals(t, gpxDoc.Routes[0].Points[0].Comment, "example cmt r")
		assertEquals(t, gpxDoc.Routes[0].Points[0].Description, "example desc r")
		assertEquals(t, gpxDoc.Routes[0].Points[0].Source, "example src r")
		// TODO
		//assertEquals(t, gpxDoc.Routes[0].Points[0].Link, "http://linkrtept")
		//assertEquals(t, gpxDoc.Routes[0].Points[0].Text, "rtept link")
		assertEquals(t, gpxDoc.Routes[0].Points[0].Type, "example type r")
		assertEquals(t, gpxDoc.Routes[0].Points[0].Symbol, "example sym r")
		assertEquals(t, gpxDoc.Routes[0].Points[0].Type, "example type r")
		assertEquals(t, gpxDoc.Routes[0].Points[0].TypeOfGpsFix, "3d")
		assertEquals(t, gpxDoc.Routes[0].Points[0].Satellites, 6)
		assertEquals(t, gpxDoc.Routes[0].Points[0].HorizontalDilution, 7.0)
		assertEquals(t, gpxDoc.Routes[0].Points[0].VerticalDilution, 8.0)
		assertEquals(t, gpxDoc.Routes[0].Points[0].PositionalDilution, 9.0)
		assertEquals(t, gpxDoc.Routes[0].Points[0].AgeOfDGpsData, 10.0)
		assertEquals(t, gpxDoc.Routes[0].Points[0].DGpsId, 99)
		// TODO: Extensions

		assertEquals(t, gpxDoc.Routes[1].Name, "second route")
		assertEquals(t, gpxDoc.Routes[1].Description, "example desc 2")
		assertEquals(t, len(gpxDoc.Routes[1].Points), 2)

		// Tracks:
		assertEquals(t, len(gpxDoc.Tracks), 2)
		assertEquals(t, gpxDoc.Tracks[0].Name, "example name t")
		assertEquals(t, gpxDoc.Tracks[0].Comment, "example cmt t")
		assertEquals(t, gpxDoc.Tracks[0].Description, "example desc t")
		assertEquals(t, gpxDoc.Tracks[0].Source, "example src t")
		assertEquals(t, gpxDoc.Tracks[0].Number, 1)
		assertEquals(t, gpxDoc.Tracks[0].Type, "t")
		// TODO link

		// TODO: segment points

		assertEquals(t, len(gpxDoc.Tracks[0].Segments), 2)
	}
}

func TestParseAndReparseGPX10(t *testing.T) {
	gpxDocuments := []*GPX{}

	{
		gpxDoc, err := ParseFile("../../test_files/gpx1.0_with_all_fields.gpx")
		if err != nil || gpxDoc == nil {
			t.Error("Error parsing:" + err.Error())
		}
		gpxDocuments = append(gpxDocuments, gpxDoc)

		// Test after reparsing
		xml, err := gpxDoc.ToXml("1.0")
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

	/* TODO Other asserts
	 */
}
