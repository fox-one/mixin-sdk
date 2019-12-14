package mixin

import (
	"context"
)

// ModifyPIN modify pin
func (user *User) ModifyPIN(ctx context.Context, oldPIN, pin string) error {
	if pin == oldPIN {
		return nil
	}

	pinEncrypted, err := user.signPIN(oldPIN)
	if err != nil {
		return err
	}

	return user.SendRequestWithPIN(ctx, "POST", "/pin/update", map[string]interface{}{"old_pin": pinEncrypted}, pin, nil)
}

// VerifyPIN verify user pin
func (user User) VerifyPIN(ctx context.Context, pin string) error {
	return user.SendRequestWithPIN(ctx, "POST", "/pin/verify", nil, pin, nil)
}
