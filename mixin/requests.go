package mixin

import (
	"context"
	"encoding/json"
	"strings"

	mixinsdk "github.com/fox-one/mixin-sdk"
)

func (user *User) paramWithPIN(payload map[string]interface{}, pin string) (map[string]interface{}, error) {
	pinToken, err := user.signPIN(pin)
	if err != nil {
		return nil, err
	}
	payload["pin"] = pinToken
	return payload, nil
}

func (user *User) SendRequest(ctx context.Context, method, uri string, body interface{}, resp interface{}) error {
	ctx = mixinsdk.WithAuth(ctx, user)
	req := mixinsdk.Request(ctx)
	if body != nil {
		req = req.SetBody(body)
	}

	r, err := req.Execute(strings.ToUpper(method), uri)
	if err != nil {
		return err
	}

	if resp != nil {
		err = mixinsdk.UnmarshalResponse(r, resp)
	} else {
		_, err = mixinsdk.DecodeResponse(r)
	}
	return err
}

func (user *User) SendRequestWithPIN(ctx context.Context, method, uri string, body map[string]interface{}, pin string, resp interface{}) error {
	if body == nil {
		body = map[string]interface{}{}
	}

	body, err := user.paramWithPIN(body, pin)
	if err != nil {
		return err
	}

	return user.SendRequest(ctx, method, uri, body, resp)
}

// TODO Deprecated, use SendRequestWithPIN instead
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

// TODO Deprecated, use SendRequest instead
// Request sign and request
func (user *User) Request(ctx context.Context, method, uri string, payload []byte) ([]byte, error) {
	ctx = mixinsdk.WithAuth(ctx, user)
	resp, err := mixinsdk.Request(ctx).SetBody(payload).Execute(method, uri)
	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}
