package command

import (
	"errors"
	"strings"
)

var (
	ErrInvalidCommand   = errors.New("invalid command")
	ErrInvalidArguments = errors.New("invalid arguments")
)

// Method represents a type of command
type Method int

func (t *Method) String() string {
	switch *t {
	case MethodSet:
		return "SET"
	case MethodGet:
		return "GET"
	case MethodDel:
		return "DEL"
	}
	return "unknown"
}

// Command types
const (
	MethodSet Method = iota
	MethodGet
	MethodDel
)

func ParseMethod(input string) (*Method, error) {
	if len(input) != 3 {
		return nil, ErrInvalidCommand
	}

	var cmd Method
	switch strings.ToUpper(input) {
	case "SET":
		cmd = MethodSet
	case "GET":
		cmd = MethodGet
	case "DEL":
		cmd = MethodDel
	default:
		return nil, ErrInvalidCommand
	}

	return &cmd, nil
}

func ParseArguments(cmd *Method, args ...string) ([]string, error) {
	cleared := make([]string, 0, len(args))
	for _, arg := range args {
		if arg != "" {
			cleared = append(cleared, arg)
		}
	}

	switch *cmd {
	case MethodSet:
		if len(cleared) != 2 {
			return nil, ErrInvalidArguments
		}
	case MethodGet, MethodDel:
		if len(cleared) != 1 {
			return nil, ErrInvalidArguments
		}
	}

	return cleared, nil
}
