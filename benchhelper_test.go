package fallback_test

import (
	"github.com/jenchik/fallback"
	"github.com/jenchik/listener"

	"sync"
	"time"
)

type (
	benchHelper struct {
		fb  fallback.Fallback
		obs *listener.Listeners
		m   map[string]int
	}

	benchClassicHelper struct {
		m  map[string]int
		mu sync.RWMutex
	}
)

func newBatchHelper(size int) *benchHelper {
	return &benchHelper{
		fb:  fallback.NewBatchFallback(size),
		obs: listener.NewListeners(),
		m:   initMap(),
	}
}

func newClassicHelper() *benchHelper {
	return &benchHelper{
		fb:  fallback.NewClassicFallback(),
		obs: listener.NewListeners(),
		m:   initMap(),
	}
}

func newClassicWithMutexHelper() *benchClassicHelper {
	return &benchClassicHelper{
		m: initMap(),
	}
}

func (b *benchHelper) lazyLoader(key string) int {
	// sample, new element
	time.Sleep(lazyLoaderTimeout)
	return -1
}

func (b *benchHelper) exclusive(c listener.Listener, key string, h fallback.Helper) {
	h.Exclusive(func() {
		b.m[key] = c.Wait().(int)
		b.obs.Delete(key)
	}, func() {
		c.Broadcast(b.lazyLoader(key))
	})
}

func (b *benchHelper) do(key string, h fallback.Helper) (res int, found bool, c listener.Listener) {
	res, found = b.m[key]
	if found {
		return
	}

	c, found = b.obs.GetOrCreate(key)
	if found {
		return
	}

	b.exclusive(c, key, h)

	return
}

func (b *benchClassicHelper) lazyLoader(key string) int {
	// sample, new element
	time.Sleep(lazyLoaderTimeout)
	return -1
}

func (b *benchClassicHelper) do(key string) (res int, found bool) {
	b.mu.RLock()
	res, found = b.m[key]
	b.mu.RUnlock()
	if found {
		return
	}

	b.mu.Lock()
	if res, found = b.m[key]; !found {
		b.m[key] = b.lazyLoader(key)
	}
	b.mu.Unlock()

	return
}
