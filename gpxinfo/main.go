package main

import (
	"flag"
	"fmt"
	"github.com/ptrv/go-gpx"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Please provide a GPX file path!")
		return
	}

	g, err := gpx.Parse(args[0])

	if err != nil {
		fmt.Println("Error opening gpx file: ", err)
		return
	}

	fmt.Println(g.Metadata.Timestamp)
	// for _, trk := range g.Tracks {
	// 	fmt.Printf("%s\n", trk.Name)
	// 	for _, trkseg := range trk.Trkseg {
	// 		for _, trkpt := range trkseg.Trkpts {
	// 			fmt.Printf("%f, %f, %f, %s\n", trkpt.Lat, trkpt.Lon, trkpt.Ele, trkpt.Timestamp)
	// 		}
	// 		fmt.Println("Length2D: ", trkseg.Length2D())
	// 	}
	// }

}
