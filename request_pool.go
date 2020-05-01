package fastjsonrpc

import (
	"sync"

	"github.com/valyala/bytebufferpool"
)

type requestPool struct {
	pool           sync.Pool
	bytebufferpool bytebufferpool.Pool
}

func (r *Request) reset() {
	r.ctx = nil
	r.arena.Reset()

	r.id = r.id[:0]
	r.method = r.method[:0]

	r.params = nil

	r.bytebufferpool.Put(r.paramsBytes)
	r.bytebufferpool.Put(r.response)
}

func (cp *requestPool) Get() *Request {
	v := cp.pool.Get()
	if v == nil {
		return &Request{
			paramsBytes: cp.bytebufferpool.Get(),
			response:    cp.bytebufferpool.Get(),

			bytebufferpool: &cp.bytebufferpool,
		}
	}

	r := v.(*Request)
	r.paramsBytes = r.bytebufferpool.Get()
	r.response = r.bytebufferpool.Get()

	return r
}

func (cp *requestPool) Put(ctx *Request) {
	ctx.reset()
	cp.pool.Put(ctx)
}
