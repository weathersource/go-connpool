# go-connpool

[![Build Status](https://travis-ci.com/weathersource/go-connpool.svg?branch=master)](https://travis-ci.com/weathersource/go-connpool)
[![Codevov](https://img.shields.io/codecov/c/github/weathersource/go-connpool.svg)](https://codecov.io/gh/weathersource/go-connpool)
[![Go Report Card](https://goreportcard.com/badge/github.com/weathersource/go-connpool)](https://goreportcard.com/report/github.com/weathersource/go-connpool)
[![GoDoc](https://img.shields.io/badge/godoc-ref-blue.svg)](https://godoc.org/github.com/weathersource/go-connpool) |

Package connpool provides GRPC connection pooling. This allows clients to distribute
their requests across (level 3 or 4) load balanced services.

This package is based on the [grpc-go-pool](https://github.com/processout/grpc-go-pool)
package, but cures a fundamental inefficiency. GRPC allows for multiplexing requests
on a single connection. But grpc-go-pool removes a connection from the pool for each
call to it's New(...) function, not returning it to the pool until it is closed. This
limits the connection to a single client at a time, and thus a single request at a time.

This package leaves connections in the pool for use by multiple clients, thus allowing
for multiplexed requests. In application, we have seen a 5-fold increase in throughput.
