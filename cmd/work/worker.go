package work

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

type Worker struct {
	logger slog.Logger
	wg     *sync.WaitGroup
}

func NewWorker(wg *sync.WaitGroup) *Worker {
	return &Worker{
		logger: *slog.
			New(slog.NewTextHandler(os.Stdout, nil)).
			With("worker", "worker"),
		wg: wg,
	}
}

func (w *Worker) Go(inputs chan string, results chan<- string, done <-chan bool) {
	log := w.logger
	log.Info("worker started")

	for {
		select {
		case <-done:
			log.Info("worker stopped")
			break

		case endpoint := <-inputs:
			w.logger.Info("worker processing item", "endpoint", endpoint)
			result, err := w.processItem(endpoint)
			if err != nil {
				log.Error("worker failed to process item", "error", err)
				log.Info("re-queueing item", "endpoint", endpoint)
				inputs <- endpoint
				continue
			}

			results <- result
		}
	}
}

func (w *Worker) processItem(endpoint string) (string, error) {
	rand := rand.Intn(100)

	if rand < 50 {
		return "", fmt.Errorf("random error")
	}

	return "result " + strconv.Itoa(rand), nil
}
