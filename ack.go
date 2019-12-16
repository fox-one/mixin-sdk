package mixin

import "context"

type AcknowledgementRequest struct {
	MessageID string `json:"message_id,omitempty"`
	Status    string `json:"status,omitempty"`
}

func (user *User) SendAcknowledgements(ctx context.Context, requests []*AcknowledgementRequest) error {
	if len(requests) == 0 {
		return nil
	}

	return user.Request(ctx, "POST", "/acknowledgements", requests, nil)
}
