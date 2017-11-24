// +build go1.9

package main

import (
	"fmt"
	"sync"
	"testing"
)

func BenchmarkConcurrentMap(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var fibonacciMemoize memoizeFunction

		fibonacciMemoize = ConcurrentMapMemoize(func(x int, xs ...int) interface{} {
			if x < 2 {
				return x
			}
			return fibonacciMemoize(x-1).(int) + fibonacciMemoize(x-2).(int)
		})
		fibonacciMemoize(fibonacciN)
	}
}

func BenchmarkThreadsConcurrentMap_x10(b *testing.B) {
	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var fibonacciMemoize memoizeFunction

		fibonacciMemoize = ConcurrentMapMemoize(func(x int, xs ...int) interface{} {
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

func BenchmarkRealThreadsConcurrentMap(b *testing.B) {
	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	var fibonacciMemoize memoizeFunction

	fibonacciMemoize = ConcurrentMapMemoize(func(x int, xs ...int) interface{} {
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

func ConcurrentMapMemoize(function memoizeFunction) memoizeFunction {
	cache := &sync.Map{}

	return func(x int, xs ...int) interface{} {
		key := fmt.Sprint(x)
		for _, i := range xs {
			key += fmt.Sprintf(",%d", i)
		}

		value, found := cache.Load(key)
		if found {
			return value
		}

		c1 := make(chan interface{}, 1)
		go func() {
			v := function(x, xs...)
			c1 <- v
			cache.Store(key, v)
		}()
		return <-c1
	}
}
