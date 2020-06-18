package mixin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fox-one/pkg/uuid"
	"github.com/go-resty/resty/v2"
)

const requestIDHeaderKey = "X-Request-ID"

var httpClient = resty.New().
	SetHeader("Content-Type", "application/json").
	SetHostURL("https://api.mixin.one").
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

		if values := r.Header.Values(requestIDHeaderKey); len(values) == 0 {
			r.Header.Set(requestIDHeaderKey, uuid.New())
		}

		return nil
	}).
	OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		return checkResponseRequestID(r)
	})

func Request(ctx context.Context) *resty.Request {
	return httpClient.R().SetContext(ctx)
}

func checkResponseRequestID(r *resty.Response) error {
	expect := r.Request.Header.Get(requestIDHeaderKey)
	got := r.Header().Get(requestIDHeaderKey)
	if expect != "" && got != "" && expect != got {
		return fmt.Errorf("%s mismatch, expect %q but got %q", requestIDHeaderKey, expect, got)
	}

	return nil
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
