package main

import (
	"fmt"
	"sync"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		res := job * 2
		results <- res
		fmt.Printf("Job done by %d and results is %d\n", id, job)
	}
}

func main() {
	fmt.Println("Hello to Worker Pool")
	jobs := make(chan int, 10)
	results := make(chan int, 10)
	wg := &sync.WaitGroup{}
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, results, wg)
	}

	for i := 1; i <= 10; i++ {
		jobs <- i
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Println(res)
	}

}
