package mixin

import (
	"context"
	"encoding/json"
	"log"
	"time"

	uuid "github.com/gofrs/uuid"
	jsoniter "github.com/json-iterator/go"
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
func (user User) Transaction(ctx context.Context, in *TransferInput, pin string) (*RawTransaction, error) {
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
	data, err := user.RequestWithPIN(ctx, "POST", "/transactions", paras, pin)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		Error Error          `json:"error"`
		Data  RawTransaction `json:"data"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, requestError(err)
	} else if resp.Error.Code > 0 {
		return nil, resp.Error
	}
	return &resp.Data, nil
}

// TransactionOutput transaction output
type TransactionOutput struct {
	Mask string   `json:"mask"`
	Keys []string `json:"keys"`
}

// MakeTransactionOutput request transaction outputs for receiving assets from main net
func (user User) MakeTransactionOutput(ctx context.Context, userIDs ...string) (*TransactionOutput, error) {
	if len(userIDs) == 0 {
		userIDs = []string{user.UserID}
	}
	bts, err := jsoniter.Marshal(userIDs)
	if err != nil {
		return nil, requestError(err)
	}

	data, err := user.Request(ctx, "POST", "/outputs", bts)
	if err != nil {
		return nil, requestError(err)
	}

	log.Println(string(data))

	var resp struct {
		Error Error              `json:"error"`
		Data  *TransactionOutput `json:"data"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, requestError(err)
	} else if resp.Error.Code > 0 {
		return nil, resp.Error
	}
	return resp.Data, nil
}
