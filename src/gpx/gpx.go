package gpx

import (
	"encoding/xml"
	"errors"
	//	"fmt"
	"gpx11"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// An array cannot be constant :( The first one if the default layout:
var TIMELAYOUTS = []string{
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05Z",
	"2006-01-02 15:04:05",
}

type GPX struct {
	Version          string
	Creator          string
	Name             string
	Description      string
	AuthorName       string
	AuthorEmail      string
	AuthorLink       string
	AuthorLinkText   string
	AuthorLinkType   string
	Copyright        string
	CopyrightYear    string
	CopyrightLicense string
	Link             string
	LinkText         string
	LinkType         string
	Time             *time.Time
	Keywords         string

	// TODO
	Extensions *[]byte
}

func (g *GPX) ToXml(version string) ([]byte, error) {
	if version == "1.0" {
		return nil, errors.New("Invalid version:" + version)
	} else if version == "1.1" {
		gpx11Doc := convertToGpx11Models(g)

		return xml.Marshal(gpx11Doc)
	} else {
		return nil, errors.New("Invalid version " + version)
	}
}

func guessGPXVersion(bytes []byte) string {
	return "1.1"
}

func parseGPXTime(timestr string) (*time.Time, error) {
	if strings.Contains(timestr, ".") {
		// Probably seconds with milliseconds
		timestr = strings.Split(timestr, ".")[0]
	}
	timestr = strings.Trim(timestr, " \t\n\r")
	for _, timeLayout := range TIMELAYOUTS {
		t, err := time.Parse(timeLayout, timestr)

		if err == nil {
			return &t, nil
		}
	}

	result := time.Now()

	return &result, errors.New("Cannot parse " + timestr)
}

func formatGPXTime(time *time.Time) string {
	if time == nil {
		return ""
	}
	return time.Format(TIMELAYOUTS[0])
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
	if version == "1.0" {
		return nil, errors.New("Invalid version:" + version)
	} else if version == "1.1" {
		g := gpx11.NewGpx()
		err := xml.Unmarshal(bytes, &g)
		if err != nil {
			return nil, err
		}

		return convertFromGpx11Models(g), nil
	} else {
		return nil, errors.New("Invalid version:" + version)
	}
}
