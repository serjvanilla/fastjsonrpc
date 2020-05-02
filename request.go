package fastjsonrpc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/valyala/bytebufferpool"

	"github.com/valyala/fastjson"
)

// RequestHandler must process incoming requests.
type RequestHandler func(ctx *RequestCtx)

var _ context.Context = &RequestCtx{}

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

// Deadline returns underlying *fasthttp.RequestCtx Deadline.
//
// This method is present to make RequestCtx implement the context interface.
func (ctx *RequestCtx) Deadline() (deadline time.Time, ok bool) {
	return ctx.fasthttpCtx.Deadline()
}

// Done returns underlying *fasthttp.RequestCtx Done channel.
//
// This method is present to make RequestCtx implement the context interface.
func (ctx *RequestCtx) Done() <-chan struct{} {
	return ctx.fasthttpCtx.Done()
}

// Err returns underlying *fasthttp.RequestCtx Err.
//
// This method is present to make RequestCtx implement the context interface.
func (ctx *RequestCtx) Err() error {
	return ctx.fasthttpCtx.Err()
}

// Value returns underlying *fasthttp.RequestCtx Value.
//
// This method is present to make RequestCtx implement the context interface.
func (ctx *RequestCtx) Value(key interface{}) interface{} {
	return ctx.fasthttpCtx.Value(key)
}
