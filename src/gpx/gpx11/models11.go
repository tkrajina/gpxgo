package gpx11

import (
	"encoding/xml"
)

type Gpx struct {
	XMLName      xml.Name  `xml:"http://www.topografix.com/GPX/1/1 gpx"`
	XmlNsXsi     string    `xml:"xmlns:xsi,attr,omitempty"`
	XmlSchemaLoc string    `xml:"xsi:schemaLocation,attr,omitempty"`
	Version      string    `xml:"version,attr"`
	Creator      string    `xml:"creator,attr"`
	Name         string    `xml:"metadata>name,omitempty"`
	Desc         string    `xml:"metadata>desc,omitempty"`
	AuthorName   string    `xml:"metadata>author>name,omitempty"`
	AuthorEmail  *GpxEmail `xml:"metadata>author>email,omitempty"`
	// TODO: There can be more than one link?
	AuthorLink *GpxLink       `xml:"metadata>author>link,omitempty"`
	Copyright  *GpxCopyright  `xml:"metadata>copyright,omitempty"`
	Link       *GpxLink       `xml:"metadata>link,omitempty"`
	Timestamp  string         `xml:"metadata>time,omitempty"`
	Keywords   string         `xml:"metadata>keywords,omitempty"`
	Bounds     *GpxBounds     `xml:"bounds"`
	Extensions *GpxExtensions `xml:"extensions"`
	Waypoints  []*GpxPoint    `xml:"wpt"`
	Routes     []*GpxRte      `xml:"rte"`
	Tracks     []*GpxTrk      `xml:"trk"`
}

type GpxBounds struct {
	//XMLName xml.Name `xml:"bounds"`
	MinLat float64 `xml:"minlat,attr"`
	MaxLat float64 `xml:"maxlat,attr"`
	MinLon float64 `xml:"minlon,attr"`
	MaxLon float64 `xml:"maxlon,attr"`
}

type GpxCopyright struct {
	XMLName xml.Name `xml:"copyright"`
	Author  string   `xml:"author,attr"`
	Year    string   `xml:"year,omitempty"`
	License string   `xml:"license,omitempty"`
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
	//	Copyright *GpxCopyright `xml:"copyright,omitempty"`
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
	Ele         float64 `xml:"ele,omitempty"`
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
