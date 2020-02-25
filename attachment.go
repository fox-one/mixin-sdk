package mixin

import (
	"context"
	"strconv"
)

// Attachment attachment
type Attachment struct {
	AttachmentID string `json:"attachment_id"`
	UploadURL    string `json:"upload_url"`
	ViewURL      string `json:"view_url"`
}

// CreateAttachment create attachment
func (user *User) CreateAttachment(ctx context.Context) (*Attachment, error) {
	var attachment Attachment
	if err := user.Request(ctx, "POST", "/attachments", nil, &attachment); err != nil {
		return nil, err
	}
	return &attachment, nil
}

// Upload upload files
func (user *User) Upload(ctx context.Context, file []byte) (string, string, error) {
	attachment, err := user.CreateAttachment(ctx)
	if err != nil {
		return "", "", err
	}

	resp, err := Request(ctx).SetBody(file).
		SetHeader("Content-Type", "application/octet-stream").
		SetHeader("x-amz-acl", "public-read").
		SetHeader("Content-Length", strconv.Itoa(len(file))).
		Put(attachment.UploadURL)
	if err != nil {
		return "", "", err
	}

	if _, err := DecodeResponse(resp); err != nil {
		return "", "", err
	}
	return attachment.AttachmentID, attachment.ViewURL, nil
}
