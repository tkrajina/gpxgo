package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/tkrajina/gpxgo/gpx"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	for _, fn := range flag.Args() {
		g, err := gpx.ParseFile(fn)
		panicIfErr(err)
		byts, err := json.MarshalIndent(g, "", "\t")
		panicIfErr(err)
		fmt.Println(string(byts))
	}
}
