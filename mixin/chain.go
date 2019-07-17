package mixin

import (
	"time"

	"github.com/shopspring/decimal"
)

// Chain chain info
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
