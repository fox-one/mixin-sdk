package mixin

import (
	"context"
	"encoding/json"

	mixin_sdk "github.com/fox-one/mixin-sdk"
)

// AssetFee asset fee info
type AssetFee struct {
	AssetID string `json:"asset_id"`
	Amount  string `json:"amount"`
}

// ReadAssetFee read asset withdraw fee
func (user *User) ReadAssetFee(ctx context.Context, assetID string) (*AssetFee, error) {
	data, err := user.Request(ctx, "GET", "/assets/"+assetID+"/fee", nil)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		AssetFee *AssetFee `json:"data,omitempty"`
		Error    *Error    `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.AssetFee, nil
}

// read asset
func readAsset(ctx context.Context,assetId string) (*Asset,error) {
	resp,err := mixin_sdk.Request(ctx).Get("/assets/"+assetId)
	if err != nil {
		return nil,err
	}

	var asset Asset
	err = mixin_sdk.UnmarshalResponse(resp,&asset)
	return &asset,err
}

// ReadAsset get asset info, including balance, address info, etc.
func (user *User) ReadAsset(ctx context.Context, assetID string) (*Asset, error) {
	ctx = mixin_sdk.WithAuth(ctx,user)
	return readAsset(ctx,assetID)
}

// ReadAsset by access token
func ReadAsset(ctx context.Context, assetID string,accessToken string) (*Asset,error) {
	ctx = mixin_sdk.WithToken(ctx,accessToken)
	return readAsset(ctx,assetID)
}

func readAssets(ctx context.Context) ([]*Asset,error) {
	resp,err := mixin_sdk.Request(ctx).Get("/assets")
	if err != nil {
		return nil,err
	}

	var assets []*Asset
	err = mixin_sdk.UnmarshalResponse(resp,&assets)
	return assets,err
}

// ReadAssets get user assets info, including balance, address info, etc.
func (user *User) ReadAssets(ctx context.Context) ([]*Asset, error) {
	ctx = mixin_sdk.WithAuth(ctx,user)
	return readAssets(ctx)
}

// ReadAssets return user assets by access token
func ReadAssets(ctx context.Context,accessToken string) ([]*Asset,error) {
	ctx = mixin_sdk.WithToken(ctx,accessToken)
	return readAssets(ctx)
}
