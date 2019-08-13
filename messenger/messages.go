package messenger

import (
	"context"
	"encoding/json"

	"github.com/fox-one/mixin-sdk/mixin"
)

type Message struct {
	ConversationID   string `json:"conversation_id"`
	RecipientID      string `json:"recipient_id,omitempty"`
	MessageID        string `json:"message_id"`
	QuoteMessageID   string `json:"quote_message_id,omitempty"`
	Category         string `json:"category"`
	Data             string `json:"data"`
	RepresentativeID string `json:"representative_id,omitempty"`
}

func (m Messenger) SendMessages(ctx context.Context, messages ...Message) error {
	var body []byte
	var err error
	switch len(messages) {
	case 0:
		return nil
	case 1:
		body, err = json.Marshal(messages[0])
	default:
		body, err = json.Marshal(messages)
	}

	if err != nil {
		return err
	}

	data, err := m.Request(ctx, "POST", "/messages", body)
	if err != nil {
		return requestError(err)
	}

	var resp struct {
		Error mixin.Error `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return requestError(err)
	}
	if resp.Error.Code != 0 {
		return resp.Error
	}
	return nil
}
