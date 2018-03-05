package rwmutex

import (
	"sync/lightswitch"
	"sync/mutex"
)

type RWMutex interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

type rwmutex struct {
	resource_mutex mutex.Mutex
	readers_switch lightswitch.LightSwitch
}

func New() *rwmutex {
	mu := mutex.New()
	return &rwmutex{
		resource_mutex: mu,
		readers_switch: lightswitch.New(mu),
	}
}

func (rw *rwmutex) Lock() {
	// Lock the resource directly.
	rw.resource_mutex.Lock()
}

func (rw *rwmutex) Unlock() {
	// Unlock resource mutex.
	rw.resource_mutex.Unlock()
}

func (rw *rwmutex) RLock() {
	// Lock resource via lightswithc.
	rw.readers_switch.Lock()
}

func (rw *rwmutex) RUnlock() {
	// Unlock resource via lightswithc.
	rw.readers_switch.Unlock()
}
