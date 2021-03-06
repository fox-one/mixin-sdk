package main

import (
	"context"
	"encoding/base64"
	"log"
	"time"

	sdk "github.com/fox-one/mixin-sdk"
	"github.com/gofrs/uuid"
	jsoniter "github.com/json-iterator/go"
)

func printJSON(prefix string, item interface{}) {
	msg, _ := jsoniter.MarshalIndent(item, "", "    ")
	log.Println(prefix, string(msg))
}

func main() {
	user, err := sdk.NewUser(ClientID, SessionID, SessionKey, PINToken)
	if err != nil {
		log.Panicln(err)
	}

	ctx := context.Background()
	publicKey := doAsset(ctx, user)

	p := "123456"
	u := doCreateUser(ctx, user, p)

	doAssetFee(ctx, u)

	publicKey1 := doAsset(ctx, u)

	doAssets(ctx, u)

	assetID := "965e5c6e-434c-3fa9-b780-c50f43cd955c"
	doTransfer(ctx, user, assetID, u.UserID, "0.0001", "ping", PIN)
	time.Sleep(time.Second * 5)
	snap := doTransfer(ctx, u, assetID, user.UserID, "0.0001", "pong", p)

	doWithdraw(ctx, user, assetID, publicKey1, "0.01", "ping", PIN)
	time.Sleep(time.Second * 5)
	doWithdraw(ctx, u, assetID, publicKey, "0.01", "pong", p)

	doReadNetwork(ctx)
	doUserReadNetwork(ctx, u)
	doReadSnapshots(ctx, user)

	doReadSnapshot(ctx, u, snap.SnapshotID)

	doReadTransfer(ctx, u, snap.TraceID)

	doReadExternal(ctx)

	doReadNetworkInfo(ctx)

	doTransaction(ctx, user, "965e5c6e-434c-3fa9-b780-c50f43cd955c", "XINT55hZYxzrtqJsWViUbyoxytJ6RoKUZfpnSCQTbgX8fjcdQ7GwjRySLxiPMWxAMhoN6KPa7SFkyv9FQXC3fGJuKHLf3est", "1", "test", PIN)

	// Messenger

	conversation := doCreateConversation(ctx, user)
	doMessage(ctx, user, &sdk.MessageRequest{
		ConversationID: conversation.ConversationID,
		MessageID:      uuid.Must(uuid.NewV4()).String(),
		Category:       "PLAIN_TEXT",
		Data:           base64.StdEncoding.EncodeToString([]byte("Just A Test")),
	})
	doReadConversation(ctx, user, conversation.ConversationID)

	Handler{}.Run(ctx, user)
}
