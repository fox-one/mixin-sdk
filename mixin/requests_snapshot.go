package mixin

import (
	"context"
	"fmt"
	"time"

	mixinsdk "github.com/fox-one/mixin-sdk"
	"github.com/fox-one/mixin-sdk/utils"
)

func readSnapshots(ctx context.Context, network bool, assetID string, offset time.Time, order bool, limit uint) ([]*Snapshot, error) {
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
	if order {
		uri = uri + "&order=ASC"
	} else {
		uri = uri + "&order=DESC"
	}
	resp, err := mixinsdk.Request(ctx).Get(uri)
	if err != nil {
		return nil, err
	}

	var snapshots []*Snapshot
	if err := mixinsdk.UnmarshalResponse(resp, &snapshots); err != nil {
		return nil, err
	}

	for _, snapshot := range snapshots {
		if snapshot.Asset != nil {
			snapshot.AssetID = snapshot.Asset.AssetID
		}
	}

	return snapshots, nil
}

func (user *User) ReadNetwork(ctx context.Context, assetID string, offset time.Time, order bool, limit uint) ([]*Snapshot, error) {
	ctx = mixinsdk.WithAuth(ctx, user)
	return readSnapshots(ctx, true, assetID, offset, order, limit)
}

func ReadNetwork(ctx context.Context, assetID string, offset time.Time, order bool, limit uint, accessToken string) ([]*Snapshot, error) {
	ctx = mixinsdk.WithToken(ctx, accessToken)
	return readSnapshots(ctx, true, assetID, offset, order, limit)
}

func (user *User) ReadSnapshots(ctx context.Context, assetID string, offset time.Time, limit uint) ([]*Snapshot, error) {
	ctx = mixinsdk.WithAuth(ctx, user)
	return readSnapshots(ctx, false, assetID, offset, false, limit)
}

func ReadSnapshots(ctx context.Context, assetID string, offset time.Time, limit uint, accessToken string) ([]*Snapshot, error) {
	ctx = mixinsdk.WithToken(ctx, accessToken)
	return readSnapshots(ctx, false, assetID, offset, false, limit)
}

func readSnapshot(ctx context.Context, network bool, snapshotID string) (*Snapshot, error) {
	uri := "/snapshots/" + snapshotID
	if network {
		uri = "/network/snapshots/" + snapshotID
	}
	resp, err := mixinsdk.Request(ctx).Get(uri)
	if err != nil {
		return nil, err
	}

	var snapshot Snapshot
	if err := mixinsdk.UnmarshalResponse(resp, &snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}

func (user *User) ReadNetworkSnapshot(ctx context.Context, snapshotID string) (*Snapshot, error) {
	ctx = mixinsdk.WithAuth(ctx, user)
	return readSnapshot(ctx, false, snapshotID)
}

func ReadNetworkSnapshot(ctx context.Context, snapshotID, accessToken string) (*Snapshot, error) {
	ctx = mixinsdk.WithToken(ctx, accessToken)
	return readSnapshot(ctx, false, snapshotID)
}

func (user *User) ReadSnapshot(ctx context.Context, snapshotID string) (*Snapshot, error) {
	ctx = mixinsdk.WithAuth(ctx, user)
	return readSnapshot(ctx, true, snapshotID)
}

func ReadSnapshot(ctx context.Context, snapshotID, accessToken string) (*Snapshot, error) {
	ctx = mixinsdk.WithToken(ctx, accessToken)
	return readSnapshot(ctx, true, snapshotID)
}

// ReadTransfer read snapshot with trace id
func (user *User) ReadTransfer(ctx context.Context, traceID string) (*Snapshot, error) {
	ctx = mixinsdk.WithAuth(ctx, user)
	resp, err := mixinsdk.Request(ctx).Get("/transfers/trace/" + traceID)
	if err != nil {
		return nil, err
	}

	var snapshot Snapshot
	if err := mixinsdk.UnmarshalResponse(resp, &snapshot); err != nil {
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

	resp, err := mixinsdk.Request(ctx).Get(uri)
	if err != nil {
		return nil, err
	}

	var snapshots []*DepositTransaction
	if err := mixinsdk.UnmarshalResponse(resp, &snapshots); err != nil {
		return nil, err
	}
	return snapshots, nil
}
