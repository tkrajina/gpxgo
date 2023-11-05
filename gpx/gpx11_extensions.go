package gpx

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type gpx11ExtensionNode struct {
	XMLName xml.Name
	Attrs   []xml.Attr           `xml:",any,attr"`
	Data    string               `xml:",chardata"`
	Nodes   []gpx11ExtensionNode `xml:",any"`
}

func (n gpx11ExtensionNode) LocalName() string    { return n.XMLName.Local }
func (n gpx11ExtensionNode) SpaceNameURL() string { return n.XMLName.Space }

func (n gpx11ExtensionNode) debugXMLChunk() []byte {
	byts, err := xml.MarshalIndent(n, "", "    ")
	if err != nil {
		return []byte("???")
	}
	return byts
}

func (n gpx11ExtensionNode) toTokens(prefix string) (tokens []xml.Token) {
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

type gpx11Extension struct {
	// XMLName xml.Name
	// Attrs   []xml.Attr `xml:",any,attr"`
	Nodes []gpx11ExtensionNode `xml:",any"`

	gpx *GPX `xml:"-"`
}

// var _ xml.Marshaler = gpx11Extension{}

func (ex gpx11Extension) debugXMLChunk() []byte {
	byts, err := xml.MarshalIndent(ex, "", "    ")
	if err != nil {
		return []byte("???")
	}
	return byts
}

func (ex gpx11Extension) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(ex.Nodes) == 0 {
		return nil
	}

	start = xml.StartElement{Name: xml.Name{Local: start.Name.Local}, Attr: nil}
	tokens := []xml.Token{start}
	for _, node := range ex.Nodes {
		nsByURLs := ex.gpx.Attrs.GetNamespacesByURLs()
		prefix := ""
		fmt.Println("find prefix from ", node.XMLName.Space)
		if ns, found := nsByURLs[node.XMLName.Space]; found {
			prefix = ns + ":"
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
