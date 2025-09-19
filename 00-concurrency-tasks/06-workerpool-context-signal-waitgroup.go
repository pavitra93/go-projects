// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"fmt"
	"sync"
)

type Job struct {
	ID    int
	Value int
}

type Result struct {
	JobID  int
	Output int
	Error  error
}

func main() {
	workers := 5
	numJobs := 10

	jobs := make(chan Job, 10)
	results := make(chan Result, 10)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go workerPool(i, jobs, results, &wg, ctx)
	}

	go producejobs(numJobs, jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.Error != nil {
			fmt.Printf("Job Printed %d", result.Error)
		}
		fmt.Printf("Job Printed %d\n", result.JobID)
	}
}

func workerPool(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, ctx context.Context) {

	defer wg.Done()

	for {
		select {

		case <-ctx.Done():

			fmt.Printf("Job cancelled")
			return

		case job, ok := <-jobs:

			if !ok {
				return

			}

			output := job.Value * 2
			select {
			case results <- Result{JobID: job.ID, Output: output, Error: nil}:
			case <-ctx.Done():
				// If context cancelled while trying to send, abandon
				fmt.Printf("worker %d: context cancelled while sending result\n", id)
				return
			}
		}
	}

}

func producejobs(n int, jobs chan<- Job) {
	for i := 1; i <= n; i++ {
		jobs <- Job{ID: i, Value: i}

	}

	close(jobs)

}
