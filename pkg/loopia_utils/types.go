package loopia_utils

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// types for XML-RPC method calls and parameters

type param interface {
	param()
}

type paramString struct {
	XMLName xml.Name `xml:"param"`
	Value   string   `xml:"value>string"`
}

func (p paramString) param() {}

type paramInt struct {
	XMLName xml.Name `xml:"param"`
	Value   int      `xml:"value>int"`
}

func (p paramInt) param() {}

type paramStruct struct {
	XMLName       xml.Name       `xml:"param"`
	StructMembers []structMember `xml:"value>struct>member"`
}

func (p paramStruct) param() {}

type structMember interface {
	structMember()
}

type structMemberString struct {
	Name  string `xml:"name"`
	Value string `xml:"value>string"`
}

func (m structMemberString) structMember() {}

type structMemberInt struct {
	Name  string `xml:"name"`
	Value int    `xml:"value>int"`
}

func (m structMemberInt) structMember() {}

type methodCall struct {
	XMLName    xml.Name `xml:"methodCall"`
	MethodName string   `xml:"methodName"`
	Params     []param  `xml:"params>param"`
}

// types for XML-RPC responses

type response interface {
	faultCode() int
	faultString() string
}

type responseString struct {
	responseFault
	Value string `xml:"params>param>value>string"`
}

type responseFault struct {
	FaultCode   int    `xml:"fault>value>struct>member>value>int"`
	FaultString string `xml:"fault>value>struct>member>value>string"`
}

func (r responseFault) faultCode() int      { return r.FaultCode }
func (r responseFault) faultString() string { return r.FaultString }

type RpcError struct {
	Code    int
	Message string
}

func (rpcError *RpcError) Error() string {
	return rpcError.Message
}

func (rpcError *RpcError) GetCode() string {
	if rpcError.Code == 0 {
		return ""
	}
	return strconv.Itoa(rpcError.Code)
}

type recordObjectsResponse struct {
	responseFault
	XMLName xml.Name  `xml:"methodResponse"`
	Params  []*Record `xml:"params>param>value>array>data>value>struct"`
}

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
