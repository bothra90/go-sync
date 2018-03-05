package barrier

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"sync/semaphore"
	"testing"
)

func TestBarrier(t *testing.T) {
	assert := assert.New(t)
	wg := sync.WaitGroup{}
	b, err := New(100)
	assert.Nil(err)
	var ops int32
	var expected int32
	expected = 100
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			t.Logf("Will wait at barier: %d\n", id)
			atomic.AddInt32(&ops, 1)
			b.Wait()
			assert.True(atomic.CompareAndSwapInt32(&ops, expected, expected))
			t.Logf("Done waiting for barrier: %d\n", id)
		}(i)
	}
	wg.Wait()
}

func ExampleDancers() {
	// assert := assert.New(t)
	followerQueue, _ := semaphore.New(0)
	leaderQueue, _ := semaphore.New(0)
	danceEndBarrier, _ := New(2)
	dancingLeaders, _ := semaphore.New(1)
	dancingFollowers, _ := semaphore.New(1)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		// Leaders.
		for i := 0; i < 5; i++ {
			followerQueue.Signal(1)
			leaderQueue.Wait(1)
			dancingLeaders.Wait(1)
			danceEndBarrier.Enter()
			fmt.Printf("Dancing Pair: %d\n", i)
			danceEndBarrier.Exit()
			dancingLeaders.Signal(1)
		}
	}()
	go func() {
		defer wg.Done()
		// Followers
		for i := 0; i < 5; i++ {
			leaderQueue.Signal(1)
			followerQueue.Wait(1)
			// Prevent future followers from proceeding.
			dancingFollowers.Wait(1)
			danceEndBarrier.Enter()
			fmt.Printf("Dancing Pair: %d\n", i)
			danceEndBarrier.Exit()
			dancingFollowers.Signal(1)
		}
	}()
	wg.Wait()
	// Output:
	// Dancing Pair: 0
	// Dancing Pair: 0
	// Dancing Pair: 1
	// Dancing Pair: 1
	// Dancing Pair: 2
	// Dancing Pair: 2
	// Dancing Pair: 3
	// Dancing Pair: 3
	// Dancing Pair: 4
	// Dancing Pair: 4
}
