package main

import (
	"context"
	"log"

	sdk "github.com/fox-one/mixin-sdk"
	"github.com/gofrs/uuid"
)

func doTransaction(ctx context.Context, user *sdk.User, assetID, opponentKey, amount, memo, pin string) {
	snapshot, err := user.Transaction(ctx, &sdk.TransferInput{
		TraceID:     uuid.Must(uuid.NewV4()).String(),
		AssetID:     assetID,
		OpponentKey: opponentKey,
		Amount:      amount,
		Memo:        memo,
	}, pin)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("do transfer", snapshot)
}

func doTransfer(ctx context.Context, user *sdk.User, assetID, opponentID, amount, memo, pin string) *sdk.Snapshot {
	snapshot, err := user.Transfer(ctx, &sdk.TransferInput{
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

func doWithdraw(ctx context.Context, user *sdk.User, assetID, publicKey, amount, memo, pin string) *sdk.Snapshot {
	addrID := doCreateAddress(ctx, user, assetID, publicKey, "Test Withdraw", pin)

	snapshot, err := user.Withdraw(ctx, &sdk.TransferInput{
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
