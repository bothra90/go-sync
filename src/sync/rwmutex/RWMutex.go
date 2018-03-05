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
// Implementation note: readSwitch should internally lock/unlock noWriters, and
// writeSwitch should lock/unlock noReaders for this rwmutex to work.
type writerPriorityRWMutex struct {
	noReaders   mutex.Mutex
	readSwitch  lightswitch.LightSwitch // Controls noWriters.
	noWriters   mutex.Mutex
	writeSwitch lightswitch.LightSwitch // Controls noReaders.
}

func (rw *writerPriorityRWMutex) Lock() {
	// Use a writeSwitch to prevent other readers from acquiring lock on noReaders
	// while allowing writers to proceed and then get blocked on noWriters.
	rw.writeSwitch.Lock()
	// Prevent other writers from proceeding beyond this point while a writer is
	// present in the critical section.
	rw.noWriters.Lock()
}

func (rw *writerPriorityRWMutex) Unlock() {
	// Allow any queued writers to go into critical section.
	rw.noWriters.Unlock()
	// Allow any queued readers to proceed if this is the last writer.
	rw.writeSwitch.Unlock()
}

func (rw *writerPriorityRWMutex) RLock() {
	// Block for all queued writers to exit. Once this returns, prevent writers
	// from entering while *this* reader has the lock.
	rw.noReaders.Lock()
	// Immediately allow writers to be queued behind this reader.
	defer rw.noReaders.Unlock()
	// Prevent any writers from acquiring writeSwitch.Lock()
	rw.readSwitch.Lock()
}

func (rw *writerPriorityRWMutex) RUnlock() {
	// Allow writers to enter the critical section if this is the last reader.
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
