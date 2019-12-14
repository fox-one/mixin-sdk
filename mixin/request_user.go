package mixin

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
)

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
	if err := user.SendRequest(ctx, "POST", "/users", paras, &u); err != nil {
		return nil, err
	}

	u.privateKey = privateKey
	return &u, nil
}
