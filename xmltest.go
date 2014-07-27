package main

import (
    "encoding/xml"
    "fmt"
    "os"
    "io/ioutil"
)

type Gpx struct {
	XMLName      xml.Name     `xml:"http://www.topografix.com/GPX/1/1 gpx"`
	XmlNsXsi     string       `xml:"xmlns:xsi,attr,omitempty"`
	XmlSchemaLoc string       `xml:"xsi:schemaLocation,attr,omitempty"`
	Version      string       `xml:"version,attr"`
	Creator      string       `xml:"creator,attr"`
	Metadata     *GpxMetadata `xml:"metadata,omitempty"`
}

type GpxExtensions struct {
    Extensions []byte        `xml:",innerxml"`
}

type GpxMetadata struct {
	Name       string         `xml:"name,omitempty"`
	Desc       string         `xml:"desc,omitempty"`
	Timestamp  string         `xml:"time,omitempty"`
	Keywords   string         `xml:"keywords,omitempty"`
    Extensions *GpxExtensions `xml:"extensions"`
}

func main() {
    f, _ := os.Open("test_files/gpx1.1_with_all_fields.gpx")
    contents, _ := ioutil.ReadAll(f)
    gpx := new(Gpx)
    xml.Unmarshal(contents, &gpx)
    fmt.Println(string(gpx.Metadata.Extensions.Extensions))
}
