package main

import (
	"context"
	"log"

	"github.com/fox-one/mixin-sdk/mixin"
)

func doAddress(ctx context.Context, user *mixin.User, assetID, publicKey, label, pin string) {
	addrID := doCreateAddress(ctx, user, assetID, publicKey, label, pin)

	doAssetAddresses(ctx, user, assetID)

	doDeleteAddress(ctx, user, addrID, pin)
}

func doCreateAddress(ctx context.Context, user *mixin.User, assetID, publicKey, label, pin string) string {
	addr, err := user.CreateWithdrawAddress(ctx, mixin.WithdrawAddress{
		AssetID:   assetID,
		PublicKey: publicKey,
		Label:     label,
	}, pin)
	if err != nil {
		log.Panicln(err)
	}

	printJSON("create withdraw address", addr)
	return addr.AddressID
}

func doAssetAddresses(ctx context.Context, user *mixin.User, assetID string) {
	addrs, err := user.ReadWithdrawAddresses(ctx, assetID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("read withdraw addresses", addrs)
}

func doDeleteAddress(ctx context.Context, user *mixin.User, addressID, pin string) {
	err := user.DeleteWithdrawAddress(ctx, addressID, pin)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("delete withdraw address")
}
