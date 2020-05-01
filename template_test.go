package fastjsonrpc

import (
	"testing"
)

func TestResponseWithError(t *testing.T) {
	result := responseWithError(
		[]byte(`"null"`),
		errorCodeInternalError,
		errorMessageInternalError,
		[]byte(`{"context":42}`),
	)
	if result != `{"jsonrpc":"2.0","error":{"code":-32603,"message":"Internal error","data":{"context":42}},"id":"null"}` {
		t.Fatalf("unexpected result: `%s`", result)
	}
}

func TestResponseWithResult(t *testing.T) {
	result := responseWithResult(
		[]byte(`"null"`),
		[]byte(`true`),
	)
	if result != `{"jsonrpc":"2.0","result":true,"id":"null"}` {
		t.Fatalf("unexpected result: `%s`", result)
	}
}
