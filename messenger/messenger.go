package messenger

import (
	"github.com/fox-one/mixin-sdk/mixin"
)

// Messenger mixin messenger
type Messenger struct {
	*mixin.User
	*BlazeClient
}

// NewMessenger create messenger
func NewMessenger(user *mixin.User) *Messenger {
	return &Messenger{
		user,
		NewBlazeClient(user.UserID, user.SessionID, user.SessionKey),
	}
}

// NewMessengerWithSession new messenger
func NewMessengerWithSession(userID, sessionID, sessionKey string) (*Messenger, error) {
	user, err := mixin.NewUser(userID, sessionID, sessionKey)
	if err != nil {
		return nil, err
	}
	return NewMessenger(user), nil
}
