// Testing notes: We rely heavily on the go race detetor to test our
// implementation.

package rwmutex

import (
	"fmt"
	"sync"
)

func ExampleReads() {
	mu := New()
	x := 5
	// Acquire read lock
	mu.RLock()
	defer mu.RUnlock()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Acquire read lock. This should be OK.
		mu.RLock()
		defer mu.RUnlock()
		fmt.Printf("x is: %d\n", x)
	}()
	wg.Wait()
	// Output:
	// x is: 5
}

func ExampleReadWrite() {
	mu := New()
	x := 5
	// Acquire read lock
	mu.Lock()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Acquire read lock. This should block till write Lock is released.
		mu.RLock()
		defer mu.RUnlock()
		fmt.Printf("x is: %d\n", x)
	}()
	x = 7
	// Release write lock.
	mu.Unlock()
	wg.Wait()
	// Output:
	// x is: 7
}
