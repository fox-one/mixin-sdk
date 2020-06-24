package mixin

import (
	"context"
	"errors"
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

// ShowAttachment show attachment by id
func (user *User) ShowAttachment(ctx context.Context, id string) (*Attachment, error) {
	var attachment Attachment
	if err := user.Request(ctx, "GET", "/attachments/"+id, nil, &attachment); err != nil {
		return nil, err
	}

	return &attachment, nil
}

// Upload upload files
func UploadAttachment(ctx context.Context, attachment *Attachment, file []byte) error {
	resp, err := Request(ctx).SetBody(file).
		SetHeader("Content-Type", "application/octet-stream").
		SetHeader("x-amz-acl", "public-read").
		SetHeader("Content-Length", strconv.Itoa(len(file))).
		Put(attachment.UploadURL)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.New(resp.Status())
	}

	return nil
}
