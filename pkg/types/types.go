package types

import (
	"encoding/xml"
	"strings"
)

type Record struct {
	Type     string
	Ttl      int
	Priority int
	Rdata    string
	Id       int
}

func (record *Record) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var name string
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		switch tt := t.(type) {
		case xml.StartElement:
			switch tt.Name.Local {
			case "name": // The name of the record object: <name>
				var s string
				if err = d.DecodeElement(&s, &start); err != nil {
					return err
				}

				name = strings.TrimSpace(s)

			case "string": // A string value of the record object: <value><string>
				if err = record.decodeValueString(name, d, start); err != nil {
					return err
				}

			case "int": // An int value of the record object: <value><int>
				if err = record.decodeValueInt(name, d, start); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}
	}
}

func (record *Record) decodeValueString(name string, d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	s = strings.TrimSpace(s)
	switch name {
	case "type":
		record.Type = s
	case "rdata":
		record.Rdata = s
	}

	return nil
}

func (record *Record) decodeValueInt(name string, d *xml.Decoder, start xml.StartElement) error {
	var i int
	if err := d.DecodeElement(&i, &start); err != nil {
		return err
	}

	switch name {
	case "record_id":
		record.Id = i
	case "ttl":
		record.Ttl = i
	case "priority":
		record.Priority = i
	}

	return nil
}
