package messenger

import (
	"context"

	mixinsdk "github.com/fox-one/mixin-sdk"
)

// Attachment attachment
type Attachment struct {
	AttachmentID string `json:"attachment"`
	UploadURL    string `json:"upload_url"`
	ViewURL      string `json:"view_url"`
}

// CreateAttachment create attachment
func (m Messenger) CreateAttachment(ctx context.Context) (*Attachment, error) {
	var attachment Attachment
	if err := m.SendRequest(ctx, "POST", "/attachments", nil, &attachment); err != nil {
		return nil, err
	}
	return &attachment, nil
}

// Upload upload files
func (m Messenger) Upload(ctx context.Context, file []byte) (string, string, error) {
	attachment, err := m.CreateAttachment(ctx)
	if err != nil {
		return "", "", err
	}

	resp, err := mixinsdk.Request(ctx).SetBody(string(file)).
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("x-amz-acl", "public-read").
		Put(attachment.UploadURL)
	if err != nil {
		return "", "", err
	}

	if _, err := mixinsdk.DecodeResponse(resp); err != nil {
		return "", "", err
	}
	return attachment.AttachmentID, attachment.ViewURL, nil
}
