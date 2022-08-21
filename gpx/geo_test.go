// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToRad(t *testing.T) {
	radVal := ToRad(360)
	if radVal != math.Pi*2 {
		t.Errorf("Test failed: %f", radVal)
	}
}

func TestElevationAngle(t *testing.T) {
	loc1 := Point{Latitude: 52.5113534275, Longitude: 13.4571944922, Elevation: *NewNullableFloat64(59.26)}
	loc2 := Point{Latitude: 52.5113568641, Longitude: 13.4571697656, Elevation: *NewNullableFloat64(65.51)}

	elevAngleA := ElevationAngle(loc1, loc2, false)
	elevAngleE := 74.65347905197362

	if elevAngleE != elevAngleA {
		t.Errorf("Elevation angle expected: %f, actual: %f", elevAngleE, elevAngleA)
	}
}

func TestMaxSpeed(t *testing.T) {
	t.Parallel()

	maxSpeed := CalcMaxSpeed([]SpeedsAndDistances{
		{Speed: 5.0, Distance: 508.674260463},
		{Speed: 4.0, Distance: 593.443625286},
		{Speed: 6.0, Distance: 523.841129461},
		{Speed: 1.0, Distance: 489.306355103},
	})
	assert.Equal(t, 6.0, maxSpeed)
}

func TestVincenty(t *testing.T) {
	t.Parallel()

	for i := 0; i < 100; i++ {
		maxAway := 2.
		lat1, lon1 := rand.Float64()*180-90, rand.Float64()*360-180
		//lat2, lon2 := rand.Float64()*180-90, rand.Float64()*360-180
		lat2, lon2 := lat1+(rand.Float64()-0.5)*maxAway, lon1+(rand.Float64()-0.5)*maxAway
		p1 := Point{Latitude: lat1, Longitude: lon1}
		p2 := Point{Latitude: lat2, Longitude: lon2}
		distHaversine := HaversineDistance(lat1, lon1, lat2, lon2)
		distVincenty, err := DistanceVincenty(p1, p2, false)
		diff := math.Abs(distHaversine - distVincenty)
		assert.Nil(t, err)
		assert.True(t, diff < math.Max(distHaversine, distVincenty)*0.0001, "p1: %#v, p2: %#v, haversine: %f, vincenty: %f, diff: %f", p1, p2, distHaversine, distVincenty, diff)
		if t.Failed() {
			t.FailNow()
		}
	}
}
