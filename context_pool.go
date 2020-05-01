package fastjsonrpc

import (
	"sync"

	"github.com/valyala/bytebufferpool"
)

type contextPool struct {
	pool           sync.Pool
	bytebufferpool bytebufferpool.Pool
}

func (ctx *RequestCtx) reset() {
	ctx.fasthttpCtx = nil
	ctx.arena.Reset()

	ctx.id = ctx.id[:0]
	ctx.method = ctx.method[:0]

	ctx.params = nil

	ctx.bytebufferpool.Put(ctx.paramsBytes)
	ctx.bytebufferpool.Put(ctx.response)
}

func (cp *contextPool) Get() *RequestCtx {
	v := cp.pool.Get()
	if v == nil {
		return &RequestCtx{
			paramsBytes: cp.bytebufferpool.Get(),
			response:    cp.bytebufferpool.Get(),

			bytebufferpool: &cp.bytebufferpool,
		}
	}

	r := v.(*RequestCtx)
	r.paramsBytes = r.bytebufferpool.Get()
	r.response = r.bytebufferpool.Get()

	return r
}

func (cp *contextPool) Put(ctx *RequestCtx) {
	ctx.reset()
	cp.pool.Put(ctx)
}
