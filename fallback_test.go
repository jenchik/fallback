package fallback_test

import (
	. "github.com/jenchik/fallback"
	"github.com/jenchik/listener"
	"github.com/stretchr/testify/assert"

	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func testFallback(t *testing.T, f Fallback) {
	m := initMap()

	var cntNew, hits int32
	retriesMap := make(map[string]int, len(m))
	start := make(chan struct{})
	wg := &sync.WaitGroup{}
	listeners := listener.NewListeners()

	wg.Add(testWorks)
	for i := 0; i < testWorks; i++ {
		go func(wn int) {
			defer func() {
				wg.Done()
			}()
			var completed, hit, n, result int
			var found bool
			var key string
			var l listener.Listener

			<-start
			for i := 0; i < steps; i++ {
				n = rand.Intn(items)
				if rand.Intn(2) == 1 {
					key = newKeys[n]
				} else {
					key = keys[n]
				}

				f.Do(func(h Helper) {
					l = nil
					if result, found = m[key]; found {
						hit++
						return
					}

					if l, found = listeners.GetOrCreate(key); found {
						return
					}

					k := key
					var r int
					h.Exclusive(func() { // executing with exclusive lock
						m[k] = r
						l.Broadcast(r)
						listeners.Delete(k)

						// for test
						atomic.AddInt32(&cntNew, 1)
						retriesMap[k]++
					}, func() {
						// for example we load from DB
						time.Sleep(time.Millisecond)
						r = n
					})
				})

				if l != nil {
					result = l.Wait().(int)
				}

				if result != n {
					assert.FailNow(t, "Results of recording and subsequent reading are not equal")
					break
				}
				completed++
			}

			atomic.AddInt32(&hits, int32(hit))
			assert.Equal(t, steps, completed)
		}(i)
	}
	close(start)

	wg.Wait()
	maxRetry := 0
	for _, retry := range retriesMap {
		if retry > maxRetry {
			maxRetry = retry
		}
	}

	if testing.Verbose() {
		avgHits := float32(hits) / float32(testWorks)
		t.Log("New elements:", cntNew)
		t.Logf("Avg hits per worker: %f (%f of %d)", float32(avgHits/(float32(steps)/100)), avgHits, steps)
	}

	assert.Equal(t, len(m), len(keys)+int(cntNew))
	assert.Equal(t, 1, maxRetry)
}
