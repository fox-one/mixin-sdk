package mixin

import (
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
	PublicKey      string  `json:"public_key,omitempty"`
	AccountName    string  `json:"account_name,omitempty"`
	AccountTag     string  `json:"account_tag,omitempty"`
	Confirmations  int     `json:"confirmations,omitempty"`
	Capitalization float64 `json:"capitalization,omitempty"`

	PriceUsd  *decimal.Decimal `json:"price_usd,omitempty"`
	ChangeUsd *decimal.Decimal `json:"change_usd,omitempty"`
	Balance   *decimal.Decimal `json:"balance,omitempty"`
}
