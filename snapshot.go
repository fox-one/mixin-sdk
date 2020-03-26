package mixin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fox-one/mixin-sdk/utils"
	"github.com/shopspring/decimal"
)

// Snapshot transfer records
type Snapshot struct {
	SnapshotID string `json:"snapshot_id"`
	TraceID    string `json:"trace_id,omitempty"`

	UserID     string    `json:"user_id,omitempty"`
	AssetID    string    `json:"asset_id,omitempty"`
	ChainID    string    `json:"chain_id,omitempty"`
	OpponentID string    `json:"opponent_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`

	Source         string          `json:"source"` // Source DEPOSIT_CONFIRMED, TRANSFER_INITIALIZED, WITHDRAWAL_INITIALIZED, WITHDRAWAL_FEE_CHARGED, WITHDRAWAL_FAILED
	Amount         string          `json:"amount"`
	OpeningBalance decimal.Decimal `json:"opening_balance"`
	ClosingBalance decimal.Decimal `json:"closing_balance"`
	Data           string          `json:"data,omitempty"`
	Type           string          `json:"type,omitempty"`

	Sender          string `json:"sender,omitempty"`
	Receiver        string `json:"receiver,omitempty"`
	TransactionHash string `json:"transaction_hash,omitempty"`

	Asset *Asset `json:"asset,omitempty"`
}

func (snapshot *Snapshot) Memo() string {
	return snapshot.Data
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
	Destination string `json:"destination"`
	Tag         string `json:"tag"`

	// TODO Deprecated
	PublicKey   string `json:"public_key"`
	AccountName string `json:"account_name"`
	AccountTag  string `json:"account_tag"`
}

const (
	OrderDESC = "DESC"
	OrderASC  = "ASC"
)

func readSnapshots(ctx context.Context, network bool, assetID string, offset time.Time, order string, limit uint) ([]*Snapshot, error) {
	uri := fmt.Sprintf("/snapshots?limit=%d", limit)
	if network {
		uri = fmt.Sprintf("/network/snapshots?limit=%d", limit)
	}
	if !offset.IsZero() {
		uri = uri + "&offset=" + offset.UTC().Format(time.RFC3339Nano)
	}
	if len(assetID) > 0 {
		uri = uri + "&asset=" + assetID
	}

	switch order {
	case OrderASC:
	case OrderDESC, "":
		order = OrderDESC
	default:
		return nil, errors.New("order must be ASC or DESC")
	}

	uri = uri + "&order=" + order

	resp, err := Request(ctx).Get(uri)
	if err != nil {
		return nil, err
	}

	var snapshots []*Snapshot
	if err := UnmarshalResponse(resp, &snapshots); err != nil {
		return nil, err
	}

	for _, snapshot := range snapshots {
		if snapshot.Asset != nil {
			snapshot.AssetID = snapshot.Asset.AssetID
		}
	}

	return snapshots, nil
}

func (user *User) ReadNetwork(ctx context.Context, assetID string, offset time.Time, order string, limit uint) ([]*Snapshot, error) {
	ctx = WithAuth(ctx, user)
	return readSnapshots(ctx, true, assetID, offset, order, limit)
}

func ReadNetwork(ctx context.Context, assetID string, offset time.Time, order string, limit uint, accessToken string) ([]*Snapshot, error) {
	ctx = WithToken(ctx, accessToken)
	return readSnapshots(ctx, true, assetID, offset, order, limit)
}

func (user *User) ReadSnapshots(ctx context.Context, assetID string, offset time.Time, order string, limit uint) ([]*Snapshot, error) {
	ctx = WithAuth(ctx, user)
	return readSnapshots(ctx, false, assetID, offset, order, limit)
}

func ReadSnapshots(ctx context.Context, assetID string, offset time.Time, order string, limit uint, accessToken string) ([]*Snapshot, error) {
	ctx = WithToken(ctx, accessToken)
	return readSnapshots(ctx, false, assetID, offset, order, limit)
}

func readSnapshot(ctx context.Context, network bool, snapshotID string) (*Snapshot, error) {
	uri := "/snapshots/" + snapshotID
	if network {
		uri = "/network/snapshots/" + snapshotID
	}
	resp, err := Request(ctx).Get(uri)
	if err != nil {
		return nil, err
	}

	var snapshot Snapshot
	if err := UnmarshalResponse(resp, &snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}

func (user *User) ReadNetworkSnapshot(ctx context.Context, snapshotID string) (*Snapshot, error) {
	ctx = WithAuth(ctx, user)
	return readSnapshot(ctx, true, snapshotID)
}

func ReadNetworkSnapshot(ctx context.Context, snapshotID, accessToken string) (*Snapshot, error) {
	ctx = WithToken(ctx, accessToken)
	return readSnapshot(ctx, true, snapshotID)
}

func (user *User) ReadSnapshot(ctx context.Context, snapshotID string) (*Snapshot, error) {
	ctx = WithAuth(ctx, user)
	return readSnapshot(ctx, false, snapshotID)
}

func ReadSnapshot(ctx context.Context, snapshotID, accessToken string) (*Snapshot, error) {
	ctx = WithToken(ctx, accessToken)
	return readSnapshot(ctx, false, snapshotID)
}

// ReadTransfer read snapshot with trace id
func (user *User) ReadTransfer(ctx context.Context, traceID string) (*Snapshot, error) {
	var snapshot Snapshot
	if err := user.Request(ctx, "GET", "/transfers/trace/"+traceID, nil, &snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}

// read external snapshots
func ReadExternal(ctx context.Context, assetID, destination, tag string, offset time.Time, limit int) ([]*DepositTransaction, error) {
	var paras = make([]string, 0, 12)
	if len(assetID) > 0 {
		paras = append(paras, "asset", assetID)
	}
	if destination != "" {
		paras = append(paras, "destination", destination)
		if tag != "" {
			paras = append(paras, "tag", tag)
		}
	}
	if !offset.IsZero() {
		paras = append(paras, "offset", offset.Format(time.RFC3339Nano))
	}
	if limit > 0 {
		paras = append(paras, "limit", fmt.Sprint(limit))
	}
	uri, err := utils.BuildURL("/external/transactions", paras...)
	if err != nil {
		return nil, err
	}

	resp, err := Request(ctx).Get(uri)
	if err != nil {
		return nil, err
	}

	var snapshots []*DepositTransaction
	if err := UnmarshalResponse(resp, &snapshots); err != nil {
		return nil, err
	}
	return snapshots, nil
}
