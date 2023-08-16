package interview

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
)

var wg sync.WaitGroup

//Sometimes we’d like our Go programs to intelligently handle Unix signals.
//For example, we might want a server to gracefully shutdown when it receives a SIGTERM, or a command-line tool to
//stop processing input if it receives a SIGINT. Here’s how to handle signals in Go with channels.

func TestSignal(t *testing.T) {
	wg.Add(1)

	go func() {
		//Go signal notification works by sending os.Signal values on a channel.
		//We’ll create a channel to receive these notifications. Note that this channel should be buffered.
		chanel := make(chan os.Signal, 1)

		//signal.Notify registers the given channel to receive notifications of the specified signals
		signal.Notify(chanel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		fmt.Printf("goroutine 1 receive a signal : %v\n\n", <-chanel)
		wg.Done()
	}()

	wg.Wait()
}
