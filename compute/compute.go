package compute

import (
	"context"

	"github.com/sattellite/bcdb/compute/impl/repl"
	"github.com/sattellite/bcdb/compute/query"
	"github.com/sattellite/bcdb/compute/result"
	"github.com/sattellite/bcdb/logger"
	"github.com/sattellite/bcdb/storage"
)

type Computer interface {
	Run(ctx context.Context)
	Parse(input string) (*query.Query, error)
	Handle(ctx context.Context, q query.Query) (result.Result, error)
	Print(r result.Result) error
}

func New(eng storage.Engine) Computer {
	return repl.New(logger.WithScope("compute"), eng)
}
