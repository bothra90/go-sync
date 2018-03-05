package semaphore

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestLoadedSemaphoreCreation(t *testing.T) {
	assert := assert.New(t)
	sem, err := New(1)
	assert.Nil(err)
	sem.Wait(1)
}

func ExampleSignaling() {
	sem, _ := New(0)
	wg := sync.WaitGroup{}
	wg.Add(2)
	// Both waiters should wait for the signal.
	go func() {
		defer wg.Done()
		sem.Wait(1)
		fmt.Println("Done waiting")
	}()
	go func() {
		defer wg.Done()
		sem.Wait(1)
		fmt.Println("Done waiting")
	}()
	sem.Signal(2)
	fmt.Println("Done signaling")
	wg.Wait()
	// Output:
	// Done signaling
	// Done waiting
	// Done waiting
}
