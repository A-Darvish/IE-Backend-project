package endpointWatch

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/arvnd79/ie-proj/model"
	"github.com/arvnd79/ie-proj/store"
	"github.com/gammazero/workerpool"
)

type Watch struct {
	store      *store.Store
	URLs       []model.URL
	wp         *workerpool.WorkerPool
	workerSize int
}

// NewWatch creates a Watch instance with 'store' and 'url'
// it also creates a worker pool of size 'workerSize'
// if 'urls' is set to nil it will be initialized with an empty slice
func NewWatch(store *store.Store, urls []model.URL, workerSize int) *Watch {
	w := new(Watch)
	if urls == nil {
		w.URLs = make([]model.URL, 0)
	}
	w.URLs = urls
	w.store = store
	w.workerSize = workerSize
	// max number of workers (max number of threads)
	w.wp = workerpool.New(workerSize) //handles concurrency with abstraction
	return w
}

// LoadFromDatabase loads all urls(the endpoints that we want to watch them)
// from database into our watch to start working on them
// this function will replace all of saved URLs with the ones from database
func (w *Watch) LoadFromDatabase() error {
	urls, err := w.store.GetAllURLs()
	if err != nil {
		return err
	}
	w.URLs = urls
	return nil
}

// AddURL appends a slice of urls to the current list of urls
func (w *Watch) AddURL(urls []model.URL) {
	w.URLs = append(w.URLs, urls...)
}

// Cancel stops all tasks of fetching urls
// it will wait for current running jobs to finish
func (w *Watch) Cancel() error {
	w.wp.Stop()
	if !w.wp.Stopped() {
		return errors.New("could not stop the watch")
	}
	return nil
}

// Do ranges over URLs currently inside Watch instance
// and save each one's request inside database
func (w *Watch) Do() {
	var wg sync.WaitGroup

	for urlIndex := range w.URLs {
		url := w.URLs[urlIndex]
		wg.Add(1)
		w.wp.Submit(func() {
			defer wg.Done()
			w.watchURL(url)
		})
	}
	wg.Wait()
}

func (w *Watch) watchURL(url model.URL) {
	// sending request
	req, err := url.SendRequest()
	if err != nil {
		fmt.Println(err, "could not send the request")
		req = new(model.Request)
		req.UrlId = url.UrlId
		req.Result = http.StatusBadRequest
	}
	// add request to database
	if err = w.store.AddRequest(req); err != nil {
		fmt.Println(err, "could not save request to database")
	}
	// status code was other than 2XX
	if req.Result/100 != 2 {
		if err = w.store.IncrementFailed(&url); err != nil {
			fmt.Println(err, "could not increment failed times for url")
		}
	}
}
