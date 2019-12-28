package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	sdk "github.com/fox-one/mixin-sdk"
)

type Handler struct{}

func (h Handler) OnMessage(ctx context.Context, msgView *sdk.MessageView, userID string) error {
	log.Println("I received a msg", msgView)

	if msgView.Category != sdk.MessageCategoryPlainText {
		return nil
	}

	data, err := base64.StdEncoding.DecodeString(msgView.Data)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("I got your message, you said: %s", string(data))
	log.Println(msg)

	return nil
}

func (h Handler) OnAckReceipt(ctx context.Context, msg *sdk.MessageView, userID string) error {
	return nil
}

func (h Handler) Run(ctx context.Context, user *sdk.User) {
	for {
		select {
		case <-ctx.Done():
			break

		default:
		}
		if err := sdk.NewBlazeClient(user).Loop(ctx, h); err != nil {
			log.Println("something is wrong", err)
			time.Sleep(1 * time.Second)
		}
	}
}
