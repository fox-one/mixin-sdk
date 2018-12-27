package mixin

import (
	"encoding/json"
	"log"
)

const (
	// RequestFailed request failed
	RequestFailed = 1000000

	// InvalidTraceID invalid trace
	InvalidTraceID
)

// Error mixin request error
type Error struct {
	Status      int    `json:"status"`
	Code        int    `json:"code"`
	Description string `json:"description"`
	trace       string
}

// Error error message
func (sessionError Error) Error() string {
	str, err := json.Marshal(sessionError)
	if err != nil {
		log.Panicln(err)
	}
	return string(str)
}

func requestError(err error) *Error {
	return &Error{
		Status:      RequestFailed,
		Code:        504,
		Description: err.Error(),
	}
}

func traceError() *Error {
	return &Error{
		Status:      InvalidTraceID,
		Code:        400,
		Description: "invalid trace id",
	}
}
