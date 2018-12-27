package mixin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fox-one/mixin-sdk/utils"
)

// CreateWithdrawAddress create withdraw address
func (user User) CreateWithdrawAddress(ctx context.Context, address WithdrawAddress, pin string) (*WithdrawAddress, *Error) {
	paras := utils.UnselectFields(address, "fee", "dust")
	data, err := user.RequestWithPIN(ctx, "POST", "/addresses", paras, pin)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		Address *WithdrawAddress `json:"data,omitempty"`
		Error   *Error           `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Address, nil
}

// ReadWithdrawAddresses read withdraw addresses
func (user User) ReadWithdrawAddresses(ctx context.Context, assetID string) ([]*WithdrawAddress, error) {
	data, err := user.Request(ctx, "GET", fmt.Sprintf("/assets/%s/addresses", assetID), nil)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		Addrs []*WithdrawAddress `json:"data,omitempty"`
		Error *Error             `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Addrs, nil
}

// DeleteWithdrawAddress delete withdraw address
func (user User) DeleteWithdrawAddress(ctx context.Context, addressID, pin string) *Error {
	data, err := user.RequestWithPIN(ctx, "POST", fmt.Sprintf("/addresses/%s/delete", addressID), nil, pin)
	if err != nil {
		return requestError(err)
	}

	var resp struct {
		Addrs []*WithdrawAddress `json:"data,omitempty"`
		Error *Error             `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return requestError(err)
	} else if resp.Error != nil {
		return resp.Error
	}

	return nil
}
