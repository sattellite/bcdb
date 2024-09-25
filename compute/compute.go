package compute

import (
	"context"
	"io"

	"github.com/sattellite/bcdb/compute/impl/network"
	"github.com/sattellite/bcdb/compute/impl/repl"
	"github.com/sattellite/bcdb/compute/query"
	"github.com/sattellite/bcdb/compute/result"
	"github.com/sattellite/bcdb/config"
	"github.com/sattellite/bcdb/logger"
	"github.com/sattellite/bcdb/storage"
)

type Computer interface {
	Run(ctx context.Context)
	Parse(input string) (*query.Query, error)
	Handle(ctx context.Context, q query.Query) result.Result
	Print(w io.Writer, r result.Result) error
}

func New(eng storage.Engine, cfg *config.Config) (Computer, error) {
	l := logger.WithScope("compute")
	if cfg.Server.Address != "" {
		return network.New(l, eng, cfg)
	}
	return repl.New(l, eng)
}
