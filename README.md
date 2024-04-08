# go-int128

[![GoDoc](https://godoc.org/github.com/agmt/go-int128?status.svg)](https://godoc.org/github.com/agmt/go-int128) 
[![Go Report Card](https://goreportcard.com/badge/github.com/agmt/go-int128)](https://goreportcard.com/report/github.com/agmt/go-int128)

128-bit integer arithmetics in Go

## Install
```
go get github.com/agmt/go-int128
```

## Benchmark
```
goos: linux
goarch: arm64
pkg: github.com/agmt/go-int128

BenchmarkAdd/int64+int64-6                                                    	1000000000	         0.3134 ns/op	       0 B/op	       0 allocs/op
BenchmarkBigAdd/big.int+big.int-6                                             	182255571	         6.578 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint128Add/uint128+uint128-6                                         	1000000000	         0.3510 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt128Add/int128+iint128-6                                           	1000000000	         0.3514 ns/op	       0 B/op	       0 allocs/op
```
