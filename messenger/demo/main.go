package main

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"log"

	"github.com/fox-one/mixin-sdk/messenger"
	"github.com/fox-one/mixin-sdk/mixin"
	"github.com/gofrs/uuid"
)

func printJSON(prefix string, item interface{}) {
	msg, _ := json.MarshalIndent(item, "", "    ")
	log.Println(prefix, string(msg))
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
	ctx := context.Background()

	conversation := doCreateConversation(ctx, m)
	doMessage(ctx, m, messenger.Message{
		ConversationID: conversation.ConversationID,
		MessageID:      uuid.Must(uuid.NewV4()).String(),
		Category:       "PLAIN_TEXT",
		Data:           base64.StdEncoding.EncodeToString([]byte("Just A Test")),
	})
	doReadConversation(ctx, m, conversation.ConversationID)

	// doUpload(ctx, m)

	Handler{m}.Run(ctx)
}
