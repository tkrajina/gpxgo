package gpx

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"
)

const timeLayout = "2006-01-02T15:04:05Z"

// type GeoLocation interface {
// 	Distance2D(l GeoLocation) float64
// 	Distance3D(l GeoLocation) float64
// 	ElevationAngle(l GeoLocation, radians bool) float64
// }

type GPXTrackpoint struct {
	// GeoLocation
	Lat       float64 `xml:"lat,attr"`
	Lon       float64 `xml:"lon,attr"`
	Ele       float64 `xml:"ele"`
	Timestamp string  `xml:"time"`
	Cmt       string  `xml:"cmt"`
	Speed     float64 `xml:"speed"`
}
type GPXTrackseg struct {
	Trkpts []GPXTrackpoint `xml:"trkpt"`
}
type GPXTrack struct {
	Name   string        `xml:"name"`
	Cmt    string        `xml:"cmt"`
	Desc   string        `xml:"desc"`
	Src    string        `xml:"src"`
	Trkseg []GPXTrackseg `xml:"trkseg"`
	Bnds   Bounds
}

type GPXMetadata struct {
	Name      string `xml:"name"`
	Desc      string `xml:"desc"`
	Author    string `xml:"author"`
	Copyright string `xml:"copyright"`
	Timestamp string `xml:"time"`
}
type GPX struct {
	Metadata GPXMetadata `xml:"metadata"`
	Tracks   []GPXTrack  `xml:"trk"`
}

type TimeBounds struct {
	StartTime time.Time
	EndTime   time.Time
}

type Bounds struct {
	MinLat float64
	MaxLat float64
	MinLon float64
	MaxLon float64
}

type UphillDownhill struct {
	Uphill   float64
	Downhill float64
}

//==========================================================
func Parse(gpxPath string) (GPX, error) {
	gpxFile, err := os.Open(gpxPath)

	if err != nil {
		// fmt.Println("Error opening file: ", err)
		return GPX{}, err
	}
	defer gpxFile.Close()

	b, err := ioutil.ReadAll(gpxFile)

	if err != nil {
		// fmt.Println("Error reading file: ", err)
		return GPX{}, err
	}
	var g GPX
	xml.Unmarshal(b, &g)

	return g, nil
}

//==========================================================
func getTime(timestr string) time.Time {
	t, err := time.Parse(timeLayout, timestr)
	if err != nil {
		return time.Time{}
	}
	return t
}

//==========================================================
func (tb *TimeBounds) Equals(tb2 TimeBounds) bool {
	if tb.StartTime == tb2.StartTime && tb.EndTime == tb2.EndTime {
		return true
	}
	return false
}

func (tb *TimeBounds) String() string {
	return fmt.Sprintf("%+v, %+v", tb.StartTime, tb.EndTime)
}

func (b *Bounds) Equals(b2 Bounds) bool {
	if b.MinLon == b2.MinLon && b.MaxLat == b2.MaxLat && b.MinLon == b2.MinLon && b.MaxLon == b.MaxLon {
		return true
	}
	return false
}

func (b *Bounds) String() string {
	return fmt.Sprintf("%+v, %+v, %+v, %+v", b.MinLat, b.MaxLat, b.MinLat, b.MaxLon)
}

//==========================================================

func (trk *GPXTrack) Length2D() float64 {
	var l float64
	for _, seg := range trk.Trkseg {
		d := seg.Length2D()
		l += d
	}
	return l
}

//==========================================================
// func (seg *GPXTrackseg) Locs() []GeoLocation {
// 	locs := make([]GeoLocation, len(seg.Trkpts))
// 	for i, pt := range seg.Trkpts {
// 		var loc GeoLocation
// 		// loc.Lat = pt.Lat
// 		// loc.Lon = pt.Lon
// 		// loc.Ele = pt.Ele
// 		locs[i] = pt
// 	}
// 	return locs
// }
func (seg *GPXTrackseg) Length2D() float64 {
	return Length2D(seg.Trkpts)
}

func (seg *GPXTrackseg) Length3D() float64 {
	return Length3D(seg.Trkpts)
}

func (seg *GPXTrackseg) TimeBounds() TimeBounds {
	timeTuple := make([]time.Time, 0)

	for _, trkpt := range seg.Trkpts {
		if trkpt.Timestamp != "" {
			if len(timeTuple) < 2 {
				timeTuple = append(timeTuple, trkpt.Time())
			} else {
				timeTuple[1] = trkpt.Time()
			}
		}
	}
	if len(timeTuple) == 2 {
		return TimeBounds{StartTime: timeTuple[0], EndTime: timeTuple[1]}
	}
	return TimeBounds{}
}

func (seg *GPXTrackseg) Bounds() Bounds {

	maxLat := -math.MaxFloat64
	minLat := math.MaxFloat64
	maxLon := -math.MaxFloat64
	minLon := math.MaxFloat64

	for _, trkpt := range seg.Trkpts {
		maxLat = math.Max(trkpt.Lat, maxLat)
		minLat = math.Min(trkpt.Lat, minLat)
		maxLon = math.Max(trkpt.Lon, maxLon)
		minLon = math.Min(trkpt.Lon, minLon)
	}

	return Bounds{
		MaxLat: maxLat, MinLat: minLat,
		MaxLon: maxLon, MinLon: minLon,
	}
}

// Get speed at point
func (seg *GPXTrackseg) Speed(pointIdx int) float64 {
	trkptsLen := len(seg.Trkpts)
	if pointIdx >= trkptsLen {
		pointIdx = trkptsLen - 1
	}

	point := seg.Trkpts[pointIdx]

	var prevPt GPXTrackpoint
	var nextPt GPXTrackpoint

	havePrev := false
	haveNext := false
	if 0 < pointIdx && pointIdx < trkptsLen {
		prevPt = seg.Trkpts[pointIdx-1]
		havePrev = true
	}

	if 0 < pointIdx && pointIdx < trkptsLen-1 {
		nextPt = seg.Trkpts[pointIdx+1]
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

// Duration in seconds
func (seg *GPXTrackseg) Duration() float64 {
	trksLen := len(seg.Trkpts)
	if trksLen == 0 {
		return 0.0
	}

	first := seg.Trkpts[0]
	last := seg.Trkpts[trksLen-1]

	if first.Time().Equal(last.Time()) {
		return 0.0
	}

	if last.Time().Before(first.Time()) {
		return 0.0
	}
	dur := last.Time().Sub(first.Time())

	return dur.Seconds()
}

func (seg *GPXTrackseg) Elevations() []float64 {
	elevations := make([]float64, len(seg.Trkpts))
	for i, trkpt := range seg.Trkpts {
		elevations[i] = trkpt.Ele
	}
	return elevations
}

// Return uphill and dowhill
func (seg *GPXTrackseg) UphillDownhill() UphillDownhill {
	if len(seg.Trkpts) == 0 {
		return UphillDownhill{}
	}

	elevations := seg.Elevations()

	uphill, downhill := CalcUphillDownhill(elevations)

	return UphillDownhill{Uphill: uphill, Downhill: downhill}
}

//==========================================================

// Return Timestamp string as Time object
func (pt *GPXTrackpoint) Time() time.Time {
	return getTime(pt.Timestamp)
}

// Time difference of two GPXTrackpoints in seconds
func (pt *GPXTrackpoint) TimeDiff(pt2 GPXTrackpoint) float64 {
	t1 := pt.Time()
	t2 := pt2.Time()

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

func (pt *GPXTrackpoint) SpeedBetween(pt2 GPXTrackpoint, threeD bool) float64 {
	seconds := pt.TimeDiff(pt2)
	var distLen float64
	if threeD {
		distLen = pt.Distance3D(pt2)
	} else {
		distLen = pt.Distance2D(pt2)
	}

	return distLen / seconds
}

func (pt *GPXTrackpoint) Distance2D(pt2 GPXTrackpoint) float64 {
	return Distance2D(pt.Lat, pt.Lon, pt2.Lat, pt2.Lon, false)
}
func (pt *GPXTrackpoint) Distance3D(pt2 GPXTrackpoint) float64 {
	return Distance3D(pt.Lat, pt.Lon, pt.Ele, pt2.Lat, pt2.Lon, pt2.Ele, false)
}

//==========================================================
