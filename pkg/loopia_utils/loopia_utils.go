package loopia_utils

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	motmedelHttpErrors "github.com/Motmedel/utils_go/pkg/http/errors"
	motmedelHttpTypes "github.com/Motmedel/utils_go/pkg/http/types"
	motmedelHttpUtils "github.com/Motmedel/utils_go/pkg/http/utils"
	"github.com/Motmedel/utils_go/pkg/net/domain_breakdown"
	motmedelNetErrors "github.com/Motmedel/utils_go/pkg/net/errors"
	loopiaUtilsErrors "github.com/altshiftab/loopia_utils/pkg/errors"
	loopiaUtilsTypes "github.com/altshiftab/loopia_utils/pkg/types"
	"net/http"
	"strings"
)

const BaseUrlString = "https://api.loopia.se/RPCSERV"
const DefaultTtlValue = 3600

type Client struct {
	*http.Client
	ApiUser     string
	ApiPassword string
}

func parseDomain(domain string) (string, string, error) {
	domainBreakdown := domain_breakdown.GetDomainBreakdown(domain)
	if domainBreakdown == nil {
		return "", "", motmedelErrors.NewWithTrace(motmedelNetErrors.ErrNilDomainBreakdown)
	}

	subdomain := domainBreakdown.Subdomain
	if subdomain == "" {
		subdomain = "@"
	}

	return domainBreakdown.RegisteredDomain, subdomain, nil
}

func (client *Client) AddRecord(ctx context.Context, record *loopiaUtilsTypes.Record, domain string) error {
	if record == nil {
		return nil
	}

	if domain == "" {
		return motmedelErrors.NewWithTrace(loopiaUtilsErrors.ErrEmptyDomain)
	}

	registeredDomain, subdomain, err := parseDomain(domain)
	if err != nil {
		return fmt.Errorf("parse domain: %w", err)
	}
	if registeredDomain == "" {
		return motmedelErrors.NewWithTrace(loopiaUtilsErrors.ErrEmptyRegisteredDomain)
	}
	if subdomain == "" {
		return motmedelErrors.NewWithTrace(loopiaUtilsErrors.ErrEmptySubdomain)
	}

	ttl := record.Ttl
	if ttl == 0 {
		ttl = DefaultTtlValue
	}

	call := &methodCall{
		MethodName: "addZoneRecord",
		Params: []param{
			paramString{Value: client.ApiUser},
			paramString{Value: client.ApiPassword},
			paramString{Value: registeredDomain},
			paramString{Value: subdomain},
			paramStruct{
				StructMembers: []structMember{
					structMemberString{Name: "type", Value: record.Type},
					structMemberInt{Name: "ttl", Value: ttl},
					structMemberInt{Name: "priority", Value: record.Priority},
					structMemberString{Name: "rdata", Value: record.Rdata},
					structMemberInt{Name: "record_id", Value: record.Id},
				},
			},
		},
	}
	resp := &responseString{}

	if err := client.rpcCall(ctx, call, resp); err != nil {
		return motmedelErrors.New(fmt.Errorf("rpc call: %w", err), registeredDomain, subdomain)
	}

	responseValue := resp.Value
	if err := checkResponse(responseValue); err != nil {
		return motmedelErrors.New(fmt.Errorf("check response: %w", err), responseValue)
	}

	return nil
}

func (client *Client) RemoveRecord(ctx context.Context, domain string, recordId int) error {
	if domain == "" {
		return nil
	}

	registeredDomain, subdomain, err := parseDomain(domain)
	if err != nil {
		return fmt.Errorf("parse domain: %w", err)
	}
	if registeredDomain == "" {
		return motmedelErrors.NewWithTrace(loopiaUtilsErrors.ErrEmptyRegisteredDomain)
	}
	if subdomain == "" {
		return motmedelErrors.NewWithTrace(loopiaUtilsErrors.ErrEmptySubdomain)
	}

	call := &methodCall{
		MethodName: "removeZoneRecord",
		Params: []param{
			paramString{Value: client.ApiUser},
			paramString{Value: client.ApiPassword},
			paramString{Value: registeredDomain},
			paramString{Value: subdomain},
			paramInt{Value: recordId},
		},
	}
	resp := &responseString{}

	if err := client.rpcCall(ctx, call, resp); err != nil {
		return motmedelErrors.New(fmt.Errorf("rpc call: %w", err), registeredDomain, subdomain)
	}

	responseValue := resp.Value
	if err := checkResponse(responseValue); err != nil {
		return motmedelErrors.New(fmt.Errorf("check response: %w", err), responseValue)
	}

	return nil
}

func (client *Client) GetRecords(ctx context.Context, domain string) ([]*loopiaUtilsTypes.Record, error) {
	if domain == "" {
		return nil, nil
	}

	registeredDomain, subdomain, err := parseDomain(domain)
	if err != nil {
		return nil, fmt.Errorf("parse domain: %w", err)
	}
	if registeredDomain == "" {
		return nil, motmedelErrors.NewWithTrace(loopiaUtilsErrors.ErrEmptyRegisteredDomain)
	}
	if subdomain == "" {
		return nil, motmedelErrors.NewWithTrace(loopiaUtilsErrors.ErrEmptySubdomain)
	}

	call := &methodCall{
		MethodName: "getZoneRecords",
		Params: []param{
			paramString{Value: client.ApiUser},
			paramString{Value: client.ApiPassword},
			paramString{Value: registeredDomain},
			paramString{Value: subdomain},
		},
	}
	resp := &recordObjectsResponse{}

	if err := client.rpcCall(ctx, call, resp); err != nil {
		return nil, motmedelErrors.New(fmt.Errorf("rpc call: %w", err), registeredDomain, subdomain)
	}

	return resp.Params, nil
}

func (client *Client) RemoveSubdomain(ctx context.Context, domain string) error {
	if domain == "" {
		return nil
	}

	registeredDomain, subdomain, err := parseDomain(domain)
	if err != nil {
		return fmt.Errorf("parse domain: %w", err)
	}
	if registeredDomain == "" {
		return motmedelErrors.NewWithTrace(loopiaUtilsErrors.ErrEmptyRegisteredDomain)
	}
	if subdomain == "" {
		return motmedelErrors.NewWithTrace(loopiaUtilsErrors.ErrEmptySubdomain)
	}

	call := &methodCall{
		MethodName: "removeSubdomain",
		Params: []param{
			paramString{Value: client.ApiUser},
			paramString{Value: client.ApiPassword},
			paramString{Value: registeredDomain},
			paramString{Value: subdomain},
		},
	}
	resp := &responseString{}

	if err := client.rpcCall(ctx, call, resp); err != nil {
		return motmedelErrors.New(fmt.Errorf("rpc call: %w", err), registeredDomain, subdomain)
	}

	responseValue := resp.Value
	if err := checkResponse(responseValue); err != nil {
		return motmedelErrors.New(fmt.Errorf("check response: %w", err), responseValue)
	}

	return nil
}

func (client *Client) rpcCall(
	ctx context.Context,
	call *methodCall,
	resultTarget response,
) error {
	if call == nil {
		return nil
	}

	requestBodyBuffer := new(bytes.Buffer)
	requestBodyBuffer.WriteString(xml.Header)

	encoder := xml.NewEncoder(requestBodyBuffer)
	encoder.Indent("", "  ")

	if err := encoder.Encode(call); err != nil {
		return motmedelErrors.NewWithTrace(fmt.Errorf("xml encoder encode: %w", err))
	}

	requestBody := requestBodyBuffer.Bytes()
	_, responseBody, err := motmedelHttpUtils.Fetch(
		ctx,
		BaseUrlString,
		client.Client,
		&motmedelHttpTypes.FetchOptions{
			Method:  http.MethodPost,
			Headers: map[string]string{"Content-Type": "text/xml"},
			Body:    requestBody,
		},
	)
	if err != nil {
		return motmedelErrors.New(fmt.Errorf("fetch: %w", err), requestBody)
	}
	if len(responseBody) == 0 {
		return motmedelErrors.NewWithTrace(motmedelHttpErrors.ErrEmptyResponseBody)
	}

	if err := xml.Unmarshal(responseBody, resultTarget); err != nil {
		return motmedelErrors.NewWithTrace(
			fmt.Errorf("xml unmarshal (response body): %w", err),
			responseBody,
		)
	}

	if faultCode := resultTarget.faultCode(); faultCode != 0 {
		return motmedelErrors.NewWithTrace(
			&RpcError{
				Code:    faultCode,
				Message: strings.TrimSpace(resultTarget.faultString()),
			},
		)
	}

	return nil
}

func checkResponse(value string) error {
	switch trimmedValue := strings.TrimSpace(value); trimmedValue {
	case "OK":
		return nil
	case "AUTH_ERROR":
		return motmedelErrors.NewWithTrace(loopiaUtilsErrors.ErrAuthenticationError)
	default:
		return motmedelErrors.NewWithTrace(loopiaUtilsErrors.ErrUnexpectedStatus, trimmedValue)
	}
}
