package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/sattellite/bcdb/compute/result"
	"github.com/sattellite/bcdb/config"
	"github.com/sattellite/bcdb/storage"
)

func New(logger *slog.Logger, engine storage.Engine, cfg *config.Config) (*Network, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	if cfg == nil {
		return nil, errors.New("config is required")
	}

	listener, err := net.Listen("tcp", net.JoinHostPort(cfg.Server.Address, cfg.Server.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	return &Network{
		logger: logger.With("module", "network"),
		engine: engine,
		server: listener,
		limit:  make(chan struct{}, cfg.Server.MaxClients),
	}, nil
}

type Network struct {
	logger *slog.Logger
	engine storage.Engine
	server net.Listener
	limit  chan struct{}
}

func (n *Network) Run(ctx context.Context) {
	n.logger.Info("network server starting")
	defer func() {
		n.logger.Info("network server stopped")
	}()

	for ctx.Err() == nil {
		conn, cErr := n.server.Accept()
		if cErr != nil {
			if errors.Is(cErr, net.ErrClosed) {
				break
			}

			n.logger.Error("failed accept connection", slog.Any("error", cErr))
			continue
		}

		// take a place in queue
		n.limit <- struct{}{}
		go func(c net.Conn) {
			n.handleClient(ctx, c)
			// clear a place in queue
			<-n.limit
		}(conn)

	}
}

func (n *Network) handleClient(ctx context.Context, c net.Conn) {
	defer func() {
		if msg := recover(); msg != nil {
			n.logger.Error("panic occurred", slog.Any("panic", msg))
		}

		if err := c.Close(); err != nil {
			n.logger.Error("failed close connection", slog.Any("error", err))
		}
	}()

	request := make([]byte, 1024)
	for {
		count, rErr := c.Read(request)
		if rErr != nil && errors.Is(rErr, io.EOF) {
			n.logger.Error("failed read request",
				slog.Any("error", rErr),
				slog.String("address", c.RemoteAddr().String()),
			)
			break
		}
		if count == 1024 {
			n.logger.Error("small buffer size", slog.Int("size", 1024))
			break
		}

		// process client input
		q, pErr := n.Parse(string(request[:count]))
		if pErr != nil {
			n.logger.Error("failed to parse command", slog.Any("error", pErr))
			wErr := n.Print(c, result.Result{Value: pErr.Error()})
			if wErr != nil {
				n.logger.Error("failed send data to connection", slog.Any("error", wErr), slog.Int("place", 1))
				break
			}
			continue
		}

		// handle client input
		res, hErr := n.Handle(ctx, *q)
		if hErr != nil {
			n.logger.Error("failed to handle command", slog.Any("error", hErr))
			wErr := n.Print(c, result.Result{Value: hErr.Error()})
			if wErr != nil {
				n.logger.Error("failed send data to connection", slog.Any("error", wErr), slog.Int("place", 2))
				break
			}
			continue
		}
		// send result to client
		wErr := n.Print(c, res)
		if wErr != nil {
			n.logger.Error("failed send data to connection", slog.Any("error", wErr), slog.Int("place", 3))
			break
		}
	}
}
