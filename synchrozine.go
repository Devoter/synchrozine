package synchrozine

import (
	"context"
	"sync"
	"sync/atomic"
)

// Synchrozine is an instrument for synchronization of multiple goroutines over a single channel.
// It provides the main channel (`chan error`), as well as tools for complete synchronization
// and receives a channels list to send finish signals to goroutines.
type Synchrozine struct {
	message        chan error
	receivers      []chan<- bool
	counter        int64 // goroutines counter
	startCounter   int64 // goroutines startup
	counterCh      chan bool
	startCounterCh chan bool
	mx             sync.Mutex
	startMX        sync.Mutex
	injected       int32
}

// New creates a initialized instance of Synchrozine.
func New() *Synchrozine {
	return &Synchrozine{
		message:        make(chan error, 1),
		receivers:      make([]chan<- bool, 0),
		counterCh:      make(chan bool, 1),
		startCounterCh: make(chan bool, 1),
	}
}

// Sync waits for a sync message, sends messages to receivers, and waits for goroutines completion.
func (s *Synchrozine) Sync(ctx context.Context) error {
	message := <-s.message

	for _, channel := range s.receivers {
		channel <- true
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.counterCh:
		return message
	}
}

// Inject sends a sync message.
func (s *Synchrozine) Inject(err error) {
	if atomic.LoadInt32(&s.injected) == 0 {
		atomic.StoreInt32(&s.injected, 1)
		s.message <- err
	}
}

// StartupSync waits for all appended goroutines to start.
func (s *Synchrozine) StartupSync(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.startCounterCh:
		return nil
	}
}

// Append creates a buffered receiver channel, adds it to the receivers list and returns for as read-only channel.
func (s *Synchrozine) Append() <-chan bool {
	receiver := make(chan bool, 1)
	s.receivers = append(s.receivers, receiver)
	s.startMX.Lock()
	s.startCounter--
	counter := s.startCounter
	s.startMX.Unlock()

	if counter <= 0 {
		s.startCounterCh <- true
	}

	return receiver
}

// Add increments a counter of the internal WaitGroup  (see `sync.WaitGroup`).
func (s *Synchrozine) Add() {
	s.mx.Lock()
	s.counter++
	s.mx.Unlock()

	s.startMX.Lock()
	s.startCounter++
	s.startMX.Unlock()
}

// Done decrements a counte of the internal WaitGroup (see `sync.WaitGroup`).
func (s *Synchrozine) Done() {
	s.mx.Lock()
	s.counter--
	counter := s.counter
	s.mx.Unlock()

	if counter <= 0 {
		s.counterCh <- true
	}
}
