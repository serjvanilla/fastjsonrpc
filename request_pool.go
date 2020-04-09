package fastjsonrpc

import (
	"sync"
)

type requestPool struct {
	pool sync.Pool
}

func (r *Request) reset() {
	r.ctx = nil
	r.arena.Reset()

	r.id = r.id[:0]
	r.method = r.method[:0]

	r.params = nil

	bufferpool.Put(r.paramsBytes)
	bufferpool.Put(r.response)
}

func (cp *requestPool) Get() *Request {
	v := cp.pool.Get()
	if v == nil {
		return &Request{
			paramsBytes: bufferpool.Get(),
			response:    bufferpool.Get(),
		}
	}

	ctx := v.(*Request)
	ctx.paramsBytes = bufferpool.Get()
	ctx.response = bufferpool.Get()

	return ctx
}

func (cp *requestPool) Put(ctx *Request) {
	ctx.reset()
	cp.pool.Put(ctx)
}
