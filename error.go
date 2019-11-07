package mixin_sdk

import (
	"fmt"
)

type Error struct {
	Status      int    `json:"status"`
	Code        int    `json:"code"`
	Description string `json:"description"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s [%d/%d]",e.Description,e.Status,e.Code)
}

func createError(status,code int,description string) error {
	return &Error{
		Status:      status,
		Code:        code,
		Description: description,
	}
}

// mixin error codes https://developers.mixin.one/api/alpha-mixin-network/errors/
