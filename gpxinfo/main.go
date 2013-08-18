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
	fmt.Println("\tLength 2D: ", len2d/1000.0)
	fmt.Println("\tLength 3D: ", len3d/1000.0)

	fmt.Printf("\tBounds: %+v\n", gpxFile.Bounds())

	md := gpxFile.MovingData()
	fmt.Println("\tMoving time: ", md.MovingTime)
	fmt.Println("\tStopped time: ", md.StoppedTime)

	fmt.Printf("\tMax speed: %fm/s = %fkm/h\n", md.MaxSpeed, md.MaxSpeed*60*60/1000.0)

	updo := gpxFile.UphillDownhill()
	fmt.Println("\tTotal uphill: ", updo.Uphill)
	fmt.Println("\tTotal downhill: ", updo.Downhill)

	timeBounds := gpxFile.TimeBounds()
	fmt.Println("\tStarted: ", timeBounds.StartTime)
	fmt.Println("\tEnded: ", timeBounds.EndTime)

	fmt.Println()
}
