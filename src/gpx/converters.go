package gpx

import (
	"gpx11"
	"strings"
	//    "fmt"
)

func convertToGpx11Models(g *GPX) *gpx11.Gpx {
	gpx11Doc := gpx11.NewGpx()

	gpx11Doc.Creator = g.Creator
	gpx11Doc.Version = g.Version
	gpx11Doc.Name = g.Name
	gpx11Doc.Desc = g.Description
	gpx11Doc.AuthorName = g.AuthorName

	if len(g.AuthorEmail) > 0 {
		parts := strings.Split(g.AuthorEmail, "@")
		if len(parts) == 1 {
			gpx11Doc.AuthorEmail = new(gpx11.GpxEmail)
			gpx11Doc.AuthorEmail.Id = parts[0]
		} else if len(parts) > 1 {
			gpx11Doc.AuthorEmail = new(gpx11.GpxEmail)
			gpx11Doc.AuthorEmail.Id = parts[0]
			gpx11Doc.AuthorEmail.Domain = parts[1]
		}
	}

	if len(g.AuthorLink) > 0 || len(g.AuthorLinkText) > 0 || len(g.AuthorLinkType) > 0 {
		gpx11Doc.AuthorLink = new(gpx11.GpxLink)
		gpx11Doc.AuthorLink.Href = g.AuthorLink
		gpx11Doc.AuthorLink.Text = g.AuthorLinkText
		gpx11Doc.AuthorLink.Type = g.AuthorLinkType
	}

	if len(g.Copyright) > 0 || len(g.CopyrightYear) > 0 || len(g.CopyrightLicense) > 0 {
		gpx11Doc.Copyright = new(gpx11.GpxCopyright)
		gpx11Doc.Copyright.Author = g.Copyright
		gpx11Doc.Copyright.Year = g.CopyrightYear
		gpx11Doc.Copyright.License = g.CopyrightLicense
	}

	if len(g.Link) > 0 || len(g.LinkText) > 0 || len(g.LinkType) > 0 {
		gpx11Doc.Link = new(gpx11.GpxLink)
		gpx11Doc.Link.Href = g.Link
		gpx11Doc.Link.Text = g.LinkText
		gpx11Doc.Link.Type = g.LinkType
	}

	if g.Time != nil {
		gpx11Doc.Timestamp = formatGPXTime(g.Time)
	}

	gpx11Doc.Keywords = g.Keywords

	if g.Waypoints != nil {
		waypoints := make([]*gpx11.GpxWpt, len(g.Waypoints))
		for waypointNo, waypoint := range g.Waypoints {
			w := new(gpx11.GpxWpt)
			w.Lat = waypoint.Latitude
			w.Lon = waypoint.Longitue
			w.Ele = waypoint.Elevation
			w.Timestamp = formatGPXTime(waypoint.Timestamp)
			w.MagVar = waypoint.MagneticVariation
			w.GeoIdHeight = waypoint.GeoidHeight
			w.Name = waypoint.Name
			w.Cmt = waypoint.Comment
			w.Desc = waypoint.Description
			w.Src = waypoint.Source
			// TODO
			//w.Links = waypoint.Links
			w.Sym = waypoint.Symbol
			w.Type = waypoint.Type
			w.Fix = waypoint.TypeOfGpsFix
			w.Sat = waypoint.Satellites
			w.Hdop = waypoint.HorizontalDilution
			w.Vdop = waypoint.VerticalDiluation
			w.Pdop = waypoint.PositionalDilution
			w.AgeOfDGpsData = waypoint.AgeOfDGpsData
			w.DGpsId = waypoint.DGpsId

			waypoints[waypointNo] = w
		}
		gpx11Doc.Waypoints = waypoints
	}

	return gpx11Doc
}

func convertFromGpx11Models(g *gpx11.Gpx) *GPX {
	result := new(GPX)

	result.Creator = g.Creator
	result.Version = g.Version
	result.Name = g.Name
	result.Description = g.Desc
	result.AuthorName = g.AuthorName

	if g.AuthorEmail != nil {
		result.AuthorEmail = g.AuthorEmail.Id + "@" + g.AuthorEmail.Domain
	}
	if g.AuthorLink != nil {
		result.AuthorLink = g.AuthorLink.Href
		result.AuthorLinkText = g.AuthorLink.Text
		result.AuthorLinkType = g.AuthorLink.Type
	}

	if g.Extensions != nil {
		result.Extensions = &g.Extensions.Bytes
	}

	if len(g.Timestamp) > 0 {
		result.Time, _ = parseGPXTime(g.Timestamp)
	}

	if g.Copyright != nil {
		result.Copyright = g.Copyright.Author
		result.CopyrightYear = g.Copyright.Year
		result.CopyrightLicense = g.Copyright.License
	}

	if g.Link != nil {
		result.Link = g.Link.Href
		result.LinkText = g.Link.Text
		result.LinkType = g.Link.Type
	}

	result.Keywords = g.Keywords

	if g.Waypoints != nil {
		waypoints := make([]*GPXWaypoint, len(g.Waypoints))
		for waypointNo, waypoint := range g.Waypoints {
			w := new(GPXWaypoint)
			w.Latitude = waypoint.Lat
			w.Longitue = waypoint.Lon
			w.Elevation = waypoint.Ele
			w.Timestamp, _ = parseGPXTime(waypoint.Timestamp)
			w.MagneticVariation = waypoint.MagVar
			w.GeoidHeight = waypoint.GeoIdHeight
			w.Name = waypoint.Name
			w.Comment = waypoint.Cmt
			w.Description = waypoint.Desc
			w.Source = waypoint.Src
			// TODO
			//w.Links = waypoint.Links
			w.Symbol = waypoint.Sym
			w.Type = waypoint.Type
			w.TypeOfGpsFix = waypoint.Fix
			w.Satellites = waypoint.Sat
			w.HorizontalDilution = waypoint.Hdop
			w.VerticalDiluation = waypoint.Vdop
			w.PositionalDilution = waypoint.Pdop
			w.AgeOfDGpsData = waypoint.AgeOfDGpsData
			w.DGpsId = waypoint.DGpsId

			waypoints[waypointNo] = w
		}

		result.Waypoints = waypoints
	}

	return result
}
