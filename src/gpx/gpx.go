package gpx

import (
    "os"
    "io/ioutil"
    "errors"
    "gpx11"
    "encoding/xml"
    "strings"
    "fmt"
)

type GPX struct {
    Version string
    Creator string
    Name string
    Description string
    AuthorName string
    AuthorEmail string
    AuthorLink string
    AuthorLinkText string
    AuthorLinkType string
    Copyright string
    CopyrightYear string
    CopyrightLicense string

    // TODO
    Extensions *[]byte
}

func (g *GPX) ToXml(version string) ([]byte, error) {
    if version == "1.0" {
        return nil, errors.New("Invalid version:" + version)
    } else if version == "1.1" {
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

        return xml.Marshal(gpx11Doc)
    } else {
        return nil, errors.New("Invalid version " + version)
    }
}

func guessGPXVersion(bytes []byte) (string) {
    return "1.1"
}

func ParseFile(fileName string) (*GPX, error) {
    f, err := os.Open(fileName)
    if err != nil {
        return nil, err
    }

    bytes, err := ioutil.ReadAll(f)
    if err != nil {
        return nil, err
    }

    return ParseString(bytes)
}

func ParseString(bytes []byte) (*GPX, error) {
    version := guessGPXVersion(bytes)
    result := new(GPX)
    if version == "1.0" {
        return nil, errors.New("Invalid version:" + version)
    } else if version == "1.1" {
        g := gpx11.NewGpx()
        err := xml.Unmarshal(bytes, &g)
        if err != nil {
            return nil, err
        }

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

        fmt.Println("copyright", g.Copyright)
        if g.Copyright != nil {
            result.Copyright = g.Copyright.Author
            result.CopyrightYear = g.Copyright.Year
            result.CopyrightLicense = g.Copyright.License
        }

        /*
        if g.Metadata != nil {
            result.Name = g.Metadata.Name
            result.Name = g.Metadata.Name
            result.Description = g.Metadata.Desc
            if g.Metadata.Author != nil {
                result.AuthorName = g.Metadata.Author.Name
                result.AuthorEmail = g.Metadata.Author.Email
                if g.Metadata.Author.Link != nil {
                    result.AuthorLink = g.Metadata.Author.Link.Href
                    result.AuthorLinkText = g.Metadata.Author.Link.Text
                    result.AuthorLinkType = g.Metadata.Author.Link.Type
                }
            }
        }
        */

        return result, nil
    } else {
        fmt.Println("error")
        return nil, errors.New("Invalid version:" + version)
    }
}
