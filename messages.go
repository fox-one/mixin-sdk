package sdk

import (
	"context"
)

type LiveMessagePayload struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	ThumbUrl string `json:"thumb_url"`
	Url      string `json:"url"`
}

type RecallMessagePayload struct {
	MessageID string `json:"message_id"`
}

type MessageRequest struct {
	ConversationID   string `json:"conversation_id"`
	RecipientID      string `json:"recipient_id"`
	MessageID        string `json:"message_id"`
	Category         string `json:"category"`
	Data             string `json:"data"`
	RepresentativeID string `json:"representative_id"`
	QuoteMessageID   string `json:"quote_message_id"`
}

type AcknowledgementRequest struct {
	MessageID string `json:"message_id,omitempty"`
	Status    string `json:"status,omitempty"`
}

func (user *User) SendMessages(ctx context.Context, messages ...MessageRequest) error {
	if len(messages) == 0 {
		return nil
	}

	var paras interface{} = messages
	if len(messages) == 1 {
		paras = messages[0]
	}
	return user.Request(ctx, "POST", "/messages", paras, nil)
}

func (user *User) SendAcknowledgements(ctx context.Context, requests []*AcknowledgementRequest) error {
	if len(requests) == 0 {
		return nil
	}
	return user.Request(ctx, "POST", "/acknowledgements", requests, nil)
}
