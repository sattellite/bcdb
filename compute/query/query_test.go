package query

import (
	"testing"

	"github.com/sattellite/bcdb/compute/command"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidQueryCreation(t *testing.T) {
	method := command.MethodSet
	args := []string{"key", "value"}
	query := New(method, args...)

	require.NotNil(t, query, "Query should not be nil")
	assert.Equal(t, method, query.Command(), "Method should match")
	assert.Equal(t, args, query.Arguments(), "Arguments should match")
}

func TestEmptyArguments(t *testing.T) {
	method := command.MethodGet
	query := New(method)

	require.NotNil(t, query, "Query should not be nil")
	assert.Equal(t, method, query.Command(), "Method should match")
	assert.Empty(t, query.Arguments(), "Arguments should be empty")
}

func TestMultipleArguments(t *testing.T) {
	method := command.MethodDel
	args := []string{"key1", "key2", "key3"}
	query := New(method, args...)

	require.NotNil(t, query, "Query should not be nil")
	assert.Equal(t, method, query.Command(), "Method should match")
	assert.Equal(t, args, query.Arguments(), "Arguments should match")
}

func TestNilMethod(t *testing.T) {
	var method command.Method
	args := []string{"key"}
	query := New(method, args...)

	require.NotNil(t, query, "Query should not be nil")
	assert.Equal(t, method, query.Command(), "Method should match")
	assert.Equal(t, args, query.Arguments(), "Arguments should match")
}
