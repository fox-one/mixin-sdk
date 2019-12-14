package mixin

import (
	mixinsdk "github.com/fox-one/mixin-sdk"
)

const (
	// RequestFailed request failed
	RequestFailed = 1000000

	// InvalidTraceID invalid trace
	InvalidTraceID
)

func traceError() error {
	return &mixinsdk.Error{
		Status:      202,
		Code:        InvalidTraceID,
		Description: "invalid trace id",
	}
}
