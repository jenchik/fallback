package main

import (
	"github.com/jenchik/fallback"
	"github.com/jenchik/listener"

	"fmt"
	"sync"
	"testing"
)

const (
	fibonacciN = 45
	benchWorks = 1000
)

func BenchmarkFibonacciMemoize(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var fibonacciMemoize memoizeFunction

		fibonacciMemoize = Memoize(func(x int, xs ...int) interface{} {
			if x < 2 {
				return x
			}
			return fibonacciMemoize(x-1).(int) + fibonacciMemoize(x-2).(int)
		})
		fibonacciMemoize(fibonacciN)
	}
}

func BenchmarkThreadsFibonacciMemoize_x10(b *testing.B) {
	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var fibonacciMemoize memoizeFunction
		fibonacciMemoize = Memoize(func(x int, xs ...int) interface{} {
			if x < 2 {
				return x
			}
			return fibonacciMemoize(x-1).(int) + fibonacciMemoize(x-2).(int)
		})
		for pb.Next() {
			fibonacciMemoize(fibonacciN * 10)
		}
	})
}

func BenchmarkRealThreadsFibonacciMemoize(b *testing.B) {
	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	var fibonacciMemoize memoizeFunction
	fibonacciMemoize = ConcurrentMemoize(func(x int, xs ...int) interface{} {
		if x < 2 {
			return x
		}
		return fibonacciMemoize(x-1).(int) + fibonacciMemoize(x-2).(int)
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fibonacciMemoize(fibonacciN)
		}
	})
}

// ***

func BenchmarkFibonacciFallback(b *testing.B) {
	var fibonacciFallback func(int, map[int]int, fallback.Fallback) int

	fibonacciFallback = func(x int, fbCache map[int]int, fb fallback.Fallback) int {
		if x < 2 {
			return x
		}

		var c1 chan int
		fb.Do(func(h fallback.Helper) {
			if xx, found := fbCache[x]; found {
				x = xx
				return
			}

			var xx int
			c1 = make(chan int, 1)
			h.Exclusive(func() {
				fbCache[x] = xx
			}, func() {
				xx = fibonacciFallback(x-1, fbCache, fb)
				xx += fibonacciFallback(x-2, fbCache, fb)
				c1 <- xx
			})
		})

		if c1 != nil {
			return <-c1
		}

		return x
	}

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fb := fallback.NewClassicFallback()
		fbCache := make(map[int]int, 32)
		fibonacciFallback(fibonacciN, fbCache, fb)
	}
}

func BenchmarkThreadsFibonacciFallback_x10(b *testing.B) {
	var fibonacciFallback func(int, map[int]int, fallback.Fallback) int

	fibonacciFallback = func(x int, fbCache map[int]int, fb fallback.Fallback) int {
		if x < 2 {
			return x
		}

		var c1 chan int
		fb.Do(func(h fallback.Helper) {
			if xx, found := fbCache[x]; found {
				x = xx
				return
			}

			var xx int
			c1 = make(chan int, 1)
			h.Exclusive(func() {
				fbCache[x] = xx
			}, func() {
				xx = fibonacciFallback(x-1, fbCache, fb)
				xx += fibonacciFallback(x-2, fbCache, fb)
				c1 <- xx
			})
		})

		if c1 != nil {
			return <-c1
		}

		return x
	}

	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		fb := fallback.NewClassicFallback()
		fbCache := make(map[int]int, 32)
		for pb.Next() {
			fibonacciFallback(fibonacciN*10, fbCache, fb)
		}
	})
}

func BenchmarkRealThreadsFibonacciFallback(b *testing.B) {
	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	var fibonacciFallback func(int) int

	fb := fallback.NewClassicFallback()
	fbCache := make(map[int]int, 32)

	fibonacciFallback = func(x int) int {
		if x < 2 {
			return x
		}

		var c1 chan int
		fb.Do(func(h fallback.Helper) {
			if xx, found := fbCache[x]; found {
				x = xx
				return
			}

			var xx int
			c1 = make(chan int, 1)
			h.Exclusive(func() {
				fbCache[x] = xx
			}, func() {
				xx = fibonacciFallback(x - 1)
				xx += fibonacciFallback(x - 2)
				c1 <- xx
			})
		})

		if c1 != nil {
			return <-c1
		}

		return x
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fibonacciFallback(fibonacciN)
		}
	})
}

// ***

func BenchmarkFibonacciFallbackWithListener(b *testing.B) {
	var fibonacciFallback func(int, map[int]int, fallback.Fallback) int
	listeners := listener.NewListeners()

	fibonacciFallback = func(x int, fbCache map[int]int, fb fallback.Fallback) int {
		if x < 2 {
			return x
		}

		var l listener.Listener
		fb.Do(func(h fallback.Helper) {
			xx, found := fbCache[x]
			if found {
				x = xx
				return
			}

			l, found = listeners.GetOrCreate(x)
			if found {
				return
			}

			h.Exclusive(func() {
				fbCache[x] = xx
			}, func() {
				xx = fibonacciFallback(x-1, fbCache, fb)
				xx += fibonacciFallback(x-2, fbCache, fb)
				l.Broadcast(xx)
			})
		})

		if l != nil {
			return l.Wait().(int)
		}

		return x
	}

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fb := fallback.NewClassicFallback()
		fbCache := make(map[int]int, 32)
		fibonacciFallback(fibonacciN, fbCache, fb)
	}
}

func BenchmarkThreadsFallbackWithListener_x10(b *testing.B) {
	var fibonacciFallback func(int, map[int]int, fallback.Fallback) int
	listeners := listener.NewListeners()

	fibonacciFallback = func(x int, fbCache map[int]int, fb fallback.Fallback) int {
		if x < 2 {
			return x
		}

		var l listener.Listener
		fb.Do(func(h fallback.Helper) {
			xx, found := fbCache[x]
			if found {
				x = xx
				return
			}

			l, found = listeners.GetOrCreate(x)
			if found {
				return
			}

			h.Exclusive(func() {
				fbCache[x] = xx
			}, func() {
				xx = fibonacciFallback(x-1, fbCache, fb)
				xx += fibonacciFallback(x-2, fbCache, fb)
				l.Broadcast(xx)
			})
		})

		if l != nil {
			return l.Wait().(int)
		}

		return x
	}

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		fb := fallback.NewClassicFallback()
		fbCache := make(map[int]int, 32)
		for pb.Next() {
			fibonacciFallback(fibonacciN*10, fbCache, fb)
		}
	})
}

func BenchmarkRealThreadsFallbackWithListener(b *testing.B) {
	var fibonacciFallback func(int, map[int]int, fallback.Fallback) int
	fb := fallback.NewClassicFallback()
	fbCache := make(map[int]int, 32)
	listeners := listener.NewListeners()

	fibonacciFallback = func(x int, fbCache map[int]int, fb fallback.Fallback) int {
		if x < 2 {
			return x
		}

		var l listener.Listener
		fb.Do(func(h fallback.Helper) {
			xx, found := fbCache[x]
			if found {
				x = xx
				return
			}

			l, found = listeners.GetOrCreate(x)
			if found {
				return
			}

			h.Exclusive(func() {
				fbCache[x] = xx
			}, func() {
				xx = fibonacciFallback(x-1, fbCache, fb)
				xx += fibonacciFallback(x-2, fbCache, fb)
				l.Broadcast(xx)
			})
		})

		if l != nil {
			return l.Wait().(int)
		}

		return x
	}

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fibonacciFallback(fibonacciN, fbCache, fb)
		}
	})
}

func ConcurrentMemoize(function memoizeFunction) memoizeFunction {
	cache := make(map[string]interface{})
	mu := &sync.RWMutex{}

	return func(x int, xs ...int) interface{} {
		key := fmt.Sprint(x)
		for _, i := range xs {
			key += fmt.Sprintf(",%d", i)
		}
		mu.RLock()
		value, found := cache[key]
		if found {
			mu.RUnlock()
			return value
		}
		mu.RUnlock()
		mu.Lock()
		value, found = cache[key]
		if found {
			mu.Unlock()
			return value
		}
		c1 := make(chan interface{}, 1)
		go func() {
			v := function(x, xs...)
			c1 <- v
			mu.Lock()
			cache[key] = v
			mu.Unlock()
		}()
		mu.Unlock()
		value = <-c1
		return value
	}
}
