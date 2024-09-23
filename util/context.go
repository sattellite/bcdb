package util

import "context"

// Key to use when setting the userid.
type ctxKeyServiceEntities int

const (
	// RequestIDKey is the key that holds the request id in a context.
	RequestIDKey ctxKeyServiceEntities = iota

	// ClientIDKey is the key that holds the client id in a context.
	ClientIDKey
)

func SetContextRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetContextRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if rid, ok := ctx.Value(RequestIDKey).(string); ok {
		return rid
	}
	return ""
}

func SetContextClientID(ctx context.Context, clientID string) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, ClientIDKey, clientID)
}

func GetContextClientID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if cid, ok := ctx.Value(ClientIDKey).(string); ok {
		return cid
	}
	return ""
}
