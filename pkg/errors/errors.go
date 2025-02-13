package errors

import (
	"errors"
	"strconv"
)

var (
	ErrNilDecoder            = errors.New("nil decoder")
	ErrInvalidDomain         = errors.New("invalid domain")
	ErrEmptyRegisteredDomain = errors.New("empty registered domain")
	ErrEmptySubdomain        = errors.New("empty subdomain")
	ErrAuthenticationError   = errors.New("authentication error")
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
