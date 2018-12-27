package mixin

import (
	"context"
	"encoding/json"
)

// ModifyPIN modify pin
func (user User) ModifyPIN(ctx context.Context, oldPIN, pin string) *Error {
	if pin == oldPIN {
		return nil
	}

	pinEncrypted, err := user.signPIN(oldPIN)
	if err != nil {
		return requestError(err)
	}

	paras := map[string]interface{}{
		"old_pin": pinEncrypted,
	}

	data, err := user.RequestWithPIN(ctx, "POST", "/pin/update", paras, pin)
	if err != nil {
		return requestError(err)
	}

	var resp struct {
		User  *User  `json:"data"`
		Error *Error `json:"error"`
	}

	if err = json.Unmarshal(data, &resp); err != nil {
		return requestError(err)
	} else if resp.Error != nil {
		return resp.Error
	}

	return nil
}

// VerifyPIN verify user pin
func (user User) VerifyPIN(ctx context.Context, pin string) *Error {
	data, err := user.RequestWithPIN(ctx, "POST", "/pin/verify", nil, pin)

	var resp struct {
		User  *User  `json:"data"`
		Error *Error `json:"error"`
	}

	if err = json.Unmarshal(data, &resp); err != nil {
		return requestError(err)
	} else if resp.Error != nil {
		return resp.Error
	}
	return nil
}
