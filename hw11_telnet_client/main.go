package main

import (
	"fmt"
	"os"
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

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	go func() {
		if err := client.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "Receive error: %v\n", err)
			os.Exit(1)
		}
	}()

	if err := client.Send(); err != nil {
		fmt.Fprintf(os.Stderr, "Send error: %v\n", err)
	}
	fmt.Fprintln(os.Stderr, "...Connection closed")
}
