package mixin_sdk

import (
	"context"
	"net/http"
)

type Authentication interface {
	Auth(r *http.Request) (string,error)
}

type contextKey int

const (
	authKey = iota
)

func WithAuth(ctx context.Context,auth Authentication) context.Context {
	return context.WithValue(ctx,authKey,auth)
}

type accessToken string

func (token accessToken) Auth(r *http.Request) (string,error) {
	return string(token),nil
}

func WithToken(ctx context.Context,token string) context.Context {
	return WithAuth(ctx,accessToken(token))
}
