package fastjsonrpc

const (
	renderedParseError     = `{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}`
	renderedInvalidRequest = `{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null}`
)

type ErrorCode int

const (
	errorCodeParseError    ErrorCode = -32700
	errorMessageParseError string    = `Parse error`

	errorCodeInvalidRequest    ErrorCode = -32600
	errorMessageInvalidRequest string    = `Invalid Request`

	errorCodeMethodNotFound    ErrorCode = -32601
	errorMessageMethodNotFound string    = `Method not found`

	errorCodeInvalidParams    ErrorCode = -32602
	errorMessageInvalidParams string    = `Invalid params`

	errorCodeInternalError    ErrorCode = -32603
	errorMessageInternalError string    = `Internal error`

	errorMessageServerError string = `Server error`
)

type Error struct {
	Code    ErrorCode
	Message string
	Data    interface{}
}

func (e *Error) WithMessage(message string) *Error {
	e.Message = message

	return e
}

func (e *Error) WithData(data interface{}) *Error {
	e.Data = data

	return e
}

func ErrParseError() *Error {
	return &Error{
		Code:    errorCodeParseError,
		Message: errorMessageParseError,
	}
}

func ErrInvalidRequest() *Error {
	return &Error{
		Code:    errorCodeInvalidRequest,
		Message: errorMessageInvalidRequest,
	}
}

func ErrMethodNotFound() *Error {
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
