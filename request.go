package mixin_sdk

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty/v2"
)

var httpClient = resty.New().
	SetHeader("Content-Type", "application/json").
	SetHostURL("https://mixin-api.zeromesh.net").
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
