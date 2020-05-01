package fastjsonrpc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/valyala/bytebufferpool"

	"github.com/valyala/fastjson"
)

type RequestHandler func(ctx *RequestCtx)

var _ context.Context = &RequestCtx{}

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

func (ctx *RequestCtx) ParamsUnmarshal(v interface{}) error {
	if json.Unmarshal(ctx.ParamsBytes(), v) != nil {
		return ErrInvalidParams()
	}

	return nil
}

func (ctx *RequestCtx) Deadline() (deadline time.Time, ok bool) {
	return ctx.fasthttpCtx.Deadline()
}

func (ctx *RequestCtx) Done() <-chan struct{} {
	return ctx.fasthttpCtx.Done()
}

func (ctx *RequestCtx) Err() error {
	return ctx.fasthttpCtx.Err()
}

func (ctx *RequestCtx) Value(key interface{}) interface{} {
	return ctx.fasthttpCtx.Value(key)
}
