# synchrozine

Synchrozine is a complex solution for synchronization of multiple goroutines over a single channel.
It contains the main channel (`chan error`), [WaitGroup](https://golang.org/pkg/sync/#WaitGroup) for complete synchronization and receivers channels list to send finish signals to goroutines.

## Usage example

```go
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Devoter/synchrozine"
)

func main() {
	listen := flag.String("listen", ":8080", "HTTP listen address")

	flag.Parse()

	synchro := synchrozine.NewSynchrozine()

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

		finishChan := make(chan bool)
		synchro.Append(finishChan)

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

	log.Println("exit: ", synchro.Sync())
}
```