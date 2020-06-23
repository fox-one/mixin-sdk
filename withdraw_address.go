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

func readWithdrawAddresses(ctx context.Context, assetID string) ([]*WithdrawAddress, error) {
	resp, err := Request(ctx).Get(fmt.Sprintf("/assets/%s/addresses", assetID))
	if err != nil {
		return nil, err
	}

	var addresses []*WithdrawAddress
	if err = UnmarshalResponse(resp, &addresses); err != nil {
		return nil, err
	}

	// TODO just fix old version tag, remove in later version
	for _, addr := range addresses {
		addr.AccountTag = addr.Tag
	}

	return addresses, err
}

// ReadWithdrawAddresses read withdraw addresses
func (user *User) ReadWithdrawAddresses(ctx context.Context, assetID string) ([]*WithdrawAddress, error) {
	ctx = WithAuth(ctx, user)
	return readWithdrawAddresses(ctx, assetID)
}

func (ed *EdOToken) ReadWithdrawAddresses(ctx context.Context, assetID string) ([]*WithdrawAddress, error) {
	ctx = WithAuth(ctx, ed)
	return readWithdrawAddresses(ctx, assetID)
}

// ReadWithdrawAddresses read addresses with asset id
func ReadWithdrawAddresses(ctx context.Context, assetID, accessToken string) ([]*WithdrawAddress, error) {
	ctx = WithToken(ctx, accessToken)
	return readWithdrawAddresses(ctx, assetID)
}

// DeleteWithdrawAddress delete withdraw address
func (user *User) DeleteWithdrawAddress(ctx context.Context, addressID, pin string) error {
	return user.RequestWithPIN(ctx, "POST", fmt.Sprintf("/addresses/%s/delete", addressID), nil, pin, nil)
}
