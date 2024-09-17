package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func methodRef(m Method) *Method {
	return &m
}

func TestParseMethod(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Method
		err      error
	}{
		{"Valid SET command", "SET", methodRef(MethodSet), nil},
		{"Valid GET command", "GET", methodRef(MethodGet), nil},
		{"Valid DEL command", "DEL", methodRef(MethodDel), nil},
		{"Invalid command", "ABC", nil, ErrInvalidCommand},
		{"Empty command", "", nil, ErrInvalidCommand},
		{"Lowercase command", "set", methodRef(MethodSet), nil},
		{"Mixed case command", "gEt", methodRef(MethodGet), nil},
		{"Command with extra spaces", " SET ", nil, ErrInvalidCommand},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := ParseMethod(tt.input)
			if tt.err != nil {
				require.Error(t, err)
				assert.Equal(t, tt.err, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, cmd)
			}
		})
	}
}

func TestParseArguments(t *testing.T) {
	tests := []struct {
		name     string
		cmd      Method
		args     []string
		expected []string
		err      error
	}{
		{"Valid SET command with arguments", MethodSet, []string{"key", "value"}, []string{"key", "value"}, nil},
		{"Valid GET command with argument", MethodGet, []string{"key"}, []string{"key"}, nil},
		{"Valid DEL command with argument", MethodDel, []string{"key"}, []string{"key"}, nil},
		{"SET command with missing arguments", MethodSet, []string{"key"}, nil, ErrInvalidArguments},
		{"GET command with extra arguments", MethodGet, []string{"key", "extra"}, nil, ErrInvalidArguments},
		{"DEL command with extra arguments", MethodDel, []string{"key", "extra"}, nil, ErrInvalidArguments},
		{"SET command with empty arguments", MethodSet, []string{"", "value"}, []string{"value"}, ErrInvalidArguments},
		{"GET command with empty argument", MethodGet, []string{""}, nil, ErrInvalidArguments},
		{"DEL command with empty argument", MethodDel, []string{""}, nil, ErrInvalidArguments},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := ParseArguments(&tt.cmd, tt.args...)
			if tt.err != nil {
				require.Error(t, err)
				assert.Equal(t, tt.err, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, args)
			}
		})
	}
}
