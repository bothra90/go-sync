package lightswitch

// A lightswitch allows any number of go-routines to share a resource. The idea
// is that the first go-routine to call Lock() acquires the resource, which can
// then be shared by other go-routines calling Lock(). The last go-routine to
// call Unlock() releases its lock on the resource. The metaphor we can use
// here, is that of a room. The first person to enter a room can switch on a
// light; any other persons to follow do not need to switch on the light if
// someone is already in the room. The last person to leave the room switches
// off the light.

import (
	"sync"
	"sync/mutex"
)

type LightSwitch interface {
	Lock()
	Unlock()
}

type lightswitch struct {
	counter       int
	counter_lock  mutex.Mutex
	resource_lock sync.Locker
}

func New(resource_lock sync.Locker) LightSwitch {
	return &lightswitch{
		counter:       0,
		counter_lock:  mutex.New(),
		resource_lock: resource_lock,
	}
}

func (ls *lightswitch) Lock() {
	ls.counter_lock.Lock()
	defer ls.counter_lock.Unlock()
	if ls.counter == 0 {
		// If we are the first one in, lock the resource.
		ls.resource_lock.Lock()
	}
	ls.counter += 1
}

func (ls *lightswitch) Unlock() {
	ls.counter_lock.Lock()
	defer ls.counter_lock.Unlock()
	ls.counter -= 1
	if ls.counter == 0 {
		// If we are the last one out, unlock the resource.
		ls.resource_lock.Unlock()
	}
}
