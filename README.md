# Go GPX library

gpxgo is a golang library for parsing and manipulating GPX files. GPX (GPS eXchange Format) is a XML based file format for GPS track logs. 

## Example:

    import (
        ...
        "github.com/tkrajina/gpxgo/gpx"
        ...
    )

    gpxBytes := ...
	gpxFile, err := gpx.ParseBytes(gpxBytes)
	if err != nil {
        ...
	}

    // Analyize/manipulate your track data here...
	for _, track := range gpxFile.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				fmt.Print(point)
			}
		}
	}

    // (Check the API for GPX manipulation and analyzing utility methods)

    // When ready, you can write the resulting GPX file:
	xmlBytes, err := gpxFile.ToXml(gpx.ToXmlParams{Version: "1.1", Indent: true})
    ...

## GPX Compatibility

Gpxgo can read/write both GPX 1.0 and GPX 1.1 files. The only not-yet-supported part of the GPX 1.1 specification are extensions.

## gpxinfo

`gpxinfo` is a command line utility for writing basic stats from gpx files:

    $ go run gpxinfo.go test_files/Mojstrovka.gpx
    File: /Users/puzz/golang/src/github.com/tkrajina/gpxgo/test_files/Mojstrovka.gpx
    GPX name:
    GPX desctiption:
    GPX version: 1.0
    Author:
    Email:


    Global stats:
     Points: 184
     Length 2D: 2.6958067369682577
     Length 3D: 3.00439590990862
     Bounds: 46.430350, 46.435641, 13.738842, 13.748333
     Moving time: 0
     Stopped time: 0
     Max speed: 0.000000m/s = 0.000000km/h
     Total uphill: 446.4893280000001
     Total downhill: 417.65524800000026
     Started: 1901-12-13 20:45:52 +0000 UTC
     Ended: 1901-12-13 20:45:52 +0000 UTC


    Track #1:
         Points: 184
         Length 2D: 2.6958067369682577
         Length 3D: 3.00439590990862
         Bounds: 46.430350, 46.435641, 13.738842, 13.748333
         Moving time: 0
         Stopped time: 0
         Max speed: 0.000000m/s = 0.000000km/h
         Total uphill: 446.4893280000001
    ...etc...

## History

Gpxgo is based on:

 * https://github.com/tkrajina/gpxpy (python gpx library)
 * https://github.com/ptrv/go-gpx (an earlier port of gpxpy)

# License

gpxgo is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)
