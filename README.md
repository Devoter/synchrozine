# synchrozine

[![Build Status](https://travis-ci.com/Devoter/synchrozine.svg?branch=master)](https://travis-ci.com/Devoter/synchrozine)

Synchrozine is an instrument for synchronization of multiple goroutines over a single channel.
It provides the main channel (`chan error`), as well as tools for complete synchronization and receivers a channels list to send finish signals to goroutines.

Synchrozine supports the startup synchronization and thread-safe injections.

## Installation

```sh
go mod init github.com/my/repo
go get github.com/Devoter/synchrozine/v3
```

## Usage

It is easy to use Synchrozine and simular to sync.WaitGroup. We have five lifecycle steps:

1. Goroutines declaration
2. Goroutines registration
3. Startup synchronization
4. Injection
5. Final synchronization

### Goroutines declaration

At first you should declare that a goroutine will be registered. Call `Add()` before start the goroutine.

```go
synchro.Add()
go func() {
	// do something
}()
```

### Goroutines registration

The next step is to register the declared goroutine. Call `Done()` before exit from the goroutine, get and listen a finish channel (`Append()`).

```go
go func() {
	defer synchro.Done()

	finishChan := synchro.Append()

	for {
		select {
		case <-finishChan:
			return
		case <-time.After(100 * time.Millisecond): // any event
			// do something
		}
	}
}
```

### Startup synchronization

Startup synchronization is optional, but if you want to be sure that all controlled goroutines have registered successfully call `StartSync()` in the parent goroutine.

```go
synchro.StartSync(func() context.Context { return context.TODO() })
```

### Injection

It is important that `Inject()` does not stops controlled goroutines. It just sends a message to the `Synchrozine` instance. This is a non-blocking method regardless of number of calls.

```go
synchro.Inject(fmt.Errorf("stop all"))
```

### Final synchronization

In your parent goroutine insert a `Sync()` call. This method will wait for the injection. After that, it will have broadcast finish signals to controlled goroutines and wait for synchronization. This is a blocking method.

```go
synchro.Sync(func() context.Context { return context.TODO() })
```

### Example

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Devoter/synchrozine/v4"
)

func main() {
	listen := flag.String("listen", ":8080", "HTTP listen address")

	flag.Parse()

	synchro := synchrozine.New()

	// Waiting for sigint or sigterm
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		synchro.Inject(fmt.Errorf("%s", <-c))
		log.Println("signal received")
	}()

	go func() {
		log.Println("starting server...")
		log.Println("addr", *listen)
		synchro.Inject(http.ListenAndServe(*listen, nil))
		log.Println("server stopped")
	}()

	synchro.Add()
	go func() {
		defer synchro.Done()

		finishChan := synchro.Append()

		const increment = 5
		var counter int

		for {
			// do something
			select {
			case <-finishChan:
				log.Println("finish something")
				return
			case <-time.After(increment * time.Second):
				counter += increment
				log.Printf("%d seconds left\n", counter)
			}
		}
	}()

	ctx := context.Background()
	var cancel context.CancelFunc

	makeStartupCtx := func() context.Context {
		ctx, cancel = context.WithTimeout(ctx, 10 * time.Second) // wait 10 seconds for startup

		return ctx
	}
	
	defer func() { cancel() }()

	err := synchro.StartupSync(makeStartupCtx) // optional, just if you want to be sure that all goroutines have started
	if err != nil {
		log.Printf("Startup operation failed by the reason: [%v]\n", err)
		os.Exit(1)
	}

	syncCtx = context.Background()
	var syncCancel context.CancelFunc

	makeSyncCtx := func() context.Context {
		syncCtx, syncCancel = context.WithTimeout(ctx, 10 * time.Second) // wait 10 seconds for sync

		return syncCtx
	}
	
	defer func() { syncCancel() }()

	log.Println("exit: ", synchro.Sync(makeSyncCtx))
}
```

## License

[MIT](LICENSE)
