package fallback

import (
	"github.com/jenchik/thread"

	"sync"
)

type (
	fallbackBatch struct {
		h  chan func()
		mu sync.RWMutex
		w  *thread.Worker
	}

	helperBatch fallbackBatch
)

var (
	_ Fallback = &fallbackBatch{}
	_ Helper   = &helperBatch{}
)

func handlers(f *fallbackBatch) (thread.WorkerHandler, func()) {
	var fn func()
	var levellock int
	var loop bool
	return func(w *thread.Worker) {
			stop := w.StopC()
			for {
				select {
				case <-stop:
					return
				case fn = <-f.h:
					f.mu.Lock() // exclusive lock
					levellock = 2
					fn()
					loop = true
					for loop {
						select {
						case fn = <-f.h:
							fn()
						default:
							loop = false
						}
					}
					f.mu.Unlock()
					levellock = 0
				}
			}
		}, func() {
			if levellock == 2 {
				f.mu.Unlock()
			}
		}
}

func NewBatchFallback(size int) *fallbackBatch {
	if size < 1 {
		size = 1
	}
	f := &fallbackBatch{
		h: make(chan func(), size),
	}
	f.w = thread.NewWorker(handlers(f))

	return f
}

func (f *fallbackBatch) Close() error {
	return f.w.Close()
}

func (f *fallbackBatch) Wait() {
	f.w.Wait()
}

func (f *fallbackBatch) Do(sharedHandler func(Helper)) {
	f.mu.RLock() // shared lock
	defer f.mu.RUnlock()

	sharedHandler((*helperBatch)(f))
}

func (h *helperBatch) Exclusive(exclusiveHandler func(), slowAsync func()) {
	if exclusiveHandler == nil {
		return
	}

	go func() {
		if slowAsync != nil {
			slowAsync()
		}
		select {
		case <-h.w.StopC():
		case h.h <- exclusiveHandler:
		}
	}()
}
