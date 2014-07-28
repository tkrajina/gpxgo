package gpx11

import (
    "encoding/xml"
)

type Gpx struct {
	XMLName      xml.Name     `xml:"http://www.topografix.com/GPX/1/1 gpx"`
	XmlNsXsi     string       `xml:"xmlns:xsi,attr,omitempty"`
	XmlSchemaLoc string       `xml:"xsi:schemaLocation,attr,omitempty"`
	Version      string       `xml:"version,attr"`
	Creator      string       `xml:"creator,attr"`
	Name      string            `xml:"metadata>name,omitempty"`
	Desc      string            `xml:"metadata>desc,omitempty"`
    AuthorName string           `xml:"metadata>author>name,omitempty"`
    AuthorEmail *GpxEmail `xml:"metadata>author>email,omitempty"`
    // TODO: There can be more than one link?
    AuthorLink *GpxLink       `xml:"metadata>author>link,omitempty"`
	Copyright *GpxCopyright `xml:"metadata>copyright,omitempty"`
	Extensions   *GpxExtensions `xml:"extensions"`
	//Metadata     *GpxMetadata `xml:"metadata,omitempty"`
//	Waypoints    []GpxWpt     `xml:"wpt"`
//	Routes       []GpxRte     `xml:"rte"`
//	Tracks       []GpxTrk     `xml:"trk"`
}

type GpxCopyright struct {
	XMLName xml.Name `xml:"copyright"`
	Author  string   `xml:"author,attr"`
	Year    string   `xml:"year,omitempty"`
	License string   `xml:"license,omitempty"`
}

type GpxAuthor struct {
	Name      string       `xml:"name,omitempty"`
	Email      string       `xml:"email,omitempty"`
	Link      *GpxLink       `xml:"link"`
}

type GpxEmail struct {
    Id string `xml:"id,attr"`
    Domain string `xml:"domain,attr"`
}

type GpxLink struct {
    Href string `xml:"href,attr"`
    Text string `xml:"text,omitempty"`
    Type string `xml:"type,omitempty"`
}

type GpxMetadata struct {
	XMLName   xml.Name      `xml:"metadata"`
	Name      string        `xml:"name,omitempty"`
	Desc      string        `xml:"desc,omitempty"`
	Author    *GpxAuthor    `xml:"author,omitempty"`
//	Copyright *GpxCopyright `xml:"copyright,omitempty"`
//	Links     []GpxLink     `xml:"link"`
	Timestamp string        `xml:"time,omitempty"`
	Keywords  string        `xml:"keywords,omitempty"`
//	Bounds    *GpxBounds    `xml:"bounds"`
}

type GpxExtensions struct {
    Bytes []byte        `xml:",innerxml"`
}

func NewGpx() (*Gpx) {
    return new(Gpx)
}
