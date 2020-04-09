# fastjsonrpc [![GoDoc](https://godoc.org/github.com/serjvanilla/fastjsonrpc?status.svg)](http://godoc.org/github.com/serjvanilla/fastjsonrpc)
Fast JSON-RPC 2.0 implementation for [fasthttp](https://github.com/valyala/fasthttp) server.

```
$ GOMAXPROCS=1 go test -bench=. -benchmem -benchtime 10s
BenchmarkEchoHandler            20206782               588 ns/op               0 B/op          0 allocs/op
BenchmarkSumHandler             16310700               732 ns/op               0 B/op          0 allocs/op
BenchmarkBatchSumHandler         7480480              1591 ns/op               0 B/op          0 allocs/op
```

## Install
```
go get -u github.com/serjvanilla/fastjsonrpc
```

## TODO
- [ ] Documentation
- [ ] Examples
- [ ] End-to-end benchmarks
- [ ] Migration from https://github.com/osamingo/jsonrpc examples