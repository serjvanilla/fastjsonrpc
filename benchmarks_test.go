package fastjsonrpc_test

import (
	"testing"

	. "github.com/serjvanilla/fastjsonrpc"
	"github.com/valyala/fasthttp"
)

func BenchmarkEchoHandler(b *testing.B) {
	r := NewRepository()
	r.Register("echo", func(ctx *Request) {
		ctx.Result(ctx.Params())
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{"jsonrpc":"2.0","method":"echo","params":"hello","id":1}`)

	handler := r.RequestHandler()

	for i := 0; i < b.N; i++ {
		ctx.Response.ResetBody()
		handler(ctx)
	}
}

func BenchmarkSumHandler(b *testing.B) {
	r := NewRepository()
	r.Register("sum", func(req *Request) {
		params := req.Params()

		a := params.GetInt("a")
		b := params.GetInt("b")

		req.Result(req.Arena().NewNumberInt(a + b))
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{"jsonrpc":"2.0","method":"sum","params":{"a":7,"b":42},"id":1}`)

	handler := r.RequestHandler()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx.Response.ResetBody()
		handler(ctx)
	}
}

func BenchmarkBatchSumHandler(b *testing.B) {
	r := NewRepository()
	r.Register("sum", func(ctx *Request) {
		params := ctx.Params()

		a := params.GetInt("a")
		b := params.GetInt("b")

		ctx.Result(ctx.Arena().NewNumberInt(a + b))
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(
		`[{"jsonrpc":"2.0","method":"sum","params":{"a":7,"b":42},"id":1},
		{"jsonrpc":"2.0","method":"sum","params":{"a":42,"b":7},"id":2}]`)

	handler := r.RequestHandler()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx.Response.ResetBody()
		handler(ctx)
	}
}

func BenchmarkErrorHandler(b *testing.B) {
	r := NewRepository()
	r.Register("error", func(req *Request) {
		req.Error(ErrServerError(ErrorCode(-32000)).WithMessage("Server defined error"))
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{"jsonrpc":"2.0","method":"error"}`)

	handler := r.RequestHandler()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx.Response.ResetBody()
		handler(ctx)
	}
}
