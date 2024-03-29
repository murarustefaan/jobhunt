package main

import (
	"github.com/murarustefaan/jobhunt/cmd/work"
	"sync"
)

const (
	NumWorkers = 2
)

func main() {
	//var wq queue.WorkScheduler[string] = queue.NewMemoryQueue[string]()
	inputs := make(chan string)
	results := make(chan string, NumWorkers)
	done := make(chan bool, NumWorkers)
	wg := &sync.WaitGroup{}

	go func() {
		for {
			select {
			case result := <-results:
				println("result:", result)
			}
		}
	}()

	wg.Add(1)
	for i := 0; i < NumWorkers; i++ {
		w := work.NewWorker(wg)
		go w.Go(inputs, results, done)
	}

	for _, endpoint := range STARTING_ENDPOINTS {
		inputs <- endpoint
	}
	wg.Done()

	wg.Wait()
	close(inputs)
	for i := 0; i < NumWorkers; i++ {
		done <- true
	}
}
