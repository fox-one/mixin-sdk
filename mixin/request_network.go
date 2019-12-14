package mixin

import (
	"context"

	mixinsdk "github.com/fox-one/mixin-sdk"
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
func ReadNetworkInfo(ctx context.Context) (*NetworkInfo, error) {
	resp, err := mixinsdk.Request(ctx).Get("/network")
	if err != nil {
		return nil, err
	}

	var info NetworkInfo
	if err := mixinsdk.UnmarshalResponse(resp, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
