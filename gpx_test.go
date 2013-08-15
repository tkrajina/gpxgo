package gpx

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"testing"
	"time"
)

var g GPX

func init() {
	fmt.Println("gpx test init")

	gpxFile, err := os.Open("testdata/file.gpx")

	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer gpxFile.Close()

	b, _ := ioutil.ReadAll(gpxFile)

	xml.Unmarshal(b, &g)
}

func TestToRad(t *testing.T) {
	radVal := ToRad(360)
	if radVal != math.Pi*2 {
		t.Errorf("Test failed: %f", radVal)
	}
}

func TestGpxParsing(t *testing.T) {
	// t.Log("Test parser")
	timestampA := g.Metadata.Timestamp
	timestampE := "2012-03-17T15:44:18Z"
	if timestampA != timestampE {
		t.Errorf("timestamp expected: %s, actual: %s", timestampE, timestampA)
	}

	trknameA := g.Tracks[0].Name
	trknameE := "17-MRZ-12 16:44:12"
	if trknameA != trknameE {
		t.Errorf("Trackname expected: %s, actual: %s", trknameE, trknameA)
	}
}

func TestLength2D(t *testing.T) {
	lengthA := g.Tracks[0].Trkseg[0].Length2D()
	lengthE := 56.77577732775905

	if lengthA != lengthE {
		t.Errorf("Length 2d expected: %f, actual %f", lengthE, lengthA)
	}
}

func TestLength3D(t *testing.T) {
	lengthA := g.Tracks[0].Trkseg[0].Length3D()
	lengthE := 61.76815317436073

	if lengthA != lengthE {
		t.Errorf("Length 3d expected: %f, actual %f", lengthE, lengthA)
	}
}

func TestGetTime(t *testing.T) {
	timestampA := getTime("2012-03-17T12:46:19Z")
	timestampE := time.Date(2012, 3, 17, 12, 46, 19, 0, time.UTC)

	if timestampA != timestampE {
		t.Errorf("Time expected: %s, actual: %s", timestampE.String(), timestampA.String())
	}
}

func TestPtTime(t *testing.T) {
	timeA := g.Tracks[0].Trkseg[0].Trkpts[0].Time()
	//2012-03-17T12:46:19Z
	timeE := time.Date(2012, 3, 17, 12, 46, 19, 0, time.UTC)

	if timeA != timeE {
		t.Errorf("Time expected: %s, actual: %s", timeE.String(), timeA.String())
	}
}

func TestTimeBounds(t *testing.T) {
	timeBoundsA := g.Tracks[0].Trkseg[0].TimeBounds()
	timeBoundsE := TimeBounds{
		StartTime: time.Date(2012, 3, 17, 12, 46, 19, 0, time.UTC),
		EndTime:   time.Date(2012, 3, 17, 12, 47, 23, 0, time.UTC),
	}

	if !timeBoundsE.Equals(timeBoundsA) {
		t.Errorf("TimeBounds expected: %s, actual: %s", timeBoundsE.String(), timeBoundsA.String())
	}
}

func TestBounds(t *testing.T) {
	boundsA := g.Tracks[0].Trkseg[0].Bounds()
	boundsE := Bounds{
		MaxLat: 52.5117189623, MinLat: 52.5113534275,
		MaxLon: 13.4571944922, MinLon: 13.4567520116,
	}

	if !boundsE.Equals(boundsA) {
		t.Errorf("Bounds expected: %s, actual: %s", boundsE.String(), boundsA.String())
	}
}

func TestSpeed(t *testing.T) {
	speedA := g.Tracks[0].Trkseg[0].Speed(2)
	speedE := 1.5386074011963367

	if speedE != speedA {
		t.Errorf("Speed expected: %f, actual: %f", speedE, speedA)
	}
}

func TestDuration(t *testing.T) {
	durA := g.Tracks[0].Trkseg[0].Duration()
	durE := 64.0

	if durE != durA {
		t.Errorf("Duration expected: %f, actual: %f", durE, durA)
	}
}

func TestUphillDownHill(t *testing.T) {
	updoA := g.Tracks[0].Trkseg[0].UphillDownhill()
	updoE := UphillDownhill{
		Uphill:   5.863000000000007,
		Downhill: 1.5430000000000064}

	fmt.Println(updoE.Uphill, updoE.Downhill, updoA.Uphill, updoA.Downhill)

	if updoE.Uphill != updoA.Uphill || updoE.Downhill != updoA.Downhill {
		t.Errorf("UphillDownhill expected: %f, %f, actual: %f, %f",
			updoE.Uphill, updoE.Downhill, updoA.Uphill, updoA.Downhill)
	}
}
