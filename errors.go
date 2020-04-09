package fastjsonrpc

type errorCode int

const (
	errorCodeMethodNotFound errorCode = -32601
	ErrorCodeInvalidParams  errorCode = -32602
	ErrorCodeInternalError  errorCode = -32603
)

var (
	renderedParseError     = []byte(`{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}`)
	renderedInvalidRequest = []byte(`{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null}`)
)
