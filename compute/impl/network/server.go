package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"

	"github.com/sattellite/bcdb/util"

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

	listener, err := net.Listen("tcp", net.JoinHostPort(cfg.Server.Address, strconv.Itoa(cfg.Server.Port)))
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
	defer func() {
		n.logger.Info("network server stopped")
	}()

	n.logger.Info("network server listening", slog.String("address", n.server.Addr().String()))

	for ctx.Err() == nil {
		conn, cErr := n.server.Accept()
		if cErr != nil {
			if errors.Is(cErr, net.ErrClosed) {
				break
			}

			n.logger.Error("failed accept connection", slog.Any("error", cErr))
			continue
		}

		clientID := util.ClientID(conn.RemoteAddr().String())
		n.logger.Debug("accepted connection", slog.String("client", clientID))
		ctx = util.SetContextClientID(ctx, clientID)

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
	clientID := util.GetContextClientID(ctx)
	defer func() {
		if msg := recover(); msg != nil {
			n.logger.Error("panic occurred",
				slog.Any("panic", msg),
				slog.String("client", clientID),
			)
		}

		if err := c.Close(); err != nil {
			n.logger.Error("failed close connection",
				slog.Any("error", err),
				slog.String("client", clientID),
			)
		}
	}()

	n.logger.Debug("handling connection", slog.String("client", clientID))

	request := make([]byte, 1024)
	for {
		count, rErr := c.Read(request)
		if rErr != nil && errors.Is(rErr, io.EOF) {
			n.logger.Error("failed read request",
				slog.Any("error", rErr),
				slog.String("address", c.RemoteAddr().String()),
				slog.String("client", clientID),
			)
			break
		}
		if count == 1024 {
			n.logger.Error("small buffer size",
				slog.Int("size", 1024),
				slog.String("client", clientID),
			)
			break
		}

		reqID := util.RequestID()
		ctx = util.SetContextRequestID(ctx, reqID)
		n.logger.Debug("got request", slog.String("client", clientID), slog.String("request", reqID))

		// process client input
		q, pErr := n.Parse(string(request[:count]))
		if pErr != nil {
			n.logger.Error("failed to parse command",
				slog.Any("error", pErr),
				slog.String("client", clientID),
				slog.String("request", reqID),
			)
			wErr := n.Print(c, result.Result{Value: pErr.Error()})
			if wErr != nil {
				n.logger.Error("failed send data to connection",
					slog.Any("error", wErr),
					slog.String("client", clientID),
					slog.String("request", reqID),
					slog.Int("place", 1),
				)
				break
			}
			continue
		}

		// handle client input
		res := n.Handle(ctx, *q)
		if res.Error != nil {
			n.logger.Error("failed to handle command",
				slog.Any("error", res.Error),
				slog.String("client", clientID),
				slog.String("request", reqID),
			)
			wErr := n.Print(c, res)
			if wErr != nil {
				n.logger.Error("failed send data to connection",
					slog.Any("error", wErr),
					slog.String("client", clientID),
					slog.String("request", reqID),
					slog.Int("place", 2),
				)
				break
			}
			continue
		}
		// send result to client
		wErr := n.Print(c, res)
		if wErr != nil {
			n.logger.Error("failed send data to connection",
				slog.Any("error", wErr),
				slog.String("client", clientID),
				slog.String("request", reqID),
				slog.Int("place", 3),
			)
			break
		}

		n.logger.Debug("request handled",
			slog.String("client", clientID),
			slog.String("request", reqID),
		)
	}
}
