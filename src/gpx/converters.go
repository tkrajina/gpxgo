package gpx

import (
	"gpx11"
	"strings"
	//    "fmt"
)

func convertToGpx11Models(gpxDoc *GPX) *gpx11.Gpx {
	gpx11Doc := gpx11.NewGpx()

	gpx11Doc.Creator = gpxDoc.Creator
	gpx11Doc.Version = gpxDoc.Version
	gpx11Doc.Name = gpxDoc.Name
	gpx11Doc.Desc = gpxDoc.Description
	gpx11Doc.AuthorName = gpxDoc.AuthorName

	if len(gpxDoc.AuthorEmail) > 0 {
		parts := strings.Split(gpxDoc.AuthorEmail, "@")
		if len(parts) == 1 {
			gpx11Doc.AuthorEmail = new(gpx11.GpxEmail)
			gpx11Doc.AuthorEmail.Id = parts[0]
		} else if len(parts) > 1 {
			gpx11Doc.AuthorEmail = new(gpx11.GpxEmail)
			gpx11Doc.AuthorEmail.Id = parts[0]
			gpx11Doc.AuthorEmail.Domain = parts[1]
		}
	}

	if len(gpxDoc.AuthorLink) > 0 || len(gpxDoc.AuthorLinkText) > 0 || len(gpxDoc.AuthorLinkType) > 0 {
		gpx11Doc.AuthorLink = new(gpx11.GpxLink)
		gpx11Doc.AuthorLink.Href = gpxDoc.AuthorLink
		gpx11Doc.AuthorLink.Text = gpxDoc.AuthorLinkText
		gpx11Doc.AuthorLink.Type = gpxDoc.AuthorLinkType
	}

	if len(gpxDoc.Copyright) > 0 || len(gpxDoc.CopyrightYear) > 0 || len(gpxDoc.CopyrightLicense) > 0 {
		gpx11Doc.Copyright = new(gpx11.GpxCopyright)
		gpx11Doc.Copyright.Author = gpxDoc.Copyright
		gpx11Doc.Copyright.Year = gpxDoc.CopyrightYear
		gpx11Doc.Copyright.License = gpxDoc.CopyrightLicense
	}

	if len(gpxDoc.Link) > 0 || len(gpxDoc.LinkText) > 0 || len(gpxDoc.LinkType) > 0 {
		gpx11Doc.Link = new(gpx11.GpxLink)
		gpx11Doc.Link.Href = gpxDoc.Link
		gpx11Doc.Link.Text = gpxDoc.LinkText
		gpx11Doc.Link.Type = gpxDoc.LinkType
	}

	if gpxDoc.Time != nil {
		gpx11Doc.Timestamp = formatGPXTime(gpxDoc.Time)
	}

	gpx11Doc.Keywords = gpxDoc.Keywords

	if gpxDoc.Waypoints != nil {
		gpx11Doc.Waypoints = make([]*gpx11.GpxPoint, len(gpxDoc.Waypoints))
		for waypointNo, waypoint := range gpxDoc.Waypoints {
			gpx11Doc.Waypoints[waypointNo] = convertPointToGpx11(waypoint)
		}
	}

	if gpxDoc.Routes != nil {
		gpx11Doc.Routes = make([]*gpx11.GpxRte, len(gpxDoc.Routes))
		for routeNo, route := range gpxDoc.Routes {
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

func convertFromGpx11Models(gpx11Doc *gpx11.Gpx) *GPX {
	gpxDoc := new(GPX)

	gpxDoc.Creator = gpx11Doc.Creator
	gpxDoc.Version = gpx11Doc.Version
	gpxDoc.Name = gpx11Doc.Name
	gpxDoc.Description = gpx11Doc.Desc
	gpxDoc.AuthorName = gpx11Doc.AuthorName

	if gpx11Doc.AuthorEmail != nil {
		gpxDoc.AuthorEmail = gpx11Doc.AuthorEmail.Id + "@" + gpx11Doc.AuthorEmail.Domain
	}
	if gpx11Doc.AuthorLink != nil {
		gpxDoc.AuthorLink = gpx11Doc.AuthorLink.Href
		gpxDoc.AuthorLinkText = gpx11Doc.AuthorLink.Text
		gpxDoc.AuthorLinkType = gpx11Doc.AuthorLink.Type
	}

	if gpx11Doc.Extensions != nil {
		gpxDoc.Extensions = &gpx11Doc.Extensions.Bytes
	}

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
		waypoints := make([]*GPXPoint, len(gpx11Doc.Waypoints))
		for waypointNo, waypoint := range gpx11Doc.Waypoints {
			waypoints[waypointNo] = convertPointFromGpx11(waypoint)
		}
		gpxDoc.Waypoints = waypoints
	}

	if gpx11Doc.Routes != nil {
		gpxDoc.Routes = make([]*GPXRoute, len(gpx11Doc.Routes))
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

			if route.Points != nil {
				r.Points = make([]*GPXPoint, len(route.Points))
				for pointNo, point := range route.Points {
					r.Points[pointNo] = convertPointFromGpx11(point)
				}
			}

			gpxDoc.Routes[routeNo] = r
		}
	}

	return gpxDoc
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
