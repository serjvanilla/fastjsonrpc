package fastjsonrpc

import (
	"encoding/json"

	"github.com/valyala/fastjson"
)

//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=.

func (ctx *RequestCtx) writeString(s string) {
	_, _ = ctx.response.WriteString(s)
}

// SetError writes JSON-RPC response with error.
func (ctx *RequestCtx) SetError(err *Error) {
	if len(ctx.id) == 0 {
		return
	}

	ctx.response.Reset()

	if err.Data == nil {
		writeresponseWithError(ctx.response, ctx.id, err.Code, err.Message, nil)
		return
	}

	switch v := err.Data.(type) {
	case *fastjson.Value:
		buf := ctx.bytebufferpool.Get()
		writeresponseWithError(ctx.response, ctx.id, err.Code, err.Message, v.MarshalTo(buf.B))
		ctx.bytebufferpool.Put(buf)
	case []byte:
		writeresponseWithError(ctx.response, ctx.id, err.Code, err.Message, v)
	default:
		buf := ctx.bytebufferpool.Get()
		_ = json.NewEncoder(buf).Encode(err.Data)
		writeresponseWithError(ctx.response, ctx.id, err.Code, err.Message, buf.B)
		ctx.bytebufferpool.Put(buf)
	}
}

// SetResult writes JSON-RPC response with result.
//
// result may be *fastjson.Value, []byte, or interface{} (slower).
func (ctx *RequestCtx) SetResult(result interface{}) {
	if len(ctx.id) == 0 {
		return
	}

	ctx.response.Reset()

	switch v := result.(type) {
	case *fastjson.Value:
		buf := ctx.bytebufferpool.Get()
		writeresponseWithResult(ctx.response, ctx.id, v.MarshalTo(buf.B))
		ctx.bytebufferpool.Put(buf)
	case []byte:
		writeresponseWithResult(ctx.response, ctx.id, v)
	default:
		buf := ctx.bytebufferpool.Get()
		_ = json.NewEncoder(buf).Encode(result)
		writeresponseWithResult(ctx.response, ctx.id, buf.B)
		ctx.bytebufferpool.Put(buf)
	}
}
