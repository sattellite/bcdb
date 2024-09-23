package network

import (
	"errors"
	"strings"

	"github.com/sattellite/bcdb/compute/command"
	"github.com/sattellite/bcdb/compute/query"
)

var ErrInvalidQuery = errors.New("invalid query")

func (n *Network) Parse(input string) (*query.Query, error) {
	parts := strings.Split(input, " ")
	if len(parts) < 2 {
		return nil, ErrInvalidQuery
	}

	cmd, err := command.ParseMethod(parts[0])
	if err != nil {
		return nil, err
	}

	args, aErr := command.ParseArguments(cmd, parts[1:]...)
	if aErr != nil {
		return nil, aErr
	}

	return query.New(*cmd, args...), nil
}
