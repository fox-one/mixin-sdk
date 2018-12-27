package main

import (
	"context"
	"log"
	"time"

	"github.com/fox-one/mixin-sdk/mixin"
)

func doReadNetwork(ctx context.Context, user *mixin.User) {
	snapshots, err := user.ReadNetwork(ctx, "965e5c6e-434c-3fa9-b780-c50f43cd955c", time.Time{}, false, 10)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("read network", snapshots)
}

func doReadSnapshot(ctx context.Context, user *mixin.User) {
	snapshot, err := user.ReadSnapshot(ctx, "cb329d89-def3-4ddc-b3eb-d5450199eff4")
	if err != nil {
		log.Panicln(err)
	}
	printJSON("read snapshot", snapshot)
}

func doReadTransfer(ctx context.Context, user *mixin.User) {
	snapshot, err := user.ReadTransfer(ctx, "01663326-4cfd-407a-997b-33f297c716e2")
	if err != nil {
		log.Panicln(err)
	}
	printJSON("read transfer", snapshot)
}
