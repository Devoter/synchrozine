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
	message           chan error
	receiverCh        chan struct{}
	counter           int64 // goroutines counter
	startCounter      int64 // goroutines startup
	counterCh         chan struct{}
	startCounterCh    chan struct{}
	mx                sync.RWMutex
	startMX           sync.Mutex
	injected          int32
	waitingForStartup bool
}

// New creates a initialized instance of Synchrozine.
func New() *Synchrozine {
	return &Synchrozine{
		message:        make(chan error, 1),
		receiverCh:     make(chan struct{}, 0),
		counterCh:      make(chan struct{}, 1),
		startCounterCh: make(chan struct{}, 1),
	}
}

// Sync waits for a sync message, sends messages to receivers, and waits for goroutines completion.
func (s *Synchrozine) Sync(ctxFactory func() context.Context) error {
	message := <-s.message

	close(s.receiverCh)

	ctx := ctxFactory()

	s.mx.RLock()
	counter := s.counter
	s.mx.RUnlock()

	if counter <= 0 {
		return message
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
func (s *Synchrozine) StartupSync(ctxFactory func() context.Context) error {
	s.startMX.Lock()

	if s.startCounter <= 0 {
		s.startMX.Unlock()
		ctxFactory() // initialize context to prevent runtime errors of canceling of nil context
		return nil
	}

	s.waitingForStartup = true
	s.startMX.Unlock()

	ctx := ctxFactory()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.startCounterCh:
		return nil
	}
}

// Append decrements a startup counter and returns a read-only channel.
func (s *Synchrozine) Append() <-chan struct{} {
	s.startMX.Lock()
	s.startCounter--
	counter := s.startCounter
	waitingForStartup := s.waitingForStartup
	s.startMX.Unlock()

	if counter <= 0 && waitingForStartup {
		s.startCounterCh <- struct{}{}
	}

	return s.receiverCh
}

// Add increments a counter of controlled goroutines.
func (s *Synchrozine) Add() {
	s.mx.Lock()
	s.counter++
	s.mx.Unlock()

	s.startMX.Lock()
	s.startCounter++
	s.startMX.Unlock()
}

// AddMany increments a counter of controlled goroutines on the specified value.
func (s *Synchrozine) AddMany(count int) {
	s.mx.Lock()
	s.counter += int64(count)
	s.mx.Unlock()

	s.startMX.Lock()
	s.startCounter += int64(count)
	s.startMX.Unlock()
}

// Done decrements a counter of controlled goroutines.
func (s *Synchrozine) Done() {
	s.mx.Lock()
	s.counter--
	counter := s.counter
	s.mx.Unlock()

	if counter <= 0 {
		s.counterCh <- struct{}{}
	}
}
