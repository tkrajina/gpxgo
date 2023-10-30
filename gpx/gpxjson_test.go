package gpx

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

	reparsedFromXML, err := reparse(g)
	assert.Nil(t, err)

	assert.Equal(t, cleanReparsed(g), cleanReparsed(*reparsedFromXML))

	if t.Failed() {
		t.FailNow()
	}

	jsn, err := json.MarshalIndent(g, "", "  ")
	assert.Nil(t, err, "%#v", err)
	fmt.Println(string(jsn))

	var unmarshaled GPX
	err = json.Unmarshal(jsn, &unmarshaled)
	assert.Nil(t, err)

	assert.Equal(t, cleanReparsed(g), cleanReparsed(unmarshaled))
	if t.Failed() {
		t.FailNow()
	}
}

func cleanReparsed(g GPX) GPX {
	if len(g.Attrs.NamespaceAttributes) == 0 {
		g.Attrs.NamespaceAttributes = nil
		return g
	}
	for k, v := range g.Attrs.NamespaceAttributes {
		for k2, v2 := range v {
			v2.replacement = ""
			g.Attrs.NamespaceAttributes[k][k2] = v2
		}
	}
	return g
}
