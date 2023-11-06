// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"encoding/xml"
	"sort"
	"strings"
)

// defaultCreator contains the original repo path
const defaultCreator = "https://github.com/tkrajina/gpxgo"

// ----------------------------------------------------------------------------------------------------
// Gpx 1.0 Stuff
// ----------------------------------------------------------------------------------------------------

func convertToGpx10Models(gpxDoc *GPX) *gpx10Gpx {
	gpx10Doc := &gpx10Gpx{}
	//gpx10Doc.Attrs = namespacesMapToAttrs(gpxDoc.Namespaces)
	//gpx10Doc.XMLNs = gpxDoc.XMLNs
	attrs := gpxDoc.Attrs.ToXMLAttrs()
	gpx10Doc.Attrs = &attrs
	gpx10Doc.Attrs.Update("", "xmlns", "http://www.topografix.com/GPX/1/0")

	gpx10Doc.Version = "1.0"
	if len(gpxDoc.Creator) == 0 {
		gpx10Doc.Creator = defaultCreator
	} else {
		gpx10Doc.Creator = gpxDoc.Creator
	}
	gpx10Doc.Name = gpxDoc.Name
	gpx10Doc.Desc = gpxDoc.Description
	gpx10Doc.Author = gpxDoc.AuthorName
	gpx10Doc.Email = gpxDoc.AuthorEmail

	if len(gpxDoc.AuthorLink) > 0 || len(gpxDoc.AuthorLinkText) > 0 {
		// TODO
	}

	if len(gpxDoc.Link) > 0 || len(gpxDoc.LinkText) > 0 {
		gpx10Doc.Url = gpxDoc.Link
		gpx10Doc.UrlName = gpxDoc.LinkText
	}

	if gpxDoc.Time != nil {
		gpx10Doc.Time = formatGPXTime(gpxDoc.Time)
	}

	gpx10Doc.Keywords = gpxDoc.Keywords

	if gpxDoc.Waypoints != nil {
		gpx10Doc.Waypoints = make([]*gpx10GpxPoint, len(gpxDoc.Waypoints))
		for waypointNo, waypoint := range gpxDoc.Waypoints {
			gpx10Doc.Waypoints[waypointNo] = convertPointToGpx10(&waypoint)
		}
	}

	if gpxDoc.Routes != nil {
		gpx10Doc.Routes = make([]*gpx10GpxRte, len(gpxDoc.Routes))
		for routeNo, route := range gpxDoc.Routes {
			r := new(gpx10GpxRte)
			r.Name = route.Name
			r.Cmt = route.Comment
			r.Desc = route.Description
			r.Src = route.Source
			// TODO
			//r.Links = route.Links
			r.Number = route.Number
			r.Type = route.Type
			// TODO
			//r.RoutePoints = route.RoutePoints

			gpx10Doc.Routes[routeNo] = r

			if route.Points != nil {
				r.Points = make([]*gpx10GpxPoint, len(route.Points))
				for pointNo, point := range route.Points {
					r.Points[pointNo] = convertPointToGpx10(&point)
				}
			}
		}
	}

	if gpxDoc.Tracks != nil {
		gpx10Doc.Tracks = make([]*gpx10GpxTrk, len(gpxDoc.Tracks))
		for trackNo, track := range gpxDoc.Tracks {
			gpx10Track := new(gpx10GpxTrk)
			gpx10Track.Name = track.Name
			gpx10Track.Cmt = track.Comment
			gpx10Track.Desc = track.Description
			gpx10Track.Src = track.Source
			gpx10Track.Number = track.Number
			gpx10Track.Type = track.Type

			if track.Segments != nil {
				gpx10Track.Segments = make([]*gpx10GpxTrkSeg, len(track.Segments))
				for segmentNo, segment := range track.Segments {
					gpx10Segment := new(gpx10GpxTrkSeg)
					if segment.Points != nil {
						gpx10Segment.Points = make([]*gpx10GpxPoint, len(segment.Points))
						for pointNo, point := range segment.Points {
							gpx10Point := convertPointToGpx10(&point)
							// TODO
							//gpx10Point.Speed = point.Speed
							//gpx10Point.Speed = point.Speed
							gpx10Segment.Points[pointNo] = gpx10Point
						}
					}
					gpx10Track.Segments[segmentNo] = gpx10Segment
				}
			}
			gpx10Doc.Tracks[trackNo] = gpx10Track
		}
	}

	return gpx10Doc
}

func convertFromGpx10Models(gpx10Doc *gpx10Gpx) *GPX {
	gpxDoc := new(GPX)
	if gpx10Doc.Attrs != nil {
		gpxDoc.Attrs = NewGPXAttributes(*gpx10Doc.Attrs)
	}

	gpxDoc.Creator = gpx10Doc.Creator
	gpxDoc.Version = gpx10Doc.Version
	gpxDoc.Name = gpx10Doc.Name
	gpxDoc.Description = gpx10Doc.Desc
	gpxDoc.AuthorName = gpx10Doc.Author
	gpxDoc.AuthorEmail = gpx10Doc.Email

	if len(gpx10Doc.Url) > 0 || len(gpx10Doc.UrlName) > 0 {
		gpxDoc.Link = gpx10Doc.Url
		gpxDoc.LinkText = gpx10Doc.UrlName
	}

	if len(gpx10Doc.Time) > 0 {
		gpxDoc.Time, _ = parseGPXTime(gpx10Doc.Time)
	}

	gpxDoc.Keywords = gpx10Doc.Keywords

	if gpx10Doc.Waypoints != nil {
		waypoints := make([]GPXPoint, len(gpx10Doc.Waypoints))
		for waypointNo, waypoint := range gpx10Doc.Waypoints {
			waypoints[waypointNo] = *convertPointFromGpx10(waypoint)
		}
		gpxDoc.Waypoints = waypoints
	}

	if gpx10Doc.Routes != nil {
		gpxDoc.Routes = make([]GPXRoute, len(gpx10Doc.Routes))
		for routeNo, route := range gpx10Doc.Routes {
			r := new(GPXRoute)

			r.Name = route.Name
			r.Comment = route.Cmt
			r.Description = route.Desc
			r.Source = route.Src
			// TODO
			//r.Links = route.Links
			r.Number = route.Number
			r.Type = route.Type
			// TODO
			//r.RoutePoints = route.RoutePoints

			if route.Points != nil {
				r.Points = make([]GPXPoint, len(route.Points))
				for pointNo, point := range route.Points {
					r.Points[pointNo] = *convertPointFromGpx10(point)
				}
			}

			gpxDoc.Routes[routeNo] = *r
		}
	}

	if gpx10Doc.Tracks != nil {
		gpxDoc.Tracks = make([]GPXTrack, len(gpx10Doc.Tracks))
		for trackNo, track := range gpx10Doc.Tracks {
			gpxTrack := new(GPXTrack)
			gpxTrack.Name = track.Name
			gpxTrack.Comment = track.Cmt
			gpxTrack.Description = track.Desc
			gpxTrack.Source = track.Src
			gpxTrack.Number = track.Number
			gpxTrack.Type = track.Type

			if track.Segments != nil {
				gpxTrack.Segments = make([]GPXTrackSegment, len(track.Segments))
				for segmentNo, segment := range track.Segments {
					gpxSegment := GPXTrackSegment{}
					if segment.Points != nil {
						gpxSegment.Points = make([]GPXPoint, len(segment.Points))
						for pointNo, point := range segment.Points {
							gpxSegment.Points[pointNo] = *convertPointFromGpx10(point)
						}
					}
					gpxTrack.Segments[segmentNo] = gpxSegment
				}
			}
			gpxDoc.Tracks[trackNo] = *gpxTrack
		}
	}

	return gpxDoc
}

func convertPointToGpx10(original *GPXPoint) *gpx10GpxPoint {
	result := new(gpx10GpxPoint)
	result.Lat = formattedFloat(original.Latitude)
	result.Lon = formattedFloat(original.Longitude)
	result.Ele = original.Elevation
	result.Timestamp = formatGPXTime(&original.Timestamp)
	result.MagVar = original.MagneticVariation
	result.GeoIdHeight = original.GeoidHeight
	result.Name = original.Name
	result.Cmt = original.Comment
	result.Desc = original.Description
	result.Src = original.Source
	// TODO
	//w.Links = original.Links
	result.Sym = original.Symbol
	result.Type = original.Type
	result.Fix = original.TypeOfGpsFix
	if original.Satellites.NotNil() {
		value := original.Satellites.Value()
		result.Sat = &value
	}
	result.Hdop = original.HorizontalDilution
	result.Vdop = original.VerticalDilution
	result.Pdop = original.PositionalDilution
	result.AgeOfDGpsData = original.AgeOfDGpsData
	result.DGpsId = original.DGpsId
	return result
}

func convertPointFromGpx10(original *gpx10GpxPoint) *GPXPoint {
	result := new(GPXPoint)
	result.Latitude = float64(original.Lat)
	result.Longitude = float64(original.Lon)
	result.Elevation = original.Ele
	time, _ := parseGPXTime(original.Timestamp)
	if time != nil {
		result.Timestamp = *time
	}
	result.MagneticVariation = original.MagVar
	result.GeoidHeight = original.GeoIdHeight
	result.Name = original.Name
	result.Comment = original.Cmt
	result.Description = original.Desc
	result.Source = original.Src
	// TODO
	//w.Links = original.Links
	result.Symbol = original.Sym
	result.Type = original.Type
	result.TypeOfGpsFix = original.Fix
	if original.Sat != nil {
		result.Satellites = NewNilableint(*original.Sat)
	}
	result.HorizontalDilution = original.Hdop
	result.VerticalDilution = original.Vdop
	result.PositionalDilution = original.Pdop
	result.AgeOfDGpsData = original.AgeOfDGpsData
	result.DGpsId = original.DGpsId
	return result
}

// ----------------------------------------------------------------------------------------------------
// Gpx 1.1 Stuff
// ----------------------------------------------------------------------------------------------------

type NamespaceAttribute struct {
	Space       string `json:"space,omitempty"`
	Local       string `json:"local,omitempty"`
	Value       string `json:"value,omitempty"`
	replacement string `json:"-"`
}

type GPXAttributes map[string]map[string]string

func NewGPXAttributes(attrs XMLAttrs) GPXAttributes {
	res := GPXAttributes(make(map[string]map[string]string))
	for _, attr := range attrs {
		res.Set(attr.Name.Space, attr.Name.Local, attr.Value)
	}
	return GPXAttributes(res)
}

func (ga GPXAttributes) Set(space, local, value string) {
	if ga == nil {
		res := GPXAttributes(make(map[string]map[string]string))
		ga = res
	}
	if _, found := (ga)[space]; !found {
		(ga)[space] = make(map[string]string)
	}
	(ga)[space][local] = value
}

func (ga GPXAttributes) RegisterNamespace(ns, url string) {
	ga.Set("xmlns", ns, url)
}

func (ga GPXAttributes) GetNamespacesByURLs() map[string]string {
	xmlns, found := ga["xmlns"]
	if !found {
		return nil
	}
	return xmlns
}

func (ga GPXAttributes) ToXMLAttrs() XMLAttrs {
	var attrs []xml.Attr
	for space, entry := range ga {
		for local, value := range entry {
			attrs = append(attrs, xml.Attr{Name: xml.Name{Space: space, Local: local}, Value: value})
		}
	}
	res := XMLAttrs(attrs)
	sort.Sort(res)
	return res
}

func convertToGpx11Extension(ext Extension, gpxDoc GPX) *gpx11Extension {
	if ext == nil {
		return nil
	}
	res := gpx11Extension{gpx: &gpxDoc}
	for _, n := range ext {
		res.Nodes = append(res.Nodes, convertToGpx11ExtensionNode(n))
	}
	return &res
}
func convertToGpx11ExtensionNode(n ExtensionNode) gpx11ExtensionNode {
	node := gpx11ExtensionNode{
		XMLName: xml.Name{
			Space: n.NameSpace,
			Local: n.NameLocal,
		},
		Data: strings.TrimSpace(n.Data),
	}
	for _, attr := range n.Attrs {
		node.Attrs = append(node.Attrs, xml.Attr{
			Name: xml.Name{
				Space: attr.NameSpace,
				Local: attr.NameLocal,
			},
			Value: attr.Value,
		})
	}
	for subn := range n.Nodes {
		node.Nodes = append(node.Nodes, convertToGpx11ExtensionNode(n.Nodes[subn]))
	}
	return node
}

func convertFromGpx11Extension(ext *gpx11Extension) Extension {
	if ext == nil {
		return nil
	}
	nodes := []ExtensionNode{}
	for n := range ext.Nodes {
		nodes = append(nodes, convertFromGpx11ExtensionNode(ext.Nodes[n]))
	}
	return Extension(nodes)
}

func convertFromGpx11ExtensionNode(n gpx11ExtensionNode) ExtensionNode {
	res := ExtensionNode{
		NameSpace: n.XMLName.Space,
		NameLocal: n.XMLName.Local,
		Data:      strings.TrimSpace(n.Data),
	}
	for _, attr := range n.Attrs {
		res.Attrs = append(res.Attrs, ExtensionNodeAttr{
			NameSpace: attr.Name.Space,
			NameLocal: attr.Name.Local,
			Value:     attr.Value,
		})
	}
	for subn := range n.Nodes {
		res.Nodes = append(res.Nodes, convertFromGpx11ExtensionNode(n.Nodes[subn]))
	}
	return res
}

func convertToGpx11Models(gpxDoc *GPX) *gpx11Gpx {
	gpx11Doc := &gpx11Gpx{}

	gpx11Doc.Version = "1.1"

	attrs := gpxDoc.Attrs.ToXMLAttrs()
	gpx11Doc.Attrs = &attrs
	gpx11Doc.Attrs.Update("", "xmlns", "http://www.topografix.com/GPX/1/1")

	gpx11Doc.Extensions = convertToGpx11Extension(gpxDoc.Extensions, *gpxDoc)
	gpx11Doc.MetadataExtensions = convertToGpx11Extension(gpxDoc.MetadataExtensions, *gpxDoc)

	if len(gpxDoc.Creator) == 0 {
		gpx11Doc.Creator = defaultCreator
	} else {
		gpx11Doc.Creator = gpxDoc.Creator
	}
	gpx11Doc.Name = gpxDoc.Name
	gpx11Doc.Desc = gpxDoc.Description
	gpx11Doc.AuthorName = gpxDoc.AuthorName

	if len(gpxDoc.AuthorEmail) > 0 {
		parts := strings.Split(gpxDoc.AuthorEmail, "@")
		if len(parts) == 1 {
			gpx11Doc.AuthorEmail = new(gpx11GpxEmail)
			gpx11Doc.AuthorEmail.Id = parts[0]
		} else if len(parts) > 1 {
			gpx11Doc.AuthorEmail = new(gpx11GpxEmail)
			gpx11Doc.AuthorEmail.Id = parts[0]
			gpx11Doc.AuthorEmail.Domain = parts[1]
		}
	}

	if len(gpxDoc.AuthorLink) > 0 || len(gpxDoc.AuthorLinkText) > 0 || len(gpxDoc.AuthorLinkType) > 0 {
		gpx11Doc.AuthorLink = new(gpx11GpxLink)
		gpx11Doc.AuthorLink.Href = gpxDoc.AuthorLink
		gpx11Doc.AuthorLink.Text = gpxDoc.AuthorLinkText
		gpx11Doc.AuthorLink.Type = gpxDoc.AuthorLinkType
	}

	if len(gpxDoc.Copyright) > 0 || len(gpxDoc.CopyrightYear) > 0 || len(gpxDoc.CopyrightLicense) > 0 {
		gpx11Doc.Copyright = new(gpx11GpxCopyright)
		gpx11Doc.Copyright.Author = gpxDoc.Copyright
		gpx11Doc.Copyright.Year = gpxDoc.CopyrightYear
		gpx11Doc.Copyright.License = gpxDoc.CopyrightLicense
	}

	if len(gpxDoc.Link) > 0 || len(gpxDoc.LinkText) > 0 || len(gpxDoc.LinkType) > 0 {
		gpx11Doc.Link = new(gpx11GpxLink)
		gpx11Doc.Link.Href = gpxDoc.Link
		gpx11Doc.Link.Text = gpxDoc.LinkText
		gpx11Doc.Link.Type = gpxDoc.LinkType
	}

	if gpxDoc.Time != nil {
		gpx11Doc.Timestamp = formatGPXTime(gpxDoc.Time)
	}

	gpx11Doc.Keywords = gpxDoc.Keywords

	if gpxDoc.Waypoints != nil {
		gpx11Doc.Waypoints = make([]*gpx11GpxPoint, len(gpxDoc.Waypoints))
		for waypointNo, waypoint := range gpxDoc.Waypoints {
			gpx11Doc.Waypoints[waypointNo] = convertPointToGpx11(&waypoint, *gpxDoc)
			gpx11Doc.Waypoints[waypointNo].Extensions = convertToGpx11Extension(waypoint.Extensions, *gpxDoc)
		}
	}

	if gpxDoc.Routes != nil {
		gpx11Doc.Routes = make([]*gpx11GpxRte, len(gpxDoc.Routes))
		for routeNo, route := range gpxDoc.Routes {
			r := new(gpx11GpxRte)
			r.Name = route.Name
			r.Cmt = route.Comment
			r.Desc = route.Description
			r.Src = route.Source
			// TODO
			//r.Links = route.Links
			r.Number = route.Number
			r.Type = route.Type
			r.Extensions = convertToGpx11Extension(route.Extensions, *gpxDoc)
			gpx11Doc.Routes[routeNo] = r

			if route.Points != nil {
				r.Points = make([]*gpx11GpxPoint, len(route.Points))
				for pointNo, point := range route.Points {
					r.Points[pointNo] = convertPointToGpx11(&point, *gpxDoc)
					r.Points[pointNo].Extensions = convertToGpx11Extension(point.Extensions, *gpxDoc)
				}
			}
		}
	}

	if gpxDoc.Tracks != nil {
		gpx11Doc.Tracks = make([]*gpx11GpxTrk, len(gpxDoc.Tracks))
		for trackNo, track := range gpxDoc.Tracks {
			gpx11Track := new(gpx11GpxTrk)
			gpx11Track.Name = track.Name
			gpx11Track.Cmt = track.Comment
			gpx11Track.Desc = track.Description
			gpx11Track.Src = track.Source
			gpx11Track.Number = track.Number
			gpx11Track.Type = track.Type
			gpx11Track.Extensions = convertToGpx11Extension(track.Extensions, *gpxDoc)

			if track.Segments != nil {
				gpx11Track.Segments = make([]*gpx11GpxTrkSeg, len(track.Segments))
				for segmentNo, segment := range track.Segments {
					gpx11Segment := new(gpx11GpxTrkSeg)
					gpx11Segment.Extensions = convertToGpx11Extension(segment.Extensions, *gpxDoc)
					if segment.Points != nil {
						gpx11Segment.Points = make([]*gpx11GpxPoint, len(segment.Points))
						for pointNo, point := range segment.Points {
							gpx11Segment.Points[pointNo] = convertPointToGpx11(&point, *gpxDoc)
							gpx11Segment.Points[pointNo].Extensions = convertToGpx11Extension(point.Extensions, *gpxDoc)
						}
					}
					gpx11Track.Segments[segmentNo] = gpx11Segment
				}
			}
			gpx11Doc.Tracks[trackNo] = gpx11Track
		}
	}

	return gpx11Doc
}

func convertFromGpx11Models(gpx11Doc *gpx11Gpx) *GPX {
	gpxDoc := new(GPX)

	if gpx11Doc.Attrs != nil {
		gpxDoc.Attrs = NewGPXAttributes(*gpx11Doc.Attrs)
	}

	gpxDoc.Creator = gpx11Doc.Creator
	gpxDoc.Version = gpx11Doc.Version
	gpxDoc.Name = gpx11Doc.Name
	gpxDoc.Description = gpx11Doc.Desc
	gpxDoc.AuthorName = gpx11Doc.AuthorName
	gpxDoc.Extensions = convertFromGpx11Extension(gpx11Doc.Extensions)
	gpxDoc.MetadataExtensions = convertFromGpx11Extension(gpx11Doc.MetadataExtensions)

	if gpx11Doc.AuthorEmail != nil {
		gpxDoc.AuthorEmail = gpx11Doc.AuthorEmail.Id + "@" + gpx11Doc.AuthorEmail.Domain
	}
	if gpx11Doc.AuthorLink != nil {
		gpxDoc.AuthorLink = gpx11Doc.AuthorLink.Href
		gpxDoc.AuthorLinkText = gpx11Doc.AuthorLink.Text
		gpxDoc.AuthorLinkType = gpx11Doc.AuthorLink.Type
	}

	/* TODO
	if gpx11Doc.Extensions != nil {
		gpxDoc.Extensions = &gpx11Doc.Extensions.Bytes
	}
	*/

	if len(gpx11Doc.Timestamp) > 0 {
		gpxDoc.Time, _ = parseGPXTime(gpx11Doc.Timestamp)
	}

	if gpx11Doc.Copyright != nil {
		gpxDoc.Copyright = gpx11Doc.Copyright.Author
		gpxDoc.CopyrightYear = gpx11Doc.Copyright.Year
		gpxDoc.CopyrightLicense = gpx11Doc.Copyright.License
	}

	if gpx11Doc.Link != nil {
		gpxDoc.Link = gpx11Doc.Link.Href
		gpxDoc.LinkText = gpx11Doc.Link.Text
		gpxDoc.LinkType = gpx11Doc.Link.Type
	}

	gpxDoc.Keywords = gpx11Doc.Keywords

	if gpx11Doc.Waypoints != nil {
		waypoints := make([]GPXPoint, len(gpx11Doc.Waypoints))
		for waypointNo, waypoint := range gpx11Doc.Waypoints {
			waypoints[waypointNo] = *convertPointFromGpx11(waypoint)
		}
		gpxDoc.Waypoints = waypoints
	}

	if gpx11Doc.Routes != nil {
		gpxDoc.Routes = make([]GPXRoute, len(gpx11Doc.Routes))
		for routeNo, route := range gpx11Doc.Routes {
			r := new(GPXRoute)

			r.Name = route.Name
			r.Comment = route.Cmt
			r.Description = route.Desc
			r.Source = route.Src
			// TODO
			//r.Links = route.Links
			r.Number = route.Number
			r.Type = route.Type
			// TODO
			//r.RoutePoints = route.RoutePoints
			r.Extensions = convertFromGpx11Extension(route.Extensions)

			if route.Points != nil {
				r.Points = make([]GPXPoint, len(route.Points))
				for pointNo, point := range route.Points {
					r.Points[pointNo] = *convertPointFromGpx11(point)
				}
			}

			gpxDoc.Routes[routeNo] = *r
		}
	}

	if gpx11Doc.Tracks != nil {
		gpxDoc.Tracks = make([]GPXTrack, len(gpx11Doc.Tracks))
		for trackNo, track := range gpx11Doc.Tracks {
			gpxTrack := new(GPXTrack)
			gpxTrack.Name = track.Name
			gpxTrack.Comment = track.Cmt
			gpxTrack.Description = track.Desc
			gpxTrack.Source = track.Src
			gpxTrack.Number = track.Number
			gpxTrack.Type = track.Type
			if track.Extensions != nil {
				gpxTrack.Extensions = convertFromGpx11Extension(track.Extensions)
			}

			if track.Segments != nil {
				gpxTrack.Segments = make([]GPXTrackSegment, len(track.Segments))
				if track.Extensions != nil {
					gpxTrack.Extensions = convertFromGpx11Extension(track.Extensions)
				}
				for segmentNo, segment := range track.Segments {
					gpxSegment := GPXTrackSegment{}
					gpxSegment.Extensions = convertFromGpx11Extension(segment.Extensions)
					if segment.Points != nil {
						gpxSegment.Points = make([]GPXPoint, len(segment.Points))
						for pointNo, point := range segment.Points {
							gpxSegment.Points[pointNo] = *convertPointFromGpx11(point)
						}
					}
					gpxTrack.Segments[segmentNo] = gpxSegment
				}
			}
			gpxDoc.Tracks[trackNo] = *gpxTrack
		}
	}

	return gpxDoc
}

func convertPointToGpx11(original *GPXPoint, gpxdoc GPX) *gpx11GpxPoint {
	result := new(gpx11GpxPoint)
	result.Lat = formattedFloat(original.Latitude)
	result.Lon = formattedFloat(original.Longitude)
	result.Ele = original.Elevation
	result.Timestamp = formatGPXTime(&original.Timestamp)
	result.MagVar = original.MagneticVariation
	result.GeoIdHeight = original.GeoidHeight
	result.Name = original.Name
	result.Cmt = original.Comment
	result.Desc = original.Description
	result.Src = original.Source
	// TODO
	//w.Links = original.Links
	result.Sym = original.Symbol
	result.Type = original.Type
	result.Fix = original.TypeOfGpsFix
	result.Extensions = convertToGpx11Extension(original.Extensions, gpxdoc)
	if original.Satellites.NotNil() {
		value := original.Satellites.Value()
		result.Sat = &value
	}
	result.Hdop = original.HorizontalDilution
	result.Vdop = original.VerticalDilution
	result.Pdop = original.PositionalDilution
	result.AgeOfDGpsData = original.AgeOfDGpsData
	result.DGpsId = original.DGpsId
	return result
}

func convertPointFromGpx11(original *gpx11GpxPoint) *GPXPoint {
	result := new(GPXPoint)
	result.Latitude = float64(original.Lat)
	result.Longitude = float64(original.Lon)
	result.Elevation = original.Ele
	time, _ := parseGPXTime(original.Timestamp)
	if time != nil {
		result.Timestamp = *time
	}
	result.MagneticVariation = original.MagVar
	result.GeoidHeight = original.GeoIdHeight
	result.Name = original.Name
	result.Comment = original.Cmt
	result.Description = original.Desc
	result.Source = original.Src
	// TODO
	//w.Links = original.Links
	result.Symbol = original.Sym
	result.Type = original.Type
	result.TypeOfGpsFix = original.Fix
	result.Extensions = convertFromGpx11Extension(original.Extensions)
	if original.Sat != nil {
		result.Satellites = NewNilableint(*original.Sat)
	}
	result.HorizontalDilution = original.Hdop
	result.VerticalDilution = original.Vdop
	result.PositionalDilution = original.Pdop
	result.AgeOfDGpsData = original.AgeOfDGpsData
	result.DGpsId = original.DGpsId
	return result
}
