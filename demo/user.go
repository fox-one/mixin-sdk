package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"log"

	sdk "github.com/fox-one/mixin-sdk"
)

func doCreateUser(ctx context.Context, user *sdk.User, pin string) *sdk.User {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Panicln(err)
	}

	u, e := user.CreateUser(ctx, privateKey, "test tset")
	if e != nil {
		log.Panicln(e)
	}
	printJSON("created user", u)

	doModifyPIN(ctx, u, "", pin)

	doVerifyPIN(ctx, u, pin)

	return u
}

func doModifyPIN(ctx context.Context, user *sdk.User, oldPIN, pin string) {
	err := user.ModifyPIN(ctx, oldPIN, pin)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("modify PIN succ")
}

func doVerifyPIN(ctx context.Context, user *sdk.User, pin string) {
	err := user.VerifyPIN(ctx, pin)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("verify PIN succ")
}
