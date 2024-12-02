package main

import (
	"context"
	"flag"
	"fmt"
	"io"
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
	var timeout *int64
	value := int64(10 * time.Second)
	timeout = &value
	address := ""
	if len(os.Args) == 3 { // передан timeout
		address = fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])
	} else {
		address = fmt.Sprintf("%s:%s", os.Args[2], os.Args[3])
		timeout = flag.Int64("timeout", int64(10*time.Second), "any int")
	}
	client := NewTelnetClient(address, time.Duration(*timeout), os.Stdin, os.Stdout)

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

	// Добавляем горутину для завершения по CTRL+D
	go func() {
		io.Copy(io.Discard, os.Stdin) // Блокирует чтение до конца потока
		cancelFunc()
	}()

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
