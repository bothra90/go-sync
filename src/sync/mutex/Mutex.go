package mutex

import "sync/semaphore"

type Mutex interface {
	Lock()
	Unlock()
}

// Implementation notes:
// A simple implementation of a mutex using only 1 semaphore is susceptible to
// thread starvation. The following technique was provided by Joseph Morris in
// "A starvation-free solution to the mutual exclusion problem".
// The implementation divides the mutex locking into 3 phases: room1, room2 and
// the critical section. Threads go through "turnstiles" from one room to the
// next (t1 for room1 to room2, and t2 for room2 to critical section). When
// a thread arrives, it goes into room1. If t1 is unlocked, it then goes into
// room2. When a thread arrives in room2, it checks if there are any other
// threads in room1. If there are, it loads t1 to let it enter room1. Otherwise,
// it loads t2 to let itself enter the critical section. When a thread exits the
// critical section, it first checks if there are threads in room2 waiting to
// enter the critical section. If there are, it loads t2 to let one of them in.
// Otherwise it loads t1 to allow threads in room1 to enter room2.

type mutex struct {
	room1 int
	room2 int
	t1    semaphore.Semaphore // used to control entry into room2
	t2    semaphore.Semaphore // used to control entry into the critical section.
	guard semaphore.Semaphore
}

func New() *mutex {
	// Ignore error.
	t1, _ := semaphore.New(1)
	t2, _ := semaphore.New(0)
	guard, _ := semaphore.New(1)
	return &mutex{
		room1: 0,
		room2: 0,
		t1:    t1,
		t2:    t2,
		guard: guard,
	}
}

func (mu *mutex) Lock() {
	mu.guard.Wait(1)
	// Enter room1.
	mu.room1 += 1
	mu.guard.Signal(1)
	// Enter the second room if t1 is open. Close t1.
	mu.t1.Wait(1)
	mu.room2 += 1
	mu.guard.Wait(1)
	mu.room1 -= 1
	if mu.room1 == 0 {
		// If we were the last thread in room1, we keep t1 closed and open t2.
		// This allows threads in room2 to start entering the critical section.
		mu.t2.Signal(1)
	} else {
		// Re-open t1 to allow threads to wait in room2.
		mu.t1.Signal(1)
	}
	mu.guard.Signal(1)
	// Enter the critical section if t2 is open.
	mu.t2.Wait(1)
	mu.room2 -= 1
}

func (mu *mutex) Unlock() {
	if mu.room2 == 0 {
		// If we were the last thread in room2, keep t2 closed, and open t1 to allow
		// threads in room1 to enter room2.
		mu.t1.Signal(1)
	} else {
		mu.t2.Signal(1)
	}
}
