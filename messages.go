package mixin

import (
	"context"
)

const (
	MessageCategoryPlainText             = "PLAIN_TEXT"
	MessageCategoryPlainImage            = "PLAIN_IMAGE"
	MessageCategoryPlainData             = "PLAIN_DATA"
	MessageCategoryPlainSticker          = "PLAIN_STICKER"
	MessageCategoryPlainLive             = "PLAIN_LIVE"
	MessageCategoryPlainContact          = "PLAIN_CONTACT"
	MessageCategorySystemConversation    = "SYSTEM_CONVERSATION"
	MessageCategorySystemAccountSnapshot = "SYSTEM_ACCOUNT_SNAPSHOT"
	MessageCategoryMessageRecall         = "MESSAGE_RECALL"
	MessageCategoryAppButtonGroup        = "APP_BUTTON_GROUP"
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

func (user *User) SendMessages(ctx context.Context, messages []*MessageRequest) error {
	switch len(messages) {
	case 0:
		return nil
	case 1:
		return user.SendMessage(ctx, messages[0])
	default:
		return user.Request(ctx, "POST", "/messages", messages, nil)
	}
}

func (user *User) SendMessage(ctx context.Context, message *MessageRequest) error {
	return user.Request(ctx, "POST", "/messages", message, nil)
}
