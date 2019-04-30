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

func requestError(err error) error {
	return &Error{
		Status:      504,
		Code:        RequestFailed,
		Description: err.Error(),
	}
}

func traceError() error {
	return &Error{
		Status:      202,
		Code:        InvalidTraceID,
		Description: "invalid trace id",
	}
}
