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
	contextPool contextPool
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
		rCtx.fasthttpCtx = ctx

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

func (r *Repository) handleRequest(rCtx *RequestCtx, request *fastjson.Value) {
	jsonrpc := request.GetStringBytes("jsonrpc")
	method := request.GetStringBytes("method")

	if !bytes.Equal(jsonrpc, []byte(`2.0`)) || len(method) == 0 {
		rCtx.writeString(renderedInvalidRequest)
		return
	}

	if id := request.Get("id"); id != nil {
		rCtx.id = id.MarshalTo(rCtx.id)
	}

	handler, ok := r.handlers[string(method)]
	if !ok {
		rCtx.SetError(ErrMethodNotFound())
		return
	}

	rCtx.method = method
	rCtx.params = request.Get("params")

	if rCtx.params != nil {
		rCtx.paramsBytes.B = rCtx.params.MarshalTo(rCtx.paramsBytes.B)
	}

	defer func() {
		if recover() != nil {
			rCtx.SetError(ErrInternalError())
		}
	}()

	handler(rCtx)
}

func (r *Repository) handleBatchRequest(batchCtx *RequestCtx, requests *fastjson.Value) {
	requestsArr := requests.GetArray()
	if len(requestsArr) == 0 {
		batchCtx.writeString(renderedInvalidRequest)
		return
	}

	_ = batchCtx.response.WriteByte('[')

	var needComma bool

	for _, request := range requestsArr {
		rCtx := r.contextPool.Get()
		rCtx.fasthttpCtx = batchCtx.fasthttpCtx
		r.handleRequest(rCtx, request)

		if rCtx.response.Len() > 0 {
			if needComma {
				_ = batchCtx.response.WriteByte(',')
				needComma = false
			}

			if n, _ := batchCtx.response.Write(rCtx.response.B); n != 0 {
				needComma = true
			}
		}

		r.contextPool.Put(rCtx)
	}

	if batchCtx.response.Len() > 1 {
		_ = batchCtx.response.WriteByte(']')
	} else {
		batchCtx.response.Reset()
	}
}
