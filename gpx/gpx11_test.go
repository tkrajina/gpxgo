package gpx

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJustGPX11(t *testing.T) {
	t.Parallel()

	byts := []byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</gpx>`)

	gpxDoc, err := ParseBytes(byts)
	assert.Nil(t, err)

	var g gpx11Gpx
	assert.Nil(t, xml.Unmarshal(byts, &g))

	g.Waypoints[0].Extensions.gpx = gpxDoc

	byts2, err := xml.MarshalIndent(g, "", "\t")
	assert.Nil(t, err)
	fmt.Println("----------------------------------------------------------------------------------------------------")
	fmt.Println(string(byts))
	fmt.Println("----------------------------------------------------------------------------------------------------")
	fmt.Println(string(byts2))
	fmt.Println("----------------------------------------------------------------------------------------------------")

	for _, a := range g.Attrs {
		fmt.Printf("attr %s.%s: %s\n", a.Name.Space, a.Name.Local, a.Value)
	}
	fmt.Println("XMLNs=", g.XMLNs)
	fmt.Println("XmlNsXsi=", g.XmlNsXsi)
	fmt.Println("XmlSchemaLoc=", g.XmlSchemaLoc)
	fmt.Println("Version=", g.Version)
	fmt.Println("Creator=", g.Creator)

	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
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
	</gpx>`, string(byts))
}
