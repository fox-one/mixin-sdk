package mixin

import (
	"context"

	"github.com/fox-one/mixin-sdk/utils"
	uuid "github.com/gofrs/uuid"
)

const (
	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"
)

type (
	VerifyPaymentInput struct {
		AssetID    string `json:"asset_id,omitempty"`
		OpponentID string `json:"opponent_id,omitempty"`
		Amount     string `json:"amount,omitempty"`
		TraceID    string `json:"trace_id,omitempty"`
	}

	VerifyPaymentResponse struct {
		Recipient *Profile `json:"recipient,omitempty"`
		Asset     *Asset   `json:"asset,omitempty"`
		Amount    string   `json:"amount,omitempty"`
		Statue    string   `json:"statue,omitempty"`
	}
)

// VerifyPayment verify payment
//	asset_id, opponent_id, amount, trace_id
func (user *User) VerifyPayment(ctx context.Context, input *VerifyPaymentInput) (*VerifyPaymentResponse, error) {
	var resp VerifyPaymentResponse
	err := user.Request(ctx, "POST", "/payments", input, &resp)
	return &resp, err
}

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
	if err := user.RequestWithPIN(ctx, "POST", "/transfers", utils.UnselectFields(input), pin, &resp); err != nil {
		return nil, err
	}

	resp.Snapshot.Data = resp.Memo
	return resp.Snapshot, nil
}

// Withdraw withdraw to address
//	address_id, opponent_id, amount, traceID, memo
func (user *User) Withdraw(ctx context.Context, input *TransferInput, pin string) (*Snapshot, error) {
	if len(input.TraceID) == 0 {
		input.TraceID = uuid.Must(uuid.NewV4()).String()
	}

	var resp struct {
		*Snapshot
		Memo string `json:"memo,omitempty"`
	}
	if err := user.RequestWithPIN(ctx, "POST", "/withdrawals", utils.UnselectFields(input), pin, &resp); err != nil {
		return nil, err
	}

	resp.Snapshot.Data = resp.Memo
	return resp.Snapshot, nil
}
