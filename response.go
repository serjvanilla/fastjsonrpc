package fastjsonrpc

import (
	"encoding/json"

	"github.com/valyala/fastjson"
)

//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=.

func (r *Request) writeString(s string) {
	_, _ = r.response.WriteString(s)
}

// Error writes JSON-RPC response with error.
func (r *Request) Error(err *Error) {
	if len(r.id) == 0 {
		return
	}

	r.response.Reset()

	if err.Data == nil {
		writeresponseWithError(r.response, r.id, err.Code, err.Message, nil)
		return
	}

	switch v := err.Data.(type) {
	case *fastjson.Value:
		buf := bufferpool.Get()
		writeresponseWithError(r.response, r.id, err.Code, err.Message, v.MarshalTo(buf.B))
		bufferpool.Put(buf)
	case []byte:
		writeresponseWithError(r.response, r.id, err.Code, err.Message, v)
	default:
		out, _ := json.Marshal(err.Data)
		writeresponseWithError(r.response, r.id, err.Code, err.Message, out)
	}
}

// Result writes JSON-RPC response with result.
//
// result may be *fastjson.Value, []byte, or interface{} (slower).
func (r *Request) Result(result interface{}) {
	if len(r.id) == 0 {
		return
	}

	r.response.Reset()

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
