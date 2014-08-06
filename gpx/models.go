package gpx

import (
	"fmt"
	"math"
	"time"
)

const DEFAULT_STOPPED_SPEED_THRESHOLD = 1.0

// ----------------------------------------------------------------------------------------------------

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
	//Extensions []byte
	Waypoints []*GPXPoint
	Routes    []*GPXRoute
	Tracks    []*GPXTrack
}

/*
 * Params are optional, you can set null to use GPXs Version and no indentation.
 */
func (g *GPX) ToXml(params ToXmlParams) ([]byte, error) {
	return ToXml(g, params)
}

// Length2D returns the 2D length of all tracks in a Gpx.
func (g *GPX) Length2D() float64 {
	var length2d float64
	for _, trk := range g.Tracks {
		length2d += trk.Length2D()
	}
	return length2d
}

// Length3D returns the 3D length of all tracks,
func (g *GPX) Length3D() float64 {
	var length3d float64
	for _, trk := range g.Tracks {
		length3d += trk.Length3D()
	}
	return length3d
}

// TimeBounds returns the time bounds of all tacks in a Gpx.
func (g *GPX) TimeBounds() *TimeBounds {
	var tbGpx *TimeBounds
	for i, trk := range g.Tracks {
		tbTrk := trk.TimeBounds()
		if i == 0 {
			tbGpx = trk.TimeBounds()
		} else {
			tbGpx.EndTime = tbTrk.EndTime
		}
	}
	return tbGpx
}

// Bounds returns the bounds of all tracks in a Gpx.
func (g *GPX) Bounds() *GpxBounds {
	minmax := getMinimaMaximaStart()
	for _, trk := range g.Tracks {
		bnds := trk.Bounds()
		minmax.MaxLat = math.Max(bnds.MaxLat, minmax.MaxLat)
		minmax.MinLat = math.Min(bnds.MinLat, minmax.MinLat)
		minmax.MaxLon = math.Max(bnds.MaxLon, minmax.MaxLon)
		minmax.MinLon = math.Min(bnds.MinLon, minmax.MinLon)
	}
	return minmax
}

// MovingData returns the moving data for all tracks in a Gpx.
func (g *GPX) MovingData() *MovingData {
	var (
		movingTime      float64
		stoppedTime     float64
		movingDistance  float64
		stoppedDistance float64
		maxSpeed        float64
	)

	for _, trk := range g.Tracks {
		md := trk.MovingData()
		movingTime += md.MovingTime
		stoppedTime += md.StoppedTime
		movingDistance += md.MovingDistance
		stoppedDistance += md.StoppedDistance

		if md.MaxSpeed > maxSpeed {
			maxSpeed = md.MaxSpeed
		}
	}
	return &MovingData{
		MovingTime:      movingTime,
		MovingDistance:  movingDistance,
		StoppedTime:     stoppedTime,
		StoppedDistance: stoppedDistance,
		MaxSpeed:        maxSpeed,
	}
}

// Split splits the Gpx segment segNo in a given track trackNo at
// pointNo.
func (g *GPX) Split(trackNo, segNo, pointNo int) {
	if trackNo >= len(g.Tracks) {
		return
	}

	track := &g.Tracks[trackNo]

	track.Split(segNo, pointNo)
}

// Duration returns the duration of all tracks in a Gpx in seconds.
func (g *GPX) Duration() float64 {
	if len(g.Tracks) == 0 {
		return 0.0
	}
	var result float64
	for _, trk := range g.Tracks {
		result += trk.Duration()
	}

	return result
}

// UphillDownhill returns uphill and downhill values for all tracks in a
// Gpx.
func (g *GPX) UphillDownhill() *UphillDownhill {
	if len(g.Tracks) == 0 {
		return nil
	}

	var (
		uphill   float64
		downhill float64
	)

	for _, trk := range g.Tracks {
		updo := trk.UphillDownhill()

		uphill += updo.Uphill
		downhill += updo.Downhill
	}

	return &UphillDownhill{
		Uphill:   uphill,
		Downhill: downhill,
	}
}

// Checks if *tracks* and segments have time information. Routes and Waypoints are ignored.
func (g *GPX) HasTimes() bool {
    result := true
    for _, track := range g.Tracks {
        result = result && track.HasTimes()
    }
    return result
}

// LocationAt returns a LocationResultsPair consisting the segment index
// and the GpxWpt at a certain time.
func (g *GPX) LocationAt(t time.Time) []LocationsResultPair {
	results := make([]LocationsResultPair, 0)

	for _, trk := range g.Tracks {
		locs := trk.LocationAt(t)
		if len(locs) > 0 {
			results = append(results, locs...)
		}
	}
	return results
}

func (g *GPX) ExecuteOnAllPoints(executor func(*GPXPoint)) {
	g.ExecuteOnWaypoints(executor)
	g.ExecuteOnRoutePoints(executor)
	g.ExecuteOnTrackPoints(executor)
}

func (g *GPX) ExecuteOnWaypoints(executor func(*GPXPoint)) {
	for _, waypoint := range g.Waypoints {
		executor(waypoint)
	}
}

func (g *GPX) ExecuteOnRoutePoints(executor func(*GPXPoint)) {
	for _, route := range g.Routes {
		route.ExecuteOnPoints(executor)
	}
}

func (g *GPX) ExecuteOnTrackPoints(executor func(*GPXPoint)) {
	for _, track := range g.Tracks {
		track.ExecuteOnPoints(executor)
	}
}

func (g *GPX) AddElevation(elevation float64) {
	g.ExecuteOnAllPoints(func(point *GPXPoint) {
		point.Elevation += elevation
	})
}

func (g *GPX) RemoveElevation() {
	for _, waypoint := range g.Waypoints {
		waypoint.RemoveElevation()
	}
	for _, route := range g.Routes {
		route.RemoveElevation()
	}
	for _, track := range g.Tracks {
		track.RemoveElevation()
	}
}

func (g *GPX) AppendTrack(t *GPXTrack) {
	g.Tracks = append(g.Tracks, t)
}

func (g *GPX) AppendRoute(r *GPXRoute) {
	g.Routes = append(g.Routes, r)
}

func (g *GPX) AppendWaypoint(w *GPXPoint) {
	g.Waypoints = append(g.Waypoints, w)
}

// ----------------------------------------------------------------------------------------------------

type GpxBounds struct {
	MinLat float64
	MaxLat float64
	MinLon float64
	MaxLon float64
}

// Equals returns true if two Bounds objects are equal
func (b *GpxBounds) Equals(b2 *GpxBounds) bool {
	return b.MinLon == b2.MinLon && b.MaxLat == b2.MaxLat && b.MinLon == b2.MinLon && b.MaxLon == b.MaxLon
}

func (b *GpxBounds) String() string {
	return fmt.Sprintf("Max: %+v, %+v Min: %+v, %+v", b.MinLat, b.MinLon, b.MaxLat, b.MaxLon)
}

// ----------------------------------------------------------------------------------------------------

// Generic point data
type Point struct {
	Latitude  float64
	Longitude float64
	Elevation float64
}

// Distance2D returns the 2D distance of two GpxWpts.
func (pt *Point) Distance2D(pt2 *Point) float64 {
	return Distance2D(pt.Latitude, pt.Longitude, pt2.Latitude, pt2.Longitude, false)
}

// Distance3D returns the 3D distance of two GpxWpts.
func (pt *Point) Distance3D(pt2 *Point) float64 {
	return Distance3D(pt.Latitude, pt.Longitude, pt.Elevation, pt2.Latitude, pt2.Longitude, pt2.Elevation, false)
}

func (pt *Point) RemoveElevation() {
	// TODO: This should be nil!
	pt.Elevation = 0
}

// ----------------------------------------------------------------------------------------------------

type TimeBounds struct {
	StartTime time.Time
	EndTime   time.Time
}

func (tb *TimeBounds) Equals(tb2 *TimeBounds) bool {
	fmt.Println(tb.StartTime)
	fmt.Println(tb2.StartTime)
	fmt.Println(tb.EndTime.Equal(tb2.EndTime))
	if tb.StartTime == tb2.StartTime && tb.EndTime == tb2.EndTime {
		return true
	}
	return false
}

func (tb *TimeBounds) String() string {
	return fmt.Sprintf("%+v, %+v", tb.StartTime, tb.EndTime)
}

// ----------------------------------------------------------------------------------------------------

type UphillDownhill struct {
	Uphill   float64
	Downhill float64
}

func (ud *UphillDownhill) Equals(ud2 *UphillDownhill) bool {
	if ud.Uphill == ud2.Uphill && ud.Downhill == ud2.Downhill {
		return true
	}
	return false
}

// ----------------------------------------------------------------------------------------------------

type LocationsResultPair struct {
	SegmentNo int
	PointNo   int
}

// ----------------------------------------------------------------------------------------------------

type GPXPoint struct {
	Point
	// TODO
	Timestamp time.Time
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

// SpeedBetween calculates the speed between two GpxWpts.
func (pt *GPXPoint) SpeedBetween(pt2 *GPXPoint, threeD bool) float64 {
	seconds := pt.TimeDiff(pt2)
	var distLen float64
	if threeD {
		distLen = pt.Distance3D(&pt2.Point)
	} else {
		distLen = pt.Distance2D(&pt2.Point)
	}

	return distLen / seconds
}

// TimeDiff returns the time difference of two GpxWpts in seconds.
func (pt *GPXPoint) TimeDiff(pt2 *GPXPoint) float64 {
	t1 := pt.Timestamp
	t2 := pt2.Timestamp

	if t1.Equal(t2) {
		return 0.0
	}

	var delta time.Duration
	if t1.After(t2) {
		delta = t1.Sub(t2)
	} else {
		delta = t2.Sub(t1)
	}

	return delta.Seconds()
}

// MaxDilutionOfPrecision returns the dilution precision of a GpxWpt.
func (pt *GPXPoint) MaxDilutionOfPrecision() float64 {
	return math.Max(pt.HorizontalDilution, math.Max(pt.VerticalDilution, pt.PositionalDilution))
}

// ----------------------------------------------------------------------------------------------------

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

// Length returns the length of a GPX route.
func (rte *GPXRoute) Length() float64 {
	// TODO: npe check
	points := make([]*Point, len(rte.Points))
	for pointNo, point := range rte.Points {
		points[pointNo] = &point.Point
	}
	return Length2D(points)
}

// Center returns the center of a GPX route.
func (rte *GPXRoute) Center() (float64, float64) {
	lenRtePts := len(rte.Points)
	if lenRtePts == 0 {
		return 0.0, 0.0
	}

	var (
		sumLat float64
		sumLon float64
	)

	for _, pt := range rte.Points {
		sumLat += pt.Latitude
		sumLon += pt.Longitude
	}

	n := float64(lenRtePts)
	return sumLat / n, sumLon / n
}

func (rte *GPXRoute) ExecuteOnPoints(executor func(*GPXPoint)) {
	for _, point := range rte.Points {
		executor(point)
	}
}

func (rte *GPXRoute) RemoveElevation() {
	for _, point := range rte.Points {
		point.RemoveElevation()
	}
}

// ----------------------------------------------------------------------------------------------------

type GPXTrackSegment struct {
	Points []*GPXPoint
	// TODO extensions
}

// Length2D returns the 2D length of a GPX segment.
func (seg *GPXTrackSegment) Length2D() float64 {
	// TODO: There should be a better way to do this:
	points := make([]*Point, len(seg.Points))
	for pointNo, point := range seg.Points {
		points[pointNo] = &point.Point
	}
	return Length2D(points)
}

// Length3D returns the 3D length of a GPX segment.
func (seg *GPXTrackSegment) Length3D() float64 {
	// TODO: There should be a better way to do this:
	points := make([]*Point, len(seg.Points))
	for pointNo, point := range seg.Points {
		points[pointNo] = &point.Point
	}
	return Length3D(points)
}

// TimeBounds returns the time bounds of a GPX segment.
func (seg *GPXTrackSegment) TimeBounds() *TimeBounds {
	timeTuple := make([]time.Time, 0)

	for _, trkpt := range seg.Points {
		if len(timeTuple) < 2 {
			timeTuple = append(timeTuple, trkpt.Timestamp)
		} else {
			timeTuple[1] = trkpt.Timestamp
		}
	}
	if len(timeTuple) == 2 {
		return &TimeBounds{StartTime: timeTuple[0], EndTime: timeTuple[1]}
	}
	return nil
}

// Bounds returns the bounds of a GPX segment.
func (seg *GPXTrackSegment) Bounds() *GpxBounds {
	minmax := getMinimaMaximaStart()
	for _, pt := range seg.Points {
		minmax.MaxLat = math.Max(pt.Latitude, minmax.MaxLat)
		minmax.MinLat = math.Min(pt.Latitude, minmax.MinLat)
		minmax.MaxLon = math.Max(pt.Longitude, minmax.MaxLon)
		minmax.MinLon = math.Min(pt.Longitude, minmax.MinLon)
	}
	return minmax
}

func (seg *GPXTrackSegment) HasTimes() bool {
    return false
    /*
    withTimes := 0
    for _, point := range seg.Points {
        if point.Timestamp != nil {
            withTimes += 1
        }
    }
    return withTimes / len(seg.Points) >= 0.75
    */
}

// Speed returns the speed at point number in a GPX segment.
func (seg *GPXTrackSegment) Speed(pointIdx int) float64 {
	trkptsLen := len(seg.Points)
	if pointIdx >= trkptsLen {
		pointIdx = trkptsLen - 1
	}

	point := seg.Points[pointIdx]

	var prevPt *GPXPoint
	var nextPt *GPXPoint

	havePrev := false
	haveNext := false
	if 0 < pointIdx && pointIdx < trkptsLen {
		prevPt = seg.Points[pointIdx-1]
		havePrev = true
	}

	if 0 < pointIdx && pointIdx < trkptsLen-1 {
		nextPt = seg.Points[pointIdx+1]
		haveNext = true
	}

	haveSpeed1 := false
	haveSpeed2 := false

	var speed1 float64
	var speed2 float64
	if havePrev {
		speed1 = math.Abs(point.SpeedBetween(prevPt, true))
		haveSpeed1 = true
	}
	if haveNext {
		speed2 = math.Abs(point.SpeedBetween(nextPt, true))
		haveSpeed2 = true
	}

	if haveSpeed1 && haveSpeed2 {
		return (speed1 + speed2) / 2.0
	}

	if haveSpeed1 {
		return speed1
	}
	return speed2
}

// Duration returns the duration in seconds in a GPX segment.
func (seg *GPXTrackSegment) Duration() float64 {
	trksLen := len(seg.Points)
	if trksLen == 0 {
		return 0.0
	}

	first := seg.Points[0]
	last := seg.Points[trksLen-1]

	firstTimestamp := first.Timestamp
	lastTimestamp := last.Timestamp

	if firstTimestamp.Equal(lastTimestamp) {
		return 0.0
	}

	if lastTimestamp.Before(firstTimestamp) {
		return 0.0
	}
	dur := lastTimestamp.Sub(firstTimestamp)

	return dur.Seconds()
}

// Elevations returns a slice with the elevations in a GPX segment.
func (seg *GPXTrackSegment) Elevations() []float64 {
	elevations := make([]float64, len(seg.Points))
	for i, trkpt := range seg.Points {
		elevations[i] = trkpt.Elevation
	}
	return elevations
}

// UphillDownhill returns uphill and dowhill in a GPX segment.
func (seg *GPXTrackSegment) UphillDownhill() *UphillDownhill {
	if len(seg.Points) == 0 {
		return nil
	}

	elevations := seg.Elevations()

	uphill, downhill := CalcUphillDownhill(elevations)

	return &UphillDownhill{Uphill: uphill, Downhill: downhill}
}

func (seg *GPXTrackSegment) ExecuteOnPoints(executor func(*GPXPoint)) {
	for _, point := range seg.Points {
		executor(point)
	}
}

func (seg *GPXTrackSegment) AddElevation(elevation float64) {
	for _, point := range seg.Points {
		point.Elevation += elevation
	}
}

func (seg *GPXTrackSegment) RemoveElevation() {
	for _, point := range seg.Points {
		point.RemoveElevation()
	}
}

// Split splits a GPX segment at point index pt. Point pt remains in
// first part.
func (seg *GPXTrackSegment) Split(pt int) (*GPXTrackSegment, *GPXTrackSegment) {
	pts1 := seg.Points[:pt+1]
	pts2 := seg.Points[pt+1:]

	return &GPXTrackSegment{Points: pts1}, &GPXTrackSegment{Points: pts2}
}

// Join concatenates to GPX segments.
func (seg *GPXTrackSegment) Join(seg2 *GPXTrackSegment) {
	seg.Points = append(seg.Points, seg2.Points...)
}

// LocationAt returns the GpxWpt at a given time.
func (seg *GPXTrackSegment) LocationAt(t time.Time) int {
	lenPts := len(seg.Points)
	if lenPts == 0 {
		return -1
	}
	firstT := seg.Points[0]
	lastT := seg.Points[lenPts-1]

	firstTimestamp := firstT.Timestamp
	lastTimestamp := lastT.Timestamp

	if firstTimestamp.Equal(lastTimestamp) || firstTimestamp.After(lastTimestamp) {
		return -1
	}

	for i := 0; i < len(seg.Points); i++ {
		pt := seg.Points[i]
		if t.Before(pt.Timestamp) {
			return i
		}
	}

	return -1
}

// MovingData returns the moving data of a GPX segment.
func (seg *GPXTrackSegment) MovingData() *MovingData {
	var (
		movingTime      float64
		stoppedTime     float64
		movingDistance  float64
		stoppedDistance float64
	)

	speedsDistances := make([]SpeedsAndDistances, 0)

	for i := 1; i < len(seg.Points); i++ {
		prev := seg.Points[i-1]
		pt := seg.Points[i]

		dist := pt.Distance3D(&prev.Point)

		timedelta := pt.Timestamp.Sub(prev.Timestamp)
		seconds := timedelta.Seconds()
		var speedKmh float64

		if seconds > 0 {
			speedKmh = (dist / 1000.0) / (timedelta.Seconds() / math.Pow(60, 2))
		}

		if speedKmh <= DEFAULT_STOPPED_SPEED_THRESHOLD {
			stoppedTime += timedelta.Seconds()
			stoppedDistance += dist
		} else {
			movingTime += timedelta.Seconds()
			movingDistance += dist

			sd := SpeedsAndDistances{dist / timedelta.Seconds(), dist}
			speedsDistances = append(speedsDistances, sd)
		}
	}

	var maxSpeed float64
	if len(speedsDistances) > 0 {
		maxSpeed = CalcMaxSpeed(speedsDistances)
	}

	return &MovingData{
		movingTime,
		stoppedTime,
		movingDistance,
		stoppedDistance,
		maxSpeed,
	}
}

func (seg *GPXTrackSegment) AppendPoint(p *GPXPoint) {
	seg.Points = append(seg.Points, p)
}

// ----------------------------------------------------------------------------------------------------

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

// Length2D returns the 2D length of a GPX track.
func (trk *GPXTrack) Length2D() float64 {
	var l float64
	for _, seg := range trk.Segments {
		d := seg.Length2D()
		l += d
	}
	return l
}

// Length3D returns the 3D length of a GPX track.
func (trk *GPXTrack) Length3D() float64 {
	var l float64
	for _, seg := range trk.Segments {
		d := seg.Length3D()
		l += d
	}
	return l
}

// TimeBounds returns the time bounds of a GPX track.
func (trk *GPXTrack) TimeBounds() *TimeBounds {
	var tbTrk *TimeBounds

	for i, seg := range trk.Segments {
		tbSeg := seg.TimeBounds()
		if i == 0 {
			tbTrk = tbSeg
		} else {
			tbTrk.EndTime = tbSeg.EndTime
		}
	}
	return tbTrk
}

// Bounds returns the bounds of a GPX track.
func (trk *GPXTrack) Bounds() *GpxBounds {
	minmax := getMinimaMaximaStart()
	for _, seg := range trk.Segments {
		bnds := seg.Bounds()
		minmax.MaxLat = math.Max(bnds.MaxLat, minmax.MaxLat)
		minmax.MinLat = math.Min(bnds.MinLat, minmax.MinLat)
		minmax.MaxLon = math.Max(bnds.MaxLon, minmax.MaxLon)
		minmax.MinLon = math.Min(bnds.MinLon, minmax.MinLon)
	}
	return minmax
}

func (trk *GPXTrack) HasTimes() bool {
    result := true
    for _, segment := range trk.Segments {
        result = result && segment.HasTimes()
    }
    return result
}

// Split splits a GPX segment at a point number ptNo in a GPX track.
func (trk *GPXTrack) Split(segNo, ptNo int) {
	lenSegs := len(trk.Segments)
	if segNo >= lenSegs {
		return
	}

	newSegs := make([]*GPXTrackSegment, 0)
	for i := 0; i < lenSegs; i++ {
		seg := trk.Segments[i]

		if i == segNo && ptNo < len(seg.Points) {
			seg1, seg2 := seg.Split(ptNo)
			newSegs = append(newSegs, seg1, seg2)
		} else {
			newSegs = append(newSegs, seg)
		}
	}
	trk.Segments = newSegs
}

func (trk *GPXTrack) ExecuteOnPoints(executor func(*GPXPoint)) {
	for _, segment := range trk.Segments {
		segment.ExecuteOnPoints(executor)
	}
}

func (trk *GPXTrack) AddElevation(elevation float64) {
	for _, segment := range trk.Segments {
		segment.AddElevation(elevation)
	}
}

func (trk *GPXTrack) RemoveElevation() {
	for _, segment := range trk.Segments {
		segment.RemoveElevation()
	}
}

// Join joins two GPX segments in a GPX track.
func (trk *GPXTrack) Join(segNo, segNo2 int) {
	lenSegs := len(trk.Segments)
	if segNo >= lenSegs && segNo2 >= lenSegs {
		return
	}
	newSegs := make([]*GPXTrackSegment, 0)
	for i := 0; i < lenSegs; i++ {
		seg := trk.Segments[i]
		if i == segNo {
			secondSeg := trk.Segments[segNo2]
			seg.Join(secondSeg)
			newSegs = append(newSegs, seg)
		} else if i == segNo2 {
			// do nothing, its already joined
		} else {
			newSegs = append(newSegs, seg)
		}
	}
	trk.Segments = newSegs
}

// JoinNext joins a GPX segment with the next segment in the current GPX
// track.
func (trk *GPXTrack) JoinNext(segNo int) {
	trk.Join(segNo, segNo+1)
}

// MovingData returns the moving data of a GPX track.
func (trk *GPXTrack) MovingData() *MovingData {
	var (
		movingTime      float64
		stoppedTime     float64
		movingDistance  float64
		stoppedDistance float64
		maxSpeed        float64
	)

	for _, seg := range trk.Segments {
		md := seg.MovingData()
		movingTime += md.MovingTime
		stoppedTime += md.StoppedTime
		movingDistance += md.MovingDistance
		stoppedDistance += md.StoppedDistance

		if md.MaxSpeed > maxSpeed {
			maxSpeed = md.MaxSpeed
		}
	}
	return &MovingData{
		MovingTime:      movingTime,
		MovingDistance:  movingDistance,
		StoppedTime:     stoppedTime,
		StoppedDistance: stoppedDistance,
		MaxSpeed:        maxSpeed,
	}
}

// Duration returns the duration of a GPX track.
func (trk *GPXTrack) Duration() float64 {
	if len(trk.Segments) == 0 {
		return 0.0
	}

	var result float64
	for _, seg := range trk.Segments {
		result += seg.Duration()
	}
	return result
}

// UphillDownhill return the uphill and downhill values of a GPX track.
func (trk *GPXTrack) UphillDownhill() *UphillDownhill {
	if len(trk.Segments) == 0 {
		return nil
	}

	var (
		uphill   float64
		downhill float64
	)

	for _, seg := range trk.Segments {
		updo := seg.UphillDownhill()

		uphill += updo.Uphill
		downhill += updo.Downhill
	}

	return &UphillDownhill{
		Uphill:   uphill,
		Downhill: downhill,
	}
}

// LocationAt returns a LocationResultsPair for a given time.
func (trk *GPXTrack) LocationAt(t time.Time) []LocationsResultPair {
	results := make([]LocationsResultPair, 0)

	for i := 0; i < len(trk.Segments); i++ {
		seg := trk.Segments[i]
		loc := seg.LocationAt(t)
		if loc != -1 {
			results = append(results, LocationsResultPair{i, loc})
		}
	}
	return results
}

func (trk *GPXTrack) AppendSegment(s *GPXTrackSegment) {
	trk.Segments = append(trk.Segments, s)
}

// ----------------------------------------------------------------------------------------------------

/**
 * Useful when looking for smaller bounds
 *
 * TODO does it work is region is between 179E and 179W?
 */
func getMinimaMaximaStart() *GpxBounds {
	return &GpxBounds{
		MaxLat: -math.MaxFloat64,
		MinLat: math.MaxFloat64,
		MaxLon: -math.MaxFloat64,
		MinLon: math.MaxFloat64,
	}
}
