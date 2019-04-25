package holder

import "sync"

type mutexHolder struct {
	val string
	mu  sync.Mutex
}

// NewMutexHolder returns a Holder backed by a sync.Mutex.
func NewMutexHolder() Holder {
	return &mutexHolder{}
}

func (h *mutexHolder) Get() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.val
}

func (h *mutexHolder) Set(s string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.val = s
}
