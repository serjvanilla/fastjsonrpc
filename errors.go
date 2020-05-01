package fastjsonrpc

import (
	"fmt"
)

const (
	renderedParseError     = `{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}`
	renderedInvalidRequest = `{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null}`
)

type ErrorCode int

const (
	errorCodeMethodNotFound    ErrorCode = -32601
	errorMessageMethodNotFound string    = `Method not found`

	errorCodeInvalidParams    ErrorCode = -32602
	errorMessageInvalidParams string    = `Invalid params`

	errorCodeInternalError    ErrorCode = -32603
	errorMessageInternalError string    = `Internal error`

	errorCodeServerError    ErrorCode = -32000
	errorMessageServerError string    = `Server error`
)

var _ error = &Error{}

type Error struct {
	Code    ErrorCode
	Message string
	Data    interface{}
}

func (e *Error) Error() string {
	if e.Data == nil {
		return fmt.Sprintf("json-rpc error: [%d] %s", e.Code, e.Message)
	}

	return fmt.Sprintf("json-rpc error: [%d] %s (%+v)", e.Code, e.Message, e.Data)
}

func (e *Error) WithData(data interface{}) *Error {
	e.Data = data

	return e
}

func errMethodNotFound() *Error {
	return &Error{
		Code:    errorCodeMethodNotFound,
		Message: errorMessageMethodNotFound,
	}
}

func ErrInvalidParams() *Error {
	return &Error{
		Code:    errorCodeInvalidParams,
		Message: errorMessageInvalidParams,
	}
}

func ErrInternalError() *Error {
	return &Error{
		Code:    errorCodeInternalError,
		Message: errorMessageInternalError,
	}
}

func ErrServerError(code ErrorCode) *Error {
	return &Error{
		Code:    code,
		Message: errorMessageServerError,
	}
}
