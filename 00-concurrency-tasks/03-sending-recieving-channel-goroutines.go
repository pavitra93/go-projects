package main

import (
	"fmt"
	"sync"
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
	Channel := make(chan int)

	//Producer
	go producer(Channel, wg)

	//Consumer
	consumer(Channel)

	wg.Wait()
	fmt.Println("Done from main thread")
}
