package main

import (
	"context"
	"log"

	"github.com/fox-one/mixin-sdk/messenger"
)

func doUpload(ctx context.Context, m *messenger.Messenger) {
	attatchmentID, viewURL, err := m.Upload(ctx, []byte("hahaaaaaaaaa"))
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("attatchment id: %s; view url: %s\n", attatchmentID, viewURL)
}
