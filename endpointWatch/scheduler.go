package endpointWatch

import (
	"errors"
	"fmt"
	"time"

	"github.com/gammazero/workerpool"
)

type Scheduler struct {
	W    *Watch
	Quit chan struct{}
}

// NewScheduler creates a new scheduler instance with w as watch
// it also creates a quit signal channel for emergency exits
func NewScheduler(w *Watch) (*Scheduler, error) {
	sch := &Scheduler{Quit: make(chan struct{})}
	if w != nil {
		sch.W = w
		return sch, nil
	}
	return nil, errors.New("cannot create a scheduler with nil monitor")
}

// DoWithIntervals creates a ticker to the execute w.Do() every d duration
// it listens to a quit channel as well for termination signal.
func (sch *Scheduler) DoWithIntervals(d time.Duration) {
	ticker := time.NewTicker(d)
	go func() {
		for {
			select {
			case <-ticker.C:
				sch.W.Do()
			case <-sch.Quit:
				// stopping worker pool from accepting anymore jobs
				err := sch.W.Cancel()

				if err != nil {
					fmt.Println("error canceling watch on quit signal in DoWithIntervals()")
				}

				// since out mnt's worker pool is useless after cancel we instantiate another one
				sch.W.wp = workerpool.New(sch.W.workerSize)

				ticker.Stop()
				return
			}
		}
	}()
}
