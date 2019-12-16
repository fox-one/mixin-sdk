package mixin

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

var httpClient = resty.New().
	SetHeader("Content-Type", "application/json").
	SetHostURL("https://mixin-api.zeromesh.net").
	SetTimeout(10 * time.Second).
	SetPreRequestHook(func(c *resty.Client, r *http.Request) error {
		ctx := r.Context()
		if auth, ok := ctx.Value(authKey).(Authentication); ok {
			token, err := auth.Auth(r)
			if err != nil {
				return err
			}

			r.Header.Set("Authorization", "Bearer "+token)
		}

		return nil
	})

func Request(ctx context.Context) *resty.Request {
	return httpClient.R().SetContext(ctx)
}

func DecodeResponse(resp *resty.Response) ([]byte, error) {
	var body struct {
		Error *Error          `json:"error,omitempty"`
		Data  json.RawMessage `json:"data,omitempty"`
	}

	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		if resp.IsError() {
			return nil, createError(resp.StatusCode(), resp.StatusCode(), resp.Status())
		}

		return nil, createError(resp.StatusCode(), resp.StatusCode(), string(resp.Body()))
	}

	if body.Error != nil && body.Error.Code > 0 {
		return nil, body.Error
	}

	return body.Data, nil
}

func UnmarshalResponse(resp *resty.Response, v interface{}) error {
	data, err := DecodeResponse(resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// User API Request

func (user *User) Request(ctx context.Context, method, uri string, body interface{}, resp interface{}) error {
	ctx = WithAuth(ctx, user)
	req := Request(ctx)
	if body != nil {
		req = req.SetBody(body)
	}

	r, err := req.Execute(strings.ToUpper(method), uri)
	if err != nil {
		return err
	}

	if resp != nil {
		err = UnmarshalResponse(r, resp)
	} else {
		_, err = DecodeResponse(r)
	}
	return err
}

func (user *User) RequestWithPIN(ctx context.Context, method, uri string, body map[string]interface{}, pin string, resp interface{}) error {
	if body == nil {
		body = map[string]interface{}{}
	}

	pinToken, err := user.EncryptPIN(pin)
	if err != nil {
		return err
	}

	body["pin"] = pinToken
	return user.Request(ctx, method, uri, body, resp)
}
