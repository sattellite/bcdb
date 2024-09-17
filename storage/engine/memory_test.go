package engine

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var noopLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestMemory_SetGetDel(t *testing.T) {
	ctx := context.Background()
	done := make(chan struct{})
	defer close(done)
	mem, _ := NewMemory(noopLogger, done)

	// Test Set and Get
	key := "testKey"
	value := []byte("testValue")

	sErr := mem.Set(ctx, key, value)
	require.NoError(t, sErr, "Set failed")

	got1, gErr := mem.Get(ctx, key)
	require.NoError(t, gErr, "Get failed")
	assert.Equal(t, value, got1, "Get returned wrong value")

	// Test Del
	dErr := mem.Del(ctx, key)
	require.NoError(t, dErr, "Del failed")

	got2, gdErr := mem.Get(ctx, key)
	require.Error(t, gdErr, "Get after Del failed")
	assert.Nil(t, got2, "Get after Del returned non-nil value")
}

func TestMemory_ContextCancellation(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	mem, _ := NewMemory(noopLogger, done)
	key := "testKey"
	value := []byte("testValue")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(2 * time.Millisecond)

	err := mem.Set(ctx, key, value)
	assert.Error(t, err, "Set should have failed due to context cancellation")

	_, err = mem.Get(ctx, key)
	assert.Error(t, err, "Get should have failed due to context cancellation")

	err = mem.Del(ctx, key)
	assert.Error(t, err, "Del should have failed due to context cancellation")
}

func TestMemory_SetEmptyKey(t *testing.T) {
	ctx := context.Background()
	done := make(chan struct{})
	defer close(done)
	mem, _ := NewMemory(noopLogger, done)

	err := mem.Set(ctx, "", "value")
	assert.Error(t, err, "Set should fail for empty key")
	assert.Equal(t, ErrEmptyKey, err, "Set should return ErrEmptyKey for empty key")
}

func TestMemory_GetEmptyKey(t *testing.T) {
	ctx := context.Background()
	done := make(chan struct{})
	defer close(done)
	mem, _ := NewMemory(noopLogger, done)

	_, err := mem.Get(ctx, "")
	assert.Error(t, err, "Get should fail for empty key")
	assert.Equal(t, ErrEmptyKey, err, "Get should return ErrEmptyKey for empty key")
}

func TestMemory_DelEmptyKey(t *testing.T) {
	ctx := context.Background()
	done := make(chan struct{})
	defer close(done)
	mem, _ := NewMemory(noopLogger, done)

	err := mem.Del(ctx, "")
	assert.Error(t, err, "Del should fail for empty key")
	assert.Equal(t, ErrEmptyKey, err, "Del should return ErrEmptyKey for empty key")
}

func TestMemory_DoneClose(t *testing.T) {
	done := make(chan struct{})
	mem, _ := NewMemory(noopLogger, done)

	select {
	case _, ok := <-mem.Done():
		assert.True(t, ok, "Done channel should be opened")
	default:
	}

	go func() {
		<-mem.Done()
	}()

	ctx := context.Background()
	mem.Close(ctx)

	select {
	case _, ok := <-mem.Done():
		assert.False(t, ok, "Done channel should be closed")
	default:
	}

	// close again
	mem.Close(ctx)

	select {
	case _, ok := <-mem.Done():
		assert.False(t, ok, "Done channel should be closed")
	default:
	}
}

func TestNewMemoryInitialization(t *testing.T) {
	tests := []struct {
		name    string
		logger  *slog.Logger
		done    chan struct{}
		wantErr bool
	}{
		{"ValidLoggerAndDoneChannel", noopLogger, make(chan struct{}), false},
		{"NilLogger", nil, make(chan struct{}), true},
		{"NilDoneChannel", noopLogger, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem, err := NewMemory(tt.logger, tt.done)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, mem)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, mem)
			}
		})
	}
}
