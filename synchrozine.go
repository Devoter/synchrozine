package synchrozine

import "sync"

import "sync/atomic"

// Synchrozine is a complex solution for synchronization of multiple goroutines over a single channel.
// It contains the main channel (`chan error`), WaitGroup for complete synchronization and receivers channels list
// to send finish signals to goroutines.
type Synchrozine struct {
	message   chan error
	receivers []chan<- bool
	wg        *sync.WaitGroup
	startWG   *sync.WaitGroup
	injected  int32
}

// New creates a initialized instance of Synchrozine.
func New() *Synchrozine {
	return &Synchrozine{
		message:   make(chan error, 1),
		receivers: make([]chan<- bool, 0),
		wg:        &sync.WaitGroup{},
		startWG:   &sync.WaitGroup{},
	}
}

// Sync waits for a sync message, sends messages to receivers, and waits for goroutines completion.
func (s *Synchrozine) Sync() error {
	message := <-s.message

	for _, channel := range s.receivers {
		channel <- true
	}

	s.wg.Wait()

	return message
}

// Inject sends a sync message.
func (s *Synchrozine) Inject(err error) {
	if atomic.LoadInt32(&s.injected) == 0 {
		atomic.StoreInt32(&s.injected, 1)
		s.message <- err
	}
}

// StartupSync waits for all appended goroutines to start.
func (s *Synchrozine) StartupSync() {
	s.startWG.Wait()
}

// Append creates a buffered receiver channel, adds it to the receivers list and returns for as read-only channel.
func (s *Synchrozine) Append() <-chan bool {
	receiver := make(chan bool, 1)
	s.receivers = append(s.receivers, receiver)
	s.startWG.Done()

	return receiver
}

// Add increments a counter of the internal WaitGroup  (see `sync.WaitGroup`).
func (s *Synchrozine) Add() {
	s.wg.Add(1)
	s.startWG.Add(1)
}

// Done decrements a counte of the internal WaitGroup (see `sync.WaitGroup`).
func (s *Synchrozine) Done() {
	s.wg.Done()
}
