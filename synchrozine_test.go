package synchrozine_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/Devoter/synchrozine/v2"
)

func TestSyncTwoGoroutines(t *testing.T) {
	synchro := New()

	dur := 50 * time.Millisecond
	f1Failed := false

	synchro.Add()
	go func() {
		defer synchro.Done()

		finishChan := synchro.Append()

		select {
		case <-finishChan:
			return
		case <-time.After(dur):
			f1Failed = true
			return
		}
	}()

	f2Failed := false

	synchro.Add()
	go func() {
		defer synchro.Done()

		finishChan := synchro.Append()

		select {
		case <-finishChan:
			return
		case <-time.After(dur):
			f2Failed = true
			return
		}
	}()

	synchro.StartupSync()

	injectErr := fmt.Errorf("stop goroutines")

	synchro.Inject(injectErr)

	err := synchro.Sync()
	if injectErr != err {
		t.Fatalf("Expected error is \"%v\", got \"%v\"\n", injectErr, err)
	}

	if f1Failed {
		t.Fatalf("First goroutine synchonization failed\n")
	}

	if f2Failed {
		t.Fatalf("Second goroutine synchonization failed\n")
	}
}

func TestMultipleInjections(t *testing.T) {
	synchro := New()

	dur := 50 * time.Millisecond
	f1Failed := false

	synchro.Add()
	go func() {
		defer synchro.Done()

		finishChan := synchro.Append()

		select {
		case <-finishChan:
			return
		case <-time.After(dur):
			f1Failed = true
			return
		}
	}()

	synchro.StartupSync()

	injectErr := fmt.Errorf("stop goroutines")
	injectErr2 := fmt.Errorf("stop 2")

	synchro.Inject(injectErr)
	synchro.Inject(injectErr2)

	err := synchro.Sync()
	if injectErr != err {
		t.Fatalf("Expected error is \"%v\", got \"%v\"\n", injectErr, err)
	}

	if f1Failed {
		t.Fatalf("First goroutine synchonization failed\n")
	}
}

func TestInjectBeforeStartupSync(t *testing.T) {
	synchro := New()

	dur := 50 * time.Millisecond
	f1Failed := false

	synchro.Add()
	go func() {
		defer synchro.Done()

		finishChan := synchro.Append()

		select {
		case <-finishChan:
			return
		case <-time.After(dur):
			f1Failed = true
			return
		}
	}()

	injectErr := fmt.Errorf("stop goroutines")

	synchro.Inject(injectErr)

	synchro.StartupSync()

	err := synchro.Sync()
	if injectErr != err {
		t.Fatalf("Expected error is \"%v\", got \"%v\"\n", injectErr, err)
	}

	if f1Failed {
		t.Fatalf("First goroutine synchonization failed\n")
	}
}

func TestInjectFromAnotherGoroutine(t *testing.T) {
	synchro := New()

	dur := 50 * time.Millisecond
	f1Failed := false

	synchro.Add()
	go func() {
		defer synchro.Done()

		finishChan := synchro.Append()

		select {
		case <-finishChan:
			return
		case <-time.After(dur):
			f1Failed = true
			return
		}
	}()

	f2Failed := false

	synchro.Add()
	go func() {
		defer synchro.Done()

		finishChan := synchro.Append()

		select {
		case <-finishChan:
			return
		case <-time.After(dur):
			f2Failed = true
			return
		}
	}()

	injectErr := fmt.Errorf("stop goroutines")

	synchro.StartupSync()

	go func() {
		synchro.Inject(injectErr)
	}()

	err := synchro.Sync()
	if injectErr != err {
		t.Fatalf("Expected error is \"%v\", got \"%v\"\n", injectErr, err)
	}

	if f1Failed {
		t.Fatalf("First goroutine synchonization failed\n")
	}

	if f2Failed {
		t.Fatalf("Second goroutine synchonization failed\n")
	}
}

func TestInjectNil(t *testing.T) {
	synchro := New()

	dur := 50 * time.Millisecond
	f1Failed := false

	synchro.Add()
	go func() {
		defer synchro.Done()

		finishChan := synchro.Append()

		select {
		case <-finishChan:
			return
		case <-time.After(dur):
			f1Failed = true
			return
		}
	}()

	synchro.StartupSync()

	synchro.Inject(nil)

	err := synchro.Sync()
	if err != nil {
		t.Fatalf("Expected error is \"%v\", got \"%v\"\n", nil, err)
	}

	if f1Failed {
		t.Fatalf("First goroutine synchonization failed\n")
	}
}
