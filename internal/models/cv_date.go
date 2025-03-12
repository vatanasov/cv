package models

import (
	"encoding/xml"
	"errors"
	"time"
)

type CVDate time.Time

func (x *CVDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var rawDescription string
	err := d.DecodeElement(&rawDescription, &start)
	if err != nil {
		return err
	}

	layouts := []string{"2006", "2006-01", "2006-01-02"}
	for _, layout := range layouts {
		t, err := time.Parse(layout, rawDescription)
		if err == nil {
			*x = CVDate(t)
			return nil
		}
	}

	return errors.New("invalid Employment Date")
}

func (x CVDate) String() string {
	return time.Time(x).Format("2006")
}
