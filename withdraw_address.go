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

	Label       string `json:"label,omitempty"`
	Destination string `json:"destination,omitempty"`
	Tag         string `json:"tag,omitempty"`

	Fee  string `json:"fee,omitempty"`
	Dust string `json:"dust,omitempty"`

	// TODO Deprecated
	PublicKey   string `json:"public_key,omitempty"`
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
	addr.AccountTag = addr.Tag
	return &addr, nil
}

// ReadWithdrawAddresses read withdraw addresses
func (user *User) ReadWithdrawAddresses(ctx context.Context, assetID string) ([]*WithdrawAddress, error) {
	var addresses []*WithdrawAddress
	if err := user.Request(ctx, "GET", fmt.Sprintf("/assets/%s/addresses", assetID), nil, &addresses); err != nil {
		return nil, err
	}
	for _, addr := range addresses {
		addr.AccountTag = addr.Tag
	}
	return addresses, nil
}

// DeleteWithdrawAddress delete withdraw address
func (user *User) DeleteWithdrawAddress(ctx context.Context, addressID, pin string) error {
	return user.RequestWithPIN(ctx, "POST", fmt.Sprintf("/addresses/%s/delete", addressID), nil, pin, nil)
}

// ReadWithdrawAddresses read addresses with asset id
func ReadWithdrawAddresses(ctx context.Context, assetID, accessToken string) ([]*WithdrawAddress, error) {
	ctx = WithToken(ctx, accessToken)
	resp, err := Request(ctx).Get(fmt.Sprintf("/assets/%s/addresses", assetID))
	if err != nil {
		return nil, err
	}

	var addresses []*WithdrawAddress
	err = UnmarshalResponse(resp, &addresses)

	return addresses, err
}

// ReadWithdrawAddress read address with address id
func ReadWithdrawAddress(ctx context.Context, addressID, accessToken string) (*WithdrawAddress, error) {
	ctx = WithToken(ctx, accessToken)
	resp, err := Request(ctx).Get(fmt.Sprintf("/addresses/%s", addressID))
	if err != nil {
		return nil, err
	}

	var address *WithdrawAddress
	err = UnmarshalResponse(resp, &address)

	return address, err
}
