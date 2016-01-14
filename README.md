# Go GPX library

gpgpx is a golang library for parsing and manipulating GPX files. GPX (GPS eXchange Format) is a XML based file format for GPS track logs. 

## A simple example:

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

## History

Gpxgo is based on 

 * https://github.com/tkrajina/gpxpy (python gpx library)
 * https://github.com/ptrv/go-gpx (an earlier port of gpxgo)

# License

gpxgo is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)
