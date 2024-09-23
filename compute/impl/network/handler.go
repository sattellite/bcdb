package network

import (
	"context"
	"errors"
	"fmt"

	"github.com/sattellite/bcdb/compute/command"
	"github.com/sattellite/bcdb/compute/query"
	"github.com/sattellite/bcdb/compute/result"
)

func (n *Network) Handle(ctx context.Context, q query.Query) (result.Result, error) {
	switch q.Command() {
	case command.MethodSet:
		err := n.engine.Set(ctx, q.Arguments()[0], q.Arguments()[1])
		if err != nil {
			return result.Result{}, err
		}
		return result.Result{Value: fmt.Sprintf("saved key %q", q.Arguments()[0])}, err
	case command.MethodGet:
		v, err := n.engine.Get(ctx, q.Arguments()[0])
		if err != nil {
			return result.Result{}, err
		}
		return result.Result{Value: fmt.Sprintf("value: %v", v)}, err
	case command.MethodDel:
		err := n.engine.Del(ctx, q.Arguments()[0])
		if err != nil {
			return result.Result{}, err
		}
		return result.Result{Value: fmt.Sprintf("deleted key %q", q.Arguments()[0])}, err
	}
	return result.Result{}, errors.New("unknown command")
}
