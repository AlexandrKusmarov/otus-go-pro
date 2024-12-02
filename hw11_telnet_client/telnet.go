package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	c.conn = conn
	return nil
}

func (c *telnetClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *telnetClient) Send() error {
	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		_, err := fmt.Fprintln(c.conn, scanner.Text())
		if err != nil {
			return fmt.Errorf("send error: %w", err)
		}
	}
	return scanner.Err()
}

func (c *telnetClient) Receive() error {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		_, err := fmt.Fprintln(c.out, scanner.Text())
		if err != nil {
			return fmt.Errorf("receive error: %w", err)
		}
	}
	return scanner.Err()
}
