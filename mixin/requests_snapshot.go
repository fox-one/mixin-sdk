package mixin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fox-one/mixin-sdk/utils"
)

// ReadNetwork read network snapshots
func (user User) ReadNetwork(ctx context.Context, assetID string, offset time.Time, order bool, limit uint) ([]*Snapshot, *Error) {
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
	return resp.Snapshots, nil
}

// ReadSnapshot read snapshot with snapshot id
func (user User) ReadSnapshot(ctx context.Context, snapshotID string) (*Snapshot, *Error) {
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

// ReadTransfer read snapshot with trace id
func (user User) ReadTransfer(ctx context.Context, traceID string) (*Snapshot, *Error) {
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
func (user User) ReadExternal(ctx context.Context, assetID, publicKey, accountName, accountTag string, offset time.Time, limit int) ([]*PendingSnapshot, *Error) {
	var paras = make([]string, 0, 12)
	if len(assetID) > 0 {
		paras = append(paras, "asset", assetID)
	}
	if len(publicKey) > 0 {
		paras = append(paras, "public_key", publicKey)
	} else if len(accountName) > 0 {
		paras = append(paras, "account_name", accountName)
		if len(accountName) > 0 {
			paras = append(paras, "account_tag", accountTag)
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
		Snapshots []*PendingSnapshot `json:"data,omitempty"`
		Error     *Error             `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Snapshots, nil
}
