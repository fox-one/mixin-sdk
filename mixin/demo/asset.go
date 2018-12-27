package main

import (
	"context"
	"log"

	"github.com/fox-one/mixin-sdk/mixin"
)

func doAssetFee(ctx context.Context, user *mixin.User) {
	assetID := "43d61dcd-e413-450d-80b8-101d5e903357"
	fee, err := user.ReadAssetFee(ctx, assetID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("asset fee", fee)

	if fee.AssetID != assetID {
		log.Panicf("fee asset should be %s but get %s", assetID, fee.AssetID)
	}

	if len(fee.Amount) == 0 {
		log.Panicln("empty fee amount")
	}
}

func validateAsset(asset *mixin.Asset) {
	if len(asset.PublicKey)+len(asset.AccountName) == 0 {
		log.Panicln("empty public key and account name", asset)
	}

	if len(asset.Balance) == 0 {
		log.Panicln("empty balance")
	}
}

func doAsset(ctx context.Context, user *mixin.User) string {
	assetID := "965e5c6e-434c-3fa9-b780-c50f43cd955c"
	asset, err := user.ReadAsset(ctx, assetID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("asset", asset)

	if asset.AssetID != assetID {
		log.Panicf("asset should be %s but get %s\n", assetID, asset.AssetID)
	}

	validateAsset(asset)
	return asset.PublicKey
}

func doAssets(ctx context.Context, user *mixin.User) {
	assets, err := user.ReadAssets(ctx)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("assets", assets)

	for _, asset := range assets {
		validateAsset(asset)
	}
}
