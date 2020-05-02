package fastjsonrpc

import (
	"fmt"
)

const (
	renderedParseError     = `{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}`
	renderedInvalidRequest = `{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null}`
)

// ErrorCode is JSON-RPC 2.0 spec defined error code.
//
// For user defined errors it should be in range from -32099 to -32000.
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

// Error is wrapper for JSON-RPC 2.0 Error Object.
type Error struct {
	Code    ErrorCode
	Message string
	Data    interface{}
}

// Error implements standard error interface.
func (e *Error) Error() string {
	if e.Data == nil {
		return fmt.Sprintf("json-rpc error: [%d] %s", e.Code, e.Message)
	}

	return fmt.Sprintf("json-rpc error: [%d] %s (%+v)", e.Code, e.Message, e.Data)
}

// WithData sets error's data value.
//
// Useful with ErrServerError.
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

// ErrInvalidParams returns pre-built JSON-RPC error with code -32602 and message "Invalid params".
func ErrInvalidParams() *Error {
	return &Error{
		Code:    errorCodeInvalidParams,
		Message: errorMessageInvalidParams,
	}
}

// ErrInternalError returns pre-built JSON-RPC error with code -32603 and message "Internal error".
func ErrInternalError() *Error {
	return &Error{
		Code:    errorCodeInternalError,
		Message: errorMessageInternalError,
	}
}

// ErrServerError returns pre-built JSON-RPC error with provided code and message "Server error".
func ErrServerError(code ErrorCode) *Error {
	return &Error{
		Code:    code,
		Message: errorMessageServerError,
	}
}
