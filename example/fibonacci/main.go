package main

import (
	"github.com/jenchik/fallback"

	"fmt"
	"runtime"
	"time"
)

type (
	memoizeFunction func(int, ...int) interface{}
)

var (
	FibonacciMemoize   memoizeFunction
	FibonacciFallback  func(int) int
	FibonacciRecursive func(int) int

	FB      = fallback.NewClassicFallback()
	fbCache = make(map[int]int, 32)
)

func init() {
	FibonacciMemoize = Memoize(func(x int, xs ...int) interface{} {
		if x < 2 {
			return x
		}
		return FibonacciMemoize(x-1).(int) + FibonacciMemoize(x-2).(int)
	})

	FibonacciFallback = func(x int) int {
		if x < 2 {
			return x
		}

		var c1 chan int
		FB.Do(func(h fallback.Helper) {
			if xx, found := fbCache[x]; found {
				x = xx
				return
			}

			var xx int
			c1 = make(chan int, 1)
			h.Exclusive(func() {
				fbCache[x] = xx
			}, func() {
				xx = FibonacciFallback(x - 1)
				xx += FibonacciFallback(x - 2)
				c1 <- xx
			})
		})

		if c1 != nil {
			return <-c1
		}

		return x
	}

	FibonacciRecursive = func(x int) int {
		if x < 2 {
			return x
		}
		return FibonacciRecursive(x-1) + FibonacciRecursive(x-2)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Starting...")

	t1 := time.Now()
	fmt.Println("FibonacciRecursive(40) =", FibonacciRecursive(40), time.Now().Sub(t1).String())

	t1 = time.Now()
	fmt.Println("FibonacciMemoize(45) =", FibonacciMemoize(45), time.Now().Sub(t1).String())
	t1 = time.Now()
	fmt.Println("FibonacciFallback(45) =", FibonacciFallback(45), time.Now().Sub(t1).String())

	t1 = time.Now()
	fmt.Println("repeat FibonacciMemoize(45) =", FibonacciMemoize(45), time.Now().Sub(t1).String())
	t1 = time.Now()
	fmt.Println("repeat FibonacciFallback(45) =", FibonacciFallback(45), time.Now().Sub(t1).String())

	t1 = time.Now()
	fmt.Println("large FibonacciMemoize(400) =", FibonacciMemoize(400), time.Now().Sub(t1).String())
	t1 = time.Now()
	fmt.Println("large FibonacciFallback(400) =", FibonacciFallback(400), time.Now().Sub(t1).String())
}

func Memoize(function memoizeFunction) memoizeFunction {
	cache := make(map[string]interface{})
	return func(x int, xs ...int) interface{} {
		key := fmt.Sprint(x)
		for _, i := range xs {
			key += fmt.Sprintf(",%d", i)
		}
		if value, found := cache[key]; found {
			return value
		}
		value := function(x, xs...)
		cache[key] = value
		return value
	}
}
