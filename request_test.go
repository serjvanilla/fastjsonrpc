package fastjsonrpc_test

import (
	"bytes"
	"reflect"
	"testing"

	. "github.com/serjvanilla/fastjsonrpc"
	"github.com/valyala/fasthttp"
)

func TestRequestID(t *testing.T) {
	r := NewRepository()
	r.Register("ping", func(ctx *RequestCtx) {
		if !bytes.Equal(ctx.ID(), []byte("1")) {
			t.Fatalf("unexpected id: `%s`", ctx.ID())
		}
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{"jsonrpc":"2.0","method":"ping","id":1}`)

	r.RequestHandler()(ctx)
}

func TestRequestMethod(t *testing.T) {
	r := NewRepository()
	r.Register("ping", func(ctx *RequestCtx) {
		if !bytes.Equal(ctx.Method(), []byte("ping")) {
			t.Fatalf("unexpected method: `%s`", ctx.Method())
		}
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{"jsonrpc":"2.0","method":"ping","id":1}`)

	r.RequestHandler()(ctx)
}

func TestRequestParams(t *testing.T) {
	r := NewRepository()
	r.Register("echo", func(ctx *RequestCtx) {
		params := ctx.Params()
		if !bytes.Equal(params.MarshalTo(nil), []byte(`"ping"`)) {
			t.Fatalf("unexpected params: `%s`", ctx.Params())
		}
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{"jsonrpc":"2.0","method":"echo","params":"ping","id":1}`)

	r.RequestHandler()(ctx)
}

func TestRequestParamsBytes(t *testing.T) {
	r := NewRepository()
	r.Register("echo", func(ctx *RequestCtx) {
		params := ctx.ParamsBytes()
		if !bytes.Equal(params, []byte(`"ping"`)) {
			t.Fatalf("unexpected params: `%s`", ctx.ParamsBytes())
		}
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{"jsonrpc":"2.0","method":"echo","params":"ping","id":1}`)

	r.RequestHandler()(ctx)
}

func TestRequestParamsUnmarshal(t *testing.T) {
	r := NewRepository()
	r.Register("echo", func(ctx *RequestCtx) {
		var param string

		err := ctx.ParamsUnmarshal(&param)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if param != "ping" {
			t.Fatalf("unexpected param: `%s`", param)
		}
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{"jsonrpc":"2.0","method":"echo","params":"ping","id":1}`)

	r.RequestHandler()(ctx)
}

func TestRequestParamsUnmarshalEmpty(t *testing.T) {
	r := NewRepository()
	r.Register("ping", func(ctx *RequestCtx) {
		var param string

		err := ctx.ParamsUnmarshal(&param)
		if !reflect.DeepEqual(err, ErrInvalidParams()) {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{"jsonrpc":"2.0","method":"ping","id":1}`)

	r.RequestHandler()(ctx)
}
