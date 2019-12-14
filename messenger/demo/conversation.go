package main

import (
	"context"
	"log"

	"github.com/fox-one/mixin-sdk/messenger"
	"github.com/fox-one/mixin-sdk/utils"
)

func doCreateConversation(ctx context.Context, m *messenger.Messenger) *messenger.Conversation {
	receiver := "170e40f0-627f-4af2-acf5-0f25c009e523"
	conversationID := utils.UniqueConversationID(m.UserID, receiver)
	participants := []*messenger.Participant{
		&messenger.Participant{
			UserID: receiver,
			Role:   "ADMIN",
		},
	}
	conversation, err := m.CreateConversation(ctx, "GROUP", conversationID, "Haha", "", "", "", participants)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("create conversation", conversation)
	return conversation
}

func doReadConversation(ctx context.Context, m *messenger.Messenger, conversationID string) *messenger.Conversation {
	conversation, err := m.ReadConversation(ctx, conversationID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("read conversation", conversation)
	return conversation
}

func doMessage(ctx context.Context, m *messenger.Messenger, message messenger.Message) {
	err := m.SendMessages(ctx, message)
	if err != nil {
		log.Panicln(err)
	}
}
