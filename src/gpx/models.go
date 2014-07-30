package gpx

import (
    "time"
)

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
	Waypoints  []*GPXPoint
	Routes     []*GPXRoute
	Tracks     []*GPXTrack
}

type GPXPoint struct {
	Latitude float64
	Longitue float64
	// Position info
	Elevation float64
	// TODO
	Timestamp *time.Time
	// TODO: Type
	MagneticVariation string
	// TODO: Type
	GeoidHeight string
	// Description info
	Name        string
	Comment     string
	Description string
	Source      string
	// TODO
	// Links       []GpxLink
	Symbol string
	Type   string
	// Accuracy info
	TypeOfGpsFix       string
	Satellites         int
	HorizontalDilution float64
	VerticalDilution   float64
	PositionalDilution float64
	AgeOfDGpsData      float64
	DGpsId             int
}

type GPXRoute struct {
	Name        string
	Comment     string
	Description string
	Source      string
	// TODO
	//Links       []Link
	Number int
	Type   string
	// TODO
	Points []*GPXPoint
}

type GPXTrackSegment struct {
	Points []*GPXPoint
	// TODO extensions
}

type GPXTrack struct {
	Name        string
	Comment     string
	Description string
	Source      string
	// TODO
	//Links    []Link
	Number   int
	Type     string
	Segments []*GPXTrackSegment
}

