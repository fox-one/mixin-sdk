package mixin

import (
	"context"
	"fmt"

	"github.com/fox-one/mixin-sdk/utils"
)

// WithdrawAddress withdraw address
type WithdrawAddress struct {
	AddressID string `json:"address_id,omitempty"`
	AssetID   string `json:"asset_id"`

	Destination string `json:"destination,omitempty"`
	Tag         string `json:"tag,omitempty"`

	Fee  string `json:"fee,omitempty"`
	Dust string `json:"dust,omitempty"`

	// TODO Deprecated
	PublicKey   string `json:"public_key,omitempty"`
	Label       string `json:"label,omitempty"`
	AccountName string `json:"account_name,omitempty"`
	AccountTag  string `json:"account_tag,omitempty"`
}

// API

// CreateWithdrawAddress create withdraw address
func (user *User) CreateWithdrawAddress(ctx context.Context, address WithdrawAddress, pin string) (*WithdrawAddress, error) {
	if len(address.Label) == 0 {
		address.Label = "Created by Fox.ONE Mixin SDK"
	}

	var addr WithdrawAddress
	if err := user.RequestWithPIN(ctx, "POST", "/addresses", utils.UnselectFields(address, "fee", "dust"), pin, &addr); err != nil {
		return nil, err
	}
	return &addr, nil
}

// ReadWithdrawAddresses read withdraw addresses
func (user *User) ReadWithdrawAddresses(ctx context.Context, assetID string) ([]*WithdrawAddress, error) {
	var addresses []*WithdrawAddress
	if err := user.Request(ctx, "GET", fmt.Sprintf("/assets/%s/addresses", assetID), nil, &addresses); err != nil {
		return nil, err
	}
	return addresses, nil
}

// DeleteWithdrawAddress delete withdraw address
func (user *User) DeleteWithdrawAddress(ctx context.Context, addressID, pin string) error {
	return user.RequestWithPIN(ctx, "POST", fmt.Sprintf("/addresses/%s/delete", addressID), nil, pin, nil)
}
