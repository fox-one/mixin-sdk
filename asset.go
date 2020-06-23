package mixin

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

// Asset asset info
type Asset struct {
	AssetID        string  `json:"asset_id"`
	ChainID        string  `json:"chain_id"`
	AssetKey       string  `json:"asset_key,omitempty"`
	Symbol         string  `json:"symbol,omitempty"`
	Name           string  `json:"name,omitempty"`
	IconURL        string  `json:"icon_url,omitempty"`
	Confirmations  int     `json:"confirmations,omitempty"`
	Capitalization float64 `json:"capitalization,omitempty"`

	PriceUsd  decimal.Decimal `json:"price_usd,omitempty"`
	ChangeUsd decimal.Decimal `json:"change_usd,omitempty"`
	Balance   decimal.Decimal `json:"balance,omitempty"`

	Destination string `json:"destination,omitempty"`
	Tag         string `json:"tag,omitempty"`

	// TODO Deprecated
	PublicKey   string `json:"public_key,omitempty"`
	AccountName string `json:"account_name,omitempty"`
	AccountTag  string `json:"account_tag,omitempty"`
}

// AssetFee asset fee info
type AssetFee struct {
	AssetID string `json:"asset_id"`
	Amount  string `json:"amount"`
}

// API

// read asset
func readAsset(ctx context.Context, assetId string) (*Asset, error) {
	resp, err := Request(ctx).Get("/assets/" + assetId)
	if err != nil {
		return nil, err
	}

	var asset Asset
	err = UnmarshalResponse(resp, &asset)
	return &asset, err
}

// ReadAsset get asset info, including balance, address info, etc.
func (user *User) ReadAsset(ctx context.Context, assetID string) (*Asset, error) {
	ctx = WithAuth(ctx, user)
	return readAsset(ctx, assetID)
}

func (ed *EdOToken) ReadAsset(ctx context.Context, assetID string) (*Asset, error) {
	ctx = WithAuth(ctx, ed)
	return readAsset(ctx, assetID)
}

// ReadAsset by access token
func ReadAsset(ctx context.Context, assetID string, accessToken string) (*Asset, error) {
	ctx = WithToken(ctx, accessToken)
	return readAsset(ctx, assetID)
}

func readAssets(ctx context.Context) ([]*Asset, error) {
	resp, err := Request(ctx).Get("/assets")
	if err != nil {
		return nil, err
	}

	var assets []*Asset
	err = UnmarshalResponse(resp, &assets)
	return assets, err
}

// ReadAssets get user assets info, including balance, address info, etc.
func (user *User) ReadAssets(ctx context.Context) ([]*Asset, error) {
	ctx = WithAuth(ctx, user)
	return readAssets(ctx)
}

func (ed *EdOToken) ReadAssets(ctx context.Context) ([]*Asset, error) {
	ctx = WithAuth(ctx, ed)
	return readAssets(ctx)
}

// ReadAssets return user assets by access token
func ReadAssets(ctx context.Context, accessToken string) ([]*Asset, error) {
	ctx = WithToken(ctx, accessToken)
	return readAssets(ctx)
}

func (user *User) ReadAssetFee(ctx context.Context, assetID string) (*AssetFee, error) {
	var fee AssetFee
	if err := user.Request(ctx, "GET", fmt.Sprintf("/assets/%s/fee", assetID), nil, &fee); err != nil {
		return nil, err
	}
	return &fee, nil
}
