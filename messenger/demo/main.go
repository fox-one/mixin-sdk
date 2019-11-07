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

	return nil
}

func (h Handler) Run(ctx context.Context) {
	for {
		if err := h.Loop(ctx, h); err != nil {
			log.Println("something is wrong", err)
			time.Sleep(1 * time.Second)
		}
	}
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

	h.Run(ctx)
}
