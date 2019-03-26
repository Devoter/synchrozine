package synchrozine

import "sync"

// Synchrozine is a complex solution for synchronization of multiple goroutines over a single channel.
// It contains the main channel (`chan error`), WaitGroup for complete synchronization and receivers channels list
// to send finish signals to goroutines.
type Synchrozine struct {
	message   chan error
	receivers []chan bool
	wg        *sync.WaitGroup
}

// NewSynchrozine is Synchrozine constructor.
func NewSynchrozine() *Synchrozine {
	return &Synchrozine{message: make(chan error), receivers: make([]chan bool, 0), wg: &sync.WaitGroup{}}
}

// Sync waits for a sync message, sends messages to receivers, and waits for goroutines completion.
func (synchrozine *Synchrozine) Sync() error {
	message := <-synchrozine.message

	for _, channel := range synchrozine.receivers {
		channel <- true
	}

	synchrozine.wg.Wait()

	return message
}

// Inject sends a sync message.
func (synchrozine *Synchrozine) Inject(err error) {
	synchrozine.message <- err
}

// Append adds a receiver channel to the receivers list.
func (synchrozine *Synchrozine) Append(receiver chan bool) {
	synchrozine.receivers = append(synchrozine.receivers, receiver)
}

// Add increments a counter of the internal WaitGroup  (see `sync.WaitGroup`).
func (synchrozine *Synchrozine) Add() {
	synchrozine.wg.Add(1)
}

// Done decrements a counte of the internal WaitGroup (see `sync.WaitGroup`).
func (synchrozine *Synchrozine) Done() {
	synchrozine.wg.Done()
}
