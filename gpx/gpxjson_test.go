package gpx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html/charset"
)

func TestJSONEmptyGPX(t *testing.T) {
	t.Parallel()

	var g GPX

	jsn, err := json.Marshal(g)
	assert.Nil(t, err, "%#v", err)

	var unmarshaled GPX
	json.Unmarshal(jsn, &unmarshaled)

	fmt.Println(unmarshaled)
}

func TestJSONEmptyGPXFromString(t *testing.T) {
	t.Parallel()

	xml := `<?xml version="1.0" encoding="UTF-8"?>
<gpx  xmlns="http://www.topografix.com/GPX/1/1" version="1.1" creator="https://github.com/tkrajina/gpxgo">
</gpx>`
	g, err := ParseString(xml)
	assert.Nil(t, err)

	jsn, err := json.MarshalIndent(g, "", "  ")
	assert.Nil(t, err, "%#v", err)

	fmt.Println(string(jsn))

	if err != nil {
		if me, is := err.(*json.MarshalerError); is {
			fmt.Println("type ", me.Type.Name())
		}
	}

	var unmarshaled GPX
	err = json.Unmarshal(jsn, &unmarshaled)
	assert.Nil(t, err)

	fmt.Println(unmarshaled)
	assert.Equal(t, cleanReparsed(*g), cleanReparsed(unmarshaled))
}

func TestJSON(t *testing.T) {
	t.Parallel()

	xml := `<?xml version="1.0" encoding="UTF-8"?>
<gpx  xmlns="http://www.topografix.com/GPX/1/1" version="1.1" creator="https://github.com/tkrajina/gpxgo">
	<metadata>
			<author></author>
	</metadata>
	<trk>
		<trkseg>
			<trkpt lat="2.1234" lon="5.1234">
				<ele>1234</ele>
			</trkpt>
			<trkpt lat="2.1233" lon="5.1235">
				<ele>1235</ele>
			</trkpt>
			<trkpt lat="2.1235" lon="5.1236">
				<ele>1236</ele>
			</trkpt>
		</trkseg>
	</trk>
</gpx>`
	g, err := ParseString(xml)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(g.Tracks))
	assert.Equal(t, 1, len(g.Tracks[0].Segments))
	assert.Equal(t, 3, len(g.Tracks[0].Segments[0].Points))

	jsn, err := json.MarshalIndent(g, "", "  ")
	assert.Nil(t, err, "%#v", err)

	fmt.Println(string(jsn))

	if err != nil {
		if me, is := err.(*json.MarshalerError); is {
			fmt.Println("type ", me.Type.Name())
		}
	}

	var unmarshaled GPX
	err = json.Unmarshal(jsn, &unmarshaled)
	assert.Nil(t, err)

	fmt.Println(unmarshaled)
	assert.Equal(t, 1, len(unmarshaled.Tracks))
	assert.Equal(t, 1, len(unmarshaled.Tracks[0].Segments))
	assert.Equal(t, 3, len(unmarshaled.Tracks[0].Segments[0].Points))
	assert.Equal(t, cleanReparsed(*g), cleanReparsed(unmarshaled))
}

func TestNullableInt(t *testing.T) {
	t.Parallel()

	var person struct {
		Name string      `json:"name"`
		Age  *NilableInt `json:"age"`
	}
	byts, err := json.Marshal(person)
	assert.Nil(t, err)
	if err != nil {
		if me, is := err.(*json.MarshalerError); is {
			fmt.Println("type ", me.Type.Name())
			fmt.Println("err:", me.Err.Error())
		}
	}
	assert.Equal(t, `{"name":"","age":null}`, string(byts))
}

func TestWithExtension(t *testing.T) {
	t.Parallel()

	xml := `<?xml version="1.0" encoding="UTF-8"?>
<gpx
  version="1.1"
  creator="Runkeeper - http://www.runkeeper.com"
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xmlns="http://www.topografix.com/GPX/1/1"
  xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd"
  xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1">
<wpt lat="37.778259000" lon="-122.391386000">
<ele>3.4</ele><time>2016-06-17T23:41:03Z</time><extensions><gpxtpx:TrackPointExtension><gpxtpx:hr>171</gpxtpx:hr></gpxtpx:TrackPointExtension></extensions>
</wpt>
</gpx>`

	g, err := ParseString(xml)
	assert.NotNil(t, g)
	assert.Nil(t, err)
	if t.Failed() {
		t.FailNow()
	}

	testGPXJSON(t, *g)
}

func TestWithGPXAttrs(t *testing.T) {
	t.Parallel()

	xmlStr := `<?xml version="1.0" encoding="UTF-8"?>
<gpx
  version="1.0"
  creator="GPSBabel - http://www.gpsbabel.org"
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xmlns="http://www.topografix.com/GPX/1/0"
  xsi:schemaLocation="http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd"></gpx>`

	{
		g := &gpx10Gpx{}
		decoder := xml.NewDecoder(bytes.NewBufferString(xmlStr))
		decoder.CharsetReader = charset.NewReaderLabel
		assert.Nil(t, decoder.Decode(&g))
		fmt.Println(g.Attrs)
	}

	g, err := ParseString(xmlStr)
	assert.NotNil(t, g)
	assert.Nil(t, err)
	if t.Failed() {
		t.FailNow()
	}

	testGPXJSON(t, *g)
}

func TestParseGPXAttrs(t *testing.T) {
	t.Parallel()

	xmlStr := `<?xml version="1.0" encoding="UTF-8"?>
<gpx
  version="1.0"
  creator="GPSBabel - http://www.gpsbabel.org"
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xmlns="http://www.topografix.com/GPX/1/0"
  xsi:schemaLocation="http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd"></gpx>`

	g, err := ParseString(xmlStr)
	assert.Nil(t, err)
	assert.NotNil(t, g)

	xml, _ := g.ToXml(ToXmlParams{})
	fmt.Println("xml=", string(xml))
}

func TestFileGPX(t *testing.T) {
	t.Parallel()

	g, err := ParseString(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<gpx xmlns="http://www.topografix.com/GPX/1/1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3" xmlns:wptx1="http://www.garmin.com/xmlschemas/WaypointExtension/v1" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" creator="eTrex 10" version="1.1" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www8.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/WaypointExtension/v1 http://www8.garmin.com/xmlschemas/WaypointExtensionv1.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd">
	  <trk>
		<name>17-MRZ-12 16:44:12</name>
		<extensions>
		  <gpxx:TrackExtension>
			<gpxx:DisplayColor>Cyan</gpxx:DisplayColor>
		  </gpxx:TrackExtension>
		</extensions>
	  </trk>
	</gpx>
	`)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(g.Tracks[0].Extensions))

	xml, _ := g.ToXml(ToXmlParams{})
	g2, err := ParseBytes(xml)
	assert.Nil(t, err)
	assert.Equal(t, "1.1", g2.Version)
	assert.Equal(t, 1, len(g2.Tracks[0].Extensions), jsonize(g2.Tracks[0].Extensions))

	// assert.Nil(t, err)
	// testGPXJSON(t, *g)
}

func TestGPXJSONForAllTestFiles(t *testing.T) {
	t.Parallel()

	testFilesDir := "../test_files/"

	list, err := os.ReadDir(testFilesDir)
	assert.Nil(t, err, "%v", err)
	for _, entry := range list {
		if strings.HasSuffix(entry.Name(), ".gz") {
			continue
		}
		fmt.Println("Testing", entry.Name())
		g, err := ParseFile(testFilesDir + entry.Name())
		assert.Nil(t, err)
		testGPXJSON(t, *g)
	}
}

func testGPXJSON(t *testing.T, g GPX) {
	fmt.Println("gpx:", jsonizeFormatted(g))

	fmt.Println("attrs1", jsonizeFormatted(g.Attrs))

	reparsedFromXML, err := reparse(g)
	assert.Nil(t, err)

	fmt.Println("reparsed:", jsonizeFormatted(reparsedFromXML.Attrs))
	assert.Equal(t, jsonizeFormatted(cleanReparsed(g)), jsonizeFormatted(cleanReparsed(*reparsedFromXML)))
	if t.Failed() {
		t.FailNow()
	}

	assert.Equal(t, cleanReparsedAttrs(g.Attrs), cleanReparsedAttrs(reparsedFromXML.Attrs))
	if t.Failed() {
		t.FailNow()
	}

	assert.Equal(t, cleanReparsed(g), cleanReparsed(*reparsedFromXML))

	if t.Failed() {
		t.FailNow()
	}

	var unmarshaled GPX
	err = json.Unmarshal([]byte(jsonizeFormatted(reparsedFromXML)), &unmarshaled)
	assert.Nil(t, err)

	fmt.Println("attrs unmarshalled", jsonizeFormatted(g.Attrs))

	if t.Failed() {
		t.FailNow()
	}

	fmt.Println("unmarshalled from reparsed:", jsonizeFormatted(unmarshaled))

	assert.Equal(t, cleanReparsed(g), cleanReparsed(unmarshaled))
	if t.Failed() {
		t.FailNow()
	}
}

func TestUnmarshalGPXJSONWithNSAttrs(t *testing.T) {
	t.Parallel()

	var g GPX
	assert.Nil(t, json.Unmarshal([]byte(`{
		"xmlns": "http://www.topografix.com/GPX/1/1",
		"attrs": {
			"nsattrs": {
				"xmlns": {
					"gpxtpx": {
						"space": "xmlns",
						"local": "gpxtpx",
						"value": "http://www.garmin.com/xmlschemas/TrackPointExtension/v1"
					},
					"xsi": {
						"space": "xmlns",
						"local": "xsi",
						"value": "http://www.w3.org/2001/XMLSchema-instance"
					}
				},
				"xsi": {
					"schemaLocation": {
						"space": "http://www.w3.org/2001/XMLSchema-instance",
						"local": "schemaLocation",
						"value": "http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd"
					}
				}
			}
		},
		"waypoints": [
			{
				"lat": 37.778259,
				"lon": -122.391386,
				"ext": [
					{
						"ns": "http://www.garmin.com/xmlschemas/TrackPointExtension/v1",
						"name": "TrackPointExtension",
						"nodes": [
							{
								"ns": "http://www.garmin.com/xmlschemas/TrackPointExtension/v1",
								"name": "hr",
								"data": "171"
							}
						]
					}
				]
			}
		],
		"trk": null
	}`), &g))

	assert.Equal(t, 1, len(g.Waypoints))
	assert.Equal(t, 1, len(g.Waypoints[0].Extensions))
	{
		node, found := g.Waypoints[0].Extensions.GetNode(AnyNamespace, "TrackPointExtension")
		assert.True(t, found)
		hrNode, found := node.GetNode("hr")
		assert.True(t, found)
		assert.Equal(t, hrNode.Data, "171")
	}
	{
		node, found := g.Waypoints[0].Extensions.GetNode(NamespaceURL("http://www.garmin.com/xmlschemas/TrackPointExtension/v1"), "TrackPointExtension")
		assert.True(t, found)
		hrNode, found := node.GetNode("hr")
		assert.True(t, found)
		assert.Equal(t, hrNode.Data, "171")
	}
	{
		_, found := g.Waypoints[0].Extensions.GetNode(NamespaceURL("jkl"), "TrackPointExtension")
		assert.False(t, found)
	}

	spew.Dump(g)

	xml, err := g.ToXml(ToXmlParams{Indent: true})
	assert.Nil(t, err)
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<gpx
version="1.1"
creator="Runkeeper - http://www.runkeeper.com"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns="http://www.topografix.com/GPX/1/1"
xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd"
xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1">
	<wpt lat="37.778259000" lon="-122.391386000">
		<ele>3.4</ele>
		<time>2016-06-17T23:41:03Z</time>
		<extensions>
			<gpxtpx:TrackPointExtension>
				<gpxtpx:hr>171</gpxtpx:hr>
			</gpxtpx:TrackPointExtension>
		</extensions>
	</wpt>
</gpx>`, string(xml))

	if t.Failed() {
		t.FailNow()
	}
}

func TestUnmarshalGPXJSON(t *testing.T) {
	t.Parallel()

	var g GPX
	assert.Nil(t, json.Unmarshal([]byte(`{
		"xmlns": "http://www.topografix.com/GPX/1/1",
		"attrs": {
			"nsattrs": {
				"xmlns": {
					"gpxtpx": {
						"space": "xmlns",
						"local": "gpxtpx",
						"value": "http://www.garmin.com/xmlschemas/TrackPointExtension/v1"
					},
					"xsi": {
						"space": "xmlns",
						"local": "xsi",
						"value": "http://www.w3.org/2001/XMLSchema-instance"
					}
				},
				"xsi": {
					"schemaLocation": {
						"space": "http://www.w3.org/2001/XMLSchema-instance",
						"local": "schemaLocation",
						"value": "http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd"
					}
				}
			}
		},
		"version": "1.1",
		"creator": "Runkeeper - http://www.runkeeper.com",
		"time": null,
		"waypoints": [
			{
				"lat": 37.778259,
				"lon": -122.391386,
				"ele": 3.4,
				"time": "2016-06-17T23:41:03Z",
				"ext": [
					{
						"ns": "http://www.garmin.com/xmlschemas/TrackPointExtension/v1",
						"name": "TrackPointExtension",
						"nodes": [
							{
								"ns": "http://www.garmin.com/xmlschemas/TrackPointExtension/v1",
								"name": "hr",
								"data": "171"
							}
						]
					}
				]
			}
		],
		"trk": null
	}`), &g))

	xml, err := g.ToXml(ToXmlParams{Indent: true})
	assert.Nil(t, err)
	assert.Equal(t, ``, string(xml))

	if t.Failed() {
		t.FailNow()
	}

	assert.Equal(t, 1, len(g.Waypoints))
	assert.Equal(t, 1, len(g.Waypoints[0].Extensions))
	{
		node, found := g.Waypoints[0].Extensions.GetNode(AnyNamespace, "TrackPointExtension")
		assert.True(t, found)
		hrNode, found := node.GetNode("hr")
		assert.True(t, found)
		assert.Equal(t, hrNode.Data, "171")
	}
	{
		node, found := g.Waypoints[0].Extensions.GetNode(NamespaceURL("http://www.garmin.com/xmlschemas/TrackPointExtension/v1"), "TrackPointExtension")
		assert.True(t, found)
		hrNode, found := node.GetNode("hr")
		assert.True(t, found)
		assert.Equal(t, hrNode.Data, "171")
	}
	{
		_, found := g.Waypoints[0].Extensions.GetNode(NamespaceURL("jkl"), "TrackPointExtension")
		assert.False(t, found)
	}
}

func cleanReparsed(g GPX) GPX {
	g.Attrs = cleanReparsedAttrs(g.Attrs)
	if g.Creator == "" {
		g.Creator = defaultCreator
	}
	if g.Version == "" {
		g.Version = "1.1"
	}
	return g
}

func cleanReparsedAttrs(attrs GPXAttributes) GPXAttributes {
	return attrs
}

func jsonizeFormatted(a any) string {
	jsn, _ := json.MarshalIndent(a, "", "  ")
	return string(jsn)
}

func jsonize(a any) string {
	jsn, _ := json.Marshal(a)
	return string(jsn)
}
