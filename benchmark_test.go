package fallback_test

import (
	"github.com/jenchik/fallback"
	"github.com/jenchik/listener"

	"sync/atomic"
	"testing"
)

func BenchmarkClassicWithMutex(b *testing.B) {
	v := newClassicWithMutexHelper()
	var key string

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key = range dispersion {
			v.do(key)
		}
	}
}

func BenchmarkDoClassic(b *testing.B) {
	v := newClassicHelper()
	var cRes listener.Listener
	var key string

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key = range dispersion {
			v.fb.Do(func(h fallback.Helper) {
				_, _, cRes = v.do(key, h)
			})

			if cRes != nil {
				_ = cRes.Wait().(int)
			}
		}
	}
}

func BenchmarkDoBatchSizeN1(b *testing.B) {
	v := newBatchHelper(1)
	var cRes listener.Listener
	var key string

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key = range dispersion {
			v.fb.Do(func(h fallback.Helper) {
				_, _, cRes = v.do(key, h)
			})

			if cRes != nil {
				_ = cRes.Wait().(int)
			}
		}
	}
	v.fb.Close()
}

func BenchmarkDoBatchSizeN100(b *testing.B) {
	v := newBatchHelper(100)
	var cRes listener.Listener
	var key string

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key = range dispersion {
			v.fb.Do(func(h fallback.Helper) {
				_, _, cRes = v.do(key, h)
			})

			if cRes != nil {
				_ = cRes.Wait().(int)
			}
		}
	}
	v.fb.Close()
}

func BenchmarkThreadsClassicWithMutex(b *testing.B) {
	v := newClassicWithMutexHelper()
	var d uint32

	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		dd := atomic.AddUint32(&d, 1)
		disp := dispersions[int(dd)%benchWorks]
		var i int
		var key string
		for pb.Next() {
			key = disp[i%steps]
			v.do(key)
			i++
		}
	})
}

func BenchmarkThreadsDoClassic(b *testing.B) {
	var d uint32
	v := newClassicHelper()

	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		dd := atomic.AddUint32(&d, 1)
		disp := dispersions[int(dd)%benchWorks]
		var cRes listener.Listener
		var i int
		var key string
		for pb.Next() {
			key = disp[i%steps]
			v.fb.Do(func(h fallback.Helper) {
				_, _, cRes = v.do(key, h)
			})

			if cRes != nil {
				_ = cRes.Wait().(int)
			}
			i++
		}
	})
}

func BenchmarkThreadsDoBatchSizeN1(b *testing.B) {
	var d uint32
	v := newBatchHelper(1)

	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		dd := atomic.AddUint32(&d, 1)
		disp := dispersions[int(dd)%benchWorks]
		var cRes listener.Listener
		var i int
		var key string
		for pb.Next() {
			key = disp[i%steps]
			v.fb.Do(func(h fallback.Helper) {
				_, _, cRes = v.do(key, h)
			})

			if cRes != nil {
				_ = cRes.Wait().(int)
			}
			i++
		}
	})
	v.fb.Close()
}

func BenchmarkThreadsDoBatchSizeN1000(b *testing.B) {
	var d uint32
	v := newBatchHelper(1000)

	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		dd := atomic.AddUint32(&d, 1)
		disp := dispersions[int(dd)%benchWorks]
		var cRes listener.Listener
		var i int
		var key string
		for pb.Next() {
			key = disp[i%steps]
			v.fb.Do(func(h fallback.Helper) {
				_, _, cRes = v.do(key, h)
			})

			if cRes != nil {
				_ = cRes.Wait().(int)
			}
			i++
		}
	})
	v.fb.Close()
}
