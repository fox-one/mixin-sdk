package mixin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fox-one/mixin-sdk/utils"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

// TransferInput input for transfer/verify payment request
type TransferInput struct {
	AddressID  string `json:"address_id,omitempty"`
	AssetID    string `json:"asset_id,omitempty"`
	OpponentID string `json:"opponent_id,omitempty"`
	Amount     string `json:"amount,omitempty"`
	TraceID    string `json:"trace_id,omitempty"`
	Memo       string `json:"memo,omitempty"`
}

func (input TransferInput) verify(snapshot Snapshot) bool {
	mins := time.Now().Sub(snapshot.CreatedAt).Minutes()
	if mins < -10 || mins > 10 {
		log.Debugf("request send %d minutes ago\n", mins)
		return false
	}
	if len(input.OpponentID) > 0 {
		if snapshot.AssetID != input.AssetID {
			log.Debugln("asset id doses not match", snapshot.AssetID, input.AssetID)
			return false
		}
		if snapshot.Data != input.Memo {
			log.Debugln("asset id doses not match", snapshot.AssetID, input.AssetID)
			return false
		}
	}

	iAmount, _ := decimal.NewFromString(input.Amount)
	oAmount, _ := decimal.NewFromString(snapshot.Amount)
	diff := iAmount.Add(oAmount).Truncate(8)
	if !diff.IsZero() {
		log.Debugln("amount does not match", input.Amount, snapshot.Amount, diff.IsZero())
		return false
	}

	return true
}

// VerifyPayment verify payment
//	asset_id, opponent_id, amount, trace_id
func (user User) VerifyPayment(ctx context.Context, input *TransferInput) (bool, error) {
	paras, err := json.Marshal(input)
	if err != nil {
		return false, requestError(err)
	}
	data, err := user.Request(ctx, "POST", "/payments", paras)
	if err != nil {
		return false, requestError(err)
	}

	var resp struct {
		Data *struct {
			User   *User  `json:"receipient"`
			Amount string `json:"amount"`
			Status string `json:"status"`
		} `json:"data,omitempty"`
		Error *Error `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return false, requestError(err)
	} else if resp.Error != nil {
		return false, resp.Error
	}

	if resp.Data.Amount != input.Amount || resp.Data.Status != "paid" {
		return false, nil
	}

	return true, nil
}

// Transfer transfer to account
//	asset_id, opponent_id, amount, traceID, memo
func (user User) Transfer(ctx context.Context, input *TransferInput, pin string) (*Snapshot, *Error) {
	paras := utils.UnselectFields(input)
	data, err := user.RequestWithPIN(ctx, "POST", "/transfers", paras, pin)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		Snapshot *struct {
			*Snapshot
			Memo string `json:"memo,omitempty"`
		} `json:"data,omitempty"`
		Error *Error `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}

	resp.Snapshot.Data = resp.Snapshot.Memo
	if !input.verify(*resp.Snapshot.Snapshot) {
		return nil, traceError()
	}

	return resp.Snapshot.Snapshot, nil
}

// Withdraw withdraw to address
//	address_id, opponent_id, amount, traceID, memo
func (user User) Withdraw(ctx context.Context, input *TransferInput, pin string) (*Snapshot, *Error) {
	paras := utils.UnselectFields(input)
	data, err := user.RequestWithPIN(ctx, "POST", "/withdrawals", paras, pin)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		Snapshot *struct {
			*Snapshot
			Memo string `json:"memo,omitempty"`
		} `json:"data,omitempty"`
		Error *Error `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}

	resp.Snapshot.Data = resp.Snapshot.Memo
	if !input.verify(*resp.Snapshot.Snapshot) {
		return nil, traceError()
	}

	return resp.Snapshot.Snapshot, nil
}
