# fallback
The pattern, similar to RWMutex. With low latency between RLock and WLock, as well as calling a slow lazy asynchronous data loader. The code looks like a WLock call inside RLock. Please, readme https://github.com/golang/go/issues/4026

[![Build Status](https://travis-ci.org/jenchik/fallback.svg)](https://travis-ci.org/jenchik/fallback)
[![GoDoc](https://godoc.org/github.com/jenchik/fallback?status.svg)](https://godoc.org/github.com/jenchik/fallback)
[![Go Report Card](https://goreportcard.com/badge/github.com/jenchik/fallback)](https://goreportcard.com/report/github.com/jenchik/fallback)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjenchik%2Ffallback.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjenchik%2Ffallback?ref=badge_shield)

Installation
------------

```bash
go get github.com/jenchik/fallback
```

Example
-------
```go
package main

import (
	"github.com/jenchik/fallback"

	"fmt"
	"time"
)

var cacheDB map[int]string = make(map[int]string)

func getUserName(f fallback.Fallback, userId int) (name string, found bool) {
	f.Do(func(h fallback.Helper) {
		// sample code in shared lock
		name, found = cacheDB[userId]
		if found {
			return
		}

		load := make(chan struct{})
		h.Exclusive(func() {
			// sample code in exclusive lock
			cacheDB[userId] = name
		}, func() {
			// for example we load from DB or HTTP
			time.Sleep(time.Millisecond * 100)
			name = fmt.Sprintf("User ID: %d", userId)
			close(load)
		})

		<-load
	})

	return
}

func main() {
	start := time.Now()
	f := fallback.NewBatchFallback(1)

	for i := 0; i < 100; i++ {
		go func(id int) {
			time.Sleep(time.Millisecond)
			getUserName(f, id)
		}(i)
	}

	for i := 100; i > 0; i-- {
		fmt.Println(getUserName(f, i))
	}

	fmt.Println("Duration:", time.Since(start).String())
}
```

Best usage
----------
```go
package main

import (
	"github.com/jenchik/fallback"
	"github.com/jenchik/listener"

	// ...
)

type (
	// T ...
)

// ...

// GetSampleResult with competitive access for each 'conditionKey' value will be called only once AsyncLoader and ExclusiveHandler
func GetSampleResult(f fallback.Fallback, listeners listener.Listeners, conditionKey interface{}) (result T, found bool) {
	var l listener.Listener

	f.Do(func(h fallback.Helper) {
		// sample code in shared lock
		if /* some condition */ {
			// result = ...
			return
		}

		if l, found = listeners.GetOrCreate(conditionKey); found {
			return
		}

		h.Exclusive(func() {
			// sample code in exclusive lock

			// ... = l.Wait().(T)

			listeners.Delete(conditionKey)
		}, func() {
			// for example may be async load from DB or HTTP

			var val T
			// val = ...

			l.Broadcast(val)
		})
	})

	if l != nil {
		return l.Wait().(T)
	}

	return
}
```

Benchmarks
----------
with lazyLoaderTimeout = 500ms
```
BenchmarkThreadsDoClassic-4                 10000        100202 ns/op       0.02 MB/s         166 B/op           4 allocs/op
BenchmarkThreadsDoBatchSizeN1-4             10000        100485 ns/op       0.02 MB/s         161 B/op           4 allocs/op
BenchmarkThreadsDoBatchSizeN1000-4          10000        100200 ns/op       0.02 MB/s         155 B/op           3 allocs/op
BenchmarkThreadsConcurrentMap-4             10000        350299 ns/op       0.01 MB/s          18 B/op           0 allocs/op
```

with lazyLoaderTimeout = 100ms
```
BenchmarkThreadsDoClassic-4              10000000           173 ns/op      11.55 MB/s          32 B/op           1 allocs/op
BenchmarkThreadsDoBatchSizeN1-4          10000000           176 ns/op      11.30 MB/s          32 B/op           1 allocs/op
BenchmarkThreadsDoBatchSizeN1000-4       10000000           154 ns/op      12.95 MB/s          32 B/op           1 allocs/op
BenchmarkThreadsConcurrentMap-4          10000000           108 ns/op      18.51 MB/s           0 B/op           0 allocs/op
```

with lazyLoaderTimeout = 50ms
```
BenchmarkThreadsDoClassic-4              10000000           179 ns/op      11.14 MB/s          32 B/op           1 allocs/op
BenchmarkThreadsDoBatchSizeN1-4          10000000           154 ns/op      12.94 MB/s          32 B/op           1 allocs/op
BenchmarkThreadsDoBatchSizeN1000-4       10000000           168 ns/op      11.88 MB/s          32 B/op           1 allocs/op
BenchmarkThreadsConcurrentMap-4          50000000           24.6 ns/op     81.30 MB/s           0 B/op           0 allocs/op
```

with lazyLoaderTimeout = 10ms
```
BenchmarkThreadsDoClassic-4              10000000           153 ns/op      13.04 MB/s          32 B/op           1 allocs/op
BenchmarkThreadsDoBatchSizeN1-4          10000000           159 ns/op      12.56 MB/s          32 B/op           1 allocs/op
BenchmarkThreadsDoBatchSizeN1000-4       10000000           155 ns/op      12.86 MB/s          32 B/op           1 allocs/op
BenchmarkThreadsConcurrentMap-4          100000000         20.0 ns/op     100.21 MB/s           0 B/op           0 allocs/op
```

with lazyLoaderTimeout = 1ms
```
BenchmarkThreadsClassicWithMutex-4       20000            54138 ns/op       0.04 MB/s          39 B/op           0 allocs/op
BenchmarkThreadsDoClassic-4              10000000           188 ns/op      10.59 MB/s          36 B/op           1 allocs/op
BenchmarkThreadsDoBatchSizeN1-4          10000000           147 ns/op      13.55 MB/s          32 B/op           1 allocs/op
BenchmarkThreadsDoBatchSizeN1000-4       10000000           141 ns/op      14.17 MB/s          32 B/op           1 allocs/op
BenchmarkThreadsConcurrentMap-4          100000000         18.3 ns/op     109.36 MB/s           0 B/op           0 allocs/op
```


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjenchik%2Ffallback.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjenchik%2Ffallback?ref=badge_large)