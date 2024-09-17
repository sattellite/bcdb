package repl

import (
	"testing"

	"github.com/sattellite/bcdb/compute/command"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedMethod command.Method
		expectedArgs   []string
		expectError    bool
	}{
		{name: "Valid SET command", input: "SET key value", expectedMethod: command.MethodSet, expectedArgs: []string{"key", "value"}, expectError: false},
		{name: "Valid GET command", input: "GET key", expectedMethod: command.MethodGet, expectedArgs: []string{"key"}, expectError: false},
		{name: "Invalid command", input: "INVALID key", expectError: true},
		{name: "Empty input", input: "", expectError: true},
		{name: "SET command with missing arguments", input: "SET key", expectError: true},
		{name: "GET command with extra arguments", input: "GET key extra", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &REPL{}
			q, err := r.Parse(tt.input)
			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, q)
			} else {
				require.NoError(t, err)
				require.NotNil(t, q)
				assert.Equal(t, tt.expectedMethod, q.Command())
				assert.Equal(t, tt.expectedArgs, q.Arguments())
			}
		})
	}
}

func methodRef(m command.Method) *command.Method {
	return &m
}
