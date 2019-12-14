package messenger

import (
	"context"
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
	if len(messages) == 0 {
		return nil
	}

	var paras interface{} = messages
	if len(messages) == 1 {
		paras = messages[0]
	}
	return m.SendRequest(ctx, "POST", "/messages", paras, nil)
}
