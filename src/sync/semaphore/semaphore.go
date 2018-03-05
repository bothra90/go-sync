package semaphore

import (
	"fmt"
	"sync"
)

type Semaphore interface {
	Signal(i int)
	Wait(i int)
}

// Implementation note:
// In the initial implementation, we tried using a buffered channel as
// recommended by http://www.golangpatterns.info/concurrency/semaphores.
// However, we quickly realized that it provides the wrong abstraction, since
// semaphores don't have their initial value as their limit, which would be the
// case if we use a buffered channel. We therefore chose to go with an int value
// combined with a sync.Cond which is used by Signal() to notify any go routines
// blocked on Wait().
type semaphore struct {
	value int
	cond  *sync.Cond
}

func New(n int) (*semaphore, error) {
	if n < 0 {
		return nil, fmt.Errorf("Cannot create a semaphore of initial value: %d", n)
	}
	return &semaphore{
		value: n,
		cond:  sync.NewCond(&sync.Mutex{}),
	}, nil
}

func (s *semaphore) Signal(i int) {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()
	s.value += i
	// This might not be very efficient if the value of i is enough to satisfy
	// only 1 out of N possible waiters.
	s.cond.Broadcast()
}

func (s *semaphore) Wait(i int) {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()
	for s.value < i {
		s.cond.Wait()
	}
	s.value -= i
}
