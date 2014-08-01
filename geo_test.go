// Copyright 2013 Peter Vasil. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gpx

import (
	"math"
	"testing"
)

func TestToRad(t *testing.T) {
	radVal := ToRad(360)
	if radVal != math.Pi*2 {
		t.Errorf("Test failed: %f", radVal)
	}
}

func TestElevationAngle(t *testing.T) {
	loc1 := Point{Latitude: 52.5113534275, Longitue: 13.4571944922, Elevation: 59.26}
	loc2 := Point{Latitude: 52.5113568641, Longitue: 13.4571697656, Elevation: 65.51}

	elevAngleA := ElevationAngle(&loc1, &loc2, false)
	elevAngleE := 74.65347905197362

	if elevAngleE != elevAngleA {
		t.Errorf("Elevation angle expected: %f, actual: %f", elevAngleE, elevAngleA)
	}
}
