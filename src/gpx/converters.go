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
		gpx11Doc.Waypoints = make([]*gpx11.GpxPoint, len(g.Waypoints))
		for waypointNo, waypoint := range g.Waypoints {
			gpx11Doc.Waypoints[waypointNo] = convertPointToGpx11(waypoint)
		}
	}

	if g.Routes != nil {
		gpx11Doc.Routes = make([]*gpx11.GpxRte, len(g.Routes))
		for routeNo, route := range g.Routes {
			r := new(gpx11.GpxRte)
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

			gpx11Doc.Routes[routeNo] = r

			if route.Points != nil {
				r.Points = make([]*gpx11.GpxPoint, len(route.Points))
				for pointNo, point := range route.Points {
					r.Points[pointNo] = convertPointToGpx11(point)
				}
			}
		}
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
		waypoints := make([]*GPXPoint, len(g.Waypoints))
		for waypointNo, waypoint := range g.Waypoints {
			waypoints[waypointNo] = convertPointFromGpx11(waypoint)
		}
		result.Waypoints = waypoints
	}

	if g.Routes != nil {
		result.Routes = make([]*GPXRoute, len(g.Routes))
		for routeNo, route := range g.Routes {
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
				r.Points = make([]*GPXPoint, len(route.Points))
				for pointNo, point := range route.Points {
					r.Points[pointNo] = convertPointFromGpx11(point)
				}
			}

			result.Routes[routeNo] = r
		}
	}

	return result
}

func convertPointToGpx11(original *GPXPoint) *gpx11.GpxPoint {
	result := new(gpx11.GpxPoint)
	result.Lat = original.Latitude
	result.Lon = original.Longitue
	result.Ele = original.Elevation
	result.Timestamp = formatGPXTime(original.Timestamp)
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
	result.Sat = original.Satellites
	result.Hdop = original.HorizontalDilution
	result.Vdop = original.VerticalDilution
	result.Pdop = original.PositionalDilution
	result.AgeOfDGpsData = original.AgeOfDGpsData
	result.DGpsId = original.DGpsId
	return result
}

func convertPointFromGpx11(original *gpx11.GpxPoint) *GPXPoint {
	result := new(GPXPoint)
	result.Latitude = original.Lat
	result.Longitue = original.Lon
	result.Elevation = original.Ele
	result.Timestamp, _ = parseGPXTime(original.Timestamp)
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
	result.Satellites = original.Sat
	result.HorizontalDilution = original.Hdop
	result.VerticalDilution = original.Vdop
	result.PositionalDilution = original.Pdop
	result.AgeOfDGpsData = original.AgeOfDGpsData
	result.DGpsId = original.DGpsId
	return result
}
