package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(Channel chan int, wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		Channel <- i
	}
	close(Channel)
	wg.Done()
}

func consumer(Channel chan int) {
	for value := range Channel {
		fmt.Println(value)
	}
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	start := time.Now()
	Channel := make(chan int, 10)

	//Producer
	go producer(Channel, wg)

	//Consumer
	consumer(Channel)

	wg.Wait()
	fmt.Println("Done from main thread")
	end := time.Now()
	fmt.Println(end.Sub(start))
}
