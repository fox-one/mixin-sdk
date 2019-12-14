package main

import (
	"context"
	"log"

	sdk "github.com/fox-one/mixin-sdk"
)

func doUpload(ctx context.Context, user *sdk.User) {
	attatchmentID, viewURL, err := user.Upload(ctx, []byte("hahaaaaaaaaa"))
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("attatchment id: %s; view url: %s\n", attatchmentID, viewURL)
}
