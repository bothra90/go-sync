package lightswitch

import (
	"sync/mutex"
	"sync/semaphore"
)

type LightSwitch interface {
	Lock()
	Unlock()
}

type lightswitch struct {
	counter int
	mutex   mutex.Mutex
	status  semaphore.Semaphore // The switch.
}

func New(status semaphore.Semaphore) LightSwitch {
	return &lightswitch{
		counter: 0,
		mutex:   mutex.New(),
		status:  status,
	}
}

func (ls *lightswitch) Lock() {
	ls.mutex.Lock()
	if ls.counter == 0 {
		// If we are the first one in, turn on the light.
		ls.status.Wait(1)
	}
	ls.counter += 1
	defer ls.mutex.Unlock()
}

func (ls *lightswitch) Unlock() {
	ls.mutex.Lock()
	ls.counter -= 1
	if ls.counter == 0 {
		// If we are the last one out, turn off the light.
		ls.status.Signal(1)
	}
	defer ls.mutex.Unlock()
}
