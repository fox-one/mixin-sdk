package mixin

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// CreateUser create a new wallet user
func (user User) CreateUser(ctx context.Context, privateKey *rsa.PrivateKey, fullname string) (*User, *Error) {
	pbts, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return nil, requestError(err)
	}

	paras := map[string]interface{}{
		"session_secret": base64.StdEncoding.EncodeToString(pbts),
	}
	if len(fullname) > 0 {
		paras["full_name"] = fullname
	}

	payload, err := json.Marshal(paras)
	if err != nil {
		return nil, requestError(err)
	}

	data, err := user.Request(ctx, "POST", "/users", payload)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		User  *User  `json:"data"`
		Error *Error `json:"error"`
	}

	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}

	if len(resp.User.UserID) == 0 || len(resp.User.SessionID) == 0 || len(resp.User.PINToken) == 0 {
		return nil, requestError(fmt.Errorf("create mixin user failed"))
	}

	resp.User.privateKey = privateKey
	return resp.User, nil
}
