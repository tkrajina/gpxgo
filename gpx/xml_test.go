// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"testing"
)

func TestParseTime(t *testing.T) {
	time, err := parseGPXTime("")
	if time != nil {
		t.Errorf("Empty string should not return a nonnil time")
	}
	if err == nil {
		t.Errorf("Empty string should result in error")
	}
}
