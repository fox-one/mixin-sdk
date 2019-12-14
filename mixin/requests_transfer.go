package mixin

import (
	"context"
	"strings"

	"github.com/fox-one/mixin-sdk/utils"
	uuid "github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

// TransferInput input for transfer/verify payment request
type TransferInput struct {
	AddressID   string `json:"address_id,omitempty"`
	AssetID     string `json:"asset_id,omitempty"`
	OpponentID  string `json:"opponent_id,omitempty"`
	Amount      string `json:"amount,omitempty"`
	TraceID     string `json:"trace_id,omitempty"`
	Memo        string `json:"memo,omitempty"`
	OpponentKey string `json:"opponent_key,omitempty"`
}

func (input TransferInput) verify(snapshot Snapshot) bool {
	// transfer
	if len(input.OpponentID) > 0 {
		if snapshot.AssetID != input.AssetID {
			log.Debugln("asset id doses not match", snapshot.AssetID, input.AssetID)
			return false
		}
		if snapshot.Data != input.Memo {
			log.Debugln("asset id doses not match", snapshot.AssetID, input.AssetID)
			return false
		}

		iAmount, _ := decimal.NewFromString(input.Amount)
		oAmount, _ := decimal.NewFromString(snapshot.Amount)
		diff := iAmount.Add(oAmount).Truncate(8)
		if !diff.IsZero() {
			log.Debugln("amount does not match", input.Amount, snapshot.Amount, diff.IsZero())
			return false
		}
	}

	// withdraw
	// TODO EOS withdraw will Truncate(4)

	return true
}

// VerifyPayment verify payment
//	asset_id, opponent_id, amount, trace_id
func (user *User) VerifyPayment(ctx context.Context, input *TransferInput) (bool, error) {
	var resp struct {
		User   *User  `json:"receipient"`
		Amount string `json:"amount"`
		Status string `json:"status"`
	}
	if err := user.SendRequest(ctx, "POST", "/payments", input, &resp); err != nil {
		return false, err
	}

	if resp.Amount != input.Amount || strings.ToLower(resp.Status) != "paid" {
		return false, nil
	}

	return true, nil
}

// Transfer transfer to account
//	asset_id, opponent_id, amount, traceID, memo
func (user *User) Transfer(ctx context.Context, input *TransferInput, pin string) (*Snapshot, error) {
	if len(input.TraceID) == 0 {
		input.TraceID = uuid.Must(uuid.NewV4()).String()
	}

	var resp struct {
		*Snapshot
		Memo string `json:"memo,omitempty"`
	}
	if err := user.SendRequestWithPIN(ctx, "POST", "/transfers", utils.UnselectFields(input), pin, &resp); err != nil {
		return nil, err
	}

	resp.Snapshot.Data = resp.Memo
	if !input.verify(*resp.Snapshot) {
		return nil, traceError()
	}

	return resp.Snapshot, nil
}

// Withdraw withdraw to address
//	address_id, opponent_id, amount, traceID, memo
func (user User) Withdraw(ctx context.Context, input *TransferInput, pin string) (*Snapshot, error) {
	if len(input.TraceID) == 0 {
		input.TraceID = uuid.Must(uuid.NewV4()).String()
	}

	var resp struct {
		*Snapshot
		Memo string `json:"memo,omitempty"`
	}
	if err := user.SendRequestWithPIN(ctx, "POST", "/withdrawals", utils.UnselectFields(input), pin, &resp); err != nil {
		return nil, err
	}

	resp.Snapshot.Data = resp.Memo
	if !input.verify(*resp.Snapshot) {
		return nil, traceError()
	}

	return resp.Snapshot, nil
}
