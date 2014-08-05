package gpx10

import (
	"encoding/xml"
)

/*

The GPX XML hierarchy:

gpx
    - attr: version (xsd:string) required
    - attr: creator (xsd:string) required
    name
    desc
    author
    email
    url
    urlname
    time
    keywords
    bounds
    wpt
        - attr: lat (gpx:latitudeType) required
        - attr: lon (gpx:longitudeType) required
        ele
        time
        magvar
        geoidheight
        name
        cmt
        desc
        src
        url
        urlname
        sym
        type
        fix
        sat
        hdop
        vdop
        pdop
        ageofdgpsdata
        dgpsid
    rte
        name
        cmt
        desc
        src
        url
        urlname
        number
        rtept
            - attr: lat (gpx:latitudeType) required
            - attr: lon (gpx:longitudeType) required
            ele
            time
            magvar
            geoidheight
            name
            cmt
            desc
            src
            url
            urlname
            sym
            type
            fix
            sat
            hdop
            vdop
            pdop
            ageofdgpsdata
            dgpsid
    trk
        name
        cmt
        desc
        src
        url
        urlname
        number
        trkseg
            trkpt
                - attr: lat (gpx:latitudeType) required
                - attr: lon (gpx:longitudeType) required
                ele
                time
                course
                speed
                magvar
                geoidheight
                name
                cmt
                desc
                src
                url
                urlname
                sym
                type
                fix
                sat
                hdop
                vdop
                pdop
                ageofdgpsdata
                dgpsid
*/

type Gpx struct {
	XMLName xml.Name `xml:"http://www.topografix.com/GPX/1/0 gpx"`
	//XMLNs        string    `xml:"xmlns,attr"`
	//XmlNsXsi     string      `xml:"xmlns:xsi,attr,omitempty"`
	//XmlSchemaLoc string      `xml:"xsi:schemaLocation,attr,omitempty"`
	Version   string      `xml:"version,attr"`
	Creator   string      `xml:"creator,attr"`
	Name      string      `xml:"name,omitempty"`
	Desc      string      `xml:"desc,omitempty"`
	Author    string      `xml:"author,omitempty"`
	Email     string      `xml:"email,omitempty"`
	Url       string      `xml:"url,omitempty"`
	UrlName   string      `xml:"urlname,omitempty"`
	Time      string      `xml:"time,omitempty"`
	Keywords  string      `xml:"keywords,omitempty"`
	Bounds    *GpxBounds  `xml:"bounds"`
	Waypoints []*GpxPoint `xml:"wpt"`
	Routes    []*GpxRte   `xml:"rte"`
	Tracks    []*GpxTrk   `xml:"trk"`
}

type GpxBounds struct {
	//XMLName xml.Name `xml:"bounds"`
	MinLat float64 `xml:"minlat,attr"`
	MaxLat float64 `xml:"maxlat,attr"`
	MinLon float64 `xml:"minlon,attr"`
	MaxLon float64 `xml:"maxlon,attr"`
}

type GpxAuthor struct {
	Name  string   `xml:"name,omitempty"`
	Email string   `xml:"email,omitempty"`
	Link  *GpxLink `xml:"link"`
}

type GpxEmail struct {
	Id     string `xml:"id,attr"`
	Domain string `xml:"domain,attr"`
}

type GpxLink struct {
	Href string `xml:"href,attr"`
	Text string `xml:"text,omitempty"`
	Type string `xml:"type,omitempty"`
}

type GpxMetadata struct {
	XMLName xml.Name   `xml:"metadata"`
	Name    string     `xml:"name,omitempty"`
	Desc    string     `xml:"desc,omitempty"`
	Author  *GpxAuthor `xml:"author,omitempty"`
	//	Links     []GpxLink     `xml:"link"`
	Timestamp string `xml:"time,omitempty"`
	Keywords  string `xml:"keywords,omitempty"`
	//	Bounds    *GpxBounds    `xml:"bounds"`
}

type GpxExtensions struct {
	Bytes []byte `xml:",innerxml"`
}

/**
 * Common struct fields for all points
 */
type GpxPoint struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	// Position info
	Ele         *float64 `xml:"ele,omitempty"`
	Timestamp   string  `xml:"time,omitempty"`
	MagVar      string  `xml:"magvar,omitempty"`
	GeoIdHeight string  `xml:"geoidheight,omitempty"`
	// Description info
	Name  string    `xml:"name,omitempty"`
	Cmt   string    `xml:"cmt,omitempty"`
	Desc  string    `xml:"desc,omitempty"`
	Src   string    `xml:"src,omitempty"`
	Links []GpxLink `xml:"link"`
	Sym   string    `xml:"sym,omitempty"`
	Type  string    `xml:"type,omitempty"`
	// Accuracy info
	Fix           string  `xml:"fix,omitempty"`
	Sat           int     `xml:"sat,omitempty"`
	Hdop          float64 `xml:"hdop,omitempty"`
	Vdop          float64 `xml:"vdop,omitempty"`
	Pdop          float64 `xml:"pdop,omitempty"`
	AgeOfDGpsData float64 `xml:"ageofdgpsdata,omitempty"`
	DGpsId        int     `xml:"dgpsid,omitempty"`

	// Those two values are here for simplicity, but they are available only when this is part of a track segment (not route or waypoint)!
	Course string `xml:"course,omitempty"`
	Speed  string `speed:"fix,omitempty"`
}

type GpxRte struct {
	XMLName xml.Name `xml:"rte"`
	Name    string   `xml:"name,omitempty"`
	Cmt     string   `xml:"cmt,omitempty"`
	Desc    string   `xml:"desc,omitempty"`
	Src     string   `xml:"src,omitempty"`
	// TODO
	//Links       []Link   `xml:"link"`
	Number int         `xml:"number,omitempty"`
	Type   string      `xml:"type,omitempty"`
	Points []*GpxPoint `xml:"rtept"`
}

type GpxTrkSeg struct {
	XMLName xml.Name    `xml:"trkseg"`
	Points  []*GpxPoint `xml:"trkpt"`
}

// Trk is a GPX track
type GpxTrk struct {
	XMLName xml.Name `xml:"trk"`
	Name    string   `xml:"name,omitempty"`
	Cmt     string   `xml:"cmt,omitempty"`
	Desc    string   `xml:"desc,omitempty"`
	Src     string   `xml:"src,omitempty"`
	// TODO
	//Links    []Link   `xml:"link"`
	Number   int          `xml:"number,omitempty"`
	Type     string       `xml:"type,omitempty"`
	Segments []*GpxTrkSeg `xml:"trkseg,omitempty"`
}

func NewGpx() *Gpx {
	return new(Gpx)
}
