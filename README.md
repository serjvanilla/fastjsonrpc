# fastjsonrpc [![GoDoc](https://godoc.org/github.com/serjvanilla/fastjsonrpc?status.svg)](http://godoc.org/github.com/serjvanilla/fastjsonrpc)
Fast JSON-RPC 2.0 implementation for [fasthttp](https://github.com/valyala/fasthttp) server.

```
$ GOMAXPROCS=1 go test -bench=. -benchmem -benchtime 10s
BenchmarkEchoHandler            20473168               584 ns/op               0 B/op          0 allocs/op
BenchmarkSumHandler             16297743               729 ns/op               0 B/op          0 allocs/op
BenchmarkBatchSumHandler         7587087              1569 ns/op               0 B/op          0 allocs/op
BenchmarkErrorHandler           17734203               671 ns/op               0 B/op          0 allocs/op
```

## Install
```
go get -u github.com/serjvanilla/fastjsonrpc
```

## TODO
- [ ] Documentation
- [ ] Examples
- [ ] Parallel batch processing
- [ ] End-to-end benchmarks
- [ ] Migration from https://github.com/osamingo/jsonrpc examples
