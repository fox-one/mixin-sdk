package mixin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fox-one/mixin-sdk/utils"
)

// ReadNetwork read network snapshots
func (user *User) ReadNetwork(ctx context.Context, assetID string, offset time.Time, order bool, limit uint) ([]*Snapshot, error) {
	uri := fmt.Sprintf("/network/snapshots?limit=%d", limit)
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
	data, err := user.Request(ctx, "GET", uri, nil)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		Snapshots []*Snapshot `json:"data,omitempty"`
		Error     *Error      `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}

	for _, snapshot := range resp.Snapshots {
		if snapshot.Asset != nil {
			snapshot.AssetID = snapshot.Asset.AssetID
		}
	}

	return resp.Snapshots, nil
}

// ReadNetworkSnapshot read snapshot with snapshot id
func (user *User) ReadNetworkSnapshot(ctx context.Context, snapshotID string) (*Snapshot, error) {
	data, err := user.Request(ctx, "GET", "/network/snapshots/"+snapshotID, nil)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		Snapshot *Snapshot `json:"data,omitempty"`
		Error    *Error    `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Snapshot, nil
}

// ReadSnapshot read snapshot with snapshot id
func (user *User) ReadSnapshot(ctx context.Context, snapshotID string) (*Snapshot, error) {
	data, err := user.Request(ctx, "GET", "/snapshots/"+snapshotID, nil)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		Snapshot *Snapshot `json:"data,omitempty"`
		Error    *Error    `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Snapshot, nil
}

// ReadTransfer read snapshot with trace id
func (user *User) ReadTransfer(ctx context.Context, traceID string) (*Snapshot, error) {
	data, err := user.Request(ctx, "GET", "/transfers/trace/"+traceID, nil)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		Snapshot *Snapshot `json:"data,omitempty"`
		Error    *Error    `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Snapshot, nil
}

// ReadExternal read external snapshots
func (user *User) ReadExternal(ctx context.Context, assetID, destination, tag string, offset time.Time, limit int) ([]*DepositTransaction, error) {
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
		return nil, requestError(err)
	}
	data, err := user.Request(ctx, "GET", uri, nil)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		Snapshots []*DepositTransaction `json:"data,omitempty"`
		Error     *Error                `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Snapshots, nil
}
