package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"log"
	"time"

	"github.com/fox-one/mixin-sdk/mixin"
	jsoniter "github.com/json-iterator/go"
)

func printJSON(prefix string, item interface{}) {
	msg, err := jsoniter.MarshalToString(item)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(prefix, msg)
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

	ctx := context.Background()
	publicKey := doAsset(ctx, user)

	p := "123456"
	u := doCreateUser(ctx, user, p)

	doAssetFee(ctx, u)

	publicKey1 := doAsset(ctx, u)

	doAssets(ctx, u)

	assetID := "965e5c6e-434c-3fa9-b780-c50f43cd955c"
	doTransfer(ctx, user, assetID, u.UserID, "0.1", "ping", PIN)
	time.Sleep(time.Second * 5)
	snap := doTransfer(ctx, u, assetID, user.UserID, "0.1", "pong", p)

	doWithdraw(ctx, user, assetID, publicKey1, "0.1", "ping", PIN)
	time.Sleep(time.Second * 5)
	doWithdraw(ctx, u, assetID, publicKey, "0.1", "pong", p)

	doReadNetwork(ctx, u)

	doReadSnapshot(ctx, u, snap.SnapshotID)

	doReadTransfer(ctx, u, snap.TraceID)

	doReadExternal(ctx, u)

	doReadNetworkInfo(ctx, user)

	doTransaction(ctx, user, "965e5c6e-434c-3fa9-b780-c50f43cd955c", "XINT55hZYxzrtqJsWViUbyoxytJ6RoKUZfpnSCQTbgX8fjcdQ7GwjRySLxiPMWxAMhoN6KPa7SFkyv9FQXC3fGJuKHLf3est", "1", "test", PIN)
}
