package fastjsonrpc

import (
	"bytes"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
)

func NewRepository() *Repository {
	return &Repository{
		handlers: make(map[string]RequestHandler),
	}
}

// Repository is a JSON-RPC 2.0 methods repository.
type Repository struct {
	contextPool requestPool
	parserPool  fastjson.ParserPool

	handlers map[string]RequestHandler
}

// RequestHandler is suitable for using with fasthttp.
func (r *Repository) RequestHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsPost() {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			return
		}

		parser := r.parserPool.Get()

		request, err := parser.ParseBytes(ctx.PostBody())
		if err != nil {
			_, _ = ctx.WriteString(renderedParseError)
			return
		}

		rCtx := r.contextPool.Get()
		rCtx.ctx = ctx

		ctx.SetContentType("application/json")

		switch request.Type() {
		case fastjson.TypeObject:
			r.handleRequest(rCtx, request)
		case fastjson.TypeArray:
			r.handleBatchRequest(rCtx, request)
		default:
			_, _ = rCtx.response.WriteString(renderedInvalidRequest)
		}

		_, _ = rCtx.response.WriteTo(ctx)

		r.contextPool.Put(rCtx)
		r.parserPool.Put(parser)
	}
}

// Register registers new method handler.
func (r *Repository) Register(method string, handler RequestHandler) {
	r.handlers[method] = handler
}

func (r *Repository) handleRequest(ctx *Request, request *fastjson.Value) {
	jsonrpc := request.GetStringBytes("jsonrpc")
	method := request.GetStringBytes("method")

	if !bytes.Equal(jsonrpc, []byte(`2.0`)) || len(method) == 0 {
		ctx.writeString(renderedInvalidRequest)
		return
	}

	if id := request.Get("id"); id != nil {
		ctx.id = id.MarshalTo(ctx.id)
	}

	handler, ok := r.handlers[string(method)]
	if !ok {
		ctx.Error(ErrMethodNotFound())
		return
	}

	ctx.method = method
	ctx.params = request.Get("params")

	if ctx.params != nil {
		ctx.paramsBytes.B = ctx.params.MarshalTo(ctx.paramsBytes.B)
	}

	defer func() {
		if recover() != nil {
			ctx.Error(ErrInternalError())
		}
	}()

	handler(ctx)
}

func (r *Repository) handleBatchRequest(batchCtx *Request, requests *fastjson.Value) {
	requestsArr := requests.GetArray()
	if len(requestsArr) == 0 {
		batchCtx.writeString(renderedInvalidRequest)
		return
	}

	_ = batchCtx.response.WriteByte('[')

	var needComma bool

	for _, request := range requestsArr {
		ctx := r.contextPool.Get()
		ctx.ctx = batchCtx.ctx
		r.handleRequest(ctx, request)

		if ctx.response.Len() > 0 {
			if needComma {
				_ = batchCtx.response.WriteByte(',')
				needComma = false
			}

			if n, _ := batchCtx.response.Write(ctx.response.B); n != 0 {
				needComma = true
			}
		}

		r.contextPool.Put(ctx)
	}

	if batchCtx.response.Len() > 1 {
		_ = batchCtx.response.WriteByte(']')
	} else {
		batchCtx.response.Reset()
	}
}
