package mixin

import (
	"context"

	"github.com/fox-one/pkg/uuid"
)

type contextKey int

const (
	authKey contextKey = iota
	requestIdKey
)

func WithAuth(ctx context.Context, auth Authentication) context.Context {
	return context.WithValue(ctx, authKey, auth)
}

func WithToken(ctx context.Context, token string) context.Context {
	return WithAuth(ctx, accessToken(token))
}

// WithRequestID bind request id to context
// request id must be uuid
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIdKey, requestID)
}

func RequestIdFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIdKey).(string); ok {
		return v
	}

	return uuid.New()
}
