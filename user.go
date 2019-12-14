package sdk

import (
	"context"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"

	"github.com/fox-one/pkg/encrypt"
)

// User wallet entity
type User struct {
	UserID     string `json:"user_id"`
	SessionID  string `json:"session_id"`
	PINToken   string `json:"pin_token"`
	SessionKey string `json:"session_key"`

	FullName string `json:"full_name"`

	pinCipher  *cipher.Block
	privateKey *rsa.PrivateKey
	scopes     string
}

// NewUser new user
func NewUser(userID, sessionID, sessionKey string, pinToken ...string) (*User, error) {
	pri, _, err := encrypt.ParsePrivatePem(sessionKey)
	if err != nil {
		return nil, err
	}
	user := User{
		UserID:     userID,
		SessionID:  sessionID,
		privateKey: pri.(*rsa.PrivateKey),
	}

	if len(pinToken) > 0 && pinToken[0] != "" {
		user.PINToken = pinToken[0]
	}
	return &user, nil
}

// SetPrivateKey set private key
func (user *User) SetPrivateKey(privateKey *rsa.PrivateKey) {
	user.privateKey = privateKey
}

// SetScopes set scopes
func (user *User) SetScopes(scopes string) {
	user.scopes = scopes
}

// API

// CreateUser create a new wallet user
func (user *User) CreateUser(ctx context.Context, privateKey *rsa.PrivateKey, fullname string) (*User, error) {
	pbts, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return nil, err
	}

	paras := map[string]interface{}{
		"session_secret": base64.StdEncoding.EncodeToString(pbts),
		"full_name":      fullname,
	}

	var u User
	if err := user.Request(ctx, "POST", "/users", paras, &u); err != nil {
		return nil, err
	}

	u.privateKey = privateKey
	return &u, nil
}
