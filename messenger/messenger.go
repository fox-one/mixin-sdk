package messenger

import "github.com/fox-one/mixin-sdk/mixin"

// Messenger mixin messenger
type Messenger struct {
	*mixin.User
}

// NewMessenger create messenger
func NewMessenger(user *mixin.User) *Messenger {
	return &Messenger{
		user,
	}
}
