package mixin

import (
	"context"
	"fmt"

	"github.com/fox-one/mixin-sdk/utils"
)

// CreateWithdrawAddress create withdraw address
func (user *User) CreateWithdrawAddress(ctx context.Context, address WithdrawAddress, pin string) (*WithdrawAddress, error) {
	if len(address.Label) == 0 {
		address.Label = "Created by FoxONE"
	}

	paras := utils.UnselectFields(address, "fee", "dust")

	var addr WithdrawAddress
	if err := user.SendRequestWithPIN(ctx, "POST", "/addresses", paras, pin, &addr); err != nil {
		return nil, err
	}
	return &addr, nil
}

// ReadWithdrawAddresses read withdraw addresses
func (user *User) ReadWithdrawAddresses(ctx context.Context, assetID string) ([]*WithdrawAddress, error) {
	var addrs []*WithdrawAddress
	if err := user.SendRequest(ctx, "GET", fmt.Sprintf("/assets/%s/addresses", assetID), nil, &addrs); err != nil {
		return nil, err
	}
	return addrs, nil
}

// DeleteWithdrawAddress delete withdraw address
func (user *User) DeleteWithdrawAddress(ctx context.Context, addressID, pin string) error {
	return user.SendRequestWithPIN(ctx, "POST", fmt.Sprintf("/addresses/%s/delete", addressID), nil, pin, nil)
}
