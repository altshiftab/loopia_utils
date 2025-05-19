package errors

import (
	"errors"
	"strconv"
)

var (
	ErrEmptyDomain           = errors.New("empty domain")
	ErrEmptyRegisteredDomain = errors.New("empty registered domain")
	ErrEmptySubdomain        = errors.New("empty subdomain")
	ErrAuthenticationError   = errors.New("authentication error")
	ErrUnexpectedStatus      = errors.New("unexpected status")
)

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
