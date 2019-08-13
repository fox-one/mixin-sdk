package mixin

import "time"

// Snapshot transfer records
type Snapshot struct {
	SnapshotID string `json:"snapshot_id"`
	TraceID    string `json:"trace_id,omitempty"`

	UserID     string    `json:"user_id,omitempty"`
	AssetID    string    `json:"asset_id,omitempty"`
	ChainID    string    `json:"chain_id,omitempty"`
	OpponentID string    `json:"opponent_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`

	Source string `json:"source"` // Source DEPOSIT_CONFIRMED, TRANSFER_INITIALIZED, WITHDRAWAL_INITIALIZED, WITHDRAWAL_FEE_CHARGED, WITHDRAWAL_FAILED
	Amount string `json:"amount"`
	Data   string `json:"data,omitempty"`

	Sender          string `json:"sender,omitempty"`
	Receiver        string `json:"receiver,omitempty"`
	TransactionHash string `json:"transaction_hash,omitempty"`

	Asset *Asset `gorm:"-" json:"asset,omitempty"`
}

// DepositTransaction deposit transaction
type DepositTransaction struct {
	Type string `json:"type"`

	TransactionID   string    `json:"transaction_id"`
	TransactionHash string    `json:"transaction_hash"`
	CreatedAt       time.Time `json:"created_at"`

	AssetID       string `json:"asset_id,omitempty"`
	ChainID       string `json:"chain_id,omitempty"`
	Amount        string `json:"amount"`
	Confirmations int    `json:"confirmations"`
	Threshold     int    `json:"threshold"`

	Sender      string `json:"sender"`
	PublicKey   string `json:"public_key"`
	AccountName string `json:"account_name"`
	AccountTag  string `json:"account_tag"`
}
