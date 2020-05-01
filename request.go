package fastjsonrpc

import (
	"context"

	"github.com/valyala/bytebufferpool"

	"github.com/valyala/fastjson"
)

type RequestHandler func(ctx *Request)

type Request struct {
	ctx context.Context

	arena fastjson.Arena

	id     []byte
	method []byte

	params *fastjson.Value

	paramsBytes *bytebufferpool.ByteBuffer
	response    *bytebufferpool.ByteBuffer

	bytebufferpool *bytebufferpool.Pool
}

// Context returns request's underlying context.
func (r *Request) Context() context.Context {
	return r.ctx
}

func (r *Request) Arena() *fastjson.Arena {
	return &r.arena
}

// ID returns "id" field of JSON-RPC 2.0 request.
func (r *Request) ID() []byte {
	return r.id
}

// Method returns matched method.
func (r *Request) Method() []byte {
	return r.method
}

// Params returns request parameters already unmarshalled with valyala/fastjson.
func (r *Request) Params() *fastjson.Value {
	return r.params
}

// ParamsBytes returns raw bytes of request's "params" field.
func (r *Request) ParamsBytes() []byte {
	return r.paramsBytes.B
}
