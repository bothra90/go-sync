package rwmutex

import (
	"sync/lightswitch"
	"sync/mutex"
)

// Implementaton notes:
// We implement the RWMutex using a resource_mutex, and a lightswitch to share
// access to the "resource" b/w readers. A writer will directly lock the
// resource_mutex, preventing other readers and writers from locking it.
// Readers, on the other hand, share access to the resource using a lightswitch
// which in turn locks and unlocks the resource when the first reader enters and
// last reader leaves respectively. In addition, we use a "queue" mutex to
// prevent starvation.

type RWMutex interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

type rwmutex struct {
	resource_mutex mutex.Mutex
	readers_switch lightswitch.LightSwitch
	queue          mutex.Mutex
}

func (rw *rwmutex) Lock() {
	// Lock the queue. This queues any incoming readers and writers after this
	// writer to block behind this writer.
	rw.queue.Lock()
	defer rw.queue.Unlock()
	// Lock the resource directly.
	rw.resource_mutex.Lock()
}

func (rw *rwmutex) Unlock() {
	// Unlock resource mutex.
	rw.resource_mutex.Unlock()
}

func (rw *rwmutex) RLock() {
	// Lock the queue. This forces readers to queue behind other readers and
	// writers. If no writer is queued, the readers should be able to proceed one-
	// by-one into the readers_switch.Lock().
	rw.queue.Lock()
	defer rw.queue.Unlock()
	// Lock resource via lightswitch.
	rw.readers_switch.Lock()
}

func (rw *rwmutex) RUnlock() {
	// Unlock resource via lightswitch.
	rw.readers_switch.Unlock()
}

// A writerPriorityRWMutex gives priority to writers over readers by allowing
// any waiting writer to proceed before any waiting reader.
type writerPriorityRWMutex struct {
	noReaders   mutex.Mutex
	readSwitch  lightswitch.LightSwitch // Controls noWriters.
	noWriters   mutex.Mutex
	writeSwitch lightswitch.LightSwitch // Controls noReaders.
}

func (rw *writerPriorityRWMutex) Lock() {
	rw.writeSwitch.Lock()
	rw.noWriters.Lock()
}

func (rw *writerPriorityRWMutex) Unlock() {
	rw.noWriters.Unlock()
	rw.writeSwitch.Unlock()
}

func (rw *writerPriorityRWMutex) RLock() {
	rw.noReaders.Lock()
	defer rw.noReaders.Unlock()
	rw.readSwitch.Lock()
}

func (rw *writerPriorityRWMutex) RUnlock() {
	rw.readSwitch.Unlock()
}

func New() *rwmutex {
	mu := mutex.New()
	return &rwmutex{
		resource_mutex: mu,
		readers_switch: lightswitch.New(mu),
		queue:          mutex.New(),
	}
}

//func New() RWMutex {
//readMutex := mutex.New()
//writeMutex := mutex.New()
//return &writerPriorityRWMutex{
//noReaders:   readMutex,
//noWriters:   writeMutex,
//readSwitch:  lightswitch.New(writeMutex),
//writeSwitch: lightswitch.New(readMutex),
//}
//}
