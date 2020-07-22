package mixin

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
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
var uploadClient = &http.Client{}

func UploadAttachment(ctx context.Context, attachment *Attachment, file []byte) error {
	req, err := http.NewRequest("PUT", attachment.UploadURL, bytes.NewReader(file))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("x-amz-acl", "public-read")
	req.Header.Add("Content-Length", strconv.Itoa(len(file)))

	resp, err := uploadClient.Do(req)
	if resp != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}

	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}

	return nil
}
