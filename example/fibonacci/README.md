# fallback
Fibonacci examples

Benchmarks
----------
go run
```
Starting...
FibonacciRecursive(40) = 102334155 1.391731009s
FibonacciMemoize(45) = 1134903170 83.993µs
FibonacciFallback(45) = 1134903170 162.562µs
repeat FibonacciMemoize(45) = 1134903170 1.014µs
repeat FibonacciFallback(45) = 1134903170 544ns
large FibonacciMemoize(400) = 2121778230729308891 529.777µs
large FibonacciFallback(400) = 2121778230729308891 1.012764ms
```

goos: linux
goarch: amd64
```
BenchmarkFibonacciMemoize-4                  	   50000	     33077 ns/op	   0.06 MB/s	    6747 B/op	     231 allocs/op
BenchmarkThreadsFibonacciMemoize_x10-4       	 1000000	      1932 ns/op	   1.03 MB/s	     250 B/op	       7 allocs/op
BenchmarkRealThreadsFibonacciMemoize-4       	10000000	       246 ns/op	   8.11 MB/s	      22 B/op	       2 allocs/op
BenchmarkFibonacciFallback-4                 	   20000	     78737 ns/op	   0.03 MB/s	   17792 B/op	     442 allocs/op
BenchmarkThreadsFibonacciFallback_x10-4      	 5000000	       796 ns/op	   2.51 MB/s	     234 B/op	       6 allocs/op
BenchmarkRealThreadsFibonacciFallback-4      	10000000	       140 ns/op	  14.19 MB/s	      19 B/op	       2 allocs/op
BenchmarkFibonacciFallbackWithListener-4     	 1000000	      1540 ns/op	   1.30 MB/s	    1336 B/op	       8 allocs/op
BenchmarkThreadsFallbackWithListener_x10-4   	10000000	       168 ns/op	  11.87 MB/s	     101 B/op	       4 allocs/op
BenchmarkRealThreadsFallbackWithListener-4   	10000000	       148 ns/op	  13.51 MB/s	      96 B/op	       4 allocs/op
BenchmarkConcurrentMap-4                     	   10000	    107654 ns/op	   0.02 MB/s	   13669 B/op	     513 allocs/op
BenchmarkThreadsConcurrentMap_x10-4          	20000000	       292 ns/op	   6.83 MB/s	      50 B/op	       3 allocs/op
BenchmarkRealThreadsConcurrentMap-4          	20000000	        97.5 ns/op	  20.52 MB/s	      16 B/op	       2 allocs/op
```
