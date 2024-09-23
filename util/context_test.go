package util

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetContextRequestID(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		requestID  string
		expectedID string
	}{
		{
			name:       "ContextRequestIDIsSetCorrectly",
			ctx:        context.Background(),
			requestID:  "12345",
			expectedID: "12345",
		},
		{
			name:       "ContextRequestIDIsRetrievedCorrectly",
			ctx:        context.WithValue(context.Background(), RequestIDKey, "12345"),
			requestID:  "12345",
			expectedID: "12345",
		},
		{
			name:       "ContextRequestIDIsEmptyWhenNotSet",
			ctx:        context.Background(),
			requestID:  "",
			expectedID: "",
		},
		{
			name:       "ContextRequestIDIsEmptyWhenContextIsNil",
			ctx:        nil,
			requestID:  "",
			expectedID: "",
		},
		{
			name:       "ContextRequestIDIsEmptyWhenWrongType",
			ctx:        context.WithValue(context.Background(), RequestIDKey, 12345),
			requestID:  "",
			expectedID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := SetContextRequestID(tt.ctx, tt.requestID)
			assert.Equal(t, tt.expectedID, GetContextRequestID(ctx))
		})
	}
}

func TestGetContextRequestID(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		expectedID string
	}{
		{
			name:       "ContextRequestIDIsRetrievedCorrectly",
			ctx:        context.WithValue(context.Background(), RequestIDKey, "12345"),
			expectedID: "12345",
		},
		{
			name:       "ContextRequestIDIsEmptyWhenNotSet",
			ctx:        context.Background(),
			expectedID: "",
		},
		{
			name:       "ContextRequestIDIsEmptyWhenContextIsNil",
			ctx:        nil,
			expectedID: "",
		},
		{
			name:       "ContextRequestIDIsEmptyWhenWrongType",
			ctx:        context.WithValue(context.Background(), RequestIDKey, 12345),
			expectedID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedID, GetContextRequestID(tt.ctx))
		})
	}
}

func TestSetContextClientID(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		clientID   string
		expectedID string
	}{
		{
			name:       "ContextClientIDIsSetCorrectly",
			ctx:        context.Background(),
			clientID:   "client123",
			expectedID: "client123",
		},
		{
			name:       "ContextClientIDIsEmptyWhenContextIsNil",
			ctx:        nil,
			clientID:   "",
			expectedID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := SetContextClientID(tt.ctx, tt.clientID)
			assert.Equal(t, tt.expectedID, GetContextClientID(ctx))
		})
	}
}

func TestGetContextClientID(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		expectedID string
	}{
		{
			name:       "ContextClientIDIsRetrievedCorrectly",
			ctx:        context.WithValue(context.Background(), ClientIDKey, "client123"),
			expectedID: "client123",
		},
		{
			name:       "ContextClientIDIsEmptyWhenNotSet",
			ctx:        context.Background(),
			expectedID: "",
		},
		{
			name:       "ContextClientIDIsEmptyWhenContextIsNil",
			ctx:        nil,
			expectedID: "",
		},
		{
			name:       "ContextClientIDIsEmptyWhenWrongType",
			ctx:        context.WithValue(context.Background(), ClientIDKey, 12345),
			expectedID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedID, GetContextClientID(tt.ctx))
		})
	}
}
