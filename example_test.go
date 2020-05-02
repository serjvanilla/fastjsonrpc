package fastjsonrpc_test

import (
	"github.com/serjvanilla/fastjsonrpc"
	"github.com/valyala/fasthttp"
)

func Example() {
	repo := fastjsonrpc.NewRepository()

	repo.Register("sum", func(ctx *fastjsonrpc.RequestCtx) {
		params := ctx.Params()

		a := params.GetInt("a")
		b := params.GetInt("b")

		ctx.SetResult(ctx.Arena().NewNumberInt(a + b))
	})
	repo.Register("sum_struct", func(ctx *fastjsonrpc.RequestCtx) {
		type (
			sumRequest struct {
				A int `json:"a"`
				B int `json:"b"`
			}
			sumResponse int
		)

		var req sumRequest
		if err := ctx.ParamsUnmarshal(&req); err != nil {
			ctx.SetError(err)
			return
		}

		ctx.SetResult(sumResponse(req.A + req.B))
	})

	_ = fasthttp.ListenAndServe(":8080", repo.RequestHandler())
}

func ExampleErrServerError() {
	_ = fastjsonrpc.
		ErrServerError(fastjsonrpc.ErrorCode(-32042)).
		WithData("something went wrong")
}
