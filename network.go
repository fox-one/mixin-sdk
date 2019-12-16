package mixin

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Chain struct {
	ChainID              string          `json:"chain_id"`
	IconURL              string          `json:"icon_url"`
	Name                 string          `json:"name"`
	Type                 string          `json:"type"`
	WithdrawFee          decimal.Decimal `json:"withdrawal_fee"`
	WithdrawPendingCount int             `json:"withdrawal_pending_count"`
	WithdrawTimestamp    time.Time       `json:"withdrawal_timestamp"`

	DepositBlockHeight  int  `json:"deposit_block_height"`
	ExternalBlockHeight int  `json:"external_block_height"`
	ManagedBlockHeight  int  `json:"managed_block_height"`
	IsSynchronized      bool `json:"is_synchronized"`
}

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

// API

// ReadNetworkInfo read mixin network
func ReadNetworkInfo(ctx context.Context) (*NetworkInfo, error) {
	resp, err := Request(ctx).Get("/network")
	if err != nil {
		return nil, err
	}

	var info NetworkInfo
	if err := UnmarshalResponse(resp, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
