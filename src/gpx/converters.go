package gpx

import (
	"gpx11"
	"strings"
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

	return result
}
