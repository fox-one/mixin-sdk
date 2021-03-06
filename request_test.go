package mixin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestRequestID(t *testing.T) {
	defer gock.Off()
	// gock.Observe(gock.DumpRequest)

	requestID := "48148640-3B3D-449B-B68A-8899A13EF0a3"

	if c := httpClient.GetClient(); c != nil {
		gock.InterceptClient(c)
		defer gock.RestoreClient(c)
	}

	gock.New(httpClient.HostURL).
		Persist().
		Reply(200).
		SetHeader(requestIDHeaderKey, requestID).
		BodyString("ok")

	ctx := context.Background()
	t.Run("request id mismatch", func(t *testing.T) {
		_, err := Request(ctx).Get("/mismatch")
		if assert.NotNilf(t, err, "request should failed by request id mismatch") {
			assert.Contains(t, err.Error(), requestID)
		}
	})

	t.Run("request id match", func(t *testing.T) {
		r, err := Request(ctx).SetHeader(requestIDHeaderKey, requestID).Get("/match")
		assert.Nilf(t, err, "request should be ok")
		assert.Equalf(t, 200, r.StatusCode(), "status code should be 200")
		assert.Equalf(t, "ok", string(r.Body()), "body should be %q", "ok")
	})

	t.Run("request id in ctx", func(t *testing.T) {
		ctx := WithRequestID(ctx, requestID)
		r, err := Request(ctx).Get("/match")
		assert.Nilf(t, err, "request should be ok")
		assert.Equalf(t, 200, r.StatusCode(), "status code should be 200")
		assert.Equalf(t, "ok", string(r.Body()), "body should be %q", "ok")
	})
}
