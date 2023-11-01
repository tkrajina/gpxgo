package gpx

import (
	"encoding/xml"
)

// formattedFloat forces XML attributes to be marshalled as a fixed point decimal with 10 decimal places.
type formattedFloat float64 // TODO: Delete

func (f formattedFloat) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  xml.Name{Local: name.Local},
		Value: formatNumber(float64(f)),
	}, nil
}
