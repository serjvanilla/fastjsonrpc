package fastjsonrpc

import (
	"encoding/json"

	"github.com/valyala/fastjson"
)

//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=.

func (r *Request) writeByte(buf []byte) {
	_, _ = r.response.Write(buf)
}

// Errors writes JSON-RPC response with error.
func (r *Request) Error(code errorCode, message string, data interface{}) {
	if len(r.id) == 0 {
		return
	}

	if data == nil {
		writeresponseWithError(r.response, r.id, code, message, nil)
		return
	}

	switch v := data.(type) {
	case *fastjson.Value:
		buf := bufferpool.Get()
		writeresponseWithError(r.response, r.id, code, message, v.MarshalTo(buf.B))
		bufferpool.Put(buf)
	case []byte:
		writeresponseWithError(r.response, r.id, code, message, v)
	default:
		out, _ := json.Marshal(data)
		writeresponseWithError(r.response, r.id, code, message, out)
	}
}

// Result writes JSON-RPC response with result.
//
// result may be *fastjson.Value, []byte, or interface{} (slower).
func (r *Request) Result(result interface{}) {
	if len(r.id) == 0 {
		return
	}

	switch v := result.(type) {
	case *fastjson.Value:
		buf := bufferpool.Get()
		writeresponseWithResult(r.response, r.id, v.MarshalTo(buf.B))
		bufferpool.Put(buf)
	case []byte:
		writeresponseWithResult(r.response, r.id, v)
	default:
		out, _ := json.Marshal(result)
		writeresponseWithResult(r.response, r.id, out)
	}
}
