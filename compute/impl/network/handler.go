package network

import (
	"context"
	"errors"
	"fmt"

	"github.com/sattellite/bcdb/compute/command"
	"github.com/sattellite/bcdb/compute/query"
	"github.com/sattellite/bcdb/compute/result"
)

func (n *Network) Handle(ctx context.Context, q query.Query) result.Result {
	switch q.Command() {
	case command.MethodSet:
		err := n.engine.Set(ctx, q.Arguments()[0], q.Arguments()[1])
		if err != nil {
			return result.Result{Error: err}
		}
		return result.Result{Value: "success"}
	case command.MethodGet:
		v, err := n.engine.Get(ctx, q.Arguments()[0])
		if err != nil {
			return result.Result{Error: err}
		}
		return result.Result{Value: fmt.Sprintf("%v", v)}
	case command.MethodDel:
		err := n.engine.Del(ctx, q.Arguments()[0])
		if err != nil {
			return result.Result{Error: err}
		}
		return result.Result{Value: "success"}
	}
	return result.Result{Error: errors.New("unknown command")}
}
