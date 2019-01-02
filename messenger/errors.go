package messenger

import "github.com/fox-one/mixin-sdk/mixin"

func requestError(err error) *mixin.Error {
	return &mixin.Error{
		Status:      mixin.RequestFailed,
		Code:        504,
		Description: err.Error(),
	}
}
