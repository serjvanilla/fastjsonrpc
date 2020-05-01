package fastjsonrpc_test

import (
	"testing"

	"github.com/valyala/fasthttp"

	. "github.com/serjvanilla/fastjsonrpc"
)

func TestRepositoryMethodPost(t *testing.T) {
	r := NewRepository()

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodGet)

	r.RequestHandler()(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusMethodNotAllowed {
		t.Fatal("unexpected status code")
	}
}

func TestRepositoryBadRequest(t *testing.T) {
	r := NewRepository()

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)

	r.RequestHandler()(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatal("unexpected status code")
	}

	if !assertJSONUnordered(t, string(ctx.Response.Body()),
		`{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}`) {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}

func TestRepositoryValidJSON(t *testing.T) {
	r := NewRepository()

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`true`)

	r.RequestHandler()(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatal("unexpected status code")
	}

	if !assertJSONUnordered(
		t, string(ctx.Response.Body()),
		`{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null}`,
	) {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}

func TestRepositoryInvalidRequest(t *testing.T) {
	r := NewRepository()

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{}`)

	r.RequestHandler()(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatal("unexpected status code")
	}

	if !assertJSONUnordered(
		t, string(ctx.Response.Body()),
		`{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null}`,
	) {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}

func TestRepositoryMethodNotFound(t *testing.T) {
	r := NewRepository()

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{"jsonrpc":"2.0","method":"ping","id":1}`)

	r.RequestHandler()(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatal("unexpected status code")
	}

	if !assertJSONUnordered(
		t, string(ctx.Response.Body()),
		`{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":1}`,
	) {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}

func TestRepositoryRequestWithoutParams(t *testing.T) {
	r := NewRepository()
	r.Register("ping", func(ctx *RequestCtx) {
		ctx.SetResult(ctx.Arena().NewTrue())
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
		`{"jsonrpc":"2.0","result":true,"id":1}`,
	) {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}

func TestRepositoryRequestNotification(t *testing.T) {
	r := NewRepository()
	r.Register("ping", func(ctx *RequestCtx) {
		ctx.SetResult(ctx.Arena().NewTrue())
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`{"jsonrpc":"2.0","method":"ping"}`)

	r.RequestHandler()(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatal("unexpected status code")
	}

	if len(ctx.Response.Body()) != 0 {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}

func TestRepositoryRequestWithParams(t *testing.T) {
	r := NewRepository()
	r.Register("ping", func(ctx *RequestCtx) {
		ctx.SetResult(ctx.Arena().NewTrue())
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
		`{"jsonrpc":"2.0","result":true,"id":1}`,
	) {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}

func TestRepositoryEmptyBatchRequest(t *testing.T) {
	r := NewRepository()

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(`[]`)

	r.RequestHandler()(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatal("unexpected status code")
	}

	if !assertJSONUnordered(
		t, string(ctx.Response.Body()),
		`{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null}`,
	) {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}

func TestRepositoryBatchRequest(t *testing.T) {
	r := NewRepository()
	r.Register("echo", func(ctx *RequestCtx) {
		ctx.SetResult(ctx.Params())
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(
		`[{"jsonrpc":"2.0","method":"echo","params":"foo","id":1},{"jsonrpc":"2.0","method":"echo","params":"bar","id":2}]`,
	)

	r.RequestHandler()(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatal("unexpected status code")
	}

	if !assertJSONUnordered(
		t, string(ctx.Response.Body()),
		`[{"jsonrpc":"2.0","result":"foo","id":1},{"jsonrpc":"2.0","result":"bar","id":2}]`,
	) {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}

func TestRepositoryBatchRequestWithError(t *testing.T) {
	r := NewRepository()
	r.Register("echo", func(ctx *RequestCtx) {
		ctx.SetResult(ctx.Params())
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(
		`[{"jsonrpc":"2.0","method":"echo","params":"foo","id":1},{"jsonrpc":"2.0","method":"ping","id":2}]`,
	)

	r.RequestHandler()(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatal("unexpected status code")
	}

	if !assertJSONUnordered(
		t, string(ctx.Response.Body()),
		`[{"jsonrpc":"2.0","result":"foo","id":1},{"jsonrpc":"2.0",
		"error":{"code":-32601,"message":"Method not found"},"id":2}]`,
	) {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}

func TestRepositoryBatchRequestWithNotification(t *testing.T) {
	r := NewRepository()
	r.Register("echo", func(ctx *RequestCtx) {
		ctx.SetResult(ctx.Params())
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(
		`
        [
          {"jsonrpc":"2.0","method":"echo","params":"foo","id":1},
          {"jsonrpc":"2.0","method":"echo","params":"bar"},
          {"jsonrpc":"2.0","method":"echo","params":"baz","id":3}
        ]`,
	)

	r.RequestHandler()(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatal("unexpected status code")
	}

	if !assertJSONUnordered(
		t, string(ctx.Response.Body()),
		`[{"jsonrpc":"2.0","result":"foo","id":1},{"jsonrpc":"2.0","result":"baz","id":3}]`,
	) {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}

func TestRepositoryHandlerPanic(t *testing.T) {
	r := NewRepository()
	r.Register("ping", func(ctx *RequestCtx) {
		panic("ha-ha")
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
