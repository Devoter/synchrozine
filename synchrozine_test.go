package synchrozine_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/Devoter/synchrozine/v4"
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

	synchro.StartupSync(func() context.Context { return context.TODO() })

	injectErr := fmt.Errorf("stop goroutines")

	synchro.Inject(injectErr)

	err := synchro.Sync(func() context.Context { return context.TODO() })
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

func TestMultipleGoroutinesWithoutStartupSync(t *testing.T) {
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

	f3Failed := false

	synchro.Add()
	go func() {
		defer synchro.Done()

		finishChan := synchro.Append()

		select {
		case <-finishChan:
			return
		case <-time.After(dur):
			f3Failed = true
			return
		}
	}()

	injectErr := fmt.Errorf("stop goroutines")

	time.Sleep(20 * time.Millisecond)

	synchro.Inject(injectErr)

	err := synchro.Sync(func() context.Context { return context.TODO() })
	if injectErr != err {
		t.Fatalf("Expected error is [%v], got [%v]\n", injectErr, err)
	}

	if f1Failed {
		t.Fatalf("First goroutine synchonization failed\n")
	}

	if f2Failed {
		t.Fatalf("Second goroutine synchonization failed\n")
	}

	if f3Failed {
		t.Fatalf("Third goroutine synchonization failed\n")
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

	synchro.StartupSync(func() context.Context { return context.TODO() })

	injectErr := fmt.Errorf("stop goroutines")
	injectErr2 := fmt.Errorf("stop 2")

	synchro.Inject(injectErr)
	synchro.Inject(injectErr2)

	err := synchro.Sync(func() context.Context { return context.TODO() })
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

	synchro.StartupSync(func() context.Context { return context.TODO() })

	err := synchro.Sync(func() context.Context { return context.TODO() })
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

	synchro.StartupSync(func() context.Context { return context.TODO() })

	go func() {
		synchro.Inject(injectErr)
	}()

	err := synchro.Sync(func() context.Context { return context.TODO() })
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

	synchro.StartupSync(func() context.Context { return context.TODO() })

	synchro.Inject(nil)

	err := synchro.Sync(func() context.Context { return context.TODO() })
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

	makeCtx := func() context.Context {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		return ctx
	}

	err := synchro.StartupSync(makeCtx)
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

	synchro.StartupSync(func() context.Context { return context.TODO() })

	injectErr := fmt.Errorf("stop goroutines")

	synchro.Inject(injectErr)

	var cancel context.CancelFunc

	makeCtx := func() context.Context {
		var ctx context.Context
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)

		return ctx
	}

	defer func() { cancel() }()

	err := synchro.Sync(makeCtx)
	if context.DeadlineExceeded != err {
		t.Fatalf("Expected error is [%v], got [%v]\n", context.DeadlineExceeded, err)
	}
}
