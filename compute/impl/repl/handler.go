package repl

import (
	"context"
	"errors"
	"fmt"

	"github.com/sattellite/bcdb/compute/command"
	"github.com/sattellite/bcdb/compute/query"
	"github.com/sattellite/bcdb/compute/result"
)

func (r *REPL) Handle(ctx context.Context, q query.Query) result.Result {
	switch q.Command() {
	case command.MethodSet:
		err := r.engine.Set(ctx, q.Arguments()[0], q.Arguments()[1])
		if err != nil {
			return result.Result{Error: err}
		}
		return result.Result{Value: "success"}
	case command.MethodGet:
		v, err := r.engine.Get(ctx, q.Arguments()[0])
		if err != nil {
			return result.Result{Error: err}
		}
		return result.Result{Value: fmt.Sprintf("%v", v)}
	case command.MethodDel:
		err := r.engine.Del(ctx, q.Arguments()[0])
		if err != nil {
			return result.Result{Error: err}
		}
		return result.Result{Value: "success"}
	}
	return result.Result{Error: errors.New("unknown command")}
}
