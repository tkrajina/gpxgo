package main

import (
	"fmt"
	"os"

	"github.com/tkrajina/gpxgo/gpx"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	outfile = "gpxjson.ts"
)

func main() {
	ts := typescriptify.New().
		WithBackupDir("").
		// WithPrefix("API_").
		Add(gpx.GPX{}).
		WithConstructor(false).
		WithCreateFromMethod(false)

	str, err := ts.Convert(nil)
	panicIfErr(err)

	panicIfErr(os.WriteFile(outfile, []byte(str), 0700))
	fmt.Println("Generated", outfile)
}
