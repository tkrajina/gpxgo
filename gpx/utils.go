package gpx

import (
	"encoding/json"
	"fmt"
)

var Debug = false

func jsonizeFormatted(a any) string {
	jsn, _ := json.MarshalIndent(a, "", "  ")
	return string(jsn)
}

func jsonize(a any) string {
	jsn, _ := json.Marshal(a)
	return string(jsn)
}

func debugf(s string, a ...any) {
	if Debug {
		fmt.Printf(s, a...)
	}
}

func debugJSONized(msg string, a any) {
	byts, _ := json.MarshalIndent(a, "", "\t")
	debugln(msg + ":\n" + string(byts) + "\n")
}

func debugln(a ...any) {
	if Debug {
		fmt.Println(a...)
	}
}
