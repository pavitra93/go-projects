package main

import (
	"fmt"
	"sync"
)

type SafeCounter struct {
	counter map[string]int
	mutex   sync.Mutex
}

func (sc *SafeCounter) Inc(key string) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.counter[key]++
}

func (sc *SafeCounter) Value(key string) int {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	return sc.counter[key]
}

func main() {
	fmt.Println("Hello from Mutexes")
	wg := &sync.WaitGroup{}
	c := SafeCounter{
		counter: make(map[string]int),
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc("key")
		}()
	}

	wg.Wait()
	fmt.Println("Done from main thread")
	fmt.Println("final count: ", c.Value("key"))
}
