package mixin

import (
	"context"
	"fmt"

	mixinsdk "github.com/fox-one/mixin-sdk"
)

// AssetFee asset fee info
type AssetFee struct {
	AssetID string `json:"asset_id"`
	Amount  string `json:"amount"`
}

// ReadAssetFee read asset withdraw fee
func (user *User) ReadAssetFee(ctx context.Context, assetID string) (*AssetFee, error) {
	var fee AssetFee
	if err := user.SendRequest(ctx, "GET", fmt.Sprintf("/assets/%s/fee", assetID), nil, &fee); err != nil {
		return nil, err
	}
	return &fee, nil
}

// read asset
func readAsset(ctx context.Context, assetId string) (*Asset, error) {
	resp, err := mixinsdk.Request(ctx).Get("/assets/" + assetId)
	if err != nil {
		return nil, err
	}

	var asset Asset
	err = mixinsdk.UnmarshalResponse(resp, &asset)
	return &asset, err
}

// ReadAsset get asset info, including balance, address info, etc.
func (user *User) ReadAsset(ctx context.Context, assetID string) (*Asset, error) {
	ctx = mixinsdk.WithAuth(ctx, user)
	return readAsset(ctx, assetID)
}

// ReadAsset by access token
func ReadAsset(ctx context.Context, assetID string, accessToken string) (*Asset, error) {
	ctx = mixinsdk.WithToken(ctx, accessToken)
	return readAsset(ctx, assetID)
}

func readAssets(ctx context.Context) ([]*Asset, error) {
	resp, err := mixinsdk.Request(ctx).Get("/assets")
	if err != nil {
		return nil, err
	}

	var assets []*Asset
	err = mixinsdk.UnmarshalResponse(resp, &assets)
	return assets, err
}

// ReadAssets get user assets info, including balance, address info, etc.
func (user *User) ReadAssets(ctx context.Context) ([]*Asset, error) {
	ctx = mixinsdk.WithAuth(ctx, user)
	return readAssets(ctx)
}

// ReadAssets return user assets by access token
func ReadAssets(ctx context.Context, accessToken string) ([]*Asset, error) {
	ctx = mixinsdk.WithToken(ctx, accessToken)
	return readAssets(ctx)
}
