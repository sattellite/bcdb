package repl

import (
	"bufio"
	"context"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/sattellite/bcdb/compute/result"
	"github.com/sattellite/bcdb/storage"
)

func New(logger *slog.Logger, engine storage.Engine) *REPL {
	return &REPL{
		logger: logger.With("module", "repl"),
		engine: engine,
		in:     make(chan string),
		out:    log.New(os.Stdout, "", 0).Writer(),
	}
}

type REPL struct {
	logger *slog.Logger
	engine storage.Engine
	in     chan string
	out    io.Writer
}

func (r *REPL) Run(ctx context.Context) {
	r.logger.Info("compute started")
	defer func() {
		r.logger.Info("compute stopped")
	}()

	scanner := bufio.NewScanner(os.Stdin)

	_ = r.prompt(r.out, prefixIn)
	for scanner.Scan() && ctx.Err() == nil {
		input := scanner.Text()
		// process user input
		q, pErr := r.Parse(input)
		if pErr != nil {
			r.logger.Error("failed to parse command", slog.Any("error", pErr))
			_ = r.Print(r.out, result.Result{Value: pErr.Error()})
			_ = r.prompt(r.out, prefixIn)
			continue
		}
		// handle user input
		res, hErr := r.Handle(ctx, *q)
		if hErr != nil {
			r.logger.Error("failed to handle command", slog.Any("error", hErr))
			_ = r.Print(r.out, result.Result{Value: hErr.Error()})
			_ = r.prompt(r.out, prefixIn)
			continue
		}
		// write result to stdout
		_ = r.Print(r.out, res)
		_ = r.prompt(r.out, prefixIn)
	}
}
