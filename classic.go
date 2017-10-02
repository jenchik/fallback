package fallback

import (
	"sync"
)

type (
	classicFallback struct {
		mu sync.RWMutex
	}

	helperClassic classicFallback
)

var (
	_ Fallback = &classicFallback{}
	_ Helper   = &helperClassic{}
)

func NewClassicFallback() *classicFallback {
	return &classicFallback{}
}

func (f *classicFallback) Close() error {
	return nil
}

func (f *classicFallback) Wait() {
	f.mu.Lock()
	f.mu.Unlock()
}

func (f *classicFallback) Do(sharedHandler func(Helper)) {
	f.mu.RLock() // shared lock
	defer f.mu.RUnlock()

	sharedHandler((*helperClassic)(f))
}

func (h *helperClassic) Exclusive(exclusiveHandler func(), slowAsync func()) {
	if exclusiveHandler == nil {
		return
	}

	go func() {
		if slowAsync != nil {
			slowAsync()
		}
		h.mu.Lock() // exclusive lock
		defer h.mu.Unlock()
		exclusiveHandler()
	}()
}
