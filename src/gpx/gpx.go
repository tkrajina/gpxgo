package gpx

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"gpx/gpx10"
	"gpx/gpx11"

	//"fmt"
)

// An array cannot be constant :( The first one if the default layout:
var TIMELAYOUTS = []string{
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05Z",
	"2006-01-02 15:04:05",
}

func (g *GPX) ToXml(version string) ([]byte, error) {
	if version == "1.0" {
		gpx10Doc := convertToGpx10Models(g)
		return xml.Marshal(gpx10Doc)
	} else if version == "1.1" {
		gpx11Doc := convertToGpx11Models(g)
		return xml.Marshal(gpx11Doc)
	} else {
		return nil, errors.New("Invalid version " + version)
	}
}

func guessGPXVersion(bytes []byte) (string, error) {
	startOfDocument := string(bytes[:1000])

	parts := strings.Split(startOfDocument, "<gpx")
	if len(parts) <= 1 {
		return "", errors.New("Invalid GPX file, cannot find version")
	}
	parts = strings.Split(parts[1], "version=")

	if len(parts) <= 1 {
		return "", errors.New("Invalid GPX file, cannot find version")
	}

	if len(parts[1]) < 10 {
		return "", errors.New("Invalid GPX file, cannot find version")
	}

	result := parts[1][1:4]

	return result, nil
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
	version, _ := guessGPXVersion(bytes)
	if version == "1.0" {
		g := gpx10.NewGpx()
		err := xml.Unmarshal(bytes, &g)
		if err != nil {
			return nil, err
		}

		return convertFromGpx10Models(g), nil
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
