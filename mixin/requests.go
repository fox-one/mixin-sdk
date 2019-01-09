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
	req, err := utils.NewRequest(url, method, string(payload), "Content-Type", "application/json", "Authorization", "Bearer "+accessToken)
	if err != nil {
		return nil, err
	}

	log.Debugln("do request: ", method, uri)

	if len(timeouts) > 0 && timeouts[0] > 0 {
		c, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeouts[0]))
		defer cancel()

		ctx = c
	}

	req = req.WithContext(ctx)
	resp, _ := utils.DoRequest(req)
	return utils.ReadResponse(resp)
}
