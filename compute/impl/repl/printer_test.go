package repl

import (
	"bytes"
	"testing"

	"github.com/sattellite/bcdb/compute/result"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrinter(t *testing.T) {
	tests := []struct {
		name        string
		repl        *REPL
		res         result.Result
		expectError bool
		expectedOut string
	}{
		{"Print result successfully", &REPL{out: new(bytes.Buffer)}, result.Result{Value: "output"}, false, "< output\n"},
		{"Print result with prompt error", &REPL{out: &errorWriter{}}, result.Result{Value: "output"}, true, ""},
		{"Prompt successfully", &REPL{out: new(bytes.Buffer)}, result.Result{}, false, "> "},
		{"Prompt with error", &REPL{out: &errorWriter{}}, result.Result{}, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.name == "Prompt successfully" || tt.name == "Prompt with error" {
				err = tt.repl.prompt(tt.repl.out, prefixIn)
			} else {
				err = tt.repl.Print(tt.repl.out, tt.res)
			}
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOut, tt.repl.out.(*bytes.Buffer).String())
			}
		})
	}
}

func TestPrompt(t *testing.T) {
	tests := []struct {
		name        string
		repl        *REPL
		prefix      []byte
		expectError bool
		expectedOut string
	}{
		{"Prompt with prefixIn successfully", &REPL{out: new(bytes.Buffer)}, prefixIn, false, "> "},
		{"Prompt with prefixOut successfully", &REPL{out: new(bytes.Buffer)}, prefixOut, false, "< "},
		{"Prompt with error", &REPL{out: &errorWriter{}}, prefixIn, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.repl.prompt(tt.repl.out, tt.prefix)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOut, tt.repl.out.(*bytes.Buffer).String())
			}
		})
	}
}

type errorWriter struct{}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, assert.AnError
}
