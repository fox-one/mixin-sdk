package mixin

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

type Message struct {
	ConversationID   string    `json:"conversation_id"`
	RecipientID      string    `json:"recipient_id"`
	MessageID        string    `json:"message_id"`
	QuoteMessageID   string    `json:"quote_message_id,omitempty"`
	Category         string    `json:"category"`
	Data             string    `json:"data"`
	RepresentativeID string    `json:"representative_id,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}

func (user User) SendMessages(ctx context.Context, body []byte) error {
	data, err := user.Request(ctx, "POST", "/messages", body)
	log.Println(string(data), err)
	if err != nil {
		return requestError(err)
	}

	var resp struct {
		Error *Error `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		log.Println(resp, err)
		return requestError(err)
	}
	return nil
}
