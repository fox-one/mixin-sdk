package messenger

import (
	mixinsdk "github.com/fox-one/mixin-sdk"
)

func ServerError(err error) error {
	return &mixinsdk.Error{
		Status:      2000000,
		Code:        500,
		Description: err.Error(),
	}
}
