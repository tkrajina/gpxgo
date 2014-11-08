package gpx

import "testing"

func TestParseTime(t *testing.T) {
	time, err := parseGPXTime("")
	if time != nil {
		t.Errorf("Empty string should not return a nonnil time")
	}
	if err == nil {
		t.Errorf("Empty string should result in error")
	}
}
