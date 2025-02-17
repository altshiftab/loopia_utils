package loopia_utils

import (
	"encoding/xml"
	loopiaUtilsTypes "github.com/altshiftab/loopia_utils/pkg/types"
	"strconv"
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
	XMLName xml.Name                   `xml:"methodResponse"`
	Params  []*loopiaUtilsTypes.Record `xml:"params>param>value>array>data>value>struct"`
}
