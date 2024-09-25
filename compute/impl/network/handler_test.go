package network

import (
	"context"
	"testing"

	"github.com/sattellite/bcdb/compute/command"
	"github.com/sattellite/bcdb/compute/query"
	"github.com/sattellite/bcdb/compute/result"
	"github.com/sattellite/bcdb/storage/engine"
	storage "github.com/sattellite/bcdb/storage/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	tests := []struct {
		name        string
		query       *query.Query
		setupMock   func(*storage.Engine)
		expectedRes result.Result
		expectError bool
	}{
		{
			name:  "Handle SET command successfully",
			query: query.New(command.MethodSet, "key", "value"),
			setupMock: func(m *storage.Engine) {
				m.On("Set", mock.Anything, "key", "value").Return(nil)
			},
			expectedRes: result.Result{Value: `success`},
			expectError: false,
		},
		{
			name:  "Handle GET command successfully",
			query: query.New(command.MethodGet, "key"),
			setupMock: func(m *storage.Engine) {
				m.On("Get", mock.Anything, "key").Return("value", nil)
			},
			expectedRes: result.Result{Value: "value"},
			expectError: false,
		},
		{
			name:  "Handle DEL command successfully",
			query: query.New(command.MethodDel, "key"),
			setupMock: func(m *storage.Engine) {
				m.On("Del", mock.Anything, "key").Return(nil)
			},
			expectedRes: result.Result{Value: `success`},
			expectError: false,
		},
		{
			name:        "Handle unknown command",
			query:       query.New(command.Method(-1), "key"),
			setupMock:   func(m *storage.Engine) {},
			expectedRes: result.Result{},
			expectError: true,
		},
		{
			name:  "Handle SET command failure",
			query: query.New(command.MethodSet, "", ""),
			setupMock: func(m *storage.Engine) {
				m.On("Set", mock.Anything, "", "").Return(engine.ErrEmptyKey)
			},
			expectedRes: result.Result{},
			expectError: true,
		},
		{
			name:  "Handle GET command failure",
			query: query.New(command.MethodGet, ""),
			setupMock: func(m *storage.Engine) {
				m.On("Get", mock.Anything, "").Return("", engine.ErrEmptyKey)
			},
			expectedRes: result.Result{},
			expectError: true,
		},
		{
			name:  "Handle DEL command failure",
			query: query.New(command.MethodDel, ""),
			setupMock: func(m *storage.Engine) {
				m.On("Del", mock.Anything, "").Return(engine.ErrEmptyKey)
			},
			expectedRes: result.Result{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEngine := storage.NewEngine(t)
			tt.setupMock(mockEngine)
			r := &Network{engine: mockEngine}

			res := r.Handle(context.Background(), *tt.query)
			if tt.expectError {
				require.Error(t, res.Error)
			} else {
				require.NoError(t, res.Error)
				assert.Equal(t, tt.expectedRes, res)
			}
			mockEngine.AssertExpectations(t)
		})
	}
}
