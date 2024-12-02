package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet host port")
		os.Exit(1)
	}

	address := fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])
	client := NewTelnetClient(address, 10*time.Second, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Connection error: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	ctx, cancelFunc := context.WithCancel(context.Background())
	go worker(client.Receive, cancelFunc)
	go worker(client.Send, cancelFunc)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigCh:
		cancelFunc()
		signal.Stop(sigCh)
		return

	case <-ctx.Done():
		close(sigCh)
		return
	}
}

func worker(handler func() error, cancelFunc context.CancelFunc) {
	if err := handler(); err != nil {
		cancelFunc()
	}
}
