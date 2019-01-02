package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fox-one/mixin-sdk/mixin"
)

var httpClient *http.Client

func (m Messenger) CreateAttachment(ctx context.Context) (string, string, string, error) {
	data, err := m.Request(ctx, "POST", "/attachment", nil)
	if err != nil {
		return "", "", "", requestError(err)
	}

	var resp struct {
		Data struct {
			AttachmentId string `json:"attachment"`
			UploadUrl    string `json:"upload_url"`
			ViewUrl      string `json:"view_url"`
		} `json:"data"`
		Error *mixin.Error `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return "", "", "", requestError(err)
	} else if resp.Error != nil {
		return "", "", "", resp.Error
	}

	return resp.Data.AttachmentId, resp.Data.UploadUrl, resp.Data.ViewUrl, nil
}

func (m Messenger) Upload(ctx context.Context, file []byte) (string, string, error) {
	id, upload, view, err := m.CreateAttachment(ctx)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("PUT", upload, bytes.NewReader(file))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("x-amz-acl", "public-read")

	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	_, err = httpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	return id, view, nil
}
