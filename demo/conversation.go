package main

import (
	"context"
	"log"

	sdk "github.com/fox-one/mixin-sdk"
	"github.com/fox-one/mixin-sdk/utils"
)

func doCreateConversation(ctx context.Context, user *sdk.User) *sdk.Conversation {
	receiver := "170e40f0-627f-4af2-acf5-0f25c009e523"
	conversationID := utils.UniqueConversationID(user.UserID, receiver)
	participants := []*sdk.Participant{
		&sdk.Participant{
			UserID: receiver,
			Role:   sdk.ParticipantRoleAdmin,
		},
	}
	conversation, err := user.CreateGroupConversation(ctx, conversationID, "haha", participants)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("create conversation", conversation)
	return conversation
}

func doReadConversation(ctx context.Context, user *sdk.User, conversationID string) *sdk.Conversation {
	conversation, err := user.ReadConversation(ctx, conversationID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("read conversation", conversation)
	return conversation
}

func doMessage(ctx context.Context, user *sdk.User, message *sdk.MessageRequest) {
	err := user.SendMessage(ctx, message)
	if err != nil {
		log.Panicln(err)
	}
}
