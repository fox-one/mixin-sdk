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

// TransactionOutput transaction output
type TransactionOutput struct {
	Mask string   `json:"mask"`
	Keys []string `json:"keys"`
}

// API

// Transaction do transaction to mixin main net
func (user *User) Transaction(ctx context.Context, in *TransferInput, pin string) (*RawTransaction, error) {
	if in.TraceID == "" {
		in.TraceID = uuid.Must(uuid.NewV4()).String()
	}

	paras := map[string]interface{}{
		"asset_id": in.AssetID,
		"amount":   in.Amount,
		"trace_id": in.TraceID,
		"memo":     in.Memo,
	}

	if in.OpponentKey != "" {
		paras["opponent_key"] = in.OpponentKey
	} else {
		paras["opponent_multisig"] = map[string]interface{}{
			"receivers": in.OpponentMultisig.Receivers,
			"threshold": in.OpponentMultisig.Threshold,
		}
	}

	var resp RawTransaction
	if err := user.RequestWithPIN(ctx, "POST", "/transactions", paras, pin, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MakeTransactionOutput request transaction outputs for receiving assets from main net
func (user *User) MakeTransactionOutput(ctx context.Context, userIDs ...string) (*TransactionOutput, error) {
	if len(userIDs) == 0 {
		userIDs = []string{user.UserID}
	}

	input := map[string]interface{}{
		"receivers": userIDs,
	}

	var resp TransactionOutput
	if err := user.Request(ctx, "POST", "/outputs", input, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
