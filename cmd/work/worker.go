package work

import (
	"github.com/mvdan/xurls"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
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

func (w *Worker) Go(inputs chan Url, results chan<- CrawlResponse, done <-chan bool) {
	log := w.logger
	log.Info("worker started")

	for {
		select {
		case <-done:
			log.Info("worker stopped")
			break

		case endpoint := <-inputs:
			w.logger.Info("worker processing item", "endpoint", endpoint)
			urls, err := w.crawlEndpoint(endpoint)
			if err != nil {
				log.Error("worker failed to process item", "error", err)
				log.Info("re-queueing item", "endpoint", endpoint)

				results <- CrawlResponse{
					Urls: []Url{endpoint},
					Err:  err,
				}
				continue
			}

			nextUrls := make([]Url, len(urls))
			for _, r := range urls {
				nextUrls = append(nextUrls, Url(r))
			}

			results <- CrawlResponse{
				Urls: nextUrls,
			}
		}
	}
}

func (w *Worker) crawlEndpoint(url Url) ([]string, error) {
	defer w.wg.Done()

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Get(string(url))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	results := xurls.Strict.FindAllString(string(body), -1)
	return results, nil
}
