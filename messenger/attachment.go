package messenger

import (
	"context"
	"encoding/json"

	"github.com/fox-one/mixin-sdk/utils"

	"github.com/fox-one/mixin-sdk/mixin"
)

// Attachment attachment
type Attachment struct {
	AttachmentID string `json:"attachment"`
	UploadURL    string `json:"upload_url"`
	ViewURL      string `json:"view_url"`
}

// CreateAttachment create attachment
func (m Messenger) CreateAttachment(ctx context.Context) (*Attachment, error) {
	data, err := m.Request(ctx, "POST", "/attachment", nil)
	if err != nil {
		return nil, requestError(err)
	}

	var resp struct {
		Attachment *Attachment  `json:"data"`
		Error      *mixin.Error `json:"error,omitempty"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, requestError(err)
	} else if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Attachment, nil
}

// Upload upload files
func (m Messenger) Upload(ctx context.Context, file []byte) (string, string, error) {
	attachment, err := m.CreateAttachment(ctx)
	if err != nil {
		return "", "", err
	}

	req, err := utils.NewRequest(attachment.UploadURL, "PUT", string(file), "x-amz-acl", "public-read")
	if err != nil {
		return "", "", err
	}

	_, err = utils.DoRequest(req)
	if err != nil {
		return "", "", err
	}
	return attachment.AttachmentID, attachment.ViewURL, nil
}
