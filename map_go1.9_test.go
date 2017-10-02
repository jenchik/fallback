// +build go1.9

package fallback_test

import (
	"github.com/stretchr/testify/assert"

	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func initConcurrentMap() *sync.Map {
	m := &sync.Map{}
	for i, key := range keys {
		m.Store(key, i)
	}

	return m
}

func TestConcurrentMap(t *testing.T) {
	m := initConcurrentMap()

	start := make(chan struct{})
	wg := &sync.WaitGroup{}
	var cntNew, retries int32

	for i := 0; i < testWorks; i++ {
		wg.Add(1)
		go func(wn int) {
			defer func() {
				wg.Done()
			}()
			var completed, n int
			var key string
			var result interface{}
			var found bool
			<-start
			for i := 0; i < steps; i++ {
				completed++
				n = rand.Intn(items)
				if rand.Intn(2) == 1 {
					key = newKeys[n]
				} else {
					key = keys[n]
				}

				result, found = m.Load(key)
				if found {
					if result.(int) != n {
						assert.FailNow(t, "Results of recording and subsequent reading are not equal")
						break
					}
					continue
				}

				// for example we load from DB
				time.Sleep(lazyLoaderTimeout)
				result = n

				if _, actual := m.LoadOrStore(key, result); !actual {
					atomic.AddInt32(&cntNew, 1)
				} else {
					atomic.AddInt32(&retries, 1)
				}
			}

			assert.Equal(t, steps, completed)
		}(i)
	}
	close(start)

	wg.Wait()

	if testing.Verbose() {
		avgRetry := float32(retries) / float32(testWorks)
		t.Log("New elements:", cntNew)
		t.Logf("Repeated attempts to write for each worker: %f", avgRetry)
	}

	n := 0
	m.Range(func(k, v interface{}) bool {
		n++
		return true
	})

	assert.Equal(t, n, len(keys)+int(cntNew))
}

func BenchmarkConcurrentMap(b *testing.B) {
	m := initConcurrentMap()
	var found bool
	var key string
	var r int

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key = range dispersion {
			_, found = m.Load(key)
			if found {
				continue
			}

			time.Sleep(lazyLoaderTimeout)
			r = -1
			m.LoadOrStore(key, r)
		}
	}
}

func BenchmarkThreadsConcurrentMap(b *testing.B) {
	m := initConcurrentMap()
	var d uint32

	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		dd := atomic.AddUint32(&d, 1)
		disp := dispersions[int(dd)%benchWorks]
		var found bool
		var i int
		var key string
		var r int
		for pb.Next() {
			key = disp[i%steps]
			_, found = m.Load(key)
			if found {
				continue
			}

			time.Sleep(lazyLoaderTimeout)
			r = -1
			m.LoadOrStore(key, r)
			i++
		}
	})
}
