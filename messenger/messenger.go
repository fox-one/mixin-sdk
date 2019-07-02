package messenger

import (
	"crypto/x509"
	"encoding/pem"

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
	user := &mixin.User{
		UserID:    userID,
		SessionID: sessionID,
	}

	block, _ := pem.Decode([]byte(sessionKey))
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	user.SetPrivateKey(privateKey)
	return NewMessenger(user), nil
}
