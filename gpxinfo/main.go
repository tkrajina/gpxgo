package main

import (
	"flag"
	"fmt"
	"github.com/ptrv/go-gpx"
	"path/filepath"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Please provide a GPX file path!")
		return
	}

	gpxFileArg := args[0]
	gpxFile, err := gpx.Parse(gpxFileArg)

	if err != nil {
		fmt.Println("Error opening gpx file: ", err)
		return
	}

	gpxPath, _ := filepath.Abs(gpxFileArg)
	fmt.Println("File: ", gpxPath)

	fmt.Println("\tGPX name: ", gpxFile.Metadata.Name)
	fmt.Println("\tGPX desctiption: ", gpxFile.Metadata.Desc)
	fmt.Println("\tAuthor: ", gpxFile.Metadata.Author.Name)
	fmt.Println("\tEmail: ", gpxFile.Metadata.Author.Email)

	len2d := gpxFile.Length2D()
	len3d := gpxFile.Length3D()
	fmt.Println("Length 2D: ", len2d)
	fmt.Println("Length 3D: ", len3d)

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
