package mixin

import (
	"context"
	"time"

	uuid "github.com/gofrs/uuid"
)

// RawTransaction raw transaction
type RawTransaction struct {
	Type            string    `json:"type"`
	SnapshotID      string    `json:"snapshot"`
	OpponentKey     string    `json:"opponent_key"`
	AssetID         string    `json:"asset_id"`
	Amount          string    `json:"amount"`
	TraceID         string    `json:"trace_id"`
	Memo            string    `json:"memo"`
	State           string    `json:"state"`
	CreatedAt       time.Time `json:"created_at"`
	TransactionHash string    `json:"transaction_hash"`
	SnapshotHash    string    `json:"snapshot_hash"`
	SnapshotAt      time.Time `json:"snapshot_at"`
}

// Transaction do transaction to mixin main net
func (user *User) Transaction(ctx context.Context, in *TransferInput, pin string) (*RawTransaction, error) {
	if in.TraceID == "" {
		in.TraceID = uuid.Must(uuid.NewV4()).String()
	}

	paras := map[string]interface{}{
		"asset_id":     in.AssetID,
		"opponent_key": in.OpponentKey,
		"amount":       in.Amount,
		"trace_id":     in.TraceID,
		"memo":         in.Memo,
	}

	var resp RawTransaction
	if err := user.SendRequestWithPIN(ctx, "POST", "/transactions", paras, pin, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TransactionOutput transaction output
type TransactionOutput struct {
	Mask string   `json:"mask"`
	Keys []string `json:"keys"`
}

// MakeTransactionOutput request transaction outputs for receiving assets from main net
func (user *User) MakeTransactionOutput(ctx context.Context, userIDs ...string) (*TransactionOutput, error) {
	if len(userIDs) == 0 {
		userIDs = []string{user.UserID}
	}

	var resp TransactionOutput
	if err := user.SendRequest(ctx, "POST", "/outputs", userIDs, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
