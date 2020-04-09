package fastjsonrpc_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/valyala/fastjson"

	"github.com/valyala/fasthttp"

	"github.com/serjvanilla/fastjsonrpc"
)

func specExamplesRepository(t *testing.T) *fastjsonrpc.Repository {
	t.Helper()

	r := fastjsonrpc.NewRepository()

	r.Register("subtract", func(ctx *fastjsonrpc.Request) {
		params := ctx.Params()

		var result int64

		switch params.Type() {
		case fastjson.TypeArray:
			result = params.GetInt64("0") - params.GetInt64("1")
		case fastjson.TypeObject:
			result = params.GetInt64("minuend") - params.GetInt64("subtrahend")
		}

		ctx.Result(result)
	})

	r.Register("sum", func(ctx *fastjsonrpc.Request) {
		params := ctx.Params()

		var result int64

		for _, param := range params.GetArray() {
			result += param.GetInt64()
		}

		ctx.Result(result)
	})

	r.Register("get_data", func(ctx *fastjsonrpc.Request) {
		ctx.Result([]interface{}{"hello", 5})
	})

	return r
}

func TestSpec(t *testing.T) {
	t.Parallel()

	t.Run("PositionalParameters", func(t *testing.T) {
		testSpecExample(t,
			`{"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1}`,
			`{"jsonrpc": "2.0", "result": 19, "id": 1}`,
		)
		testSpecExample(t,
			`{"jsonrpc": "2.0", "method": "subtract", "params": [23, 42], "id": 2}`,
			`{"jsonrpc": "2.0", "result": -19, "id": 2}`,
		)
	})

	t.Run("NamedParameters", func(t *testing.T) {
		testSpecExample(t,
			`{"jsonrpc": "2.0", "method": "subtract", "params": {"subtrahend": 23, "minuend": 42}, "id": 3}`,
			`{"jsonrpc": "2.0", "result": 19, "id": 3}`,
		)
		testSpecExample(t,
			`{"jsonrpc": "2.0", "method": "subtract", "params": {"minuend": 42, "subtrahend": 23}, "id": 4}`,
			`{"jsonrpc": "2.0", "result": 19, "id": 4}`,
		)
	})

	t.Run("Notification", func(t *testing.T) {
		testSpecExample(t,
			`{"jsonrpc": "2.0", "method": "update", "params": [1,2,3,4,5]}`,
			``,
		)
		testSpecExample(t,
			`{"jsonrpc": "2.0", "method": "foobar"}`,
			``,
		)
	})

	t.Run("NonExistentMethod", func(t *testing.T) {
		testSpecExample(t,
			`{"jsonrpc": "2.0", "method": "foobar", "id": "1"}`,
			`{"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": "1"}`,
		)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		testSpecExample(t,
			`{"jsonrpc": "2.0", "method": "foobar, "params": "bar", "baz]`,
			`{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`,
		)
	})

	t.Run("InvalidRequestObject", func(t *testing.T) {
		testSpecExample(t,
			`{"jsonrpc": "2.0", "method": 1, "params": "bar"}`,
			`{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}`,
		)
	})

	t.Run("BatchInvalidJSON", func(t *testing.T) {
		testSpecExample(t,
			`
			[
				{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4], "id": "1"},
				{"jsonrpc": "2.0", "method"
			]`,
			`{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`)
	})

	t.Run("EmptyArray", func(t *testing.T) {
		testSpecExample(t,
			`[]`,
			`{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}`,
		)
	})

	t.Run("InvalidBatch", func(t *testing.T) {
		testSpecExample(t,
			`[1]`,
			`[{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}]`,
		)
		testSpecExample(t,
			`[1,2,3]`,
			`
			[
				{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null},
				{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null},
				{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}
			]`,
		)
	})

	t.Run("Batch", func(t *testing.T) {
		testSpecExample(t,
			`
			[
				{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4], "id": "1"},
				{"jsonrpc": "2.0", "method": "notify_hello", "params": [7]},
				{"jsonrpc": "2.0", "method": "subtract", "params": [42,23], "id": "2"},
				{"foo": "boo"},
				{"jsonrpc": "2.0", "method": "foo.get", "params": {"name": "myself"}, "id": "5"},
				{"jsonrpc": "2.0", "method": "get_data", "id": "9"} 
			]`,
			`
			[
				{"jsonrpc": "2.0", "result": 7, "id": "1"},
				{"jsonrpc": "2.0", "result": 19, "id": "2"},
				{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null},
				{"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": "5"},
				{"jsonrpc": "2.0", "result": ["hello", 5], "id": "9"}
			]`,
		)
	})

	t.Run("BatchNotifications", func(t *testing.T) {
		testSpecExample(t,
			`
			[
				{"jsonrpc": "2.0", "method": "notify_sum", "params": [1,2,4]},
				{"jsonrpc": "2.0", "method": "notify_hello", "params": [7]}
            ]`,
			``,
		)
	})
}

func testSpecExample(t *testing.T, request, response string) {
	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBodyString(request)

	specExamplesRepository(t).RequestHandler()(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatalf("unexpected status code: %d", ctx.Response.StatusCode())
	}

	if !assertJSON(t, response, string(ctx.Response.Body())) {
		t.Fatalf("unexpected response body: `%s`", ctx.Response.Body())
	}
}

func assertJSON(t *testing.T, s1, s2 string) bool {
	t.Helper()

	if s1 == "" && s2 == "" {
		return true
	}

	var (
		err error
		o1  interface{}
		o2  interface{}
	)

	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false
	}

	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}
