package mixin

import (
	"context"
	"encoding/json"

	"github.com/shopspring/decimal"
)

// NetworkInfo mixin network info
type NetworkInfo struct {
	Assets []struct {
		Amount  decimal.Decimal `json:"amount"`
		AssetID string          `json:"asset_id"`
		IconURL string          `json:"icon_url"`
		Symbol  string          `json:"symbol"`
		Type    string          `json:"type"`
	} `json:"assets"`
	Chains         []*Chain        `json:"chains"`
	AssetsCount    decimal.Decimal `json:"assets_count"`
	PeakThroughput decimal.Decimal `json:"peak_throughput"`
	SnapshotsCount decimal.Decimal `json:"snapshots_count"`
	Type           string          `json:"type"`
}

// ReadNetworkInfo read mixin network
func (user User) ReadNetworkInfo(ctx context.Context) (*NetworkInfo, error) {
	data, err := user.Request(ctx, "GET", "/network", nil)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		NetworkInfo *NetworkInfo `json:"data,omitempty"`
		Error       *Error       `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.NetworkInfo, nil
}
