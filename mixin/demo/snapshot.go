package main

import (
	"context"
	"log"
	"time"

	"github.com/fox-one/mixin-sdk/mixin"
)

func doReadNetwork(ctx context.Context) {
	snapshots, err := mixin.ReadNetwork(ctx, "965e5c6e-434c-3fa9-b780-c50f43cd955c", time.Time{}, false, 10, "")
	if err != nil {
		log.Panicln(err)
	}
	printJSON("read network", snapshots)
}

func doUserReadNetwork(ctx context.Context, user *mixin.User) {
	snapshots, err := user.ReadNetwork(ctx, "965e5c6e-434c-3fa9-b780-c50f43cd955c", time.Time{}, false, 10)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("read network", snapshots)
}

func doReadSnapshots(ctx context.Context, user *mixin.User) {
	snapshots, err := user.ReadSnapshots(ctx, "965e5c6e-434c-3fa9-b780-c50f43cd955c", time.Time{}, 10)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("read snapshots", snapshots)
}

func doReadSnapshot(ctx context.Context, user *mixin.User, snapshotID string) {
	snapshot, err := user.ReadSnapshot(ctx, snapshotID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("read snapshot", snapshot)
}

func doReadTransfer(ctx context.Context, user *mixin.User, traceID string) {
	snapshot, err := user.ReadTransfer(ctx, traceID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("read transfer", snapshot)
}

func doReadExternal(ctx context.Context) {
	snapshots, err := mixin.ReadExternal(ctx, "", "", "", time.Time{}, 10)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("read deposit snapshots", snapshots)
}
