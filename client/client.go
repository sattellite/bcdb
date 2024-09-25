package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/sattellite/bcdb/config"
	"github.com/sattellite/bcdb/logger"
)

func New(cfg *config.Config) (*Client, error) {
	l := logger.WithScope("client")
	conn, err := net.Dial("tcp", net.JoinHostPort(cfg.Server.Address, strconv.Itoa(cfg.Server.Port)))
	if err != nil {
		l.Error("failed to dial", slog.Any("error", err))
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return &Client{
		conn:   conn,
		logger: l,
		done:   make(chan struct{}),
	}, nil
}

type Client struct {
	conn   net.Conn
	logger *slog.Logger

	done   chan struct{}
	closed atomic.Bool
}

func (c *Client) Run(ctx context.Context) {
	defer c.Close()

	l := c.logger.With("method", "Run")
	l.Debug("client started", slog.String("address", c.conn.RemoteAddr().String()))

	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, _ = fmt.Fprint(os.Stdout, "Enter command: ")
		scanner.Scan()
		command := strings.TrimSpace(strings.Trim(scanner.Text(), "\n"))
		l.Debug("got command", slog.String("command", command))
		if command == "" {
			continue
		}

		if command == "exit" {
			return
		}

		response, err := c.Request([]byte(command))
		if err != nil {
			l.Error("failed to request", slog.Any("error", err))
			return
		}
		_, _ = fmt.Fprintf(os.Stdout, "Response: %s", response)
	}
}

func (c *Client) Request(data []byte) ([]byte, error) {
	l := c.logger.With("method", "Request")
	n, err := c.conn.Write(data)
	if err != nil {
		l.Error("failed to write", slog.Any("error", err))
		return nil, fmt.Errorf("failed to write: %w", err)
	}
	l.Debug("written", slog.Int("bytes", n))

	response := make([]byte, 1024)
	cnt, rErr := c.conn.Read(response)
	if rErr != nil && !errors.Is(rErr, io.EOF) {
		l.Error("failed to read", slog.Any("error", rErr))
		return nil, fmt.Errorf("failed to read: %w", rErr)
	}
	l.Debug("read", slog.Int("bytes", cnt))

	if cnt == 1024 {
		l.Warn("response too long")
	}

	return response[:cnt], nil
}

func (c *Client) Done() <-chan struct{} {
	return c.done
}

func (c *Client) Close() {
	_ = c.conn.Close()
	if c.closed.CompareAndSwap(false, true) {
		c.done <- struct{}{}
		close(c.done)
	}
}
