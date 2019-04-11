package mixin

import (
	"context"
	"encoding/json"
)

// AssetFee asset fee info
type AssetFee struct {
	AssetID string `json:"asset_id"`
	Amount  string `json:"amount"`
}

// ReadAssetFee read asset withdraw fee
func (user User) ReadAssetFee(ctx context.Context, assetID string) (*AssetFee, error) {
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

// ReadAsset get asset info, including balance, address info, etc.
func (user User) ReadAsset(ctx context.Context, assetID string) (*Asset, error) {
	data, err := user.Request(ctx, "GET", "/assets/"+assetID, nil)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		Asset *Asset `json:"data,omitempty"`
		Error *Error `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Asset, nil
}

// ReadAssets get user assets info, including balance, address info, etc.
func (user User) ReadAssets(ctx context.Context) ([]*Asset, error) {
	data, err := user.Request(ctx, "GET", "/assets", nil)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		Assets []*Asset `json:"data,omitempty"`
		Error  *Error   `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Assets, nil
}
