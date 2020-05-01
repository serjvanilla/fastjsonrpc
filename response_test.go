package fastjsonrpc_test

import (
	"testing"
	"testing/iotest"

	. "github.com/serjvanilla/fastjsonrpc"
	"github.com/valyala/fasthttp"
)

func TestResponseWithError(t *testing.T) {
	r := NewRepository()
	r.Register("ping", func(ctx *RequestCtx) {
		ctx.SetError(ErrInternalError())
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{"jsonrpc":"2.0","method":"ping","id":1}`)

	r.RequestHandler()(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatal("unexpected status code")
	}

	if !assertJSONUnordered(
		t, string(ctx.Response.Body()),
		`{"jsonrpc":"2.0","error":{"code":-32603,"message":"Internal error"},"id":1}`,
	) {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}

func TestResponseWithStdError(t *testing.T) {
	r := NewRepository()
	r.Register("ping", func(ctx *RequestCtx) {
		ctx.SetError(iotest.ErrTimeout)
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{"jsonrpc":"2.0","method":"ping","id":1}`)

	r.RequestHandler()(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatal("unexpected status code")
	}

	if !assertJSONUnordered(
		t, string(ctx.Response.Body()),
		`{"jsonrpc":"2.0","error":{"code":-32000,"message":"Server error","data":"timeout"},"id":1}`,
	) {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}
