package utils

import (
	"crypto/md5"
	"io"
	"strings"

	"github.com/gofrs/uuid"
)

// UUIDWithString generate uuid with text
func UUIDWithString(text string) string {
	h := md5.New()
	io.WriteString(h, text)
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	return uuid.FromBytesOrNil(sum).String()
}

// UniqueConversationID sort and generate uuid
func UniqueConversationID(userID, recipientID string) string {
	minID, maxID := userID, recipientID
	if strings.Compare(userID, recipientID) > 0 {
		maxID, minID = userID, recipientID
	}
	return UUIDWithString(minID + maxID)
}
