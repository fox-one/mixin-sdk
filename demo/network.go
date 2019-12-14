package main

import (
	"context"
	"log"

	sdk "github.com/fox-one/mixin-sdk"
)

func doReadNetworkInfo(ctx context.Context) {
	network, err := sdk.ReadNetworkInfo(ctx)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("read network info", network)
}
