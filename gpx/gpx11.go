// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"encoding/xml"
	"strings"
)

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Data    string     `xml:",chardata"`
	Nodes   []Node     `xml:",any"`
}

func (n Node) debugXMLChunk() []byte {
	byts, err := xml.MarshalIndent(n, "", "    ")
	if err != nil {
		return []byte("???")
	}
	return byts
}

func (n Node) toTokens(prefix string) (tokens []xml.Token) {
	var attrs []xml.Attr
	for _, a := range n.Attrs {
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: prefix + a.Name.Local}, Value: a.Value})
	}

	start := xml.StartElement{Name: xml.Name{Local: prefix + n.XMLName.Local, Space: ""}, Attr: attrs}
	tokens = append(tokens, start)
	data := strings.TrimSpace(n.Data)
	if len(n.Nodes) > 0 {
		for _, node := range n.Nodes {
			tokens = append(tokens, node.toTokens(prefix)...)
		}
	} else if data != "" {
		tokens = append(tokens, xml.CharData(data))
	} else {
		return nil
	}
	tokens = append(tokens, xml.EndElement{start.Name})
	return
}

func (n *Node) GetAttr(key string) (value string, found bool) {
	for i := range n.Attrs {
		if n.Attrs[i].Name.Local == key {
			value = n.Attrs[i].Value
			found = true
			return
		}
	}
	return
}

func (n *Node) SetAttr(key, value string) {
	for i := range n.Attrs {
		if n.Attrs[i].Name.Local == key {
			n.Attrs[i].Value = value
			return
		}
	}
	n.Attrs = append(n.Attrs, xml.Attr{
		Name: xml.Name{
			Space: n.SpaceNameURL(),
			Local: key,
		},
		Value: value,
	})
}

func (n *Node) GetNode(path0 string) (node *Node, found bool) {
	for subn := range n.Nodes {
		if n.Nodes[subn].LocalName() == path0 {
			node = &n.Nodes[subn]
			found = true
			return
		}
	}
	return
}

func (n *Node) getOrCreateNode(path ...string) *Node {
	if len(path) == 0 {
		return n
	}

	path0, rest := path[0], path[1:]

	subNode, found := n.GetNode(path0)
	if !found {
		n.Nodes = append(n.Nodes, Node{
			XMLName: xml.Name{
				Space: n.XMLName.Space,
				Local: path0,
			},
			Attrs: nil,
		})
		subNode = &(n.Nodes[len(n.Nodes)-1])
	}

	return subNode.getOrCreateNode(rest...)
}

func (n Node) IsEmpty() bool        { return len(n.Nodes) == 0 && len(n.Attrs) == 0 && len(n.Data) == 0 }
func (n Node) LocalName() string    { return n.XMLName.Local }
func (n Node) SpaceNameURL() string { return n.XMLName.Space }
func (n Node) GetAttrOrEmpty(attr string) string {
	val, _ := n.GetAttr(attr)
	return val
}

type Extension struct {
	// XMLName xml.Name
	// Attrs   []xml.Attr `xml:",any,attr"`
	Nodes []Node `xml:",any"`

	// Filled before deserializing:
	globalNsAttrs map[string]Attr
}

var _ xml.Marshaler = Extension{}

func (ex Extension) debugXMLChunk() []byte {
	byts, err := xml.MarshalIndent(ex, "", "    ")
	if err != nil {
		return []byte("???")
	}
	return byts
}

func (ex Extension) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(ex.Nodes) == 0 {
		return nil
	}

	start = xml.StartElement{Name: xml.Name{Local: start.Name.Local}, Attr: nil}
	tokens := []xml.Token{start}
	for _, node := range ex.Nodes {
		prefix := ""
		for _, v := range ex.globalNsAttrs {
			if node.SpaceNameURL() == v.Value || node.SpaceNameURL() == v.Name.Local {
				prefix = v.replacement
			}
		}
		tokens = append(tokens, node.toTokens(prefix)...)
	}

	tokens = append(tokens, xml.EndElement{Name: start.Name})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	err := e.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (ex *Extension) GetOrCreateNode(namespaceURL string, path ...string) *Node {
	// TODO: Check is len(nodes) == 0
	var subNode *Node
	for n := range ex.Nodes {
		if ex.Nodes[n].SpaceNameURL() == namespaceURL && ex.Nodes[n].LocalName() == path[0] {
			subNode = &ex.Nodes[n]
			break
		}
	}
	if subNode == nil {
		ex.Nodes = append(ex.Nodes, Node{
			XMLName: xml.Name{
				Space: namespaceURL,
				Local: path[0],
			},
		})
		subNode = &ex.Nodes[len(ex.Nodes)-1]
	}
	return subNode.getOrCreateNode(path[1:]...)
}

func (ex *Extension) GetNode(path0 string) (node *Node, found bool) {
	for subn := range ex.Nodes {
		if ex.Nodes[subn].LocalName() == path0 {
			node = &ex.Nodes[subn]
			found = true
			return
		}
	}
	return
}

/*

The GPX XML hierarchy:

gpx (gpxType)
    - attr: version (xsd:string) None
    - attr: creator (xsd:string) None
    metadata (metadataType)
        name (xsd:string)
        desc (xsd:string)
        author (personType)
            name (xsd:string)
            email (emailType)
                - attr: id (xsd:string) None
                - attr: domain (xsd:string) None
            link (linkType)
                - attr: href (xsd:anyURI) None
                text (xsd:string)
                type (xsd:string)
        copyright (copyrightType)
            - attr: author (xsd:string) None
            year (xsd:gYear)
            license (xsd:anyURI)
        link (linkType)
            - attr: href (xsd:anyURI) None
            text (xsd:string)
            type (xsd:string)
        time (xsd:dateTime)
        keywords (xsd:string)
        bounds (boundsType)
            - attr: minlat (latitudeType) None
            - attr: minlon (longitudeType) None
            - attr: maxlat (latitudeType) None
            - attr: maxlon (longitudeType) None
        extensions (extensionsType)
    wpt (wptType)
        - attr: lat (latitudeType) None
        - attr: lon (longitudeType) None
        ele (xsd:decimal)
        time (xsd:dateTime)
        magvar (degreesType)
        geoidheight (xsd:decimal)
        name (xsd:string)
        cmt (xsd:string)
        desc (xsd:string)
        src (xsd:string)
        link (linkType)
            - attr: href (xsd:anyURI) None
            text (xsd:string)
            type (xsd:string)
        sym (xsd:string)
        type (xsd:string)
        fix (fixType)
        sat (xsd:nonNegativeInteger)
        hdop (xsd:decimal)
        vdop (xsd:decimal)
        pdop (xsd:decimal)
        ageofdgpsdata (xsd:decimal)
        dgpsid (dgpsStationType)
        extensions (extensionsType)
    rte (rteType)
        name (xsd:string)
        cmt (xsd:string)
        desc (xsd:string)
        src (xsd:string)
        link (linkType)
            - attr: href (xsd:anyURI) None
            text (xsd:string)
            type (xsd:string)
        number (xsd:nonNegativeInteger)
        type (xsd:string)
        extensions (extensionsType)
        rtept (wptType)
            - attr: lat (latitudeType) None
            - attr: lon (longitudeType) None
            ele (xsd:decimal)
            time (xsd:dateTime)
            magvar (degreesType)
            geoidheight (xsd:decimal)
            name (xsd:string)
            cmt (xsd:string)
            desc (xsd:string)
            src (xsd:string)
            link (linkType)
                - attr: href (xsd:anyURI) None
                text (xsd:string)
                type (xsd:string)
            sym (xsd:string)
            type (xsd:string)
            fix (fixType)
            sat (xsd:nonNegativeInteger)
            hdop (xsd:decimal)
            vdop (xsd:decimal)
            pdop (xsd:decimal)
            ageofdgpsdata (xsd:decimal)
            dgpsid (dgpsStationType)
            extensions (extensionsType)
    trk (trkType)
        name (xsd:string)
        cmt (xsd:string)
        desc (xsd:string)
        src (xsd:string)
        link (linkType)
            - attr: href (xsd:anyURI) None
            text (xsd:string)
            type (xsd:string)
        number (xsd:nonNegativeInteger)
        type (xsd:string)
        extensions (extensionsType)
        trkseg (trksegType)
            trkpt (wptType)
                - attr: lat (latitudeType) None
                - attr: lon (longitudeType) None
                ele (xsd:decimal)
                time (xsd:dateTime)
                magvar (degreesType)
                geoidheight (xsd:decimal)
                name (xsd:string)
                cmt (xsd:string)
                desc (xsd:string)
                src (xsd:string)
                link (linkType)
                    - attr: href (xsd:anyURI) None
                    text (xsd:string)
                    type (xsd:string)
                sym (xsd:string)
                type (xsd:string)
                fix (fixType)
                sat (xsd:nonNegativeInteger)
                hdop (xsd:decimal)
                vdop (xsd:decimal)
                pdop (xsd:decimal)
                ageofdgpsdata (xsd:decimal)
                dgpsid (dgpsStationType)
                extensions (extensionsType)
            extensions (extensionsType)
    extensions (extensionsType)
*/

type gpx11Gpx struct {
	XMLName      xml.Name   `xml:"gpx"`
	Attrs        []xml.Attr `xml:",any,attr"`
	XMLNs        string     `xml:"xmlns,attr,omitempty"`
	XmlNsXsi     string     `xml:"xmlns:xsi,attr,omitempty"`
	XmlSchemaLoc string     `xml:"xsi:schemaLocation,attr,omitempty"`

	Version     string         `xml:"version,attr"`
	Creator     string         `xml:"creator,attr"`
	Name        string         `xml:"metadata>name,omitempty"`
	Desc        string         `xml:"metadata>desc,omitempty"`
	AuthorName  string         `xml:"metadata>author>name,omitempty"`
	AuthorEmail *gpx11GpxEmail `xml:"metadata>author>email,omitempty"`
	// TODO: There can be more than one link?
	AuthorLink *gpx11GpxLink      `xml:"metadata>author>link,omitempty"`
	Copyright  *gpx11GpxCopyright `xml:"metadata>copyright,omitempty"`
	Link       *gpx11GpxLink      `xml:"metadata>link,omitempty"`
	Timestamp  string             `xml:"metadata>time,omitempty"`
	Keywords   string             `xml:"metadata>keywords,omitempty"`
	Extensions Extension          `xml:"metadata>extensions"`
	Bounds     *gpx11GpxBounds    `xml:"bounds"`
	Waypoints  []*gpx11GpxPoint   `xml:"wpt"`
	Routes     []*gpx11GpxRte     `xml:"rte"`
	Tracks     []*gpx11GpxTrk     `xml:"trk"`
}

type gpx11GpxBounds struct {
	//XMLName xml.Name `xml:"bounds"`
	MinLat float64 `xml:"minlat,attr"`
	MaxLat float64 `xml:"maxlat,attr"`
	MinLon float64 `xml:"minlon,attr"`
	MaxLon float64 `xml:"maxlon,attr"`
}

type gpx11GpxCopyright struct {
	XMLName xml.Name `xml:"copyright"`
	Author  string   `xml:"author,attr"`
	Year    string   `xml:"year,omitempty"`
	License string   `xml:"license,omitempty"`
}

//type gpx11GpxAuthor struct {
//	Name  string        `xml:"name,omitempty"`
//	Email string        `xml:"email,omitempty"`
//	Link  *gpx11GpxLink `xml:"link"`
//}

type gpx11GpxEmail struct {
	Id     string `xml:"id,attr"`
	Domain string `xml:"domain,attr"`
}

type gpx11GpxLink struct {
	Href string `xml:"href,attr"`
	Text string `xml:"text,omitempty"`
	Type string `xml:"type,omitempty"`
}

//type gpx11GpxMetadata struct {
//	XMLName xml.Name        `xml:"metadata"`
//	Name    string          `xml:"name,omitempty"`
//	Desc    string          `xml:"desc,omitempty"`
//	Author  *gpx11GpxAuthor `xml:"author,omitempty"`
//	//	Copyright *GpxCopyright `xml:"copyright,omitempty"`
//	//	Links     []GpxLink     `xml:"link"`
//	Timestamp string `xml:"time,omitempty"`
//	Keywords  string `xml:"keywords,omitempty"`
//	//	Bounds    *GpxBounds    `xml:"bounds"`
//}

/**
 * Common struct fields for all points
 */
type gpx11GpxPoint struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	// Position info
	Ele         NullableFloat64 `xml:"ele,omitempty"`
	Timestamp   string          `xml:"time,omitempty"`
	MagVar      string          `xml:"magvar,omitempty"`
	GeoIdHeight string          `xml:"geoidheight,omitempty"`
	// Description info
	Name  string         `xml:"name,omitempty"`
	Cmt   string         `xml:"cmt,omitempty"`
	Desc  string         `xml:"desc,omitempty"`
	Src   string         `xml:"src,omitempty"`
	Links []gpx11GpxLink `xml:"link"`
	Sym   string         `xml:"sym,omitempty"`
	Type  string         `xml:"type,omitempty"`
	// Accuracy info
	Fix           string    `xml:"fix,omitempty"`
	Sat           *int      `xml:"sat,omitempty"`
	Hdop          *float64  `xml:"hdop,omitempty"`
	Vdop          *float64  `xml:"vdop,omitempty"`
	Pdop          *float64  `xml:"pdop,omitempty"`
	AgeOfDGpsData *float64  `xml:"ageofdgpsdata,omitempty"`
	DGpsId        *int      `xml:"dgpsid,omitempty"`
	Extensions    Extension `xml:"extensions"`
}

type gpx11GpxRte struct {
	XMLName xml.Name `xml:"rte"`
	Name    string   `xml:"name,omitempty"`
	Cmt     string   `xml:"cmt,omitempty"`
	Desc    string   `xml:"desc,omitempty"`
	Src     string   `xml:"src,omitempty"`
	// TODO
	//Links       []Link   `xml:"link"`
	Number     NullableInt      `xml:"number,omitempty"`
	Type       string           `xml:"type,omitempty"`
	Points     []*gpx11GpxPoint `xml:"rtept"`
	Extensions Extension        `xml:"extensions"`
}

type gpx11GpxTrkSeg struct {
	XMLName    xml.Name         `xml:"trkseg"`
	Points     []*gpx11GpxPoint `xml:"trkpt"`
	Extensions Extension        `xml:"extensions"`
}

// Trk is a GPX track
type gpx11GpxTrk struct {
	XMLName xml.Name `xml:"trk"`
	Name    string   `xml:"name,omitempty"`
	Cmt     string   `xml:"cmt,omitempty"`
	Desc    string   `xml:"desc,omitempty"`
	Src     string   `xml:"src,omitempty"`
	// TODO
	//Links    []Link   `xml:"link"`
	Number     NullableInt       `xml:"number,omitempty"`
	Type       string            `xml:"type,omitempty"`
	Segments   []*gpx11GpxTrkSeg `xml:"trkseg,omitempty"`
	Extensions Extension         `xml:"extensions"`
}
