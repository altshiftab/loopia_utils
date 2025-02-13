package loopia_utils

import (
	"bytes"
	"encoding/xml"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	motmedelHttpErrors "github.com/Motmedel/utils_go/pkg/http/errors"
	motmedelHttpTypes "github.com/Motmedel/utils_go/pkg/http/types"
	motmedelHttpUtils "github.com/Motmedel/utils_go/pkg/http/utils"
	"github.com/Motmedel/utils_go/pkg/net/domain_breakdown"
	loopiaUtilsErrors "github.com/altshiftab/loopia_utils/pkg/errors"
	"net/http"
	"strings"
)

const BaseUrlString = "https://api.loopia.se/RPCSERV"

type Client struct {
	ApiUser     string
	ApiPassword string
	HttpClient  motmedelHttpUtils.HttpClient
}

func parseDomain(domain string) (string, string, error) {
	domainBreakdown := domain_breakdown.GetDomainBreakdown(domain)
	if domainBreakdown == nil {
		return "", "", &motmedelErrors.InputError{
			Message: "The domain is invalid.",
			Input:   domain,
		}
	}

	subdomain := domainBreakdown.Subdomain
	if subdomain == "" {
		subdomain = "@"
	}

	return domainBreakdown.RegisteredDomain, subdomain, nil
}

func (c *Client) AddTxtRecord(domain string, ttl int, value string) (*motmedelHttpTypes.HttpContext, error) {
	if domain == "" {
		return nil, nil
	}

	registeredDomain, subdomain, err := parseDomain(domain)
	if err != nil {
		return nil, &motmedelErrors.InputError{
			Message: "An error occurred when parsing the domain.",
			Cause:   err,
			Input:   domain,
		}
	}
	if registeredDomain == "" {
		return nil, loopiaUtilsErrors.ErrEmptyRegisteredDomain
	}
	if subdomain == "" {
		return nil, loopiaUtilsErrors.ErrEmptySubdomain
	}

	call := &methodCall{
		MethodName: "addZoneRecord",
		Params: []param{
			paramString{Value: c.ApiUser},
			paramString{Value: c.ApiPassword},
			paramString{Value: registeredDomain},
			paramString{Value: subdomain},
			paramStruct{
				StructMembers: []structMember{
					structMemberString{Name: "type", Value: "TXT"},
					structMemberInt{Name: "ttl", Value: ttl},
					structMemberInt{Name: "priority", Value: 0},
					structMemberString{Name: "rdata", Value: value},
					structMemberInt{Name: "record_id", Value: 0},
				},
			},
		},
	}
	resp := &responseString{}

	httpContext, err := c.rpcCall(call, resp)
	if err != nil {
		return httpContext, &motmedelErrors.CauseError{
			Message: "An error occurred when making the call.",
			Cause:   err,
		}
	}

	responseValue := resp.Value
	if err := checkResponse(responseValue); err != nil {
		return httpContext, &motmedelErrors.InputError{
			Message: "An error response was received.",
			Cause:   err,
			Input:   responseValue,
		}
	}

	return httpContext, nil
}

func (c *Client) RemoveTxtRecord(domain string, recordId int) (*motmedelHttpTypes.HttpContext, error) {
	if domain == "" {
		return nil, nil
	}

	registeredDomain, subdomain, err := parseDomain(domain)
	if err != nil {
		return nil, &motmedelErrors.InputError{
			Message: "An error occurred when parsing the domain.",
			Cause:   err,
			Input:   domain,
		}
	}
	if registeredDomain == "" {
		return nil, loopiaUtilsErrors.ErrEmptyRegisteredDomain
	}
	if subdomain == "" {
		return nil, loopiaUtilsErrors.ErrEmptySubdomain
	}

	call := &methodCall{
		MethodName: "removeZoneRecord",
		Params: []param{
			paramString{Value: c.ApiUser},
			paramString{Value: c.ApiPassword},
			paramString{Value: registeredDomain},
			paramString{Value: subdomain},
			paramInt{Value: recordId},
		},
	}
	resp := &responseString{}

	httpContext, err := c.rpcCall(call, resp)
	if err != nil {
		return httpContext, &motmedelErrors.CauseError{
			Message: "An error occurred when making the call.",
			Cause:   err,
		}
	}

	responseValue := resp.Value
	if err := checkResponse(responseValue); err != nil {
		return httpContext, &motmedelErrors.InputError{
			Message: "An error response was received.",
			Cause:   err,
			Input:   responseValue,
		}
	}

	return httpContext, nil
}

func (c *Client) GetTxtRecords(domain string) ([]*Record, *motmedelHttpTypes.HttpContext, error) {
	if domain == "" {
		return nil, nil, nil
	}

	registeredDomain, subdomain, err := parseDomain(domain)
	if err != nil {
		return nil, nil, &motmedelErrors.InputError{
			Message: "An error occurred when parsing the domain.",
			Cause:   err,
			Input:   domain,
		}
	}
	if registeredDomain == "" {
		return nil, nil, loopiaUtilsErrors.ErrEmptyRegisteredDomain
	}
	if subdomain == "" {
		return nil, nil, loopiaUtilsErrors.ErrEmptySubdomain
	}

	call := &methodCall{
		MethodName: "getZoneRecords",
		Params: []param{
			paramString{Value: c.ApiUser},
			paramString{Value: c.ApiPassword},
			paramString{Value: registeredDomain},
			paramString{Value: subdomain},
		},
	}
	resp := &recordObjectsResponse{}

	httpContext, err := c.rpcCall(call, resp)
	if err != nil {
		return nil, httpContext, &motmedelErrors.CauseError{
			Message: "An error occurred when making the call.",
			Cause:   err,
		}
	}

	return resp.Params, httpContext, nil
}

// RemoveSubdomain remove a sub-domain.
func (c *Client) RemoveSubdomain(domain string) (*motmedelHttpTypes.HttpContext, error) {
	if domain == "" {
		return nil, nil
	}

	registeredDomain, subdomain, err := parseDomain(domain)
	if err != nil {
		return nil, &motmedelErrors.InputError{
			Message: "An error occurred when parsing the domain.",
			Cause:   err,
			Input:   domain,
		}
	}
	if registeredDomain == "" {
		return nil, loopiaUtilsErrors.ErrEmptyRegisteredDomain
	}
	if subdomain == "" {
		return nil, loopiaUtilsErrors.ErrEmptySubdomain
	}

	call := &methodCall{
		MethodName: "removeSubdomain",
		Params: []param{
			paramString{Value: c.ApiUser},
			paramString{Value: c.ApiPassword},
			paramString{Value: domain},
			paramString{Value: subdomain},
		},
	}
	resp := &responseString{}

	httpContext, err := c.rpcCall(call, resp)
	if err != nil {
		return httpContext, &motmedelErrors.CauseError{
			Message: "An error occurred when making the call.",
			Cause:   err,
		}
	}

	responseValue := resp.Value
	if err := checkResponse(responseValue); err != nil {
		return httpContext, &motmedelErrors.InputError{
			Message: "An error response was received.",
			Cause:   err,
			Input:   responseValue,
		}
	}

	return httpContext, nil
}

func (c *Client) rpcCall(call *methodCall, resultTarget response) (*motmedelHttpTypes.HttpContext, error) {
	if call == nil {
		return nil, nil
	}

	requestBodyBuffer := new(bytes.Buffer)
	requestBodyBuffer.WriteString(xml.Header)

	encoder := xml.NewEncoder(requestBodyBuffer)
	encoder.Indent("", "  ")

	if err := encoder.Encode(call); err != nil {
		return nil, &motmedelErrors.InputError{
			Message: "An error occurred when encoding the call.",
			Cause:   err,
			Input:   call,
		}
	}

	requestMethod := http.MethodPost
	requestUrl := BaseUrlString
	requestBody := requestBodyBuffer.Bytes()

	httpContext, err := motmedelHttpUtils.SendRequest(
		c.HttpClient,
		requestMethod,
		requestUrl,
		requestBody,
		func(request *http.Request) error {
			if request == nil {
				return motmedelHttpErrors.ErrNilHttpRequest
			}

			requestHeader := request.Header
			if requestHeader == nil {
				return motmedelHttpErrors.ErrNilHttpRequestHeader
			}

			requestHeader.Set("Content-Type", "text/xml")

			return nil
		},
	)
	if err != nil {
		return httpContext, &motmedelErrors.InputError{
			Message: "An error occurred when sending the request.",
			Cause:   err,
			Input:   []any{requestMethod, requestUrl, requestBody},
		}
	}

	if httpContext == nil {
		return nil, motmedelHttpErrors.ErrNilHttpContext
	}

	responseBody := httpContext.ResponseBody
	if len(responseBody) == 0 {
		return httpContext, motmedelHttpErrors.ErrEmptyResponseBody
	}

	if err := xml.Unmarshal(responseBody, resultTarget); err != nil {
		return httpContext, &motmedelErrors.InputError{
			Message: "An error occurred when unmarshalling the response body.",
			Cause:   err,
			Input:   responseBody,
		}
	}

	if resultTarget.faultCode() != 0 {
		return httpContext, &RpcError{
			Code:    resultTarget.faultCode(),
			Message: strings.TrimSpace(resultTarget.faultString()),
		}
	}

	return httpContext, nil
}

func checkResponse(value string) error {
	switch v := strings.TrimSpace(value); v {
	case "OK":
		return nil
	case "AUTH_ERROR":
		return loopiaUtilsErrors.ErrAuthenticationError
	default:
		return &motmedelErrors.InputError{
			Message: "An unknown status value was encountered.",
			Input:   value,
		}
	}
}
