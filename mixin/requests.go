package mixin

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fox-one/mixin-sdk/utils"
	log "github.com/sirupsen/logrus"
)

var httpClient = &http.Client{}

// RequestWithPIN sign pin and request
func (user *User) RequestWithPIN(ctx context.Context, method, uri string, payload map[string]interface{}, pin string, timeout ...int64) ([]byte, error) {
	if payload == nil {
		payload = map[string]interface{}{}
	}
	pinToken, err := user.signPIN(pin)
	if err != nil {
		return nil, err
	}
	payload["pin"] = pinToken

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return user.Request(ctx, method, uri, data, timeout...)
}

// Request sign and request
func (user *User) Request(ctx context.Context, method, uri string, payload []byte, timeouts ...int64) ([]byte, error) {
	accessToken, err := user.SignToken(method, uri, payload)
	if err != nil {
		return nil, err
	}

	url := "https://api.mixin.one" + uri
	if len(timeouts) > 0 && timeouts[0] > 0 {
		c, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeouts[0]))
		defer cancel()

		ctx = c
	}
	result := utils.SendRequest(ctx, url, method, string(payload), "Content-Type", "application/json", "Authorization", "Bearer "+accessToken)
	code, status := result.Status()
	log.Debugln("do request: ", method, uri, string(payload), code, status)
	return result.Bytes()
}

func (user *User) RequestWithToken(ctx context.Context, method, uri string, payload []byte, accessToken string) ([]byte, error) {
	url := "https://api.mixin.one" + uri
	result := utils.SendRequest(ctx, url, method, string(payload), "Content-Type", "application/json", "Authorization", "Bearer "+accessToken)
	code, status := result.Status()
	log.Debugln("do request: ", method, uri, string(payload), code, status)
	return result.Bytes()
}
