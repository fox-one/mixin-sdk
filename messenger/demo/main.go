package main

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"time"

	"github.com/fox-one/mixin-sdk/messenger"
	"github.com/fox-one/mixin-sdk/mixin"
	"github.com/fox-one/mixin-sdk/utils"
)

type Handler struct {
	*messenger.Messenger
}

func (h Handler) OnMessage(ctx context.Context, msgView messenger.MessageView, userId string) error {
	log.Println("I received a msg", msgView)

	if msgView.Category != messenger.MessageCategoryPlainText {
		return nil
	}

	data, err := base64.StdEncoding.DecodeString(msgView.Data)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("I got your message, you said: %s", string(data))
	log.Println(msg)

	return h.SendPlainText(ctx, msgView, msg)
}

func (h Handler) Run(ctx context.Context) {
	for {
		if err := h.Loop(ctx, h); err != nil {
			log.Println("something is wrong", err)
			time.Sleep(1 * time.Second)
		}
	}
}

func (h Handler) Send(ctx context.Context, userId, content string) error {
	msgView := messenger.MessageView{
		ConversationId: utils.UniqueConversationId(ClientID, userId),
		UserId:         userId,
	}
	return h.SendPlainText(ctx, msgView, content)
}

func main() {
	user := &mixin.User{
		UserID:    ClientID,
		SessionID: SessionID,
		PINToken:  PINToken,
	}

	block, _ := pem.Decode([]byte(SessionKey))
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Panicln(err)
	}
	user.SetPrivateKey(privateKey)

	m := messenger.NewMessenger(user)
	h := Handler{m}
	ctx := context.Background()

	go h.Run(ctx)
	for {
		h.Send(ctx, "7b3f0a95-3ee9-4c1b-8ae9-170e3877d909", "hello world")
		time.Sleep(5 * time.Second)
	}
}
