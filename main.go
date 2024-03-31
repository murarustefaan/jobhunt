package main

import (
	"github.com/murarustefaan/jobhunt/cmd/work"
	"log/slog"
	"sync"
)

const (
	NumWorkers = 2
)

func main() {
	//var wq queue.WorkScheduler[string] = queue.NewMemoryQueue[string]()
	inputs := make(chan work.Url, 1000)
	results := make(chan work.CrawlResponse, NumWorkers)
	done := make(chan bool, NumWorkers)
	wg := &sync.WaitGroup{}

	go func() {
		for {
			select {
			case result := <-results:
				if result.Err != nil {
					slog.Error("failed to crawl", "error", result.Err)
					inputs <- result.Urls[0]
					continue
				}

				for _, url := range result.Urls {
					slog.Info("scheduling next url", "url", url)
					inputs <- url
				}
			}

		}
	}()

	wg.Add(1)
	for i := 0; i < NumWorkers; i++ {
		w := work.NewWorker(wg)
		go w.Go(inputs, results, done)
	}

	for _, endpoint := range STARTING_ENDPOINTS {
		inputs <- work.Url(endpoint)
	}
	wg.Done()

	wg.Wait()
	close(inputs)
	for i := 0; i < NumWorkers; i++ {
		done <- true
	}
}
