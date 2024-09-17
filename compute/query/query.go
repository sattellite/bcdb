package query

import (
	"github.com/sattellite/bcdb/compute/command"
)

type Query struct {
	method    command.Method
	arguments []string
}

func (q *Query) Command() command.Method {
	return q.method
}

func (q *Query) Arguments() []string {
	return q.arguments
}

func New(method command.Method, arguments ...string) *Query {
	return &Query{
		method:    method,
		arguments: arguments,
	}
}
