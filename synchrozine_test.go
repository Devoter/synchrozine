package synchrozine_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/Devoter/synchrozine/v3"
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

	synchro.StartupSync(context.TODO())

	injectErr := fmt.Errorf("stop goroutines")

	synchro.Inject(injectErr)

	err := synchro.Sync(context.TODO())
	if injectErr != err {
		t.Fatalf("Expected error is [%v], got [%v]\n", injectErr, err)
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

	synchro.StartupSync(context.TODO())

	injectErr := fmt.Errorf("stop goroutines")
	injectErr2 := fmt.Errorf("stop 2")

	synchro.Inject(injectErr)
	synchro.Inject(injectErr2)

	err := synchro.Sync(context.TODO())
	if injectErr != err {
		t.Fatalf("Expected error is [%v], got [%v]\n", injectErr, err)
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

	synchro.StartupSync(context.TODO())

	err := synchro.Sync(context.TODO())
	if injectErr != err {
		t.Fatalf("Expected error is [%v], got [%v]\n", injectErr, err)
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

	synchro.StartupSync(context.TODO())

	go func() {
		synchro.Inject(injectErr)
	}()

	err := synchro.Sync(context.TODO())
	if injectErr != err {
		t.Fatalf("Expected error is [%v], got [%v]\n", injectErr, err)
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

	synchro.StartupSync(context.TODO())

	synchro.Inject(nil)

	err := synchro.Sync(context.TODO())
	if err != nil {
		t.Fatalf("Expected error is [%v], got [%v]\n", nil, err)
	}

	if f1Failed {
		t.Fatalf("First goroutine synchonization failed\n")
	}
}

func TestStartupCanceled(t *testing.T) {
	synchro := New()

	synchro.Add()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := synchro.StartupSync(ctx)
	if context.Canceled != err {
		t.Fatalf("Expected error is [%v], got [%v]\n", context.Canceled, err)
	}
}

func TestSyncTimeoutExceeded(t *testing.T) {
	synchro := New()

	synchro.Add()
	go func() {
		<-synchro.Append()
	}()

	synchro.StartupSync(context.TODO())

	injectErr := fmt.Errorf("stop goroutines")

	synchro.Inject(injectErr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := synchro.Sync(ctx)
	if context.DeadlineExceeded != err {
		t.Fatalf("Expected error is [%v], got [%v]\n", context.DeadlineExceeded, err)
	}
}
