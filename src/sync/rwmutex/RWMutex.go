package mutex

import (
	"sync/mutex"
	"sync/semaphore"
)

type RWMutex interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

type rwmutex struct {
	room_empty semaphore.Semaphore
	mu         mutex.Mutex
	readers    int
}

func New() RWMutex {
	lock := mutex.New()
	room_empty, _ := semaphore.New(1)
	return &rwmutex{
		readers:    0,
		room_empty: room_empty,
		mu:         lock,
	}
}

func (rw *rwmutex) Lock() {
	// Wait for all readers and writers to exit.
	rw.room_empty.Wait(1)
}

func (rw *rwmutex) Unlock() {
	// Signal room as empty
	rw.room_empty.Signal(1)
}

func (rw *rwmutex) RLock() {
	rw.mu.Lock()
	if rw.readers == 0 {
		rw.room_empty.Wait(1)
	}
	rw.readers += 1
	rw.mu.Unlock()
}

func (rw *rwmutex) RUnlock() {
	rw.mu.Lock()
	rw.readers -= 1
	if rw.readers == 0 {
		rw.room_empty.Signal(1)
	}
	rw.mu.Unlock()
}
