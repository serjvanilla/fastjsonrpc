package fastjsonrpc

import (
	"context"
	"encoding/json"

	"github.com/valyala/fasthttp"

	"github.com/valyala/bytebufferpool"

	"github.com/valyala/fastjson"
)

// RequestHandler must process incoming requests.
type RequestHandler func(ctx *RequestCtx)

// RequestCtx contains incoming request and manages outgoing response.
type RequestCtx struct {
	fasthttpCtx *fasthttp.RequestCtx

	arena fastjson.Arena

	id     []byte
	method []byte

	params *fastjson.Value

	paramsBytes *bytebufferpool.ByteBuffer
	response    *bytebufferpool.ByteBuffer

	bytebufferpool *bytebufferpool.Pool
}

// Arena returns fastjson.Arena for current request.
//
// RequestHandler should avoid holding references to Arena and/or constructed Values after the return.
func (ctx *RequestCtx) Arena() *fastjson.Arena {
	return &ctx.arena
}

// ID returns "id" field of JSON-RPC 2.0 request.
func (ctx *RequestCtx) ID() []byte {
	return ctx.id
}

// Method returns matched method.
func (ctx *RequestCtx) Method() []byte {
	return ctx.method
}

// Params returns request parameters already unmarshalled with valyala/fastjson.
func (ctx *RequestCtx) Params() *fastjson.Value {
	return ctx.params
}

// ParamsBytes returns raw bytes of request's "params" field.
func (ctx *RequestCtx) ParamsBytes() []byte {
	return ctx.paramsBytes.B
}

// ParamsUnmarshal parses request param and stores the result in the value pointed to by v.
func (ctx *RequestCtx) ParamsUnmarshal(v interface{}) *Error {
	if json.Unmarshal(ctx.ParamsBytes(), v) != nil {
		return ErrInvalidParams()
	}

	return nil
}

func (ctx *RequestCtx) Context() context.Context {
	return ctx.fasthttpCtx
}
