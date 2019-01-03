package mixin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/fox-one/mixin-sdk/utils"
)

var httpClient = &http.Client{}

// RequestWithPIN sign pin and request
func (user *User) RequestWithPIN(ctx context.Context, method, uri string, payload map[string]interface{}, pin string) ([]byte, error) {
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

	return user.Request(ctx, method, uri, data)
}

// Request sign and request
func (user *User) Request(ctx context.Context, method, uri string, payload []byte) ([]byte, error) {
	accessToken, err := user.SignToken(method, uri, payload)
	if err != nil {
		return nil, err
	}

	url := "https://api.mixin.one" + uri
	req, err := utils.NewRequest(url, method, string(payload), "Content-Type", "application/json", "Authorization", "Bearer "+accessToken)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	resp, _ := utils.DoRequest(req)
	return utils.ReadResponse(resp)
}
