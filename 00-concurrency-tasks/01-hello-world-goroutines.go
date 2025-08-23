package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Hello World from Goroutine")
	}()

	wg.Wait()
	fmt.Println("Done from main thread")
}
