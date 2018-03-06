// A barrier can be used to synchronize various go-routines in the following
// pattern:
//     | | | | |
//     | | | | |
//     | | | | |
//     --------- Enter Barrier.
//         |
//         |
//         |
//         |
//         |
//     --------- Exit Barrier.
//     | | | | |
//     | | | | |
//     | | | | |
// It is guaranteed that all go-routines will enter and exit the barrier
// "simultaneosly".

package barrier

import (
	"fmt"
	"sync/semaphore"
)

type Barrier interface {
	Wait()
	Enter()
	Exit()
}

type barrier struct {
	N              int
	entryTurnstile semaphore.Semaphore
	exitTurnstile  semaphore.Semaphore
	mu             semaphore.Semaphore
	count          int
}

func New(N int) (Barrier, error) {
	if N <= 0 {
		return nil, fmt.Errorf("Cannot create a barrier for %d go-routines", N)
	}
	// entryTurnstile is used to control entry into the barrier.
	entryTurnstile, _ := semaphore.New(0)
	// entryTurnstile is used to control exit from the barrier.
	exitTurnstile, _ := semaphore.New(0)
	// Mutex is used to get/set count.
	mu, _ := semaphore.New(1)
	return &barrier{
		N:              N,
		entryTurnstile: entryTurnstile,
		exitTurnstile:  exitTurnstile,
		mu:             mu,
		count:          0,
	}, nil
}

func (b *barrier) Enter() {
	b.mu.Wait(1)
	b.count += 1
	if b.count == b.N {
		// Load the entryTurnstile to allow all waiting threads to pass.
		b.entryTurnstile.Signal(b.N)
	}
	b.mu.Signal(1)
	b.entryTurnstile.Wait(1)
	return
}

func (b *barrier) Exit() {
	b.mu.Wait(1)
	b.count -= 1
	if b.count == 0 {
		// Load the entryTurnstile to allow all waiting threads to pass.
		b.exitTurnstile.Signal(b.N)
	}
	b.mu.Signal(1)
	b.exitTurnstile.Wait(1)
	return
}

func (b *barrier) Wait() {
	b.Enter()
	b.Exit()
}
