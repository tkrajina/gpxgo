package gpx

import (
	"encoding/xml"
	"fmt"
)

// fixedPointFloat forces XML attributes to be marshalled as a fixed point decimal with 10 decimal places.
type fixedPointFloat float64

func (f fixedPointFloat) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  xml.Name{Local: name.Local},
		Value: fmt.Sprintf("%.10f", f),
	}, nil
}
