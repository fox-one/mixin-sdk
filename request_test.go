package mixin

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestRequestID(t *testing.T) {
	defer gock.Off()
	// gock.Observe(gock.DumpRequest)

	requestID := "random request id"

	if c := httpClient.GetClient(); c != nil {
		gock.InterceptClient(c)
		defer gock.RestoreClient(c)
	}

	gock.New(httpClient.HostURL).
		Persist().
		Reply(202).
		SetHeader(requestIDHeaderKey, requestID).
		JSON(json.RawMessage(`{"data":{"foo":"bar"}}`))

	t.Run("request id mismatch", func(t *testing.T) {
		_, err := Request(context.Background()).Get("/mismatch")
		if assert.NotNilf(t, err, "request should failed by request id mismatch") {
			assert.Contains(t, err.Error(), requestID)
		}
	})

	t.Run("request id match", func(t *testing.T) {
		_, err := Request(context.Background()).SetHeader(requestIDHeaderKey, requestID).Get("/match")
		assert.Nilf(t, err, "request should be ok")
	})
}
