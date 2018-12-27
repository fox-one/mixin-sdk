package main

import (
	"context"
	"log"

	"github.com/fox-one/mixin-sdk/mixin"
	"github.com/satori/go.uuid"
)

func doTransfer(ctx context.Context, user *mixin.User, assetID, opponentID, amount, memo, pin string) *mixin.Snapshot {
	snapshot, err := user.Transfer(ctx, &mixin.TransferInput{
		TraceID:    uuid.Must(uuid.NewV4()).String(),
		AssetID:    assetID,
		OpponentID: opponentID,
		Amount:     amount,
		Memo:       memo,
	}, pin)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("do transfer", snapshot)
	return snapshot
}

func doWithdraw(ctx context.Context, user *mixin.User, assetID, publicKey, amount, memo, pin string) *mixin.Snapshot {
	addrID := doCreateAddress(ctx, user, assetID, publicKey, "Test Withdraw", pin)

	snapshot, err := user.Withdraw(ctx, &mixin.TransferInput{
		TraceID:   uuid.Must(uuid.NewV4()).String(),
		AssetID:   assetID,
		AddressID: addrID,
		Amount:    amount,
		Memo:      memo,
	}, pin)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("do withdraw", snapshot)

	doDeleteAddress(ctx, user, addrID, pin)
	return snapshot
}
